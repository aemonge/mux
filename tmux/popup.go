package tmux

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	popupWidth      = "85%"
	popupHeight     = "80%"
	minTmuxVersion  = 3.2
	DefaultBindKey  = "m"
)

// OpenPopup opens mux inside a tmux display-popup overlay.
// Must be called from inside a tmux session.
func OpenPopup() error {
	if os.Getenv("TMUX") == "" {
		return fmt.Errorf("mux popup must be run inside a tmux session")
	}

	version, err := getTmuxVersion()
	if err != nil {
		return fmt.Errorf("failed to detect tmux version: %w", err)
	}
	if version < minTmuxVersion {
		return fmt.Errorf("tmux %.1f+ required for popup (current: %.1f)", minTmuxVersion, version)
	}

	muxPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to find mux executable: %w", err)
	}

	cmd := exec.Command("tmux", "display-popup",
		"-E",
		"-w", popupWidth,
		"-h", popupHeight,
		muxPath,
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func findTmuxConf() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to find home directory: %w", err)
	}

	// Candidate paths in priority order:
	//   $XDG_CONFIG_HOME/tmux/tmux.conf → ~/.config/tmux/tmux.conf → ~/.tmux.conf
	var candidates []string
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		candidates = append(candidates, filepath.Join(xdg, "tmux", "tmux.conf"))
	}
	candidates = append(candidates,
		filepath.Join(home, ".config", "tmux", "tmux.conf"),
		filepath.Join(home, ".tmux.conf"),
	)

	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return filepath.Join(home, ".tmux.conf"), nil
}

// SetupKeybind adds a popup keybinding to the user's tmux config file.
func SetupKeybind(key string) error {
	muxPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to find mux executable: %w", err)
	}

	confPath, err := findTmuxConf()
	if err != nil {
		return err
	}
	bindLine := fmt.Sprintf(`bind %s display-popup -E -w %s -h %s "%s"`, key, popupWidth, popupHeight, muxPath)
	marker := "# mux popup keybinding"

	// Read existing config
	content, err := os.ReadFile(confPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read %s: %w", confPath, err)
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	replaced := false

	for _, line := range lines {
		if strings.Contains(line, marker) {
			// Replace existing mux keybinding
			newLines = append(newLines, bindLine+"  "+marker)
			replaced = true
			continue
		}
		newLines = append(newLines, line)
	}

	if !replaced {
		// Ensure trailing newline before appending
		if len(newLines) > 0 && newLines[len(newLines)-1] != "" {
			newLines = append(newLines, "")
		}
		newLines = append(newLines, bindLine+"  "+marker)
	}

	result := strings.Join(newLines, "\n")
	// Ensure file ends with newline
	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}

	if err := os.WriteFile(confPath, []byte(result), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", confPath, err)
	}

	fmt.Printf("Added to %s:\n  %s\n\n", confPath, bindLine)
	fmt.Printf("Reload tmux config:\n  tmux source-file %s\n\n", confPath)
	fmt.Printf("Then press: prefix + %s (default prefix: Ctrl+b)\n", key)

	return nil
}

func getTmuxVersion() (float64, error) {
	out, err := exec.Command("tmux", "-V").Output()
	if err != nil {
		return 0, err
	}
	// Output: "tmux 3.4" or "tmux 3.2a"
	s := strings.TrimSpace(string(out))
	s = strings.TrimPrefix(s, "tmux ")
	// Strip trailing letter (e.g. "3.2a" -> "3.2")
	var version float64
	fmt.Sscanf(s, "%f", &version)
	return version, nil
}
