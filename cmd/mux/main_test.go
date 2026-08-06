package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfiguredThemeNameReadsXDGConfig(t *testing.T) {
	xdg := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdg)
	t.Setenv("MUX_THEME", "")
	path := filepath.Join(xdg, "mux", "config.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"theme":"solarized-gruvbox"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := configuredThemeName()
	if err != nil {
		t.Fatalf("configuredThemeName() error = %v", err)
	}
	if got != "solarized-gruvbox" {
		t.Errorf("configuredThemeName() = %q, want solarized-gruvbox", got)
	}
}

func TestConfiguredThemeNameEnvironmentOverridesConfig(t *testing.T) {
	xdg := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdg)
	t.Setenv("MUX_THEME", "default")
	path := filepath.Join(xdg, "mux", "config.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"theme":"solarized-gruvbox"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := configuredThemeName()
	if err != nil {
		t.Fatalf("configuredThemeName() error = %v", err)
	}
	if got != "default" {
		t.Errorf("configuredThemeName() = %q, want default", got)
	}
}

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
