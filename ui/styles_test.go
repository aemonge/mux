package ui

import (
	"testing"

	"github.com/lunemis/mux/theme"
)

func TestUseThemeAppliesPaletteAndAIToolColors(t *testing.T) {
	defer UseTheme(theme.Default)

	selected, err := theme.Get("solarized-gruvbox")
	if err != nil {
		t.Fatal(err)
	}
	UseTheme(selected)

	if applyBackground {
		t.Error("background NONE must not apply a terminal background")
	}
	if got := string(colorBackground); got != "" {
		t.Errorf("colorBackground = %q, want empty", got)
	}
	if got := string(colorText); got != "#3C3836" {
		t.Errorf("colorText = %q, want #3C3836", got)
	}
	_, color := commandIconPlain("claude")
	if color != "#AF3A03" {
		t.Errorf("claude color = %q, want #AF3A03", color)
	}
}
