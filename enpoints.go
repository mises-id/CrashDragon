package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCrashreports(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"title": "Crashreports",
	})
}

func GetCrashreport(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"title": "Crashreports",
		"id":    c.Param("id"),
	})
}

func GetCrashreportFile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"title": "Crashreports",
		"id":    c.Param("id"),
		"name":  c.Param("name"),
	})
}

func GetSymfiles(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"title": "Symfiles",
	})
}

func GetSymfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"title": "Symfiles",
		"id":    c.Param("id"),
	})
}

func PostCrashreports(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"title": "Crashreports",
	})
}

func PostSymfiles(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"title": "Symfiles",
	})
}
