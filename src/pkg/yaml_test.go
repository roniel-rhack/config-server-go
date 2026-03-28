package lib

import (
	"testing"
)

func TestNewDefaultYamlPreferences(t *testing.T) {
	p := NewDefaultYamlPreferences()
	if p.Indent != 2 {
		t.Errorf("expected Indent 2, got %v", p.Indent)
	}
	if p.ColorsEnabled {
		t.Errorf("expected ColorsEnabled false")
	}
	if !p.LeadingContentPreProcessing {
		t.Errorf("expected LeadingContentPreProcessing true")
	}
	if !p.PrintDocSeparators {
		t.Errorf("expected PrintDocSeparators true")
	}
	if !p.UnwrapScalar {
		t.Errorf("expected UnwrapScalar true")
	}
	if p.EvaluateTogether {
		t.Errorf("expected EvaluateTogether false")
	}
}

func TestYamlPreferences_Copy(t *testing.T) {
	original := NewDefaultYamlPreferences()
	copied := original.Copy()

	original.Indent = 4
	original.ColorsEnabled = true
	original.UnwrapScalar = false

	if copied.Indent != 2 {
		t.Errorf("copy should be independent: Indent should be 2, got %v", copied.Indent)
	}
	if copied.ColorsEnabled {
		t.Errorf("copy should be independent: ColorsEnabled should be false")
	}
	if !copied.UnwrapScalar {
		t.Errorf("copy should be independent: UnwrapScalar should be true")
	}
}
