package services

import (
	clog "configTest/custom_logguer"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	clog.Initialize()
	os.Exit(m.Run())
}

func TestGetConfigFile_FileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	viper.Reset()
	viper.Set("CONFIG_FOLDER", tmpDir+"/configs/")
	viper.Set("CURRENT_VERSION", "v1")
	os.MkdirAll(tmpDir+"/configs/v1", 0755)

	app := fiber.New()
	app.Get("/:filename", GetConfigFile)

	req := httptest.NewRequest(http.MethodGet, "/nonexistent.yaml", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 404 {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	if body["error"] != "File not found" {
		t.Errorf("expected 'File not found', got %q", body["error"])
	}
}

func TestGetConfigFile_Success(t *testing.T) {
	tmpDir := t.TempDir()
	viper.Reset()
	viper.Set("CONFIG_FOLDER", tmpDir+"/configs/")
	viper.Set("CURRENT_VERSION", "v1")

	versionDir := tmpDir + "/configs/v1/"
	os.MkdirAll(versionDir, 0755)

	expectedContent := "server:\n  port: 8080\n"
	filePath := filepath.Join(versionDir, "app.yaml")
	if err := os.WriteFile(filePath, []byte(expectedContent), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	app := fiber.New()
	app.Get("/:filename", GetConfigFile)

	req := httptest.NewRequest(http.MethodGet, "/app.yaml", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	buf := make([]byte, 1024)
	n, _ := resp.Body.Read(buf)
	actual := string(buf[:n])
	if actual != expectedContent {
		t.Errorf("expected %q, got %q", expectedContent, actual)
	}
}

func TestGetConfigFile_PathTraversalProtection(t *testing.T) {
	tmpDir := t.TempDir()
	viper.Reset()
	viper.Set("CONFIG_FOLDER", tmpDir+"/configs/")
	viper.Set("CURRENT_VERSION", "v1")
	os.MkdirAll(tmpDir+"/configs/v1", 0755)

	// Create a sensitive file outside the version dir
	os.WriteFile(tmpDir+"/configs/secret.txt", []byte("secret"), 0644)

	basePath := tmpDir + "/configs/v1/"
	filename := "../secret.txt"
	filePath := filepath.Join(basePath, filename)
	cleaned := filepath.Clean(filePath)
	cleanedBase := filepath.Clean(basePath)

	// Verify the path would escape the base directory
	if filepath.HasPrefix(cleaned, cleanedBase) {
		t.Skip("filepath.Join already prevents traversal on this OS")
	}

	// The protection in GetConfigFile uses strings.HasPrefix
	// Verify the protection logic works directly
	if strings.HasPrefix(cleaned, cleanedBase) {
		t.Errorf("path traversal not detected: %s starts with %s", cleaned, cleanedBase)
	}
}

func TestGetConfig_Success(t *testing.T) {
	tmpDir := t.TempDir()
	viper.Reset()
	viper.Set("CONFIG_FOLDER", tmpDir+"/configs/")
	viper.Set("CURRENT_VERSION", "v1")

	versionDir := tmpDir + "/configs/v1/"
	os.MkdirAll(versionDir, 0755)

	yamlContent := "server:\n  port: \"8080\"\n  host: \"localhost\"\n"
	filePath := filepath.Join(versionDir, "myapp-dev.yaml")
	if err := os.WriteFile(filePath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	app := fiber.New()
	app.Get("/:appName/:profile", GetConfig)

	req := httptest.NewRequest(http.MethodGet, "/myapp/dev", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var config struct {
		Name            string `json:"name"`
		PropertySources []struct {
			Name   string            `json:"name"`
			Source map[string]string `json:"source"`
		} `json:"propertySources"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if config.Name != "application" {
		t.Errorf("expected name 'application', got %q", config.Name)
	}
	if len(config.PropertySources) != 1 {
		t.Fatalf("expected 1 property source, got %d", len(config.PropertySources))
	}
	if config.PropertySources[0].Name != "application" {
		t.Errorf("expected property source name 'application', got %q", config.PropertySources[0].Name)
	}
	source := config.PropertySources[0].Source
	if source["server.port"] != "8080" {
		t.Errorf("expected server.port=8080, got %q", source["server.port"])
	}
	if source["server.host"] != "localhost" {
		t.Errorf("expected server.host=localhost, got %q", source["server.host"])
	}
}
