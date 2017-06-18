package main

import (
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"

	"git.1750studios.com/GSoC/CrashDragon/config"
	"git.1750studios.com/GSoC/CrashDragon/database"

	"github.com/gin-gonic/gin"
)

func initRouter() *gin.Engine {
	router := gin.Default()
	funcMap := template.FuncMap{
		"fileAbs": func(fpath string) string {
			if fpath == "" {
				return ""
			}
			return path.Join(path.Dir(fpath), path.Base(fpath))
		},
	}

	if tmpl, err := template.New("crashdragonViews").Funcs(funcMap).ParseGlob("templates/*.html"); err == nil {
		router.SetHTMLTemplate(tmpl)
	} else {
		panic(err)
	}
	// Endpoints
	router.GET("/", GetCrashreports)
	router.GET("/crashreports", GetCrashreports)
	router.GET("/crashreports/:id", GetCrashreport)
	router.GET("/crashreports/:id/files/:name", GetCrashreportFile)
	router.GET("/symfiles", GetSymfiles)
	router.GET("/symfiles/:id", GetSymfile)
	router.POST("/crashreports", PostCrashreports)
	router.POST("/symfiles", PostSymfiles)

	router.Static("/static", config.C.AssetsDirectory)
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
