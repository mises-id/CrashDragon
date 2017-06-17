package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"git.1750studios.com/GSoC/CrashDragon/config"
	"git.1750studios.com/GSoC/CrashDragon/database"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
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
	f, err := ioutil.TempFile("", "crashdragon_")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}
	io.Copy(f, file)
	f.Close()
	var Crashreport database.Crashreport
	cmd := exec.Command("./minidump-stackwalk/stackwalker", f.Name(), path.Join(config.C.ContentDirectory, "Symfiles"))
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		os.Remove(f.Name())
		return
	}
	err = json.Unmarshal(out.Bytes(), &Crashreport.Report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		os.Remove(f.Name())
		return
	}
	if Crashreport.Report.Status != "OK" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  Crashreport.Report.Status,
		})
		os.Remove(f.Name())
		return
	}
	Crashreport.Product = c.Request.FormValue("prod")
	Crashreport.Version = c.Request.FormValue("ver")
	if err = database.Db.Create(&Crashreport).Error; err != nil {
		if err2, ok := err.(*pq.Error); ok {
			if err2.Code == "23505" {
				c.JSON(http.StatusConflict, gin.H{
					"status": http.StatusConflict,
					"error":  err2.Error(),
				})
				os.Remove(f.Name())
				return
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		os.Remove(f.Name())
		return
	}
	filepath := path.Join(config.C.ContentDirectory, "Crashreports", strconv.Itoa(int(Crashreport.ID)%100), strconv.Itoa(int(Crashreport.ID)%10))
	os.MkdirAll(filepath, 0755)
	os.Rename(f.Name(), path.Join(filepath, strconv.Itoa(int(Crashreport.ID))+".dmp"))
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
	if err = database.Db.Create(&Symfile).Error; err != nil {
		if err2, ok := err.(*pq.Error); ok {
			if err2.Code == "23505" {
				c.JSON(http.StatusConflict, gin.H{
					"status": http.StatusConflict,
					"error":  err2.Error(),
				})
				return
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}
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
	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"error":  "OK",
		"object": Symfile,
	})
	return
}
