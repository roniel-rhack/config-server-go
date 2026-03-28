package versions

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

func TestLoadAvailableVersions(t *testing.T) {
	tmpDir := t.TempDir()
	configFolder := tmpDir + "/"

	// Create version directories with files
	v1Dir := filepath.Join(tmpDir, "v1")
	v2Dir := filepath.Join(tmpDir, "v2")
	if err := os.Mkdir(v1Dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(v2Dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(v1Dir, "app.yaml"), []byte("key: val"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(v2Dir, "db.yaml"), []byte("host: localhost"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(v2Dir, "cache.yaml"), []byte("ttl: 60"), 0644); err != nil {
		t.Fatal(err)
	}

	viper.Reset()
	viper.Set("CONFIG_FOLDER", configFolder)
	viper.Set("CURRENT_VERSION", "v1")
	viper.Set("AVAILABLE_VERSIONS", []string{})

	result := LoadAvailableVersions()

	if result.Current.Version != "v1" {
		t.Errorf("Current.Version: got %q, want %q", result.Current.Version, "v1")
	}
	if len(result.Current.Files) != 1 {
		t.Errorf("Current.Files length: got %d, want 1", len(result.Current.Files))
	}
	if len(result.Available) != 2 {
		t.Fatalf("Available length: got %d, want 2", len(result.Available))
	}

	// Find v2 in available
	var foundV2 bool
	for _, v := range result.Available {
		if v.Version == "v2" {
			foundV2 = true
			if len(v.Files) != 2 {
				t.Errorf("v2 files length: got %d, want 2", len(v.Files))
			}
		}
	}
	if !foundV2 {
		t.Error("v2 not found in available versions")
	}
}

func TestLoadAvailableVersions_EmptyFolder(t *testing.T) {
	tmpDir := t.TempDir()
	configFolder := tmpDir + "/"

	viper.Reset()
	viper.Set("CONFIG_FOLDER", configFolder)
	viper.Set("CURRENT_VERSION", "v1")
	viper.Set("AVAILABLE_VERSIONS", []string{})

	// The current version folder does not exist, so getVersionStruct will log an error
	// but still return a Version struct. LoadAvailableVersions should return empty Available.
	result := LoadAvailableVersions()

	if len(result.Available) != 0 {
		t.Errorf("Available length: got %d, want 0", len(result.Available))
	}
}

func TestLoadAvailableVersions_NonExistentFolder(t *testing.T) {
	viper.Reset()
	viper.Set("CONFIG_FOLDER", "/tmp/nonexistent_config_folder_xyz_12345/")
	viper.Set("CURRENT_VERSION", "v1")
	viper.Set("AVAILABLE_VERSIONS", []string{})

	// Should not panic; the error is logged internally
	result := LoadAvailableVersions()

	if len(result.Available) != 0 {
		t.Errorf("Available length: got %d, want 0", len(result.Available))
	}
}

func TestGetVersionStruct_WithFiles(t *testing.T) {
	tmpDir := t.TempDir()
	configFolder := tmpDir + "/"

	// Create a version dir with multiple files and a subdirectory
	vDir := filepath.Join(tmpDir, "v1")
	if err := os.Mkdir(vDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(vDir, "app.yaml"), []byte("a"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(vDir, "db.properties"), []byte("b"), 0644); err != nil {
		t.Fatal(err)
	}
	// Create a subdirectory (should be excluded from files list)
	if err := os.Mkdir(filepath.Join(vDir, "subdir"), 0755); err != nil {
		t.Fatal(err)
	}

	viper.Reset()
	viper.Set("CONFIG_FOLDER", configFolder)
	viper.Set("CURRENT_VERSION", "v1")
	viper.Set("AVAILABLE_VERSIONS", []string{})

	// We test getVersionStruct indirectly via LoadAvailableVersions
	result := LoadAvailableVersions()

	if result.Current.Version != "v1" {
		t.Errorf("Current.Version: got %q, want %q", result.Current.Version, "v1")
	}
	// Should have 2 files, not the subdirectory
	if len(result.Current.Files) != 2 {
		t.Errorf("Current.Files length: got %d, want 2", len(result.Current.Files))
	}
}

func TestGetVersionStruct_EmptyFolder(t *testing.T) {
	tmpDir := t.TempDir()
	configFolder := tmpDir + "/"

	// Create an empty version directory
	vDir := filepath.Join(tmpDir, "v1")
	if err := os.Mkdir(vDir, 0755); err != nil {
		t.Fatal(err)
	}

	viper.Reset()
	viper.Set("CONFIG_FOLDER", configFolder)
	viper.Set("CURRENT_VERSION", "v1")
	viper.Set("AVAILABLE_VERSIONS", []string{})

	result := LoadAvailableVersions()

	if result.Current.Version != "v1" {
		t.Errorf("Current.Version: got %q, want %q", result.Current.Version, "v1")
	}
	// Empty folder returns empty slice (not nil), per the source code
	if result.Current.Files == nil {
		t.Error("Current.Files should not be nil for empty folder")
	}
	if len(result.Current.Files) != 0 {
		t.Errorf("Current.Files length: got %d, want 0", len(result.Current.Files))
	}
}

func TestLoadAvailableVersions_IgnoresFiles(t *testing.T) {
	tmpDir := t.TempDir()
	configFolder := tmpDir + "/"

	// Create a regular file at the config folder level (should be ignored)
	if err := os.WriteFile(filepath.Join(tmpDir, "readme.txt"), []byte("hi"), 0644); err != nil {
		t.Fatal(err)
	}
	// Create a real version directory
	vDir := filepath.Join(tmpDir, "v1")
	if err := os.Mkdir(vDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(vDir, "app.yaml"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}

	viper.Reset()
	viper.Set("CONFIG_FOLDER", configFolder)
	viper.Set("CURRENT_VERSION", "v1")
	viper.Set("AVAILABLE_VERSIONS", []string{})

	result := LoadAvailableVersions()

	// Only directories should be in Available, not regular files
	if len(result.Available) != 1 {
		t.Errorf("Available length: got %d, want 1", len(result.Available))
	}
	if len(result.Available) > 0 && result.Available[0].Version != "v1" {
		t.Errorf("Available[0].Version: got %q, want %q", result.Available[0].Version, "v1")
	}
}
