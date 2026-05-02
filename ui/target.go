package ui

import "fmt"

// formatTarget formats a tmux target string. Pass paneIdx == -1 to address the
// active pane of the given window.
func formatTarget(sessionName string, windowIdx, paneIdx int) string {
	if paneIdx < 0 {
		return fmt.Sprintf("%s:%d", sessionName, windowIdx)
	}
	return fmt.Sprintf("%s:%d.%d", sessionName, windowIdx, paneIdx)
}
