package main

import (
	"log"
	"os"
	"path/filepath"

	"code.videolan.org/videolan/CrashDragon/config"
	"code.videolan.org/videolan/CrashDragon/database"

	"github.com/gin-gonic/gin"
)

func initRouter() *gin.Engine {
	router := gin.Default()

	auth := router.Group("/", Auth)
	auth.POST("/crashes/:id/comments", PostCrashComment)
	auth.POST("/reports/:id/comments", PostReportComment)
	auth.POST("/reports/:id/crashid", PostReportCrashID)
	auth.POST("/reports/:id/reprocess", ReprocessReport)

	admin := auth.Group("/admin", IsAdmin)
	admin.GET("/", GetAdminIndex)
	admin.POST("/symfiles", PostSymfiles)

	admin.GET("/products", GetAdminProducts)
	admin.GET("/products/new", GetAdminNewProduct)
	admin.GET("/products/edit/:id", GetAdminEditProduct)
	admin.GET("/products/delete/:id", GetAdminDeleteProduct)
	admin.POST("/products/new", PostAdminNewProduct)
	admin.POST("/products/edit/:id", PostAdminEditProduct)

	admin.GET("/versions", GetAdminVersions)
	admin.GET("/versions/new", GetAdminNewVersion)
	admin.GET("/versions/edit/:id", GetAdminEditVersion)
	admin.GET("/versions/delete/:id", GetAdminDeleteVersion)
	admin.POST("/versions/new", PostAdminNewVersion)
	admin.POST("/versions/edit/:id", PostAdminEditVersion)

	admin.GET("/users", GetAdminUsers)
	admin.GET("/users/new", GetAdminNewUser)
	admin.GET("/users/edit/:id", GetAdminEditUser)
	admin.GET("/users/delete/:id", GetAdminDeleteUser)
	admin.POST("/users/new", PostAdminNewUser)
	admin.POST("/users/edit/:id", PostAdminEditUser)

	admin.GET("/symfiles", GetAdminSymfiles)
	admin.GET("/symfiles/delete/:id", GetAdminDeleteSymfile)

	// Admin JSON endpoints
	api := admin.Group("/api")
	api.GET("/products", APIGetProducts)
	api.GET("/products/:id", APIGetProduct)
	api.POST("/products", APINewProduct)
	api.PUT("/products", APIUpdateProduct)
	api.DELETE("/products/:id", APIDeleteProduct)

	api.GET("/versions", APIGetVersions)
	api.GET("/versions/:id", APIGetVersion)
	api.POST("/versions", APINewVersion)
	api.PUT("/versions", APIUpdateVersion)

	api.POST("/users", APINewUser)
	api.PUT("/users", APIUpdateUser)

	// simple-breakpad endpoints
	router.GET("/", GetIndex)
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
	err := config.GetConfig(filepath.Join(os.Getenv("HOME"), "CrashDragon/config.toml"))
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
