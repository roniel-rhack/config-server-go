package lib

import (
	"testing"
)

func TestFormat_MatchesName_FormalName(t *testing.T) {
	f := &Format{FormalName: "yaml", Names: []string{"y", "yml"}}
	if !f.MatchesName("yaml") {
		t.Errorf("expected formal name match")
	}
}

func TestFormat_MatchesName_Alias(t *testing.T) {
	f := &Format{FormalName: "yaml", Names: []string{"y", "yml"}}
	if !f.MatchesName("yml") {
		t.Errorf("expected alias match")
	}
}

func TestFormat_MatchesName_NoMatch(t *testing.T) {
	f := &Format{FormalName: "yaml", Names: []string{"y", "yml"}}
	if f.MatchesName("json") {
		t.Errorf("expected no match")
	}
}

func TestFormatFromString_Yaml(t *testing.T) {
	f, err := FormatFromString("yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.FormalName != "yaml" {
		t.Errorf("expected yaml format, got %v", f.FormalName)
	}
}

func TestFormatFromString_Props(t *testing.T) {
	f, err := FormatFromString("props")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.FormalName != "props" {
		t.Errorf("expected props format, got %v", f.FormalName)
	}
}

func TestFormatFromString_PropsAlias(t *testing.T) {
	f, err := FormatFromString("properties")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.FormalName != "props" {
		t.Errorf("expected props format, got %v", f.FormalName)
	}
}

func TestFormatFromString_Unknown(t *testing.T) {
	_, err := FormatFromString("json")
	if err == nil {
		t.Errorf("expected error for unknown format")
	}
}

func TestFormatFromString_Empty(t *testing.T) {
	_, err := FormatFromString("")
	if err == nil {
		t.Errorf("expected error for empty format")
	}
}

func TestGetAvailableOutputFormats(t *testing.T) {
	formats := GetAvailableOutputFormats()
	if len(formats) == 0 {
		t.Errorf("expected at least one output format")
	}
	// All returned formats should have an encoder factory
	for _, f := range formats {
		if f.EncoderFactory == nil {
			t.Errorf("format %v has nil encoder factory", f.FormalName)
		}
	}
}

func TestGetAvailableOutputFormatString(t *testing.T) {
	result := GetAvailableOutputFormatString()
	if result == "" {
		t.Errorf("expected non-empty format string")
	}
}

func TestGetConfiguredEncoder(t *testing.T) {
	enc := PropertiesFormat.GetConfiguredEncoder()
	if enc == nil {
		t.Errorf("expected non-nil encoder for properties format")
	}
}

func TestGetConfiguredEncoder_Yaml(t *testing.T) {
	enc := YamlFormat.GetConfiguredEncoder()
	// yaml encoder factory returns nil
	if enc != nil {
		t.Errorf("expected nil encoder for yaml format (returns nil)")
	}
}
