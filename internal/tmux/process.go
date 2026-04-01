package tmux

import (
	"fmt"
	"path/filepath"
	"strings"
)

// resolveCommand returns the logical command name for a pane.
// It always inspects child processes of the pane shell to detect AI CLIs,
// because tmux pane_current_command can be unreliable (e.g. returning
// version strings instead of the actual command name).
// Works on both Linux and macOS using pgrep/ps.
func resolveCommand(panePID int, rawCmd string) string {
	if panePID <= 0 {
		return rawCmd
	}

	// If rawCmd itself is already a known AI command, return it directly.
	if IsAICommand(rawCmd) {
		return rawCmd
	}

	// Scan child processes of the pane shell to find AI CLIs.
	out, err := runner.Output("pgrep", "-P", fmt.Sprintf("%d", panePID))
	if err != nil {
		return rawCmd
	}

	for _, pidStr := range strings.Fields(string(out)) {
		// Get the full command line of each child process
		args, err := runner.Output("ps", "-o", "args=", "-p", pidStr)
		if err != nil {
			continue
		}
		for _, part := range strings.Fields(string(args)) {
			base := filepath.Base(part)
			if IsAICommand(base) {
				return base
			}
		}
	}
	return rawCmd
}
