package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"git.1750studios.com/GSoC/CrashDragon/config"
	"git.1750studios.com/GSoC/CrashDragon/database"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// PostReports processes crashreport
func PostReports(c *gin.Context) {
	file, _, err := c.Request.FormFile("upload_file_minidump")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer file.Close()
	var Report database.Report
	Report.Processed = false
	Report.ID = uuid.NewV4()
	if err = database.Db.Create(&Report).Error; err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	Report.Product = c.Request.FormValue("prod")
	Report.Version = c.Request.FormValue("ver")
	Report.ProcessUptime, _ = strconv.Atoi(c.Request.FormValue("ptime"))
	Report.EMail = c.Request.FormValue("email")
	Report.Comment = c.Request.FormValue("comments")
	filepath := path.Join(config.C.ContentDirectory, "Reports", Report.ID.String()[0:2], Report.ID.String()[0:4])
	err = os.MkdirAll(filepath, 0755)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		database.Db.Delete(&Report)
		return
	}
	f, err := os.Create(path.Join(filepath, Report.ID.String()+".dmp"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		database.Db.Delete(&Report)
		return
	}
	io.Copy(f, file)
	f.Close()
	go processReport(Report, false)
	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"object": Report,
	})
	return
}

// ReprocessReport processes the Crashreport again with current symbols
func ReprocessReport(c *gin.Context) {
	var Report database.Report
	database.Db.Where("id = ?", c.Param("id")).First(&Report)
	processReport(Report, true)
	c.SetCookie("result", "OK", 0, "/", "", false, false)
	c.Redirect(http.StatusMovedPermanently, "/crashreports/"+Report.ID.String())
	return
}

func processReport(Report database.Report, reprocess bool) {
	file := path.Join(config.C.ContentDirectory, "Reports", Report.ID.String()[0:2], Report.ID.String()[0:4], Report.ID.String()+".dmp")
	cmdJSON := exec.Command("./minidump-stackwalk/stackwalker", file, path.Join(config.C.ContentDirectory, "Symfiles"))
	var out bytes.Buffer
	cmdJSON.Stdout = &out
	err := cmdJSON.Run()
	if err != nil {
		os.Remove(file)
		database.Db.Delete(&Report)
		return
	}
	err = json.Unmarshal(out.Bytes(), &Report.Report)
	if err != nil {
		os.Remove(file)
		database.Db.Delete(&Report)
		return
	}
	if Report.Report.Status != "OK" {
		os.Remove(file)
		database.Db.Delete(&Report)
		return
	}
	cmdTXT := exec.Command("./minidump-stackwalk/stackwalk/bin/minidump_stackwalk", file, path.Join(config.C.ContentDirectory, "Symfiles"))
	out.Reset()
	cmdTXT.Stdout = &out
	err = cmdTXT.Run()
	if err != nil {
		os.Remove(file)
		database.Db.Delete(&Report)
		return
	}
	Report.ReportContentTXT = out.String()
	Report.Processed = true
	Report.Os = Report.Report.SystemInfo.Os
	Report.OsVersion = Report.Report.SystemInfo.OsVer
	Report.Arch = Report.Report.SystemInfo.CPUArch
	for _, Frame := range Report.Report.CrashingThread.Frames {
		if Frame.File == "" && Report.Signature != "" {
			continue
		}
		Report.Signature = Frame.Function
		if Frame.File == "" {
			continue
		}
		Report.CrashLocation = path.Base(Frame.File) + ":" + strconv.Itoa(Frame.Line)
		break
	}
	if err = database.Db.Save(&Report).Error; err != nil {
		os.Remove(file)
		database.Db.Delete(&Report)
		return
	}
	var signature string
	for _, frame := range Report.Report.CrashingThread.Frames {
		if frame.Function == "" {
			continue
		} else {
			signature = frame.Function
			break
		}
	}
	if signature == "" {
		return
	}

	var Crash database.Crash
	if reprocess && Report.CrashID != uuid.Nil {
		database.Db.First(&Crash, "id = ?", Report.CrashID)
		Crash.Signature = signature
		database.Db.Save(&Crash)
	} else {
		database.Db.FirstOrInit(&Crash, "signature = ?", signature)
	}
	if Crash.ID == uuid.Nil {
		Crash.ID = uuid.NewV4()

		Crash.FirstReported = Report.CreatedAt
		Crash.Signature = signature

		Crash.AllCrashCount = 0
		Crash.WinCrashCount = 0
		Crash.MacCrashCount = 0
		Crash.LinCrashCount = 0

		database.Db.Create(&Crash)
		reprocess = false
	}
	if !reprocess || Report.CrashID == uuid.Nil {
		Crash.LastReported = Report.CreatedAt
		Crash.AllCrashCount++
		if Report.Os == "Windows" {
			Crash.WinCrashCount++
		} else if Report.Os == "Linux" {
			Crash.LinCrashCount++
		} else if Report.Os == "Mac OS X" {
			Crash.MacCrashCount++
		}
		database.Db.Save(&Crash)
	}
	Report.CrashID = Crash.ID
	database.Db.Save(&Report)
}

// PostSymfiles processes symfile
func PostSymfiles(c *gin.Context) {
	file, _, err := c.Request.FormFile("symfile")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer file.Close()
	var Symfile database.Symfile
	Symfile.Product = c.Request.FormValue("prod")
	if Symfile.Product == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("prod required"))
		return
	}
	Symfile.Version = c.Request.FormValue("ver")
	if Symfile.Version == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("ver required"))
		return
	}
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	if err = scanner.Err(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	parts := strings.Split(scanner.Text(), " ")
	if parts[0] != "MODULE" {
		c.AbortWithError(http.StatusUnprocessableEntity, errors.New("Sym-file does not start with 'MODULE'!"))
		return
	}
	updated := true
	if err = database.Db.Where("code = ?", parts[3]).First(&Symfile).Error; err != nil {
		Symfile.ID = uuid.NewV4()
		updated = false
	}
	Symfile.Os = parts[1]
	Symfile.Arch = parts[2]
	Symfile.Code = parts[3]
	Symfile.Name = parts[4]
	filepath := path.Join(config.C.ContentDirectory, "Symfiles", Symfile.Name, Symfile.Code)
	err = os.MkdirAll(filepath, 0755)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		database.Db.Delete(&Symfile)
		return
	}
	if _, existsErr := os.Stat(path.Join(filepath, Symfile.Name+".sym")); !os.IsNotExist(existsErr) {
		err = os.Remove(path.Join(filepath, Symfile.Name+".sym"))
	}
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		database.Db.Delete(&Symfile)
		return
	}
	f, err := os.Create(path.Join(filepath, Symfile.Name+".sym"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		database.Db.Delete(&Symfile)
		return
	}
	defer f.Close()
	file.Seek(0, 0)
	io.Copy(f, file)
	if updated {
		if err = database.Db.Save(&Symfile).Error; err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			database.Db.Delete(&Symfile)
			os.Remove(f.Name())
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"object": Symfile,
		})
		return
	}
	if err = database.Db.Create(&Symfile).Error; err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		os.Remove(f.Name())
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"object": Symfile,
	})
	return
}

// Auth middleware which checks the Authorization header field and looks up the user in the database
func Auth(c *gin.Context) {
	var auth string
	var user string
	// FIXME: Change the Header workaround to use the native gin function once it is stable
	//c.GetHeader("Authorization") //gin gonic develop branch
	// Workaround
	if auths, _ := c.Request.Header["Authorization"]; len(auths) > 0 {
		auth = auths[0]
	} else {
		c.Writer.Header().Set("WWW-Authenticate", "Basic realm=\"CrashDragon\"")
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	// End of workaround
	if strings.HasPrefix(auth, "Basic ") {
		base := strings.Split(auth, " ")[1]
		userpass, _ := base64.StdEncoding.DecodeString(base)
		user = strings.Split(string(userpass), ":")[0]
	}
	var User database.User
	database.Db.FirstOrInit(&User, "name = ?", user)
	if User.ID == uuid.Nil {
		User.ID = uuid.NewV4()
		User.IsAdmin = false
		User.Name = user
		database.Db.Create(&User)
	}
	c.Set("user", User)
	c.Next()
}
