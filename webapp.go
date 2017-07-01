package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"

	"git.1750studios.com/GSoC/CrashDragon/config"
	"git.1750studios.com/GSoC/CrashDragon/database"
	"github.com/gin-gonic/gin"
)

// GetCrashes returns crashes
func GetCrashes(c *gin.Context) {
	var Crashes []database.Crash
	database.Db.Order("all_crash_count DESC").Find(&Crashes)
	c.HTML(http.StatusOK, "crashes.html", gin.H{
		"title": "Crashes",
		"items": Crashes,
	})
}

// GetCrash returns details of a crash
func GetCrash(c *gin.Context) {
	var Crash database.Crash
	var Crashreports []database.Crashreport
	database.Db.First(&Crash, "id = ?", c.Param("id")).Order("created_at DESC").Related(&Crashreports)
	c.HTML(http.StatusOK, "crash.html", gin.H{
		"title": "Crash",
		"items": Crashreports,
	})
}

// GetCrashreports returns crashreports
func GetCrashreports(c *gin.Context) {
	var Reports []database.Crashreport
	var List []struct {
		ID        string
		Signature string
		Date      string
		Product   string
		Version   string
		Platform  string
		Reason    string
		Location  string
	}
	database.Db.Where("processed = true").Order("created_at DESC").Find(&Reports)
	for _, Report := range Reports {
		var Item struct {
			ID        string
			Signature string
			Date      string
			Product   string
			Version   string
			Platform  string
			Reason    string
			Location  string
		}
		Item.ID = Report.ID.String()
		Item.Date = Report.CreatedAt.Format("2006-01-02 15:04:05")
		Item.Product = Report.Product
		Item.Version = Report.Version
		Item.Platform = Report.Os
		Item.Reason = Report.Report.CrashInfo.Type
		for _, Frame := range Report.Report.CrashingThread.Frames {
			if Frame.File == "" && Item.Signature != "" {
				continue
			}
			Item.Signature = Frame.Function
			if Frame.File == "" {
				continue
			}
			Item.Location = path.Base(Frame.File) + ":" + strconv.Itoa(Frame.Line)
			break
		}
		List = append(List, Item)
	}
	c.HTML(http.StatusOK, "crashreports.html", gin.H{
		"title": "Crashreports",
		"items": List,
	})
}

// GetCrashreport returns details of crashreport
func GetCrashreport(c *gin.Context) {
	var Report database.Crashreport
	database.Db.Where("id = ?", c.Param("id")).First(&Report)
	var Item struct {
		ID        string
		Signature string
		Date      string
		Product   string
		Version   string
		Platform  string
		Arch      string
		Processor string
		Reason    string
		Location  string
		Comment   string
		Uptime    string
	}
	Item.ID = Report.ID.String()
	Item.Date = Report.CreatedAt.Format("2006-01-02 15:04:05")
	Item.Product = Report.Product
	Item.Version = Report.Version
	Item.Platform = Report.Os + " " + Report.OsVersion
	Item.Arch = Report.Arch
	Item.Processor = Report.Report.SystemInfo.CPUInfo + " (" + strconv.Itoa(Report.Report.SystemInfo.CPUCount) + " cores)"
	Item.Reason = Report.Report.CrashInfo.Type
	Item.Comment = Report.Comment
	h := (Report.ProcessUptime / 3600000) % 24
	m := (Report.ProcessUptime / 60000) % 60
	s := (Report.ProcessUptime / 1000) % 60
	Item.Uptime = fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	for _, Frame := range Report.Report.CrashingThread.Frames {
		if Frame.File == "" && Item.Signature != "" {
			continue
		}
		Item.Signature = Frame.Function
		if Frame.File == "" {
			continue
		}
		Item.Location = path.Base(Frame.File) + ":" + strconv.Itoa(Frame.Line)
		break
	}
	result, _ := c.Cookie("result")
	if result != "" {
		c.SetCookie("result", "", 1, "/", "", false, false)
	}
	c.HTML(http.StatusOK, "crashreport.html", gin.H{
		"title":  "Crashreport",
		"item":   Item,
		"report": Report.Report,
		"result": result,
	})
}

// GetCrashreportFile returns minidump file of crashreport
func GetCrashreportFile(c *gin.Context) {
	var Crashreport database.Crashreport
	if err := database.Db.Where("id = ?", c.Param("id")).First(&Crashreport).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	name := c.Param("name")
	switch name {
	case "upload_file_minidump":
		file := path.Join(config.C.ContentDirectory, "Crashreports", Crashreport.ID.String()[0:2], Crashreport.ID.String()[0:4], Crashreport.ID.String()+".dmp")
		f, err := os.Open(file)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		data, err := ioutil.ReadAll(f)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Header("Content-Disposition", "attachment; filename=\""+Crashreport.ID.String()+".dmp\"")
		c.Data(http.StatusOK, "application/octet-stream", data)
		return
	case "processed_json":
		c.Header("Content-Disposition", "attachment; filename=\""+Crashreport.ID.String()+".json\"")
		c.Data(http.StatusOK, "application/json", []byte(Crashreport.ReportContentJSON))
		return
	case "processed_txt":
		c.Header("Content-Disposition", "attachment; filename=\""+Crashreport.ID.String()+".txt\"")
		c.Data(http.StatusOK, "text/plain", []byte(Crashreport.ReportContentTXT))
		return
	default:
		c.AbortWithError(http.StatusBadRequest, errors.New(name+" is a unknwon file"))
		return
	}
}

// GetSymfiles returns symfiles
func GetSymfiles(c *gin.Context) {
	var Symfiles []database.Symfile
	database.Db.Find(&Symfiles)
	c.HTML(http.StatusOK, "symfiles.html", gin.H{
		"title": "Symfiles",
		"items": Symfiles,
	})
}

// GetSymfile returns content of symfile
func GetSymfile(c *gin.Context) {
	var Symfile database.Symfile
	database.Db.Where("id = ?", c.Param("id")).First(&Symfile)
	f, err := os.Open(path.Join(config.C.ContentDirectory, "Symfiles", Symfile.Name, Symfile.Code, Symfile.Name+".sym"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error,
		})
		return
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error,
		})
		return
	}
	c.Data(http.StatusOK, "text/plain", data)
}
