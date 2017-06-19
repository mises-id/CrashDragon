package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"git.1750studios.com/GSoC/CrashDragon/config"
	"git.1750studios.com/GSoC/CrashDragon/database"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

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
	Crashreport.ID = uuid.NewV4()
	if err = database.Db.Create(&Crashreport).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}
	Crashreport.Product = c.Request.FormValue("prod")
	Crashreport.Version = c.Request.FormValue("ver")
	filepath := path.Join(config.C.ContentDirectory, "Crashreports", Crashreport.ID.String()[0:2], Crashreport.ID.String()[0:4])
	os.MkdirAll(filepath, 0755)
	f, err := os.Create(path.Join(filepath, Crashreport.ID.String()+".dmp"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		database.Db.Delete(&Crashreport)
		return
	}
	io.Copy(f, file)
	f.Close()
	go processReport(Crashreport, f)
	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"error":  nil,
		"object": Crashreport,
	})
	return
}

func processReport(Crashreport database.Crashreport, f *os.File) {
	cmdJSON := exec.Command("./minidump-stackwalk/stackwalker", f.Name(), path.Join(config.C.ContentDirectory, "Symfiles"))
	var out bytes.Buffer
	cmdJSON.Stdout = &out
	err := cmdJSON.Run()
	if err != nil {
		os.Remove(f.Name())
		database.Db.Delete(&Crashreport)
		return
	}
	err = json.Unmarshal(out.Bytes(), &Crashreport.Report)
	if err != nil {
		os.Remove(f.Name())
		database.Db.Delete(&Crashreport)
		return
	}
	if Crashreport.Report.Status != "OK" {
		os.Remove(f.Name())
		database.Db.Delete(&Crashreport)
		return
	}
	cmdTXT := exec.Command("./minidump-stackwalk/stackwalk/bin/minidump_stackwalk", f.Name(), path.Join(config.C.ContentDirectory, "Symfiles"))
	out.Reset()
	cmdTXT.Stdout = &out
	err = cmdTXT.Run()
	if err != nil {
		os.Remove(f.Name())
		database.Db.Delete(&Crashreport)
		return
	}
	Crashreport.ReportContentTXT = out.String()
	if err = database.Db.Save(&Crashreport).Error; err != nil {
		os.Remove(f.Name())
		database.Db.Delete(&Crashreport)
		return
	}
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
	updated := true
	if err = database.Db.Where("code = ?", parts[3]).First(&Symfile).Error; err != nil {
		Symfile.ID = uuid.NewV4()
		updated = false
	}
	Symfile.Os = parts[1]
	Symfile.Arch = parts[2]
	Symfile.Code = parts[3]
	Symfile.Name = parts[4]
	filepath := path.Join(config.C.ContentDirectory, "Symfiles", Symfile.Name, Symfile.Code)
	os.MkdirAll(filepath, 0755)
	os.Remove(path.Join(filepath, Symfile.Name+".sym"))
	f, err := os.Create(path.Join(filepath, Symfile.Name+".sym"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		database.Db.Delete(&Symfile)
		return
	}
	defer f.Close()
	file.Seek(0, 0)
	io.Copy(f, file)
	if updated {
		if err = database.Db.Save(&Symfile).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"error":  err.Error(),
			})
			database.Db.Delete(&Symfile)
			os.Remove(f.Name())
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"error":  nil,
			"object": Symfile,
		})
		return
	}
	if err = database.Db.Create(&Symfile).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		os.Remove(f.Name())
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"error":  nil,
		"object": Symfile,
	})
	return
}
