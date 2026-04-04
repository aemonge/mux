#!/usr/bin/env bash
set -euo pipefail

REPO="lunemis/mux"
BINARY="mux"
INSTALL_DIR="/usr/local/bin"
HOOKS_ONLY=false

# Parse flags
for arg in "$@"; do
    case "$arg" in
        --hooks-only) HOOKS_ONLY=true ;;
        --help|-h)
            echo "Usage: install.sh [--hooks-only]"
            echo ""
            echo "Options:"
            echo "  --hooks-only   Only set up AI tool hooks (skip binary/keybind)"
            exit 0
            ;;
    esac
done

# --- Helpers ---

info()  { printf '\033[1;34m→\033[0m %s\n' "$*"; }
ok()    { printf '\033[1;32m✓\033[0m %s\n' "$*"; }
skip()  { printf '\033[1;33m⊘\033[0m %s\n' "$*"; }
warn()  { printf '\033[1;33m⚠\033[0m %s\n' "$*"; }
fail()  { printf '\033[1;31m✗\033[0m %s\n' "$*"; exit 1; }

ask() {
    local prompt="$1"
    local default="${2:-Y}"
    local yn
    if [ "$default" = "Y" ]; then
        printf '\033[1m%s\033[0m [Y/n] ' "$prompt"
    else
        printf '\033[1m%s\033[0m [y/N] ' "$prompt"
    fi
    read -r yn </dev/tty || yn=""
    yn="${yn:-$default}"
    case "$yn" in
        [Yy]*) return 0 ;;
        *) return 1 ;;
    esac
}

detect_platform() {
    local os arch
    os="$(uname -s | tr '[:upper:]' '[:lower:]')"
    arch="$(uname -m)"

    case "$os" in
        linux)  OS="linux" ;;
        darwin) OS="darwin" ;;
        *)      fail "Unsupported OS: $os" ;;
    esac

    case "$arch" in
        x86_64|amd64)   ARCH="amd64" ;;
        arm64|aarch64)  ARCH="arm64" ;;
        *)              fail "Unsupported architecture: $arch" ;;
    esac
}

latest_version() {
    curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
        | grep '"tag_name"' \
        | head -1 \
        | sed -E 's/.*"([^"]+)".*/\1/'
}

# --- Steps ---

install_binary() {
    info "Installing ${BINARY}..."

    # Try go install first
    if command -v go &>/dev/null; then
        if ask "  Go detected. Use 'go install'?"; then
            go install "github.com/${REPO}/cmd/${BINARY}@latest"
            ok "${BINARY} installed via go install"
            return
        fi
    fi

    # Fall back to GitHub Release download
    detect_platform
    local version
    version="$(latest_version)" || fail "Could not determine latest version"
    local name="${BINARY}_${version#v}_${OS}_${ARCH}"
    local url="https://github.com/${REPO}/releases/download/${version}/${name}.tar.gz"

    info "Downloading ${BINARY} ${version} (${OS}/${ARCH})..."
    local tmp
    tmp="$(mktemp -d)"
    trap 'rm -rf "$tmp"' EXIT

    curl -fsSL "$url" -o "${tmp}/${name}.tar.gz" \
        || fail "Download failed. Check https://github.com/${REPO}/releases"
    tar -xzf "${tmp}/${name}.tar.gz" -C "$tmp"

    # Determine install path
    if [ -w "$INSTALL_DIR" ]; then
        install -m 755 "${tmp}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
        ok "${BINARY} installed to ${INSTALL_DIR}/${BINARY}"
    else
        local local_bin="${HOME}/.local/bin"
        mkdir -p "$local_bin"
        install -m 755 "${tmp}/${BINARY}" "${local_bin}/${BINARY}"
        ok "${BINARY} installed to ${local_bin}/${BINARY}"
        if ! echo "$PATH" | grep -q "$local_bin"; then
            warn "Add ${local_bin} to your PATH"
        fi
    fi
}

setup_keybind() {
    info "Setting up tmux keybinding..."
    if command -v "$BINARY" &>/dev/null; then
        "$BINARY" setup-keybind m
        ok "Keybinding added: prefix + m → mux popup"
    else
        # Manual fallback
        local conf="${HOME}/.tmux.conf"
        local line='bind-key m display-popup -E -w 80% -h 80% "mux"'
        if [ -f "$conf" ] && grep -q "mux" "$conf"; then
            ok "Keybinding already exists in ${conf}"
        else
            echo "$line" >> "$conf"
            ok "Keybinding added to ${conf}"
        fi
    fi
}

setup_claude_hooks() {
    info "Setting up Claude Code hooks..."

    local settings="${HOME}/.claude/settings.json"

    if [ ! -d "${HOME}/.claude" ]; then
        warn "~/.claude directory not found. Is Claude Code installed?"
        return
    fi

    if command -v jq &>/dev/null; then
        # Use jq for safe JSON manipulation
        local tmp_file
        tmp_file="$(mktemp)"

        if [ -f "$settings" ]; then
            # Merge hooks into existing settings
            jq '.hooks = (.hooks // {}) | .hooks.PostToolUse = (.hooks.PostToolUse // []) + [{"matcher":"*","hooks":[{"type":"command","command":"touch /tmp/mux-claude-activity-$(date +%s)"}]}] | .hooks.NotificationArrived = (.hooks.NotificationArrived // []) + [{"matcher":"*","hooks":[{"type":"command","command":"touch /tmp/mux-claude-notification"}]}]' "$settings" > "$tmp_file"
        else
            # Create new settings
            cat > "$tmp_file" << 'SETTINGS_EOF'
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "*",
        "hooks": [
          {
            "type": "command",
            "command": "touch /tmp/mux-claude-activity-$(date +%s)"
          }
        ]
      }
    ]
  }
}
SETTINGS_EOF
        fi
        mv "$tmp_file" "$settings"
        ok "Claude Code hooks configured"
    else
        warn "jq not found. Skipping Claude Code hooks setup."
        warn "Install jq and re-run: install.sh --hooks-only"
    fi
}

setup_codex_hooks() {
    info "Setting up Codex CLI hooks..."
    # Codex hook setup is placeholder — the exact config location
    # depends on the Codex CLI version
    warn "Codex CLI hook setup is experimental."
    warn "Please refer to Codex CLI documentation for hook configuration."
}

# --- Main ---

echo ""
echo "  ⚡ mux installer"
echo ""

if ! "$HOOKS_ONLY"; then
    # Step 1: Install binary
    if ask "[1/4] Install ${BINARY} binary?"; then
        install_binary
    else
        skip "${BINARY} binary installation skipped"
    fi
    echo ""

    # Step 2: tmux keybinding
    if ask "[2/4] Configure tmux keybinding (prefix+m)?"; then
        setup_keybind
    else
        skip "Keybinding setup skipped"
    fi
    echo ""

    # Step 3: Claude Code hooks
    if ask "[3/4] Set up Claude Code hooks?" "n"; then
        setup_claude_hooks
    else
        skip "Claude Code hooks skipped"
    fi
    echo ""

    # Step 4: Codex CLI hooks
    if ask "[4/4] Set up Codex CLI hooks?" "n"; then
        setup_codex_hooks
    else
        skip "Codex CLI hooks skipped"
    fi
else
    # Hooks only mode
    echo "  (hooks-only mode)"
    echo ""

    if ask "[1/2] Set up Claude Code hooks?"; then
        setup_claude_hooks
    else
        skip "Claude Code hooks skipped"
    fi
    echo ""

    if ask "[2/2] Set up Codex CLI hooks?"; then
        setup_codex_hooks
    else
        skip "Codex CLI hooks skipped"
    fi
fi

echo ""
echo "  Done! Run 'mux' in tmux to start."
echo ""
