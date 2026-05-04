# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Multi-window/pane tree expansion in the session list (#14):
  - `Tab` / `→` / `l` to expand a session into its windows, then a window into its panes
  - `Shift+Tab` / `←` / `h` to collapse one level
  - Preview panel now follows the cursor — captures the targeted window or pane via `tmux capture-pane -t session:window.pane`
  - `Enter` on a window/pane row attaches and focuses that exact pane (`select-window` + `select-pane` before attach)
- `tmux.ListWindows` / `tmux.ListPanes` / `tmux.CapturePaneTarget` helpers
- MIT License
- English README with Korean translation (README.ko.md)
- CONTRIBUTING.md guide
- CODE_OF_CONDUCT.md (Contributor Covenant v2.1)
- Makefile with build, test, install, clean targets
- goreleaser configuration for automated releases
- VHS demo tape for recording demo GIFs
- Unit tests for tmux and UI packages
- `--version` flag
- `scripts/test-fixture.sh` for spinning up test sessions with multiple windows/panes

### Changed
- `Session.Windows` (int) split into `Session.WindowCount` (int) + `Session.Windows` ([]Window) — the latter is lazily populated on demand
- `AttachToSession(name)` signature extended to `AttachToSession(name, windowIdx, paneIdx)` — pass `-1` to keep tmux defaults
- Cross-platform `shortenPath` using `os.UserHomeDir()` instead of hardcoded `/Users/`
- Go version in go.mod updated to stable release

### Fixed
- `renderPreview` test call missing `captured` parameter

## [0.1.0] - 2026-03-30

### Added
- TUI session manager with list and live preview panels
- Real-time terminal output preview (500ms refresh)
- AI CLI detection (claude, codex, aider, gemini) with badge display
- Session create, delete, and rename from within the TUI
- Quick filter with `/` key
- Instant attach / switch-client
- Popup mode (`mux popup`) as tmux floating overlay
- `mux setup-keybind` for one-command tmux keybinding setup
- Cross-platform AI CLI detection (Linux/macOS)
