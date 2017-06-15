package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
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
	file, _, err := c.Request.FormFile("upload_file_minidump")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}
	defer file.Close()
	var Crashreport database.Crashreport
	database.Db.Create(&Crashreport)
	filepath := path.Join(config.C.ContentDirectory, "Crashreports")
	os.MkdirAll(filepath, 0755)
	f, err := os.Create(path.Join(filepath, strconv.Itoa(int(Crashreport.ID))+".dmp"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}
	io.Copy(f, file)
	f.Close()
	cmd := exec.Command("./minidump-stackwalk/stackwalker", path.Join(filepath, strconv.Itoa(int(Crashreport.ID))+".dmp"), path.Join(config.C.ContentDirectory, "Symfiles"))
	//cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}
	log.Print(out.String())
	Crashreport.Product = c.Request.FormValue("prod")
	Crashreport.Version = c.Request.FormValue("ver")
	database.Db.Save(&Crashreport)
	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"error":  "OK",
		"object": Crashreport,
	})
	return
}

// PostSymfiles processes symfile
func PostSymfiles(c *gin.Context) {
	file, _, err := c.Request.FormFile("symfile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}
	defer file.Close()
	var Symfile database.Symfile
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	if err = scanner.Err(); err != nil {
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
	filepath := path.Join(config.C.ContentDirectory, "Symfiles", Symfile.Name, Symfile.Code)
	os.MkdirAll(filepath, 0755)
	f, err := os.Create(path.Join(filepath, Symfile.Name+".sym"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}
	defer f.Close()
	file.Seek(0, 0)
	io.Copy(f, file)
	database.Db.Create(&Symfile)
	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"error":  "OK",
		"object": Symfile,
	})
	return
}
