package utils

import "github.com/spf13/viper"

func GetCurrentVersionPath() string {
	return viper.GetString("CONFIG_FOLDER") + viper.GetString("CURRENT_VERSION") + "/"
}
