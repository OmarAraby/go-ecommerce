#!/usr/bin/env bash
# stop_reminder.sh — Claude Code Stop hook (OPTIONAL, gentle enforcement).
# Fires when Claude finishes a response (once per turn).
#
# Purpose: nudge toward saving memory when meaningful work happened but the
# session-log hasn't been updated today. This is what makes Claude "remember"
# to log — without you having to rely on it.
#
# Behavior is intentionally SOFT: it injects a reminder via additionalContext at
# most once per session (guarded by a marker file). It does NOT hard-block, to
# avoid loops. If you want hard enforcement, see the commented block at the bottom.
#
# Input: JSON on stdin. Output: optional JSON to control the stop.

INPUT="$(cat)"
MEM=".ai-memory"
TODAY="$(date +%F)"
# Marker lives in temp (keyed by project path + date) so it never shows up in
# 'git status' and never needs gitignoring.
_KEY="$(printf '%s' "$PWD" | tr -c 'A-Za-z0-9' '_')"
MARKER="${TMPDIR:-/tmp}/.claude_stop_${_KEY}_${TODAY}"

# No memory, or already reminded today → stay silent.
[ ! -d "$MEM" ] && exit 0
[ -f "$MARKER" ] && exit 0

# Did meaningful work happen? Use uncommitted changes as a cheap proxy.
CHANGES="$(git status --short 2>/dev/null | wc -l | tr -d ' ')"
[ "${CHANGES:-0}" -eq 0 ] && exit 0

# Has the session-log already been updated today?
LOG="$MEM/session-log.md"
if [ -f "$LOG" ] && head -40 "$LOG" | grep -q "$TODAY"; then
  exit 0   # already logged today
fi

# Emit a one-time, non-blocking reminder as additional context.
touch "$MARKER"
cat <<'JSON'
{
  "hookSpecificOutput": {
    "hookEventName": "Stop",
    "additionalContext": "Reminder: there are uncommitted changes and no session-log entry for today yet. If this is a good stopping point, run /e4:end-session to update progress.md, active-tasks.md, and session-log.md so the next session can resume cleanly."
  }
}
JSON
exit 0

# ---------------------------------------------------------------------------
# HARD ENFORCEMENT (opt-in): block the stop and force Claude to log first.
# Replace the JSON above with this. Note: Stop-hook blocks are capped by
# CLAUDE_CODE_STOP_HOOK_BLOCK_CAP (default 8) to prevent infinite loops.
#
# echo '{"decision":"block","reason":"Run /e4:end-session to persist this session before stopping."}'
# exit 0
# ---------------------------------------------------------------------------
