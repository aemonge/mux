package ui

import (
	"os"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
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

func TestRenderPreviewAnchorsCapturedOutputBottomLeft(t *testing.T) {
	session := tmux.Session{Name: "pi", Directory: "/tmp"}
	item := &listItem{kind: itemSession, session: &session}

	output := ansi.Strip(renderPreview(item, "prompt", 24, 9, nil))
	lines := strings.Split(output, "\n")
	bottomInterior := lines[len(lines)-2]
	if !strings.HasPrefix(bottomInterior, "│prompt") {
		t.Errorf("bottom interior row = %q, want prompt anchored at left", bottomInterior)
	}
	if strings.Contains(strings.Join(lines[:len(lines)-2], "\n"), "prompt") {
		t.Error("prompt should render only on the bottom interior row")
	}
}

func TestRenderPreviewCropsTopAndRightWithoutWrapping(t *testing.T) {
	session := tmux.Session{Name: "pi", Directory: "/tmp"}
	item := &listItem{kind: itemSession, session: &session}

	output := ansi.Strip(renderPreview(item, "old\nleft-edge-right\nlatest", 12, 6, nil))
	if strings.Contains(output, "old") {
		t.Error("preview should crop the oldest row from the top")
	}
	if strings.Contains(output, "right") {
		t.Error("preview should crop overflowing text from the right")
	}
	if !strings.Contains(output, "left-edge-") || !strings.Contains(output, "latest") {
		t.Error("preview should preserve bottom rows and their left edge")
	}
	if lines := strings.Count(output, "\n") + 1; lines != 6 {
		t.Errorf("preview lines = %d, want 6 without wrapping", lines)
	}
}

func TestRenderPreviewNilSession(t *testing.T) {
	output := renderPreview(nil, "", 40, 10, nil)
	if !strings.Contains(output, "No session selected") {
		t.Error("nil session should show 'No session selected'")
	}
}
