package database

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type Crashreport struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Product string
	Version string
}

type Symfile struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Os       string
	Name     string
	Code     string
	Arch     string
	Contents string
}

var Db *gorm.DB

func InitDb(connection string) error {
	var err error
	Db, err = gorm.Open("postgres", connection)
	if err != nil {
		log.Fatalf("FAT Database error: %+v", err)
		return err
	}
	Db.LogMode(true)

	Db.AutoMigrate(&Crashreport{}, &Symfile{})
	return err
}
