// Package web provides the gin router and the endpoints for it
package web

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	router   *gin.Engine
	srv      *http.Server
	listener net.Listener
)

func initAuthRoutes(auth *gin.RouterGroup) {
	auth.POST("/crashes/:id/comments", PostCrashComment)
	auth.POST("/reports/:id/comments", PostReportComment)
	auth.POST("/reports/:id/crashid", PostReportCrashID)
	auth.POST("/reports/:id/reprocess", ReprocessReport)
	auth.POST("/reports/:id/delete", DeleteReport)
}

func initAdminRoutes(admin *gin.RouterGroup) {
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
}

func initAPIv1Routes(apiv1 *gin.RouterGroup) {
	apiv1.GET("/crashes", APIv1GetCrashes)
	apiv1.GET("/crashes/:id", APIv1GetCrash)
	apiv1.GET("/reports", APIv1GetReports)
	apiv1.GET("/reports/:id", APIv1GetReport)
	apiv1.GET("/symfiles", APIv1GetSymfiles)
	apiv1.GET("/symfiles/:id", APIv1GetSymfile)

	apiv1.GET("/products", APIv1GetProducts)
	apiv1.GET("/products/:id", APIv1GetProduct)
	apiv1.POST("/products", APIv1NewProduct)
	apiv1.PUT("/products/:id", APIv1UpdateProduct)
	apiv1.DELETE("/products/:id", APIv1DeleteProduct)

	apiv1.GET("/versions", APIv1GetVersions)
	apiv1.GET("/versions/:id", APIv1GetVersion)
	apiv1.POST("/versions", APIv1NewVersion)
	apiv1.PUT("/versions/:id", APIv1UpdateVersion)
	apiv1.DELETE("/versions/:id", APIv1DeleteVersion)

	apiv1.GET("/users", APIv1GetUsers)
	apiv1.GET("/users/:id", APIv1GetUser)
	apiv1.POST("/users", APIv1NewUser)
	apiv1.PUT("/users/:id", APIv1UpdateUser)
	apiv1.DELETE("/users/:id", APIv1DeleteUser)

	apiv1.GET("/comments", APIv1GetComments)
	apiv1.GET("/comments/:id", APIv1GetComment)
	apiv1.POST("/comments", APIv1NewComment)
	apiv1.PUT("/comments/:id", APIv1UpdateComment)
	apiv1.DELETE("/comments/:id", APIv1DeleteComment)
}

func initBreakpadRoutes(breakpad *gin.RouterGroup) {
	breakpad.GET("/", GetIndex)
	breakpad.GET("/crashes", GetCrashes)
	breakpad.GET("/crashes/:id", GetCrash)
	breakpad.POST("/crashes/:id/fixed", MarkCrashFixed)
	breakpad.GET("/reports", GetReports)
	breakpad.GET("/reports/:id", GetReport)
	breakpad.GET("/reports/:id/files/:name", GetReportFile)
	breakpad.GET("/symfiles", GetSymfiles)
	breakpad.GET("/symfiles/:id", GetSymfile)
	breakpad.POST("/reports", PostReports)
}

// Init initzializes the router
func Init() {
	router = gin.Default()
	srv = &http.Server{
		Handler: router,
	}

	auth := router.Group("/", Auth)
	initAuthRoutes(auth)

	admin := auth.Group("/admin", IsAdmin)
	initAdminRoutes(admin)

	// Admin JSON endpoints
	apiv1 := auth.Group("/api/v1")
	initAPIv1Routes(apiv1)

	// simple-breakpad endpoints
	breakpad := router.Group("/")
	initBreakpadRoutes(breakpad)

	// Static files and templates
	router.Static("/static", viper.GetString("Directory.Assets"))
	router.LoadHTMLGlob(filepath.Join(viper.GetString("Directory.Templates"), "*.html"))
}

func runIP(ip string) {
	srv.Addr = ip
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %+v", err)
		}
	}()
	log.Print("Listening on address ", ip)
}

func runSocket(sock string) {
	err := os.Remove(sock)
	if err != nil {
		log.Printf("There was an error removing the old socket: %+v", err)
	}
	listener, err = net.Listen("unix", sock)
	if err != nil {
		log.Fatalf("Socket error: %+v", err)
		return
	}
	go func() {
		if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %+v", err)
		}
	}()
	log.Print("Listening on socket ", sock)
}

// Run runs the web server
func Run() {
	if viper.GetBool("Web.UseSocket") {
		runSocket(viper.GetString("Web.BindSocket"))
	} else {
		runIP(viper.GetString("Web.BindAddress"))
	}
}

// Stop stops the webserver
func Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := srv.Shutdown(ctx); err != nil {
		cancel()
		log.Fatalf("Server Shutdown: %+v", err)
	}
	defer cancel()
	if viper.GetBool("Web.UseSocket") {
		defer func() {
			err := listener.Close()
			if err != nil {
				log.Printf("There was an error closing the socket: %+v", err)
			}
		}()
	}
}
