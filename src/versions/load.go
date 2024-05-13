package versions

import (
	clog "configTest/custom_logguer"
	"configTest/models"
	"github.com/spf13/viper"
	"os"
)

func LoadAvailableVersions() models.AvailableVersions {
	configFolder := viper.GetString("CONFIG_FOLDER")
	files, err := os.ReadDir(configFolder)

	if err != nil {
		clog.Error("Error getting available versions from folder:%s", configFolder)
	}

	available := viper.GetStringSlice("AVAILABLE_VERSIONS")[:0]

	for _, file := range files {
		if file.IsDir() {
			available = append(available, file.Name())
		}
	}

	viper.Set("AVAILABLE_VERSIONS", available)

	var versions []models.Version

	for _, version := range available {
		v := getVersionStruct(version)
		versions = append(versions, v)
	}

	return models.AvailableVersions{
		Current:   getVersionStruct(viper.GetString("CURRENT_VERSION")),
		Available: versions,
	}
}

func getVersionStruct(version string) models.Version {
	versionFolderPath := viper.GetString("CONFIG_FOLDER") + version
	versionFiles, err := os.ReadDir(versionFolderPath)
	if err != nil {
		clog.Error("Error getting files from folder: %s", versionFolderPath)
	}
	var filesInDir []string
	for _, file := range versionFiles {
		if !file.IsDir() {
			filesInDir = append(filesInDir, file.Name())
		}
	}

	if len(filesInDir) == 0 {
		clog.Warn("No files found in version folder: %s", versionFolderPath)
		filesInDir = []string{}
	}

	return models.Version{
		Version: version,
		Folder:  versionFolderPath,
		Files:   filesInDir,
	}
}
