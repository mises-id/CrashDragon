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
)

func main() {
	log.SetFlags(log.Lshortfile)
	log.SetOutput(os.Stderr)
	cf := flag.String("config", "../etc/crashdragon.toml", "specifies the config file to use")
	flag.Parse()

	err := config.GetConfig(*cf)
	if err != nil {
		log.Fatalf("Config error: %+v", err)
	}

	err = database.InitDb(config.C.DatabaseConnection)
	if err != nil {
		log.Fatalf("Database error: %+v", err)
	}

	migrations.RunMigrations()
	processor.StartQueue()
	web.Init()
	web.Run()

	// Wait for SIGINT
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Stopping Server...")

	web.Stop()
	log.Println("Server stopped, waiting for processor queue to empty...")
	for processor.QueueSize() > 0 {
	}

	log.Println("Queue empty, closing database...")
	err = database.Db.Close()
	if err != nil {
		log.Printf("Error closing the database!")
		os.Exit(1)
		return
	}

	log.Println("Closed database, good bye!")
	os.Exit(0)
}
