# mux

**Switch between AI CLI sessions without breaking your flow.**

Running Claude in one session, Codex in another, and a dev server in a third? Switching between them means detaching, listing sessions, remembering which is which, and reattaching. mux eliminates that friction — see every session's live output at a glance, spot which AI tools are active, and switch in a keystroke.

[한국어](README.ko.md)

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-blue.svg)

<!-- TODO: Replace with an actual demo recording
     Record with: https://github.com/charmbracelet/vhs
     vhs demo.tape  →  produces demo.gif -->
![Demo](assets/demo.gif)

## The Problem

In the age of AI-powered development, a typical workflow looks like:

- **Session 1**: Claude Code working on your feature
- **Session 2**: Codex reviewing your test suite
- **Session 3**: Dev server running your app
- **Session 4**: Another Claude session refactoring a different module

tmux's built-in `choose-session` shows you a list of names — but which session has Claude waiting for your input? Which one is still running? You end up cycling through sessions blindly.

## How mux solves it

### Live preview
See the actual terminal output of any session *before* you switch. No more guessing.

<!-- TODO: screenshot of the main view with live preview panel visible -->
<!-- ![Live preview](assets/preview.png) -->

### AI CLI detection
`claude`, `codex`, `aider`, `gemini` are automatically detected and highlighted with badges — instantly find the right session.

<!-- TODO: screenshot showing sessions with AI CLI badges (✦ claude, ◈ codex, etc.) -->
<!-- ![AI CLI badges](assets/ai-badges.png) -->

### Git branch & worktree display
Each session shows the current git branch. Worktrees are marked with `⌥⌥` so you can tell at a glance which sessions are working on isolated branches.

### Cost & token tracking
For Claude Code sessions, mux reads session logs to display real-time token usage and estimated cost — no configuration needed.

### Popup overlay
Press one key to summon mux on top of whatever you're doing — even mid-conversation with an AI CLI. Pick a session and you're there.

<!-- TODO: screenshot/gif of popup appearing over an active Claude session -->
<!-- ![Popup mode](assets/popup.gif) -->

### Vim-style navigation
`j`/`k` to browse, `/` to filter, `Enter` to attach. No mouse needed.

## Quick Start

```bash
# One-line interactive installer (recommended)
curl -sSL https://raw.githubusercontent.com/lunemis/mux/main/install.sh | bash

# Or install manually
brew install lunemis/tap/mux   # or: go install github.com/lunemis/mux/cmd/mux@latest
mux                             # launch the session manager
```

For the best experience, set up popup mode (opens mux as a floating overlay):

```bash
mux setup-keybind               # binds prefix + m
tmux source-file ~/.tmux.conf   # reload config
```

Now press `Ctrl+b` then `m` anywhere in tmux to open mux.

## Installation

### Interactive installer (recommended)

The installer guides you through binary installation and keybinding setup:

```bash
curl -sSL https://raw.githubusercontent.com/lunemis/mux/main/install.sh | bash
```

### Homebrew

```bash
brew install lunemis/tap/mux
```

### From source

```bash
git clone https://github.com/lunemis/mux.git
cd mux
make install   # builds and installs to /usr/local/bin
```

### Go install

```bash
go install github.com/lunemis/mux/cmd/mux@latest
```

## Usage

### Basic

Run `mux` to open the session manager. Use `j`/`k` to navigate, `Enter` to attach, `q` to quit.

```
⚡ tmux sessions (3)
┌──────────────────────────┐┌──────────────────────────────────────────┐
│ * my-project    2h ✦ main ││ [ my-project ]  ~/dev/project  ✦ claude  │
│   dev-server    3h        ││   45.2k in / 12.1k out  ~$0.85           │
│   dotfiles      1d   main ││ ──────────────────────────────────────── │
│                            ││ ...terminal output preview...            │
└──────────────────────────┘└──────────────────────────────────────────┘
↑↓/jk navigate  •  enter attach  •  n new  •  x kill  •  r rename  •  / filter  •  q quit
```

The left panel shows your sessions with AI badges and git branches. The right panel shows a **live preview** of the selected session's terminal output, updated every 500ms.

### Popup mode (recommended)

Open mux as a floating overlay inside tmux — works even while AI CLIs are running in the foreground.

```bash
# Set up the keybinding (one-time)
mux setup-keybind          # prefix + m (default)
mux setup-keybind Space    # or use a different key

# Reload tmux config
tmux source-file ~/.tmux.conf
```

You can also open the popup manually with `mux popup`.

> **Note:** Popup mode requires tmux 3.2+

### Statusbar widget

Show AI session icons in your tmux status bar without opening the TUI:

```bash
# Add to ~/.tmux.conf
set -g status-right '#(mux status)'
```

This runs `mux status` which outputs a compact summary like `✦ ◈` when AI sessions are active.

### Keybindings

| Key | Action |
|---|---|
| `j` / `k` | Move down / up |
| `g` / `G` | Jump to first / last |
| `Enter` | Attach to selected session |
| `n` | Create new session |
| `r` | Rename session |
| `x` | Delete session (with confirmation) |
| `/` | Filter sessions by name or path |
| `Esc` | Clear filter / cancel |
| `q` | Quit |

## Requirements

- tmux (popup mode requires 3.2+)
- Linux or macOS

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

[MIT](LICENSE)
