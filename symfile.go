package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"git.1750studios.com/GSoC/CrashDragon/config"
	"git.1750studios.com/GSoC/CrashDragon/database"
	"github.com/gin-gonic/gin"
)

// GetSymfiles returns symfiles
func GetSymfiles(c *gin.Context) {
	var Symfiles []database.Symfile
	lastDate := c.DefaultQuery("last_date", time.Time{}.Format(time.UnixDate))
	dir := c.DefaultQuery("dir", "up")
	if lastDate == "" || lastDate == config.NilDate {
		database.Db.Order("created_at ASC").Limit(50).Find(&Symfiles)
	} else if dir == "up" {
		database.Db.Where("created_at > ?", lastDate).Order("created_at ASC").Limit(50).Find(&Symfiles)
	} else {
		database.Db.Where("created_at < ?", lastDate).Order("created_at ASC").Limit(50).Find(&Symfiles)
	}
	var nextDate string
	var prevDate string
	if len(Symfiles) > 0 {
		nextDate = Symfiles[len(Symfiles)-1].CreatedAt.Format(time.UnixDate)
		prevDate = Symfiles[0].CreatedAt.Format(time.UnixDate)
	}
	c.HTML(http.StatusOK, "symfiles.html", gin.H{
		"title":    "Symfiles",
		"items":    Symfiles,
		"nextDate": nextDate,
		"prevDate": prevDate,
	})
}

// GetSymfile returns content of symfile
func GetSymfile(c *gin.Context) {
	var Symfile database.Symfile
	database.Db.Where("id = ?", c.Param("id")).First(&Symfile)
	f, err := os.Open(path.Join(config.C.ContentDirectory, "Symfiles", Symfile.Name, Symfile.Code, Symfile.Name+".sym"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.Data(http.StatusOK, "text/plain", data)
}
