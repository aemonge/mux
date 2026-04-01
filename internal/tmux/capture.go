package tmux

import (
	"fmt"
	"strings"
)

// CapturePane returns the visible content of the active pane in the given session.
func CapturePane(sessionName string) (string, error) {
	out, err := runner.Output("tmux", "capture-pane", "-t", sessionName, "-p", "-e")
	if err != nil {
		return "", fmt.Errorf("capture pane %s: %w", sessionName, err)
	}
	return strings.TrimRight(string(out), "\n"), nil
}
