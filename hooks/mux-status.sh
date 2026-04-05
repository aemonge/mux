#!/usr/bin/env bash
# mux-status.sh — Claude Code hook script that writes agent status to a file.
# Reads JSON from stdin, writes status to /tmp/mux-status-{session_id}.json
set -euo pipefail

INPUT=$(cat)
SESSION_ID=$(echo "$INPUT" | grep -o '"session_id":"[^"]*"' | head -1 | cut -d'"' -f4)
HOOK_EVENT=$(echo "$INPUT" | grep -o '"hook_event_name":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ -z "$SESSION_ID" ]; then
    exit 0
fi

STATUS_FILE="/tmp/mux-status-${SESSION_ID}.json"

case "$HOOK_EVENT" in
    PreToolUse)
        TOOL_NAME=$(echo "$INPUT" | grep -o '"tool_name":"[^"]*"' | head -1 | cut -d'"' -f4)
        printf '{"status":"thinking","tool":"%s","ts":%d}\n' "$TOOL_NAME" "$(date +%s)" > "$STATUS_FILE"
        ;;
    Stop)
        printf '{"status":"idle","ts":%d}\n' "$(date +%s)" > "$STATUS_FILE"
        ;;
    Notification)
        NTYPE=$(echo "$INPUT" | grep -o '"notification_type":"[^"]*"' | head -1 | cut -d'"' -f4)
        if [ "$NTYPE" = "permission_prompt" ]; then
            printf '{"status":"permission","ts":%d}\n' "$(date +%s)" > "$STATUS_FILE"
        fi
        ;;
esac

exit 0
