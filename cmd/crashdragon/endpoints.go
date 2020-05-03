package main

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"code.videolan.org/videolan/CrashDragon/internal/config"
	"code.videolan.org/videolan/CrashDragon/internal/database"
	"code.videolan.org/videolan/CrashDragon/internal/processor"
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

	var Product database.Product
	if err = database.Db.First(&Product, "slug = ?", c.Request.FormValue("prod")).Error; err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("the specified prod does not exist"))
		return
	}
	Report.ProductID = Product.ID
	Report.Product = Product

	var Version database.Version
	if err = database.Db.First(&Version, "slug = ? AND product_id = ? AND ignore = false", c.Request.FormValue("ver"), Report.ProductID).Error; err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("the specified ver does not exist or is ignored"))
		return
	}
	Report.VersionID = Version.ID
	Report.Version = Version

	Report.ProcessUptime, _ = strconv.Atoi(c.Request.FormValue("ptime"))
	Report.EMail = c.Request.FormValue("email")
	Report.Comment = c.Request.FormValue("comments")
	filepth := filepath.Join(config.C.ContentDirectory, "Reports", Report.ID.String()[0:2], Report.ID.String()[0:4])
	err = os.MkdirAll(filepth, 0755)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	f, err := os.Create(filepath.Join(filepth, Report.ID.String()+".dmp"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	processor.AddToQueue(Report)
	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"object": Report.ID,
	})
	return
}

// ReprocessReport processes the Crashreport again with current symbols
func ReprocessReport(c *gin.Context) {
	var Report database.Report
	database.Db.Where("id = ?", c.Param("id")).First(&Report)
	processor.Reprocess(Report)
	c.SetCookie("result", "OK", 0, "/", "", false, false)
	c.Redirect(http.StatusMovedPermanently, "/reports/"+Report.ID.String())
	return
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
	Symfile.Product = Product
	var Version database.Version
	if err = database.Db.First(&Version, "slug = ? AND product_id = ?", c.Request.FormValue("ver"), Symfile.ProductID).Error; err != nil {
		Version.ID = uuid.NewV4()
		Version.Name = c.Request.FormValue("ver")
		Version.Slug = c.Request.FormValue("ver")
		Version.Ignore = false
		Version.Product = Product
		Version.ProductID = Product.ID
		database.Db.Create(&Version)
	}
	Symfile.VersionID = Version.ID
	Symfile.Version = Version
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
		filepth := filepath.Join(config.C.ContentDirectory, "Symfiles", Symfile.Product.Slug, Symfile.Version.Slug, Symfile.Name, Symfile.Code)
		if _, existsErr := os.Stat(filepath.Join(filepth, Symfile.Name+".sym")); !os.IsNotExist(existsErr) {
			err = os.Remove(filepath.Join(filepth, Symfile.Name+".sym"))
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
	filepth := filepath.Join(config.C.ContentDirectory, "Symfiles", Symfile.Product.Slug, Symfile.Version.Slug, Symfile.Name, Symfile.Code)
	err = os.MkdirAll(filepth, 0755)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		database.Db.Delete(&Symfile)
		return
	}
	f, err := os.Create(filepath.Join(filepth, Symfile.Name+".sym"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		database.Db.Delete(&Symfile)
		return
	}
	defer f.Close()
	file.Seek(0, 0)
	_, err = io.Copy(f, file)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		database.Db.Delete(&Symfile)
		return
	}
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
