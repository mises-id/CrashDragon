package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"code.videolan.org/videolan/CrashDragon/config"
	"code.videolan.org/videolan/CrashDragon/database"
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
	if err = database.Db.FirstOrInit(&Report).Error; err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var Product database.Product
	if err = database.Db.First(&Product, "slug = ?", c.Request.FormValue("prod")).Error; err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("the specified prod does not exist"))
		return
	}
	Report.ProductID = Product.ID
	var Version database.Version
	if err = database.Db.First(&Version, "slug = ? AND product_id = ? AND ignore = false", c.Request.FormValue("ver"), Report.ProductID).Error; err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("the specified ver does not exist or is ignored"))
		return
	}
	Report.VersionID = Version.ID
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
	c.Redirect(http.StatusMovedPermanently, "/reports/"+Report.ID.String())
	return
}

func processReport(Report database.Report, reprocess bool) {
	file := path.Join(config.C.ContentDirectory, "Reports", Report.ID.String()[0:2], Report.ID.String()[0:4], Report.ID.String()+".dmp")
	cmdJSON := exec.Command("./build/bin/minidump_stackwalk", "-f", "json", file, path.Join(config.C.ContentDirectory, "Symfiles"))
	var out bytes.Buffer
	cmdJSON.Stdout = &out
	err := cmdJSON.Run()
	if err != nil {
		os.Remove(file)
		database.Db.Delete(&Report)
		return
	}
	Report.Report = database.ReportContent{}
	err = json.Unmarshal(out.Bytes(), &Report.Report)
	if err != nil {
		os.Remove(file)
		database.Db.Delete(&Report)
		return
	}
	if Report.Report.Status != "OK" {
		Report.Processed = false
	} else {
		Report.Processed = true
	}
	cmdTXT := exec.Command("./build/bin/minidump_stackwalk", "-f", "txt", file, path.Join(config.C.ContentDirectory, "Symfiles"))
	out.Reset()
	cmdTXT.Stdout = &out
	err = cmdTXT.Run()
	if err != nil {
		os.Remove(file)
		database.Db.Delete(&Report)
		return
	}
	Report.ReportContentTXT = out.String()
	Report.Os = Report.Report.SystemInfo.Os
	Report.OsVersion = Report.Report.SystemInfo.OsVer
	Report.Arch = Report.Report.SystemInfo.CPUArch

	if reprocess {
		Report.Signature = ""
		Report.CrashLocation = ""
		Report.CrashPath = ""
		Report.CrashLine = 0
	}
	for _, Frame := range Report.Report.CrashingThread.Frames {
		if Frame.File == "" && Report.Signature != "" {
			continue
		}
		Report.Signature = Frame.Function
		if Frame.File == "" {
			continue
		}
		Report.CrashLocation = Frame.File + ":" + strconv.Itoa(Frame.Line)
		Report.CrashPath = Frame.File
		Report.CrashLine = Frame.Line
		break
	}

	if !reprocess {
		Report.CreatedAt = time.Now()
	}

	var Crash database.Crash
	if reprocess && Report.CrashID != uuid.Nil {
		database.Db.First(&Crash, "id = ?", Report.CrashID)
		Crash.Signature = Report.Signature
		database.Db.Save(&Crash)
	} else {
		database.Db.FirstOrInit(&Crash, "signature = ?", Report.Signature)
	}
	if Crash.ID == uuid.Nil {
		Crash.ID = uuid.NewV4()

		Crash.FirstReported = Report.CreatedAt
		Crash.Signature = Report.Signature

		Crash.AllCrashCount = 0
		Crash.WinCrashCount = 0
		Crash.MacCrashCount = 0
		Crash.LinCrashCount = 0

		Crash.ProductID = Report.ProductID
		Crash.VersionID = Report.VersionID

		Crash.Fixed = false

		database.Db.Create(&Crash)
		reprocess = false
	}
	database.Db.Model(&Crash).Related(&Crash.Reports)
	for _, CReport := range Crash.Reports {
		if CReport.VersionID == Report.VersionID {
			break
		}
		Crash.Fixed = false
	}
	if !reprocess || Report.CrashID == uuid.Nil {
		Crash.LastReported = Report.CreatedAt
		Crash.AllCrashCount++
		if Report.Os == "Windows NT" {
			Crash.WinCrashCount++
		} else if Report.Os == "Linux" {
			Crash.LinCrashCount++
		} else if Report.Os == "Mac OS X" {
			Crash.MacCrashCount++
		}
		database.Db.Save(&Crash)
	}
	Report.CrashID = Crash.ID
	if reprocess {
		database.Db.Save(&Report)
	} else {
		database.Db.Create(&Report)
	}
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
	var Product database.Product
	if err = database.Db.First(&Product, "slug = ?", c.Request.FormValue("prod")).Error; err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("the specified prod does not exist"))
		return
	}
	Symfile.ProductID = Product.ID
	var Version database.Version
	if err = database.Db.First(&Version, "slug = ? AND product_id = ?", c.Request.FormValue("ver"), Symfile.ProductID).Error; err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("the specified ver does not exist"))
		return
	}
	Symfile.VersionID = Version.ID
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	if err = scanner.Err(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	parts := strings.Split(scanner.Text(), " ")
	if parts[0] != "MODULE" {
		c.AbortWithError(http.StatusUnprocessableEntity, errors.New("sym-file does not start with 'MODULE'"))
		return
	}
	if parts[3] == "000000000000000000000000000000000" {
		c.AbortWithError(http.StatusUnprocessableEntity, errors.New("sym-file has invalid code"))
		return
	}
	updated := true
	if err = database.Db.Where("code = ?", parts[3]).First(&Symfile).Error; err != nil {
		Symfile.ID = uuid.NewV4()
		updated = false
	} else {
		filepath := path.Join(config.C.ContentDirectory, "Symfiles", Symfile.Name, Symfile.Code)
		if _, existsErr := os.Stat(path.Join(filepath, Symfile.Name+".sym")); !os.IsNotExist(existsErr) {
			err = os.Remove(path.Join(filepath, Symfile.Name+".sym"))
		}
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			database.Db.Delete(&Symfile)
			return
		}
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
