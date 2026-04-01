package tmux

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const commandCacheTTL = 5 * time.Second

type cachedCommand struct {
	command   string
	expiresAt time.Time
}

var (
	cmdCache   = make(map[int]cachedCommand)
	cmdCacheMu sync.Mutex
)

// resolveCommand returns the logical command name for a pane.
// Results are cached by panePID with a TTL to avoid repeated pgrep/ps calls.
func resolveCommand(panePID int, rawCmd string) string {
	if panePID <= 0 {
		return rawCmd
	}

	if IsAICommand(rawCmd) {
		return rawCmd
	}

	cmdCacheMu.Lock()
	if cached, ok := cmdCache[panePID]; ok && time.Now().Before(cached.expiresAt) {
		cmdCacheMu.Unlock()
		return cached.command
	}
	cmdCacheMu.Unlock()

	result := scanChildProcesses(panePID, rawCmd)

	cmdCacheMu.Lock()
	cmdCache[panePID] = cachedCommand{
		command:   result,
		expiresAt: time.Now().Add(commandCacheTTL),
	}
	cmdCacheMu.Unlock()

	return result
}

// scanChildProcesses inspects child processes of the pane shell to detect AI CLIs.
// Works on both Linux and macOS using pgrep/ps.
func scanChildProcesses(panePID int, rawCmd string) string {
	out, err := runner.Output("pgrep", "-P", fmt.Sprintf("%d", panePID))
	if err != nil {
		return rawCmd
	}

	for _, pidStr := range strings.Fields(string(out)) {
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
