package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"

	"git.1750studios.com/GSoC/CrashDragon/config"
	"git.1750studios.com/GSoC/CrashDragon/database"
	"github.com/gin-gonic/gin"
)

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
	database.Db.Order("created_at DESC").Find(&Reports)
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
		Item.Platform = Report.Report.SystemInfo.Os
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
	}
	Item.ID = Report.ID.String()
	Item.Date = Report.CreatedAt.Format("2006-01-02 15:04:05")
	Item.Product = Report.Product
	Item.Version = Report.Version
	Item.Platform = Report.Report.SystemInfo.Os + " " + Report.Report.SystemInfo.OsVer
	Item.Arch = Report.Report.SystemInfo.CPUArch
	Item.Processor = Report.Report.SystemInfo.CPUInfo + " (" + strconv.Itoa(Report.Report.SystemInfo.CPUCount) + " cores)"
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
	c.HTML(http.StatusOK, "crashreport.html", gin.H{
		"title":  "Crashreport",
		"item":   Item,
		"report": Report.Report,
	})
}

// GetCrashreportFile returns minidump file of crashreport
func GetCrashreportFile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"title": "Crashreports",
		"id":    c.Param("id"),
		"name":  c.Param("name"),
	})
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

// GetSymfile returns details of symfile
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
