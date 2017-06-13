package main

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"git.1750studios.com/GSoC/CrashDragon/config"
	"git.1750studios.com/GSoC/CrashDragon/database"

	"github.com/gin-gonic/gin"
)

// GetCrashreports returns crashreports
func GetCrashreports(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"title": "Crashreports",
	})
}

// GetCrashreport returns details of crashreport
func GetCrashreport(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"title": "Crashreports",
		"id":    c.Param("id"),
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
	c.JSON(http.StatusOK, gin.H{
		"title": "Symfiles",
	})
}

// GetSymfile returns details of symfile
func GetSymfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"title": "Symfiles",
		"id":    c.Param("id"),
	})
}

// PostCrashreports processes crashreport
func PostCrashreports(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"title": "Crashreports",
	})
}

// PostSymfiles processes symfile
func PostSymfiles(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}
	defer file.Close()
	var Symfile database.Symfile
	database.Db.Create(&Symfile)
	filepath := path.Join(config.C.ContentDirectory, "Symfiles")
	os.MkdirAll(filepath, 0755)
	f, err := os.Create(path.Join(filepath, strconv.Itoa(int(Symfile.ID))+".sym"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}
	defer f.Close()
	io.Copy(f, file)
	f.Seek(0, 0)
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}
	parts := strings.Split(scanner.Text(), " ")
	if parts[0] != "MODULE" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  "Sym-file does not start with 'MODULE'!",
		})
		return
	}
	Symfile.Os = parts[1]
	Symfile.Arch = parts[2]
	Symfile.Code = parts[3]
	Symfile.Name = parts[4]
	database.Db.Save(&Symfile)
	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"error":  "OK",
		"object": Symfile,
	})
	return
}
