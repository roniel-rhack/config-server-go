package config

import (
	clog "configTest/custom_logguer"
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
)

func LoadConfig() {

	viper.SetDefault("SERVER.PORT", 8888)
	viper.SetDefault("CONFIG_FOLDER", "/opt/packages/config-server/configs/")
	viper.SetDefault("CONFIG_FILE", "/opt/packages/config-server/config.yaml")
	viper.SetDefault("AVAILABLE_VERSIONS", []string{})
	VerifyFolderStructure()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/opt/packages/config-server/")
	viper.AddConfigPath(".") // optionally look for config in the working directory

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			clog.Warn("Config file not found, using default values")
		} else {
			clog.Error("Error reading config file:%s", err.Error())
		}
	}

	err := viper.WriteConfigAs(viper.GetString("CONFIG_FILE"))
	if err != nil {
		clog.Error("Error writing config file: %s", err.Error())
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		clog.Info("Config file changed: %s", e.Name)
	})
	viper.WatchConfig()
}

func VerifyFolderStructure() {
	ConfigFolder := viper.GetString("CONFIG_FOLDER")
	err := os.MkdirAll(ConfigFolder, os.ModePerm)
	if err != nil {
		clog.Error("Error creating folder: %s", ConfigFolder)
	}
}
