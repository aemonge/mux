package ui

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/x/ansi"
	"github.com/lunemis/mux/tmux"
)

func TestLayoutDimensions(t *testing.T) {
	sessions := []tmux.Session{
		{Name: "claude", WindowCount: 1, Created: time.Now().Add(-2 * time.Hour), Attached: true, Directory: "/Users/test/workspace/project1"},
		{Name: "dev-server", WindowCount: 2, Created: time.Now().Add(-24 * time.Hour), Attached: false, Directory: "/Users/test/workspace/project2"},
		{Name: "deploy", WindowCount: 1, Created: time.Now().Add(-48 * time.Hour), Attached: false, Directory: "/Users/test/workspace/project3"},
	}

	widths := []int{10, 40, 80, 120, 160, 200}
	heights := []int{10, 12, 20, 30, 40, 50}

	for _, w := range widths {
		for _, h := range heights {
			t.Run("", func(t *testing.T) {
				m := NewModel()
				m.width = w
				m.height = h
				m.sessions = sessions
				m.filtered = sessions
				m.cursor = 0

				output := m.viewMain()
				lines := strings.Split(output, "\n")

				t.Logf("w=%d h=%d => output lines=%d", w, h, len(lines))

				if len(lines) > h {
					t.Errorf("w=%d h=%d: output has %d lines, exceeds terminal height %d", w, h, len(lines), h)
					// Print first and last few lines for debugging
					for i, l := range lines {
						if i < 3 || i >= len(lines)-3 {
							t.Logf("  line %d (len=%d): %q", i, len(l), truncStr(l, 80))
						}
					}
				}
			})
		}
	}
}

func TestStackedLayoutHeights(t *testing.T) {
	tests := []struct {
		height                         int
		preview, separator, list, help int
	}{
		{height: 20, preview: 9, separator: 1, list: 8, help: 2},
		{height: 30, preview: 14, separator: 1, list: 13, help: 2},
		{height: 12, preview: 5, separator: 1, list: 4, help: 2},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("height_%d", tt.height), func(t *testing.T) {
			got := calculateStackedHeights(tt.height)
			if got.preview != tt.preview || got.separator != tt.separator ||
				got.list != tt.list || got.help != tt.help {
				t.Errorf("calculateStackedHeights(%d) = %#v, want preview=%d separator=%d list=%d help=%d",
					tt.height, got, tt.preview, tt.separator, tt.list, tt.help)
			}
			if got.preview+got.separator+got.list+got.help != tt.height {
				t.Errorf("allocated height = %d, want %d", got.preview+got.separator+got.list+got.help, tt.height)
			}
		})
	}
}

func TestViewMainStacksPreviewBeforeSelectionAndHelp(t *testing.T) {
	m := NewModel()
	m.width = 100
	m.height = 20
	m.sessions = []tmux.Session{{
		Name: "selection-marker", Created: time.Now(), Directory: "/tmp/project",
	}}
	m.applyFilter()
	m.previewKey = previewKeyForItem(*m.currentItem())
	m.previewContent = "preview-marker"

	output := ansi.Strip(m.viewMain())
	previewAt := strings.Index(output, "preview-marker")
	selectionAt := strings.Index(output, "tmux sessions")
	helpAt := strings.Index(output, "navigate")
	if previewAt < 0 || selectionAt < 0 || helpAt < 0 {
		t.Fatalf("missing layout markers: preview=%d selection=%d help=%d", previewAt, selectionAt, helpAt)
	}
	if previewAt >= selectionAt || selectionAt >= helpAt {
		t.Errorf("layout order preview/selection/help = %d/%d/%d", previewAt, selectionAt, helpAt)
	}
	if lines := strings.Count(output, "\n") + 1; lines > m.height {
		t.Errorf("output lines = %d, exceeds terminal height %d", lines, m.height)
	}
}

func TestStackedLayoutFitsExtraSelectionBar(t *testing.T) {
	m := NewModel()
	m.width = 80
	m.height = minimumStackedHeight
	m.mode = modeFilter
	m.filterMod = newFilterModel("needle")

	output := m.viewMain()
	if lines := strings.Count(output, "\n") + 1; lines != m.height {
		t.Errorf("output lines = %d, want terminal height %d", lines, m.height)
	}
}

func TestRenderSeparatorFillsTerminalWidth(t *testing.T) {
	const width = 80
	got := ansi.Strip(renderSeparator(width))
	if ansi.StringWidth(got) != width {
		t.Errorf("separator width = %d, want %d", ansi.StringWidth(got), width)
	}
	if got != strings.Repeat("━", width) {
		t.Errorf("separator = %q, want heavy horizontal line", got)
	}
}

func TestSessionListScrolling(t *testing.T) {
	// Create more sessions than can fit in a small viewport
	sessions := make([]tmux.Session, 20)
	for i := range sessions {
		sessions[i] = tmux.Session{
			Name:        fmt.Sprintf("session-%02d", i),
			WindowCount: 1,
			Created:     time.Now(),
		}
	}

	width := 60
	height := 10 // innerHeight = 8, so only 8 sessions visible

	// Cursor at 0: first session should be visible
	out := renderSessionList(sessions, 0, "", width, height)
	if !strings.Contains(out, "session-00") {
		t.Error("cursor=0: expected session-00 to be visible")
	}

	// Cursor at 15: should scroll so session-15 is visible
	out = renderSessionList(sessions, 15, "", width, height)
	if !strings.Contains(out, "session-15") {
		t.Error("cursor=15: expected session-15 to be visible")
	}
	// session-00 should be scrolled out
	if strings.Contains(out, "session-00") {
		t.Error("cursor=15: expected session-00 to be scrolled out")
	}

	// Cursor at last session
	out = renderSessionList(sessions, 19, "", width, height)
	if !strings.Contains(out, "session-19") {
		t.Error("cursor=19: expected session-19 to be visible")
	}
}

func truncStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
