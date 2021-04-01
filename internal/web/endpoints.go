package web

import (
	"bufio"
	"io"
	"log"
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
//nolint:funlen
func PostReports(c *gin.Context) {
	file, _, err := c.Request.FormFile("upload_file_minidump")
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Printf("Error closing Minidump file after upload: %+v", err)
		}
	}()
	var Report database.Report
	Report.Processed = false
	Report.ID = uuid.NewV4()

	var Product database.Product
	if err = database.DB.First(&Product, "slug = ?", c.Request.FormValue("prod")).Error; err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	Report.ProductID = Product.ID
	Report.Product = Product

	var Version database.Version
	if err = database.DB.First(&Version, "slug = ? AND product_id = ? AND ignore = false", c.Request.FormValue("ver"), Report.ProductID).Error; err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	Report.VersionID = Version.ID
	Report.Version = Version

	Report.ProcessUptime, _ = strconv.ParseUint(c.Request.FormValue("ptime"), 10, 64)
	Report.EMail = c.Request.FormValue("email")
	Report.Comment = c.Request.FormValue("comments")
	filepth := filepath.Join(config.C.ContentDirectory, "Reports", Report.ID.String()[0:2], Report.ID.String()[0:4])
	err = os.MkdirAll(filepth, 0750)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	f, err := os.Create(filepath.Join(filepth, Report.ID.String()+".dmp"))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(f, file)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		err = f.Close()
		if err != nil {
			log.Printf("Error closing the local minidump: %+v", err)
		}
		return
	}
	err = f.Close()
	if err != nil {
		log.Printf("Error closing the local minidump: %+v", err)
	}
	processor.AddToQueue(Report)
	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"object": Report.ID,
	})
}

// ReprocessReport processes the Crashreport again with current symbols
func ReprocessReport(c *gin.Context) {
	var Report database.Report
	database.DB.Where("id = ?", c.Param("id")).First(&Report)
	processor.Reprocess(Report)
	c.SetCookie("result", "OK", 0, "/", "", false, false)
	c.Redirect(http.StatusMovedPermanently, "/reports/"+Report.ID.String())
}

// PostSymfiles processes symfile
//nolint:funlen,gocognit
func PostSymfiles(c *gin.Context) {
	file, _, err := c.Request.FormFile("symfile")
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Printf("Error closing Symfile after upload: %+v", err)
		}
	}()
	var Symfile database.Symfile
	var Product database.Product
	if err = database.DB.First(&Product, "slug = ?", c.Request.FormValue("prod")).Error; err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	Symfile.ProductID = Product.ID
	Symfile.Product = Product
	var Version database.Version
	if err = database.DB.First(&Version, "slug = ? AND product_id = ?", c.Request.FormValue("ver"), Symfile.ProductID).Error; err != nil {
		Version.ID = uuid.NewV4()
		Version.Name = c.Request.FormValue("ver")
		Version.Slug = c.Request.FormValue("ver")
		Version.Ignore = false
		Version.Product = Product
		Version.ProductID = Product.ID
		database.DB.Create(&Version)
	}
	Symfile.VersionID = Version.ID
	Symfile.Version = Version
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	if err = scanner.Err(); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	parts := strings.Split(scanner.Text(), " ")
	if parts[0] != "MODULE" {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	if parts[3] == "000000000000000000000000000000000" {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	updated := true
	if err = database.DB.Where("code = ?", parts[3]).First(&Symfile).Error; err != nil {
		Symfile.ID = uuid.NewV4()
		updated = false
	} else {
		filepth := filepath.Join(config.C.ContentDirectory, "Symfiles", Symfile.Product.Slug, Symfile.Version.Slug, Symfile.Name, Symfile.Code)
		if _, existsErr := os.Stat(filepath.Join(filepth, Symfile.Name+".sym")); !os.IsNotExist(existsErr) {
			err = os.Remove(filepath.Join(filepth, Symfile.Name+".sym"))
		}
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			database.DB.Delete(&Symfile)
			return
		}
	}
	Symfile.Os = parts[1]
	Symfile.Arch = parts[2]
	Symfile.Code = parts[3]
	Symfile.Name = parts[4]
	filepth := filepath.Join(config.C.ContentDirectory, "Symfiles", Symfile.Product.Slug, Symfile.Version.Slug, Symfile.Name, Symfile.Code)
	err = os.MkdirAll(filepth, 0750)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		database.DB.Delete(&Symfile)
		return
	}
	f, err := os.Create(filepath.Join(filepth, Symfile.Name+".sym"))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		database.DB.Delete(&Symfile)
		return
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		log.Printf("Error seeking to the beginning: %+v", err)
	}
	_, err = io.Copy(f, file)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		database.DB.Delete(&Symfile)
		err = f.Close()
		if err != nil {
			log.Printf("Error closing the file: %+v", err)
		}
		return
	}
	err = f.Close()
	if err != nil {
		log.Printf("Error closing the file: %+v", err)
	}
	if updated {
		if err = database.DB.Save(&Symfile).Error; err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			database.DB.Delete(&Symfile)
			err = os.Remove(f.Name())
			if err != nil {
				log.Printf("Error removing the file: %+v", err)
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"object": Symfile,
		})
		return
	}
	if err = database.DB.Create(&Symfile).Error; err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		err = os.Remove(f.Name())
		if err != nil {
			log.Printf("Error removing the file: %+v", err)
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"object": Symfile,
	})
}
