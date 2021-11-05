// Package config provides the config for CrashDragon
package config

import (
	"github.com/spf13/viper"
)

// GetConfig loads default values and overwrites them by the ones in a file, or creates a file with them if there is no file
func GetConfig() error {
	// Set default values if there is no config
	viper.SetDefault("DB.Connection", "host=localhost user=crashdragon dbname=crashdragon sslmode=disable")
	viper.SetDefault("Web.UseSocket", false)
	viper.SetDefault("Web.BindAddress", ":8080")
	viper.SetDefault("Web.BindSocket", "/var/run/crashdragon/crashdragon.sock")
	viper.SetDefault("Directory.Content", "../share/crashdragon/files")
	viper.SetDefault("Directory.Templates", "./web/templates")
	viper.SetDefault("Directory.Assets", "./web/assets")
	viper.SetDefault("Symbolicator.Executable", "./minidump_stackwalk")
	viper.SetDefault("Symbolicator.TrimModuleNames", true)
	viper.SetDefault("Housekeeping.ReportRetentionTime", "2190h") // Around 3 months (duration only supports times in hours and down due to irregular length of days/months/years)

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/crashdragon/")
	viper.AddConfigPath("$HOME/.crashdragon")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return viper.WriteConfig()
}
