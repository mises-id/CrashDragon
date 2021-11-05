package web

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"code.videolan.org/videolan/CrashDragon/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// GetSymfiles returns symfiles
func GetSymfiles(c *gin.Context) {
	var Symfiles []database.Symfile
	query := database.DB
	prod, ver := GetCookies(c)
	if prod != nil {
		query = query.Where("product_id = ?", prod.ID)
	}
	if ver != nil {
		query = query.Where("version_id = ?", ver.ID)
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		offset = 0
	}
	var count int64
	query.Model(database.Symfile{}).Count(&count)
	query.Order("created_at ASC").Offset(offset).Limit(50).Preload("Product").Preload("Version").Find(&Symfiles)
	var next int
	var prev int
	if (int64(offset) + 50) >= count {
		next = -1
	} else {
		next = offset + 50
	}
	prev = offset - 50
	if strings.HasPrefix(c.Request.Header.Get("Accept"), "text/html") {
		c.HTML(http.StatusOK, "symfiles.html", gin.H{
			"prods":      database.Products,
			"vers":       database.Versions,
			"title":      "Symfiles",
			"items":      Symfiles,
			"nextOffset": next,
			"prevOffset": prev,
		})
	} else {
		c.JSON(http.StatusOK, Symfiles)
	}
}

// GetSymfile returns content of symfile
func GetSymfile(c *gin.Context) {
	var Symfile database.Symfile
	database.DB.Where("id = ?", c.Param("id")).Preload("Product").Preload("Version").First(&Symfile)
	if !strings.HasPrefix(c.Request.Header.Get("Accept"), "text/html") {
		c.JSON(http.StatusOK, Symfile)
		return
	}
	f, err := os.Open(filepath.Join(viper.GetString("Directory.Content"), "Symfiles", Symfile.Product.Slug, Symfile.Version.Slug, Symfile.Name, Symfile.Code, Symfile.Name+".sym"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.Data(http.StatusOK, "text/plain", data)
}
