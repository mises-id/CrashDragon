package main

import (
	"log"
	"os"
	"path/filepath"

	"git.1750studios.com/GSoC/CrashDragon/config"
	"git.1750studios.com/GSoC/CrashDragon/database"

	"github.com/gin-gonic/gin"
)

func initRouter() *gin.Engine {
	router := gin.Default()

	auth := router.Group("/", Auth)
	auth.POST("/crashes/:id/comments", PostCrashComment)
	auth.POST("/reports/:id/comments", PostReportComment)
	auth.POST("/reports/:id/crashid", PostReportCrashID)
	auth.POST("/reports/:id/reprocess", ReprocessReport)
	auth.POST("/symfiles", PostSymfiles)

	// Endpoints
	router.GET("/", GetCrashes)
	router.GET("/crashes", GetCrashes)
	router.GET("/crashes/:id", GetCrash)
	router.GET("/reports", GetReports)
	router.GET("/reports/:id", GetReport)
	router.GET("/reports/:id/files/:name", GetReportFile)
	router.GET("/symfiles", GetSymfiles)
	router.GET("/symfiles/:id", GetSymfile)

	router.POST("/reports", PostReports)

	router.Static("/static", config.C.AssetsDirectory)
	router.LoadHTMLGlob("templates/*.html")
	return router
}

func main() {
	log.SetFlags(log.Lshortfile)
	log.SetOutput(os.Stderr)
	var err error
	if os.Getenv("GIN_MODE") == "release" {
		err = config.GetConfig("/etc/crashdragon/config.toml")
	} else {
		err = config.GetConfig(filepath.Join(os.Getenv("HOME"), "CrashDragon/config.toml"))
	}
	if err != nil {
		log.Fatalf("FAT Config error: %+v", err)
		return
	}
	err = database.InitDb(config.C.DatabaseConnection)
	if err != nil {
		log.Fatalf("FAT Database error: %+v", err)
		return
	}

	router := initRouter()

	if config.C.UseSocket {
		router.RunUnix(config.C.BindSocket)
	} else {
		router.Run(config.C.BindAddress)
	}
}
