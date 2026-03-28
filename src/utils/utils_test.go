package utils

import (
	"testing"

	"github.com/spf13/viper"
)

func TestGetCurrentVersionPath(t *testing.T) {
	viper.Reset()
	viper.Set("CONFIG_FOLDER", "/opt/configs/")
	viper.Set("CURRENT_VERSION", "v1")

	result := GetCurrentVersionPath()
	expected := "/opt/configs/v1/"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestGetCurrentVersionPath_EmptyValues(t *testing.T) {
	viper.Reset()
	// Both values are empty strings by default in viper

	result := GetCurrentVersionPath()
	expected := "/"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestGetCurrentVersionPath_TrailingSlash(t *testing.T) {
	viper.Reset()
	viper.Set("CONFIG_FOLDER", "/opt/configs/")
	viper.Set("CURRENT_VERSION", "v2")

	result := GetCurrentVersionPath()
	expected := "/opt/configs/v2/"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestGetCurrentVersionPath_NoTrailingSlashOnFolder(t *testing.T) {
	viper.Reset()
	viper.Set("CONFIG_FOLDER", "/opt/configs")
	viper.Set("CURRENT_VERSION", "v3")

	result := GetCurrentVersionPath()
	// The function concatenates directly, so no extra slash between folder and version
	expected := "/opt/configsv3/"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestGetCurrentVersionPath_OnlyFolderSet(t *testing.T) {
	viper.Reset()
	viper.Set("CONFIG_FOLDER", "/data/")

	result := GetCurrentVersionPath()
	expected := "/data//"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestGetCurrentVersionPath_OnlyVersionSet(t *testing.T) {
	viper.Reset()
	viper.Set("CURRENT_VERSION", "v5")

	result := GetCurrentVersionPath()
	expected := "v5/"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}
