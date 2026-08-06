package main

import (
	"strings"
	"testing"
)

func TestConfigureTheme(t *testing.T) {
	if err := configureTheme("solarized-gruvbox"); err != nil {
		t.Fatalf("configureTheme() error = %v", err)
	}
	defer configureTheme("default")

	err := configureTheme("unknown")
	if err == nil || !strings.Contains(err.Error(), "unknown") {
		t.Fatalf("configureTheme(unknown) error = %v, want unknown theme error", err)
	}
}
