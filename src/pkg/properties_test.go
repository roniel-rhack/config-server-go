package lib

import (
	"testing"
)

func TestNewDefaultPropertiesPreferences(t *testing.T) {
	p := NewDefaultPropertiesPreferences()
	if !p.UnwrapScalar {
		t.Errorf("expected UnwrapScalar true")
	}
	if p.KeyValueSeparator != ":" {
		t.Errorf("expected separator ':', got '%v'", p.KeyValueSeparator)
	}
	if !p.UseArrayBrackets {
		t.Errorf("expected UseArrayBrackets true")
	}
}

func TestPropertiesPreferences_Copy(t *testing.T) {
	original := NewDefaultPropertiesPreferences()
	copied := original.Copy()

	// Modify original
	original.UnwrapScalar = false
	original.KeyValueSeparator = "="
	original.UseArrayBrackets = false

	// Verify copy is independent
	if !copied.UnwrapScalar {
		t.Errorf("copy should be independent: UnwrapScalar should be true")
	}
	if copied.KeyValueSeparator != ":" {
		t.Errorf("copy should be independent: separator should be ':'")
	}
	if !copied.UseArrayBrackets {
		t.Errorf("copy should be independent: UseArrayBrackets should be true")
	}
}

func TestConfiguredPropertiesPreferences(t *testing.T) {
	p := ConfiguredPropertiesPreferences
	if !p.UnwrapScalar {
		t.Errorf("configured default UnwrapScalar should be true")
	}
	if p.KeyValueSeparator != ":" {
		t.Errorf("configured default separator should be ':'")
	}
}
