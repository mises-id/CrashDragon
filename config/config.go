package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	DatabaseConnection string
	UseSocket          bool
	BindAddress        string
	BindSocket         string
	ContentDirectory   string
}

var C Config

func LoadConfig(path string) error {
	_, e := toml.DecodeFile(path, &C)
	return e
}

func WriteConfig(path string) error {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(C)
	if err != nil {
		return err
	}
	os.MkdirAll(filepath.Dir(path), 0775)
	err = ioutil.WriteFile(path, buf.Bytes(), 0644)
	return err
}

func GetConfig(path string) error {
	//Set default values if there is no config
	C.DatabaseConnection = "host=localhost user=crashdragon dbname=crashdragon password=crashdragon sslmode=disable"
	C.UseSocket = false
	C.BindAddress = "0.0.0.0:8080"
	C.BindSocket = "/var/run/crashdragon/crashdragon.sock"
	C.ContentDirectory = filepath.Join(os.Getenv("HOME"), "/CrashDragon/Files")

	var cerr error
	if _, err := os.Stat(path); err == nil {
		cerr = LoadConfig(path)
	} else {
		cerr = WriteConfig(path)
	}
	return cerr
}
