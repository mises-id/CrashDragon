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

// GetSymfiles returns symfiles
func GetSymfiles(c *gin.Context) {
	var Symfiles []database.Symfile
	query := database.Db
	all, prod := GetProductCookie(c)
	if !all {
		query = query.Where("product_id = ?", prod.ID)
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		offset = 0
	}
	var count int
	query.Model(database.Symfile{}).Count(&count)
	query.Order("created_at ASC").Offset(offset).Limit(50).Preload("Product").Preload("Version").Find(&Symfiles)
	var next int
	var prev int
	if (offset + 50) >= count {
		next = -1
	} else {
		next = offset + 50
	}
	prev = offset - 50
	c.HTML(http.StatusOK, "symfiles.html", gin.H{
		"prods":      database.Products,
		"title":      "Symfiles",
		"items":      Symfiles,
		"nextOffset": next,
		"prevOffset": prev,
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
