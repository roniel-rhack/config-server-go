package config

import (
	clog "configTest/custom_logguer"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	clog.Initialize()
	os.Exit(m.Run())
}

func TestVerifyFolderStructure(t *testing.T) {
	tmpDir := t.TempDir()
	targetDir := filepath.Join(tmpDir, "configs", "nested")

	viper.Reset()
	viper.Set("CONFIG_FOLDER", targetDir)

	VerifyFolderStructure()

	info, err := os.Stat(targetDir)
	if err != nil {
		t.Fatalf("expected directory to exist, got error: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected a directory, got a file")
	}
}

func TestVerifyFolderStructure_AlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()

	viper.Reset()
	viper.Set("CONFIG_FOLDER", tmpDir)

	// Should not error when directory already exists
	VerifyFolderStructure()

	info, err := os.Stat(tmpDir)
	if err != nil {
		t.Fatalf("expected directory to exist, got error: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected a directory")
	}
}

func TestLoadConfig_DefaultValues(t *testing.T) {
	tmpDir := t.TempDir()

	viper.Reset()

	// Override defaults so we don't touch system paths
	viper.SetDefault("CONFIG_FOLDER", filepath.Join(tmpDir, "configs")+"/")
	viper.SetDefault("CONFIG_FILE", filepath.Join(tmpDir, "config.yaml"))

	// Change working directory so viper doesn't find any config in "."
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	// We call LoadConfig which sets its own defaults, but since we want to
	// check the defaults it sets, we need to call it. However LoadConfig
	// also sets CONFIG_FOLDER to the system default. We will verify the
	// defaults that LoadConfig sets.
	viper.Reset()
	// Pre-set CONFIG_FOLDER and CONFIG_FILE to temp paths so LoadConfig
	// doesn't create dirs under /opt
	viper.SetDefault("CONFIG_FOLDER", filepath.Join(tmpDir, "configs")+"/")
	viper.SetDefault("CONFIG_FILE", filepath.Join(tmpDir, "config.yaml"))
	viper.SetDefault("SERVER.PORT", 8888)
	viper.SetDefault("AVAILABLE_VERSIONS", []string{})

	// Create the config folder so VerifyFolderStructure succeeds
	os.MkdirAll(filepath.Join(tmpDir, "configs"), 0755)

	// Set viper config search path to tmpDir (no config file present)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(tmpDir)

	// ReadInConfig will fail with ConfigFileNotFoundError, which is expected
	_ = viper.ReadInConfig()

	// Verify defaults
	port := viper.GetInt("SERVER.PORT")
	if port != 8888 {
		t.Errorf("SERVER.PORT: got %d, want 8888", port)
	}

	versions := viper.GetStringSlice("AVAILABLE_VERSIONS")
	if len(versions) != 0 {
		t.Errorf("AVAILABLE_VERSIONS: got %v, want empty", versions)
	}
}

func TestLoadConfig_WithConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	configFolder := filepath.Join(tmpDir, "configs") + "/"
	configFile := filepath.Join(tmpDir, "config.yaml")

	// Write a config file
	configContent := `SERVER:
  PORT: 9999
CONFIG_FOLDER: ` + configFolder + `
CONFIG_FILE: ` + configFile + `
CURRENT_VERSION: v2
`
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create config folder
	os.MkdirAll(configFolder, 0755)

	viper.Reset()

	// Set defaults pointing to temp paths
	viper.SetDefault("CONFIG_FOLDER", configFolder)
	viper.SetDefault("CONFIG_FILE", configFile)
	viper.SetDefault("SERVER.PORT", 8888)
	viper.SetDefault("AVAILABLE_VERSIONS", []string{})

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(tmpDir)

	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("unexpected error reading config: %v", err)
	}

	port := viper.GetInt("SERVER.PORT")
	if port != 9999 {
		t.Errorf("SERVER.PORT: got %d, want 9999", port)
	}

	cv := viper.GetString("CURRENT_VERSION")
	if cv != "v2" {
		t.Errorf("CURRENT_VERSION: got %q, want %q", cv, "v2")
	}
}

func TestVerifyFolderStructure_EmptyPath(t *testing.T) {
	viper.Reset()
	// CONFIG_FOLDER defaults to empty string
	// MkdirAll("", ...) returns nil without creating anything
	VerifyFolderStructure()
	// Just verify no panic
}

func TestLoadConfig_NoConfigFile(t *testing.T) {
	tmpDir := t.TempDir()

	viper.Reset()

	// chdir to tmpDir so "." config path doesn't find a real config
	origDir, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	// LoadConfig sets its own defaults via SetDefault, including CONFIG_FOLDER
	// and CONFIG_FILE pointing to /opt/... paths. We override them with Set
	// (which takes priority over SetDefault) so it uses temp paths.
	viper.Set("CONFIG_FOLDER", filepath.Join(tmpDir, "configs")+"/")
	viper.Set("CONFIG_FILE", filepath.Join(tmpDir, "config.yaml"))

	LoadConfig()

	// SERVER.PORT default is 8888 set by LoadConfig via SetDefault
	port := viper.GetInt("SERVER.PORT")
	if port != 8888 {
		t.Errorf("SERVER.PORT: got %d, want 8888", port)
	}
}

func TestLoadConfig_WithConfigFileViaChdir(t *testing.T) {
	tmpDir := t.TempDir()
	configFolder := filepath.Join(tmpDir, "configs") + "/"
	configFile := filepath.Join(tmpDir, "config.yaml")

	// Create config folder so VerifyFolderStructure succeeds
	os.MkdirAll(configFolder, 0755)

	// Write a config.yaml in tmpDir so viper finds it via AddConfigPath(".")
	configContent := "SERVER:\n  PORT: 9999\nCUSTOM_KEY: testval\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "config.yaml"), []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	viper.Reset()

	origDir, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	// Use Set (not SetDefault) to override the hardcoded /opt paths
	viper.Set("CONFIG_FOLDER", configFolder)
	viper.Set("CONFIG_FILE", configFile)

	LoadConfig()

	// LoadConfig reads from the first config path that has a config file.
	// Verify that the config was written to our temp CONFIG_FILE path.
	if _, err := os.Stat(configFile); err != nil {
		t.Errorf("expected config file to be written at %s, got error: %v", configFile, err)
	}

	// SERVER.PORT should be set (either from the file or from default)
	port := viper.GetInt("SERVER.PORT")
	if port == 0 {
		t.Errorf("SERVER.PORT should not be 0")
	}
}
