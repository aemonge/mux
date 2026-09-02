package ui

import (
	"os"
	"strings"
	"testing"

	"github.com/lunemis/mux/tmux"
)

func TestShortenPath(t *testing.T) {
	home, _ := os.UserHomeDir()

	tests := []struct {
		input string
		want  string
	}{
		{home + "/projects/foo", "~/projects/foo"},
		{home, "~"},
		{"/tmp/other", "/tmp/other"},
		{"", ""},
	}

	for _, tt := range tests {
		got := shortenPath(tt.input)
		if got != tt.want {
			t.Errorf("shortenPath(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestShortenPathTruncatesLong(t *testing.T) {
	long := "/very/long/path/that/exceeds/thirty/five/characters/definitely"
	got := shortenPath(long)
	if len(got) > 35 {
		t.Errorf("shortenPath should truncate to 35 chars, got len=%d: %q", len(got), got)
	}
	if !strings.HasPrefix(got, "...") {
		t.Errorf("truncated path should start with '...', got %q", got)
	}
}

func TestAiLabelPlain(t *testing.T) {
	// Known commands should return non-empty
	for _, cmd := range []string{"claude", "codex", "aider", "gemini"} {
		info := aiLabelPlain(cmd)
		if info.styled == "" {
			t.Errorf("aiLabelPlain(%q) returned empty styled", cmd)
		}
		if info.text == "" {
			t.Errorf("aiLabelPlain(%q) returned empty text", cmd)
		}
		if info.extraWidth != 1 {
			t.Errorf("aiLabelPlain(%q) extraWidth = %d, want 1", cmd, info.extraWidth)
		}
	}
	// Unknown commands should return empty
	info := aiLabelPlain("bash")
	if info.styled != "" {
		t.Errorf("aiLabelPlain(%q) styled = %q, want empty", "bash", info.styled)
	}
}

func TestRenderPreviewCompactHeightOmitsOptionalTokenRow(t *testing.T) {
	session := tmux.Session{Name: "compact", Directory: "/tmp"}
	item := &listItem{kind: itemSession, session: &session}
	usage := &tmux.TokenUsage{InputTokens: 100, OutputTokens: 50, TotalCost: 0.01}

	output := renderPreview(item, "latest output", 40, 5, usage)
	if lines := strings.Count(output, "\n") + 1; lines != 5 {
		t.Errorf("preview lines = %d, want 5", lines)
	}
	if strings.Contains(output, "~$") {
		t.Error("compact preview should omit optional token row")
	}
	if !strings.Contains(output, "latest output") {
		t.Error("compact preview should preserve captured output")
	}
}

func TestRenderPreviewNilSession(t *testing.T) {
	output := renderPreview(nil, "", 40, 10, nil)
	if !strings.Contains(output, "No session selected") {
		t.Error("nil session should show 'No session selected'")
	}
}
