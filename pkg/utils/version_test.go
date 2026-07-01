package utils

import (
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	version := Version()

	// Verify version format
	expected := "4.2.0"
	if version != expected {
		t.Errorf("Expected %s, got %s", expected, version)
	}
}

func TestVersionFormat(t *testing.T) {
	version := Version()

	// Verify version contains major version
	if !strings.Contains(version, "4") {
		t.Error("Version should contain major version 4")
	}

	// Verify format has dots
	if !strings.Contains(version, ".") {
		t.Error("Version should contain dots")
	}

	// Verify it's not empty
	if len(version) == 0 {
		t.Error("Version should not be empty")
	}
}
