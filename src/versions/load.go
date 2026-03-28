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
		clog.Error("Error getting available versions from folder: %s", configFolder)
		return models.AvailableVersions{
			Current:   getVersionStruct(viper.GetString("CURRENT_VERSION")),
			Available: []models.Version{},
		}
	}

	var available []string

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

	if versions == nil {
		versions = []models.Version{}
	}

	return models.AvailableVersions{
		Current:   getVersionStruct(viper.GetString("CURRENT_VERSION")),
		Available: versions,
	}
}

func getVersionStruct(version string) models.Version {
	if version == "" {
		return models.Version{Files: []string{}}
	}

	versionFolderPath := viper.GetString("CONFIG_FOLDER") + version
	versionFiles, err := os.ReadDir(versionFolderPath)
	if err != nil {
		clog.Error("Error getting files from folder: %s", versionFolderPath)
		return models.Version{
			Version: version,
			Folder:  versionFolderPath,
			Files:   []string{},
		}
	}

	filesInDir := []string{}
	for _, file := range versionFiles {
		if !file.IsDir() {
			filesInDir = append(filesInDir, file.Name())
		}
	}

	return models.Version{
		Version: version,
		Folder:  versionFolderPath,
		Files:   filesInDir,
	}
}
