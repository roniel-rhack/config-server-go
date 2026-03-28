package services

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func setupTestEnv(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")
	os.WriteFile(configFile, []byte(""), 0644)
	viper.Reset()
	viper.SetConfigFile(configFile)
	viper.Set("CONFIG_FOLDER", tmpDir+"/configs/")
	viper.Set("CURRENT_VERSION", "v1")
	viper.Set("AVAILABLE_VERSIONS", []string{"v1", "v2"})
	os.MkdirAll(tmpDir+"/configs/v1", 0755)
	os.MkdirAll(tmpDir+"/configs/v2", 0755)
	viper.WriteConfig()
	return tmpDir
}

func TestGetVersions(t *testing.T) {
	tmpDir := setupTestEnv(t)

	// Create a file inside v1 so it appears in the response
	os.WriteFile(filepath.Join(tmpDir, "configs", "v1", "app.yaml"), []byte("test"), 0644)

	app := fiber.New()
	app.Get("/versions", GetVersions)

	req := httptest.NewRequest(http.MethodGet, "/versions", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	if body["current"] == nil {
		t.Error("expected 'current' field in response")
	}
	if body["available"] == nil {
		t.Error("expected 'available' field in response")
	}
}

func TestSetVersion_InvalidBody(t *testing.T) {
	setupTestEnv(t)

	app := fiber.New()
	app.Put("/version", SetVersion)

	req := httptest.NewRequest(http.MethodPut, "/version", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	if body["error"] != "Invalid request" {
		t.Errorf("expected 'Invalid request', got %q", body["error"])
	}
}

func TestSetVersion_EmptyVersion(t *testing.T) {
	setupTestEnv(t)

	app := fiber.New()
	app.Put("/version", SetVersion)

	payload, _ := json.Marshal(map[string]string{"version": ""})
	req := httptest.NewRequest(http.MethodPut, "/version", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	if body["error"] != "Invalid version" {
		t.Errorf("expected 'Invalid version', got %q", body["error"])
	}
}

func TestSetVersion_AlreadySet(t *testing.T) {
	setupTestEnv(t)

	app := fiber.New()
	app.Put("/version", SetVersion)

	payload, _ := json.Marshal(map[string]string{"version": "v1"})
	req := httptest.NewRequest(http.MethodPut, "/version", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	if body["success"] != "Version already set" {
		t.Errorf("expected 'Version already set', got %q", body["success"])
	}
}

func TestSetVersion_NotAvailable(t *testing.T) {
	setupTestEnv(t)

	app := fiber.New()
	app.Put("/version", SetVersion)

	payload, _ := json.Marshal(map[string]string{"version": "v999"})
	req := httptest.NewRequest(http.MethodPut, "/version", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	if body["error"] != "Version not available" {
		t.Errorf("expected 'Version not available', got %q", body["error"])
	}
}

func TestSetVersion_Success(t *testing.T) {
	setupTestEnv(t)

	app := fiber.New()
	app.Put("/version", SetVersion)

	payload, _ := json.Marshal(map[string]string{"version": "v2"})
	req := httptest.NewRequest(http.MethodPut, "/version", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	if body["success"] != "Version set" {
		t.Errorf("expected 'Version set', got %q", body["success"])
	}

	if viper.GetString("CURRENT_VERSION") != "v2" {
		t.Errorf("expected CURRENT_VERSION to be 'v2', got %q", viper.GetString("CURRENT_VERSION"))
	}
}

func TestAddVersion_InvalidBody(t *testing.T) {
	setupTestEnv(t)

	app := fiber.New()
	app.Post("/version", AddVersion)

	req := httptest.NewRequest(http.MethodPost, "/version", bytes.NewBufferString("bad json"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	if body["error"] != "Invalid request" {
		t.Errorf("expected 'Invalid request', got %q", body["error"])
	}
}

func TestAddVersion_EmptyVersion(t *testing.T) {
	setupTestEnv(t)

	app := fiber.New()
	app.Post("/version", AddVersion)

	payload, _ := json.Marshal(map[string]string{"version": ""})
	req := httptest.NewRequest(http.MethodPost, "/version", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	if body["error"] != "Invalid version" {
		t.Errorf("expected 'Invalid version', got %q", body["error"])
	}
}

func TestAddVersion_AlreadyCurrent(t *testing.T) {
	setupTestEnv(t)

	app := fiber.New()
	app.Post("/version", AddVersion)

	payload, _ := json.Marshal(map[string]string{"version": "v1"})
	req := httptest.NewRequest(http.MethodPost, "/version", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	// v1 is already in the available list, so it should return 400
	if resp.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	if body["error"] != "Version already available" {
		t.Errorf("expected 'Version already available', got %q", body["error"])
	}
}

func TestAddVersion_AlreadyAvailable(t *testing.T) {
	setupTestEnv(t)

	app := fiber.New()
	app.Post("/version", AddVersion)

	payload, _ := json.Marshal(map[string]string{"version": "v2"})
	req := httptest.NewRequest(http.MethodPost, "/version", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	if body["error"] != "Version already available" {
		t.Errorf("expected 'Version already available', got %q", body["error"])
	}
}

func TestAddVersion_Success(t *testing.T) {
	tmpDir := setupTestEnv(t)

	app := fiber.New()
	app.Post("/version", AddVersion)

	payload, _ := json.Marshal(map[string]string{"version": "v3"})
	req := httptest.NewRequest(http.MethodPost, "/version", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	if body["success"] != "Version added" {
		t.Errorf("expected 'Version added', got %q", body["success"])
	}

	// Verify the folder was created
	versionPath := filepath.Join(tmpDir, "configs", "v3")
	if _, err := os.Stat(versionPath); os.IsNotExist(err) {
		t.Error("expected version folder to be created")
	}

	// Verify version was added to available versions
	versions := viper.GetStringSlice("AVAILABLE_VERSIONS")
	found := false
	for _, v := range versions {
		if v == "v3" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'v3' in AVAILABLE_VERSIONS, got %v", versions)
	}
}

func TestDeleteVersion_EmptyVersion(t *testing.T) {
	setupTestEnv(t)

	app := fiber.New()
	// Use a route that will match empty — Fiber won't match "/:version" with empty string,
	// so we test with a route that provides an empty param via a different mechanism.
	// The realistic scenario: the route param is required, so an empty path won't match.
	// We test by registering a route with an optional-like setup.
	app.Delete("/version/:version", DeleteVersion)

	// Fiber will return 404 for /version/ with no param, so we use a workaround:
	// register a second route that passes empty string explicitly.
	app.Delete("/version/", func(c *fiber.Ctx) error {
		// Simulate empty version param
		return DeleteVersion(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/version/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	if body["error"] != "Invalid version" {
		t.Errorf("expected 'Invalid version', got %q", body["error"])
	}
}

func TestDeleteVersion_IsCurrent(t *testing.T) {
	setupTestEnv(t)

	app := fiber.New()
	app.Delete("/version/:version", DeleteVersion)

	req := httptest.NewRequest(http.MethodDelete, "/version/v1", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	if body["error"] != "Cannot delete current version" {
		t.Errorf("expected 'Cannot delete current version', got %q", body["error"])
	}
}

func TestDeleteVersion_NotAvailable(t *testing.T) {
	setupTestEnv(t)

	app := fiber.New()
	app.Delete("/version/:version", DeleteVersion)

	req := httptest.NewRequest(http.MethodDelete, "/version/v999", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	if body["error"] != "Version not available" {
		t.Errorf("expected 'Version not available', got %q", body["error"])
	}
}

func TestDeleteVersion_Success(t *testing.T) {
	tmpDir := setupTestEnv(t)

	app := fiber.New()
	app.Delete("/version/:version", DeleteVersion)

	// Verify v2 folder exists before deletion
	v2Path := filepath.Join(tmpDir, "configs", "v2")
	if _, err := os.Stat(v2Path); os.IsNotExist(err) {
		t.Fatal("expected v2 folder to exist before test")
	}

	req := httptest.NewRequest(http.MethodDelete, "/version/v2", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	if body["success"] != "Version deleted" {
		t.Errorf("expected 'Version deleted', got %q", body["success"])
	}

	// Verify the folder was removed
	if _, err := os.Stat(v2Path); !os.IsNotExist(err) {
		t.Error("expected v2 folder to be deleted")
	}

	// Verify version was removed from available versions
	versions := viper.GetStringSlice("AVAILABLE_VERSIONS")
	for _, v := range versions {
		if v == "v2" {
			t.Errorf("expected 'v2' to be removed from AVAILABLE_VERSIONS, got %v", versions)
			break
		}
	}
}
