package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"code.videolan.org/videolan/CrashDragon/config"
	"code.videolan.org/videolan/CrashDragon/database"
	"code.videolan.org/videolan/CrashDragon/migrations"
	"code.videolan.org/videolan/CrashDragon/processor"

	"github.com/gin-gonic/gin"
)

func initRouter() *gin.Engine {
	router := gin.Default()

	auth := router.Group("/", Auth)
	auth.POST("/crashes/:id/comments", PostCrashComment)
	auth.POST("/reports/:id/comments", PostReportComment)
	auth.POST("/reports/:id/crashid", PostReportCrashID)
	auth.POST("/reports/:id/reprocess", ReprocessReport)
	auth.POST("/reports/:id/delete", DeleteReport)

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
	apiv1 := auth.Group("/api/v1")
	apiv1.GET("/crashes", APIv1GetCrashes)
	apiv1.GET("/crashes/:id", APIv1GetCrash)
	apiv1.GET("/reports", APIv1GetReports)
	apiv1.GET("/reports/:id", APIv1GetReport)
	apiv1.GET("/symfiles", APIv1GetSymfiles)
	apiv1.GET("/symfiles/:id", APIv1GetSymfile)

	apiv1.GET("/products", APIGetProducts)
	apiv1.GET("/products/:id", APIGetProduct)
	apiv1.POST("/products", APINewProduct)
	apiv1.PUT("/products/:id", APIUpdateProduct)
	apiv1.DELETE("/products/:id", APIDeleteProduct)

	apiv1.GET("/versions", APIGetVersions)
	apiv1.GET("/versions/:id", APIGetVersion)
	apiv1.POST("/versions", APINewVersion)
	apiv1.PUT("/versions/:id", APIUpdateVersion)
	apiv1.DELETE("/versions/:id", APIDeleteVersion)

	apiv1.GET("/users", APIGetUsers)
	apiv1.GET("/users/:id", APIGetUser)
	apiv1.POST("/users", APINewUser)
	apiv1.PUT("/users/:id", APIUpdateUser)
	apiv1.DELETE("/users/:id", APIDeleteUser)

	apiv1.GET("/stats", GetIndex)

	// simple-breakpad endpoints
	router.GET("/", GetIndex)
	router.GET("/crashes", GetCrashes)
	router.GET("/crashes/:id", GetCrash)
	router.POST("/crashes/:id/fixed", MarkCrashFixed)
	router.GET("/reports", GetReports)
	router.GET("/reports/:id", GetReport)
	router.GET("/reports/:id/files/:name", GetReportFile)
	router.GET("/symfiles", GetSymfiles)
	router.GET("/symfiles/:id", GetSymfile)

	router.POST("/reports", PostReports)

	router.Static("/static", config.C.AssetsDirectory)
	router.LoadHTMLGlob(filepath.Join(config.C.TemplatesDirectory, "*.html"))
	return router
}

func main() {
	log.SetFlags(log.Lshortfile)
	log.SetOutput(os.Stderr)
	cf := flag.String("config", "../etc/crashdragon.toml", "specifies the config file to use")
	flag.Parse()
	err := config.GetConfig(*cf)
	if err != nil {
		log.Fatalf("FAT Config error: %+v", err)
		return
	}
	err = database.InitDb(config.C.DatabaseConnection)
	if err != nil {
		log.Fatalf("FAT Database error: %+v", err)
		return
	}
	migrations.RunMigrations()
	processor.StartQueue()

	router := initRouter()

	srv := &http.Server{
		Handler: router,
	}

	// See: https://github.com/gin-gonic/gin/blob/master/examples/graceful-shutdown/graceful-shutdown/server.go
	if config.C.UseSocket {
		os.Remove(config.C.BindSocket)
		listener, err := net.Listen("unix", config.C.BindSocket)
		if err != nil {
			log.Fatalf("FAT Socket error: %+v", err)
			return
		}
		defer listener.Close()
		go func() {
			// service connections
			if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()
		log.Print("Listening on socket ", config.C.BindSocket)
	} else {
		srv.Addr = config.C.BindAddress
		go func() {
			// service connections
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()
		log.Print("Listening on address ", config.C.BindAddress)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Closing Server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server stopped, waiting for processor queue to empty...")
	for processor.QueueSize() > 0 {
	}

	log.Println("Queue empty, closing database...")
	database.Db.Close()

	log.Println("Closed database, good bye!")
	os.Exit(0)
	return
}
