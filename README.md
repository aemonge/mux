# mux

A TUI tool for browsing and managing tmux sessions from the terminal.

[한국어](README.ko.md)

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-blue.svg)

## Features

- **Session list** — Active/inactive sessions sorted by recent activity
- **Live preview** — Real-time terminal output of the selected session in the right panel (refreshed every 500ms)
- **AI CLI detection** — Shows a badge when `claude`, `codex`, `aider`, `gemini`, or other AI CLIs are running in a session
- **Session management** — Create, delete, and rename sessions without leaving the TUI
- **Quick filter** — Press `/` to filter sessions by name or path in real time
- **Instant attach** — Press `Enter` to attach to the selected session (`switch-client` when inside tmux)

## Installation

### From source

```bash
git clone https://github.com/lunemis/mux.git
cd mux
make build
```

Move the binary to your PATH:

```bash
make install  # installs to /usr/local/bin
```

### Go install

```bash
go install github.com/lunemis/mux@latest
```

## Usage

```bash
mux
```

### Popup mode (recommended)

Open mux as a floating overlay inside tmux with a single keybinding — works even while AI CLIs (claude, codex, etc.) are running.

**Setup:**

```bash
# Bind prefix + m to open the popup (default)
mux setup-keybind

# Use a different key
mux setup-keybind Space

# Reload tmux config
tmux source-file ~/.tmux.conf
```

**Use:**

Press `Ctrl+b` then `m` (or your configured key) to open the popup. Select a session or press `q` to dismiss.

You can also open the popup manually:

```bash
mux popup
```

> **Note:** Popup mode requires tmux 3.2+ (`tmux -V` to check)

### Keybindings

| Key | Action |
|---|---|
| `Up` / `k` | Move up |
| `Down` / `j` | Move down |
| `g` / `G` | Jump to first / last |
| `Enter` | Attach to selected session |
| `n` | Create new session |
| `r` | Rename session |
| `x` | Delete session (with confirmation) |
| `/` | Filter sessions |
| `Esc` | Clear filter / cancel |
| `q` / `Ctrl+C` | Quit |

## Layout

```
⚡ tmux sessions (3)
┌─────────────────┐┌──────────────────────────────────────┐
│ ● my-project    ││ [ my-project ]  ~/dev/project  ✦ claude│
│   dev-server    ││ ─────────────────────────────────────  │
│   dotfiles      ││ ...terminal output preview...          │
└─────────────────┘└──────────────────────────────────────┘
↑↓/jk navigate  •  enter attach  •  n new  •  x kill  •  r rename  •  / filter  •  q quit
```

## Requirements

- Go 1.21+ (build only)
- tmux (popup mode requires 3.2+)
- Linux or macOS

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) — TUI components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — Styling

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

[MIT](LICENSE)
