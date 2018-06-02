package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

//NilDate stores the string a empty Time type gets printed in Unix mode
const NilDate = "Mon Jan  1 00:00:00 UTC 0001"

//Config holds the structure of the configuration
type Config struct {
	DatabaseConnection string
	UseSocket          bool
	BindAddress        string
	BindSocket         string
	ContentDirectory   string
	TemplatesDirectory string
	AssetsDirectory    string
	SymbolicatorPath   string
}

//C is the actual configuration read from the file
var C Config

//LoadConfig reads the configuration file in the specified path
func LoadConfig(path string) error {
	_, e := toml.DecodeFile(path, &C)
	return e
}

//WriteConfig writes the current configuration to the specified file
func WriteConfig(path string) error {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(C)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(path), 0775)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, buf.Bytes(), 0644)
	return err
}

//GetConfig loads default values and overwrites them by the ones in a file, or creates a file with them if there is no file
func GetConfig(path string) error {
	//Set default values if there is no config
	C.DatabaseConnection = "host=localhost user=crashdragon dbname=crashdragon sslmode=disable"
	C.UseSocket = false
	C.BindAddress = "0.0.0.0:8080"
	C.BindSocket = "/var/run/crashdragon/crashdragon.sock"
	C.ContentDirectory = filepath.Join(os.Getenv("HOME"), "/CrashDragon/Files")
	C.TemplatesDirectory = "./templates"
	C.AssetsDirectory = "./assets"
	C.SymbolicatorPath = "./build/bin/minidump_stackwalk"

	var cerr error
	if _, err := os.Stat(path); err == nil {
		cerr = LoadConfig(path)
	} else {
		cerr = WriteConfig(path)
	}
	return cerr
}
