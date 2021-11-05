// Package migrations provides the database content migrations for CrashDragon
package migrations

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"code.videolan.org/videolan/CrashDragon/internal/database"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
)

const (
	ver1_2_0 = "1.2.0"
	ver1_2_1 = "1.2.1"
	ver1_3_0 = "1.3.0"
	curVer   = ver1_3_0
)

var (
	wg  sync.WaitGroup
	sem = make(chan struct{}, 10)
)

type result struct {
	ID uuid.UUID
}

// RunMigrations runs the needed migrations
func RunMigrations() {
	dbMigrations()

	var Migration database.Migration
	database.DB.First(&Migration, "component = 'database'")
	switch Migration.Version {
	case ver1_3_0:
		log.Printf("Database migration is version 1.3.0")
		database.DB.Exec("UPDATE migrations SET version = '1.3.0' WHERE component = 'crashdragon';")
	case ver1_2_1:
		log.Printf("Database migration is version 1.2.1")
	case ver1_2_0:
		log.Print("Database migration is version 1.2.0")
		var Migration2 database.Migration
		database.DB.First(&Migration2, "component = 'crashdragon'")
		if Migration2.Version != ver1_2_0 {
			log.Print("Running crash migration, please wait...")
			migrateCrashes() // Very slow
			migrateSymfiles()
			Migration2.Version = ver1_2_0
			database.DB.Save(&Migration2)
			log.Print("Crashes migrated!")
		} else {
			log.Print("CrashDragon migration is version 1.2.0")
		}
		Migration2 = database.Migration{}
		database.DB.First(&Migration2, "component = 'database'")
		Migration2.Version = ver1_2_1
		database.DB.Save(&Migration2)
		Migration2 = database.Migration{}
		database.DB.First(&Migration2, "component = 'crashdragon'")
		Migration2.Version = ver1_2_1
		database.DB.Save(&Migration2)
	default:
		log.Fatal("Database migration version unsupported...")
	}
}

func migrateSymfiles() {
	var Symfiles []database.Symfile
	database.DB.Preload("Product").Preload("Version").Find(&Symfiles)
	for i, Symfile := range Symfiles {
		log.Printf("Moving symfile %d/%d", i+1, len(Symfiles))
		filepthnew := filepath.Join(viper.GetString("Directory.Content"), "Symfiles", Symfile.Product.Slug, Symfile.Version.Slug, Symfile.Name, Symfile.Code)
		err := os.MkdirAll(filepthnew, 0750)
		if err != nil {
			log.Fatal("Can not create directory ", err)
		}
		filepthold := filepath.Join(viper.GetString("Directory.Content"), "Symfiles", Symfile.Name, Symfile.Code)
		err = os.Rename(filepath.Join(filepthold, Symfile.Name+".sym"), filepath.Join(filepthnew, Symfile.Name+".sym"))
		if err != nil {
			log.Fatal("Could not move symfile", err)
		}
	}
}

func migrateCrashes() {
	var ccount uint
	var crashes []result
	database.DB.Model(&database.Crash{}).Select("id").Where("module IS NULL").Scan(&crashes).Count(&ccount)
	for curc, cra := range crashes {
		sem <- struct{}{}
		log.Printf("Re-reading %d/%d crashes, please wait...", curc+1, ccount)
		wg.Add(1)
		go migrateCrash(cra)
	}
	wg.Wait()
	database.DB.Exec("VACUUM ANALYZE;")
}

func migrateCrash(cra result) {
	defer wg.Done()
	tx := database.DB.Begin()
	var crash database.Crash
	tx.First(&crash, "id = ?", cra.ID)
	var reports []result
	tx.Model(&database.Report{}).Select("id").Where("crash_id = ?", crash.ID).Scan(&reports)
	existingModules := make(map[string]uuid.UUID)
	for curr, rep := range reports {
		log.Printf("\tReading reports %d/%d...", curr+1, len(reports))
		var report database.Report
		tx.First(&report, "id = ?", rep.ID)
		if report.Report.CrashInfo.CrashingThread >= len(report.Report.Threads) {
			continue
		}
	Loop:
		for _, Frame := range report.Report.Threads[report.Report.CrashInfo.CrashingThread].Frames {
			switch {
			case Frame.Function == report.Signature:
				report.Module = strings.TrimSuffix(Frame.Module, filepath.Ext(Frame.Module))
				break Loop
			case report.Module == "":
				report.Module = strings.TrimSuffix(Frame.Module, filepath.Ext(Frame.Module))
			default:
				break Loop
			}
		}
		switch {
		case crash.Module == "" && existingModules[report.Module] == uuid.Nil:
			crash.Module = report.Module
			existingModules[crash.Module] = crash.ID
			tx.Save(&crash)
			report.CrashID = crash.ID
		case crash.Module != report.Module && existingModules[report.Module] == uuid.Nil:
			crash.ID = uuid.NewV4()
			crash.Module = report.Module
			existingModules[crash.Module] = crash.ID
			tx.Create(&crash)
			report.CrashID = crash.ID
		case existingModules[report.Module] != uuid.Nil:
			report.CrashID = existingModules[report.Module]
		}

		tx.Save(&report)
	}
	err := tx.Commit().Error
	if err != nil {
		log.Fatal("Could not commit changes:", err)
	}
	<-sem
}

func dbMigrations() {
	database.DB.AutoMigrate(&database.Product{}, &database.Version{}, &database.User{}, &database.Comment{}, &database.Crash{}, &database.CrashCount{}, &database.Report{}, &database.Symfile{}, &database.Migration{})

	database.DB.Model(&database.Version{}).AddForeignKey("product_id", "products(id)", "RESTRICT", "RESTRICT")
	database.DB.Model(&database.Comment{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	database.DB.Model(&database.Report{}).AddForeignKey("crash_id", "crashes(id)", "RESTRICT", "RESTRICT")
	database.DB.Model(&database.Report{}).AddForeignKey("product_id", "products(id)", "RESTRICT", "RESTRICT")
	database.DB.Model(&database.Report{}).AddForeignKey("version_id", "versions(id)", "RESTRICT", "RESTRICT")
	database.DB.Model(&database.Symfile{}).AddForeignKey("product_id", "products(id)", "RESTRICT", "RESTRICT")
	database.DB.Model(&database.Symfile{}).AddForeignKey("version_id", "versions(id)", "RESTRICT", "RESTRICT")
	database.DB.Table("crash_versions").AddForeignKey("crash_id", "crashes(id)", "RESTRICT", "RESTRICT")
	database.DB.Table("crash_versions").AddForeignKey("version_id", "versions(id)", "RESTRICT", "RESTRICT")

	database.DB.Model(&database.Product{}).AddUniqueIndex("idx_product_slug", "slug")
	database.DB.Model(&database.Version{}).AddUniqueIndex("idx_version_slug_product", "slug", "product_id")
	database.DB.Model(&database.User{}).AddUniqueIndex("idx_user_name", "name")
	database.DB.Model(&database.Crash{}).AddUniqueIndex("idx_crash_signature_module", "signature", "module")
	database.DB.Model(&database.Symfile{}).AddUniqueIndex("idx_symfile_code", "code")

	database.DB.Model(&database.Report{}).AddIndex("idx_crash_id", "crash_id")
	database.DB.Model(&database.Report{}).AddIndex("idx_product_id", "product_id")
	database.DB.Model(&database.Report{}).AddIndex("idx_version_id", "version_id")

	database.DB.Model(&database.CrashCount{}).AddIndex("idx_crashcount_crash", "crash_id")
	database.DB.Model(&database.CrashCount{}).AddIndex("idx_crashcount_version", "version_id")
	database.DB.Model(&database.CrashCount{}).AddIndex("idx_crashcount_os", "os")

	var Migrations []database.Migration
	var cnt uint
	database.DB.Find(&Migrations).Count(&cnt)
	if cnt != 2 {
		var cd = database.Migration{ID: uuid.NewV4(), Component: "crashdragon", Version: curVer}
		database.DB.Create(&cd)
		var db = database.Migration{ID: uuid.NewV4(), Component: "database", Version: curVer}
		database.DB.Create(&db)
	}
}
