package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"code.videolan.org/videolan/CrashDragon/internal/config"
	"code.videolan.org/videolan/CrashDragon/internal/database"
	"code.videolan.org/videolan/CrashDragon/internal/migrations"
	"code.videolan.org/videolan/CrashDragon/internal/processor"
	"code.videolan.org/videolan/CrashDragon/internal/web"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

// Version holds the current version, filled by the Makefile
var Version = "Unknown"

func main() {
	log.SetFlags(log.Lshortfile)
	log.SetOutput(os.Stderr)
	flag.Bool("version", false, "")
	flag.Parse()

	if isFlagPassed("version") {
		log.Printf("Crashdragon Version: %s", Version)
		os.Exit(0)
	}

	err := config.GetConfig()
	if err != nil {
		log.Fatalf("Config error: %+v", err)
	}

	err = database.InitDB(viper.GetString("DB.Connection"))
	if err != nil {
		log.Fatalf("Database error: %+v", err)
	}

	migrations.RunMigrations()
	processor.StartQueue()
	web.Init()
	web.Run()

	c := cron.New()
	c.AddFunc("@daily", database.RemoveOldReports)
	c.Start()

	// Wait for SIGINT
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Stopping Server...")

	c.Stop()
	web.Stop()
	log.Println("Server stopped, waiting for processor queue to empty...")
	for processor.QueueSize() > 0 {
	}

	log.Println("Queue empty, closing database...")
	sqlDB, err := database.DB.DB()
	if err != nil {
		log.Printf("Error closing the database!")
		os.Exit(1)
		return
	}
	err = sqlDB.Close()
	if err != nil {
		log.Printf("Error closing the database!")
		os.Exit(1)
		return
	}

	log.Println("Closed database, good bye!")
	os.Exit(0)
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
