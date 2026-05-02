package tmux

import "testing"

func TestCapturePaneTarget_Session(t *testing.T) {
	withMock(t, func(m *mockRunner) {
		m.OnOutput([]byte("hello\n"), nil, "tmux", "capture-pane", "-t", "sess", "-p", "-e")
		got, err := CapturePaneTarget("sess")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "hello" {
			t.Errorf("got %q, want %q", got, "hello")
		}
	})
}

func TestCapturePaneTarget_Pane(t *testing.T) {
	withMock(t, func(m *mockRunner) {
		m.OnOutput([]byte("output\n"), nil, "tmux", "capture-pane", "-t", "sess:1.2", "-p", "-e")
		got, err := CapturePaneTarget("sess:1.2")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "output" {
			t.Errorf("got %q, want %q", got, "output")
		}
	})
}
