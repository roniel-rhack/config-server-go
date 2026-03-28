package models

import (
	"encoding/json"
	"reflect"
	"testing"
)

// --- WebError ---

func TestWebError_Marshal(t *testing.T) {
	we := WebError{Error: "something went wrong"}
	data, err := json.Marshal(we)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	expected := `{"error":"something went wrong"}`
	if string(data) != expected {
		t.Errorf("got %s, want %s", string(data), expected)
	}
}

func TestWebError_Unmarshal(t *testing.T) {
	raw := `{"error":"bad request"}`
	var we WebError
	if err := json.Unmarshal([]byte(raw), &we); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if we.Error != "bad request" {
		t.Errorf("got %q, want %q", we.Error, "bad request")
	}
}

func TestWebError_Empty(t *testing.T) {
	we := WebError{}
	data, err := json.Marshal(we)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	expected := `{"error":""}`
	if string(data) != expected {
		t.Errorf("got %s, want %s", string(data), expected)
	}
}

func TestWebError_JsonTag(t *testing.T) {
	field, ok := reflect.TypeOf(WebError{}).FieldByName("Error")
	if !ok {
		t.Fatal("field Error not found")
	}
	tag := field.Tag.Get("json")
	if tag != "error" {
		t.Errorf("got json tag %q, want %q", tag, "error")
	}
}

// --- WebSuccess ---

func TestWebSuccess_Marshal(t *testing.T) {
	ws := WebSuccess{Success: "operation completed"}
	data, err := json.Marshal(ws)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	expected := `{"success":"operation completed"}`
	if string(data) != expected {
		t.Errorf("got %s, want %s", string(data), expected)
	}
}

func TestWebSuccess_Unmarshal(t *testing.T) {
	raw := `{"success":"done"}`
	var ws WebSuccess
	if err := json.Unmarshal([]byte(raw), &ws); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if ws.Success != "done" {
		t.Errorf("got %q, want %q", ws.Success, "done")
	}
}

func TestWebSuccess_Empty(t *testing.T) {
	ws := WebSuccess{}
	data, err := json.Marshal(ws)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	expected := `{"success":""}`
	if string(data) != expected {
		t.Errorf("got %s, want %s", string(data), expected)
	}
}

func TestWebSuccess_JsonTag(t *testing.T) {
	field, ok := reflect.TypeOf(WebSuccess{}).FieldByName("Success")
	if !ok {
		t.Fatal("field Success not found")
	}
	tag := field.Tag.Get("json")
	if tag != "success" {
		t.Errorf("got json tag %q, want %q", tag, "success")
	}
}

// --- SetVersion ---

func TestSetVersion_Marshal(t *testing.T) {
	sv := SetVersion{Version: "v1.2.3"}
	data, err := json.Marshal(sv)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	expected := `{"version":"v1.2.3"}`
	if string(data) != expected {
		t.Errorf("got %s, want %s", string(data), expected)
	}
}

func TestSetVersion_Unmarshal(t *testing.T) {
	raw := `{"version":"v2.0.0"}`
	var sv SetVersion
	if err := json.Unmarshal([]byte(raw), &sv); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if sv.Version != "v2.0.0" {
		t.Errorf("got %q, want %q", sv.Version, "v2.0.0")
	}
}

func TestSetVersion_Empty(t *testing.T) {
	sv := SetVersion{}
	data, err := json.Marshal(sv)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	expected := `{"version":""}`
	if string(data) != expected {
		t.Errorf("got %s, want %s", string(data), expected)
	}
}

func TestSetVersion_JsonTag(t *testing.T) {
	field, ok := reflect.TypeOf(SetVersion{}).FieldByName("Version")
	if !ok {
		t.Fatal("field Version not found")
	}
	tag := field.Tag.Get("json")
	if tag != "version" {
		t.Errorf("got json tag %q, want %q", tag, "version")
	}
}

// --- Version ---

func TestVersion_Marshal(t *testing.T) {
	v := Version{
		Version: "v1",
		Folder:  "/configs/v1",
		Files:   []string{"app.yaml", "db.yaml"},
	}
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	var roundTrip Version
	if err := json.Unmarshal(data, &roundTrip); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if !reflect.DeepEqual(v, roundTrip) {
		t.Errorf("round-trip mismatch: got %+v, want %+v", roundTrip, v)
	}
}

func TestVersion_Unmarshal(t *testing.T) {
	raw := `{"version":"v2","folder":"/configs/v2","files":["a.yml"]}`
	var v Version
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if v.Version != "v2" {
		t.Errorf("Version: got %q, want %q", v.Version, "v2")
	}
	if v.Folder != "/configs/v2" {
		t.Errorf("Folder: got %q, want %q", v.Folder, "/configs/v2")
	}
	if len(v.Files) != 1 || v.Files[0] != "a.yml" {
		t.Errorf("Files: got %v, want [a.yml]", v.Files)
	}
}

func TestVersion_EmptyFiles(t *testing.T) {
	v := Version{Version: "v1", Folder: "/x", Files: []string{}}
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	var roundTrip Version
	if err := json.Unmarshal(data, &roundTrip); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if len(roundTrip.Files) != 0 {
		t.Errorf("expected empty files, got %v", roundTrip.Files)
	}
}

func TestVersion_NilFiles(t *testing.T) {
	v := Version{Version: "v1", Folder: "/x"}
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	// nil slice marshals as null
	var roundTrip Version
	if err := json.Unmarshal(data, &roundTrip); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if roundTrip.Files != nil {
		t.Errorf("expected nil files after null unmarshal, got %v", roundTrip.Files)
	}
}

func TestVersion_JsonTags(t *testing.T) {
	typ := reflect.TypeOf(Version{})
	tests := []struct {
		fieldName string
		jsonTag   string
	}{
		{"Version", "version"},
		{"Folder", "folder"},
		{"Files", "files"},
	}
	for _, tc := range tests {
		field, ok := typ.FieldByName(tc.fieldName)
		if !ok {
			t.Errorf("field %s not found", tc.fieldName)
			continue
		}
		tag := field.Tag.Get("json")
		if tag != tc.jsonTag {
			t.Errorf("field %s: got json tag %q, want %q", tc.fieldName, tag, tc.jsonTag)
		}
	}
}

// --- AvailableVersions ---

func TestAvailableVersions_Marshal(t *testing.T) {
	av := AvailableVersions{
		Current: Version{Version: "v1", Folder: "/v1", Files: []string{"a.yml"}},
		Available: []Version{
			{Version: "v1", Folder: "/v1", Files: []string{"a.yml"}},
			{Version: "v2", Folder: "/v2", Files: []string{"b.yml", "c.yml"}},
		},
	}
	data, err := json.Marshal(av)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	var roundTrip AvailableVersions
	if err := json.Unmarshal(data, &roundTrip); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if !reflect.DeepEqual(av, roundTrip) {
		t.Errorf("round-trip mismatch: got %+v, want %+v", roundTrip, av)
	}
}

func TestAvailableVersions_MultipleVersions(t *testing.T) {
	av := AvailableVersions{
		Current: Version{Version: "v3", Folder: "/v3", Files: []string{"x.yml"}},
		Available: []Version{
			{Version: "v1", Folder: "/v1", Files: []string{}},
			{Version: "v2", Folder: "/v2", Files: []string{"a.yml"}},
			{Version: "v3", Folder: "/v3", Files: []string{"x.yml"}},
		},
	}
	data, err := json.Marshal(av)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	var roundTrip AvailableVersions
	if err := json.Unmarshal(data, &roundTrip); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if len(roundTrip.Available) != 3 {
		t.Errorf("expected 3 available versions, got %d", len(roundTrip.Available))
	}
	if roundTrip.Current.Version != "v3" {
		t.Errorf("current version: got %q, want %q", roundTrip.Current.Version, "v3")
	}
}

func TestAvailableVersions_Empty(t *testing.T) {
	av := AvailableVersions{
		Available: []Version{},
	}
	data, err := json.Marshal(av)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	var roundTrip AvailableVersions
	if err := json.Unmarshal(data, &roundTrip); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if len(roundTrip.Available) != 0 {
		t.Errorf("expected 0 available versions, got %d", len(roundTrip.Available))
	}
}

func TestAvailableVersions_JsonTags(t *testing.T) {
	typ := reflect.TypeOf(AvailableVersions{})
	tests := []struct {
		fieldName string
		jsonTag   string
	}{
		{"Current", "current"},
		{"Available", "available"},
	}
	for _, tc := range tests {
		field, ok := typ.FieldByName(tc.fieldName)
		if !ok {
			t.Errorf("field %s not found", tc.fieldName)
			continue
		}
		tag := field.Tag.Get("json")
		if tag != tc.jsonTag {
			t.Errorf("field %s: got json tag %q, want %q", tc.fieldName, tag, tc.jsonTag)
		}
	}
}

// --- PropertySources ---

func TestPropertySources_Marshal(t *testing.T) {
	ps := PropertySources{
		Name:   "application",
		Source: map[string]string{"key1": "val1", "key2": "val2"},
	}
	data, err := json.Marshal(ps)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	var roundTrip PropertySources
	if err := json.Unmarshal(data, &roundTrip); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if !reflect.DeepEqual(ps, roundTrip) {
		t.Errorf("round-trip mismatch: got %+v, want %+v", roundTrip, ps)
	}
}

func TestPropertySources_Unmarshal(t *testing.T) {
	raw := `{"name":"db","source":{"host":"localhost","port":"5432"}}`
	var ps PropertySources
	if err := json.Unmarshal([]byte(raw), &ps); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if ps.Name != "db" {
		t.Errorf("Name: got %q, want %q", ps.Name, "db")
	}
	if ps.Source["host"] != "localhost" {
		t.Errorf("Source[host]: got %q, want %q", ps.Source["host"], "localhost")
	}
	if ps.Source["port"] != "5432" {
		t.Errorf("Source[port]: got %q, want %q", ps.Source["port"], "5432")
	}
}

func TestPropertySources_EmptySource(t *testing.T) {
	ps := PropertySources{Name: "empty", Source: map[string]string{}}
	data, err := json.Marshal(ps)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	var roundTrip PropertySources
	if err := json.Unmarshal(data, &roundTrip); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if len(roundTrip.Source) != 0 {
		t.Errorf("expected empty source, got %v", roundTrip.Source)
	}
}

func TestPropertySources_NilSource(t *testing.T) {
	ps := PropertySources{Name: "nil"}
	data, err := json.Marshal(ps)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	var roundTrip PropertySources
	if err := json.Unmarshal(data, &roundTrip); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if roundTrip.Source != nil {
		t.Errorf("expected nil source, got %v", roundTrip.Source)
	}
}

func TestPropertySources_JsonTags(t *testing.T) {
	typ := reflect.TypeOf(PropertySources{})
	tests := []struct {
		fieldName string
		jsonTag   string
	}{
		{"Name", "name"},
		{"Source", "source"},
	}
	for _, tc := range tests {
		field, ok := typ.FieldByName(tc.fieldName)
		if !ok {
			t.Errorf("field %s not found", tc.fieldName)
			continue
		}
		tag := field.Tag.Get("json")
		if tag != tc.jsonTag {
			t.Errorf("field %s: got json tag %q, want %q", tc.fieldName, tag, tc.jsonTag)
		}
	}
}

// --- Config ---

func TestConfig_Marshal(t *testing.T) {
	cfg := Config{
		Name: "myapp",
		PropertySources: []PropertySources{
			{
				Name:   "application",
				Source: map[string]string{"server.port": "8080"},
			},
		},
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	var roundTrip Config
	if err := json.Unmarshal(data, &roundTrip); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if !reflect.DeepEqual(cfg, roundTrip) {
		t.Errorf("round-trip mismatch: got %+v, want %+v", roundTrip, cfg)
	}
}

func TestConfig_MultiplePropertySources(t *testing.T) {
	cfg := Config{
		Name: "myapp",
		PropertySources: []PropertySources{
			{Name: "app", Source: map[string]string{"a": "1"}},
			{Name: "db", Source: map[string]string{"b": "2"}},
			{Name: "cache", Source: map[string]string{"c": "3"}},
		},
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	var roundTrip Config
	if err := json.Unmarshal(data, &roundTrip); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if len(roundTrip.PropertySources) != 3 {
		t.Errorf("expected 3 property sources, got %d", len(roundTrip.PropertySources))
	}
}

func TestConfig_EmptyPropertySources(t *testing.T) {
	cfg := Config{Name: "empty", PropertySources: []PropertySources{}}
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	var roundTrip Config
	if err := json.Unmarshal(data, &roundTrip); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if len(roundTrip.PropertySources) != 0 {
		t.Errorf("expected 0 property sources, got %d", len(roundTrip.PropertySources))
	}
}

func TestConfig_Unmarshal(t *testing.T) {
	raw := `{"name":"svc","propertySources":[{"name":"defaults","source":{"k":"v"}}]}`
	var cfg Config
	if err := json.Unmarshal([]byte(raw), &cfg); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if cfg.Name != "svc" {
		t.Errorf("Name: got %q, want %q", cfg.Name, "svc")
	}
	if len(cfg.PropertySources) != 1 {
		t.Fatalf("expected 1 property source, got %d", len(cfg.PropertySources))
	}
	if cfg.PropertySources[0].Source["k"] != "v" {
		t.Errorf("Source[k]: got %q, want %q", cfg.PropertySources[0].Source["k"], "v")
	}
}

func TestConfig_JsonTags(t *testing.T) {
	typ := reflect.TypeOf(Config{})
	tests := []struct {
		fieldName string
		jsonTag   string
	}{
		{"Name", "name"},
		{"PropertySources", "propertySources"},
	}
	for _, tc := range tests {
		field, ok := typ.FieldByName(tc.fieldName)
		if !ok {
			t.Errorf("field %s not found", tc.fieldName)
			continue
		}
		tag := field.Tag.Get("json")
		if tag != tc.jsonTag {
			t.Errorf("field %s: got json tag %q, want %q", tc.fieldName, tag, tc.jsonTag)
		}
	}
}
