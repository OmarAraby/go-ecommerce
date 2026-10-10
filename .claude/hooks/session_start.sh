#!/usr/bin/env bash
# session_start.sh — Claude Code SessionStart hook.
# Whatever this script prints to STDOUT is injected into Claude's context at the
# start of every session/resume. This is the "doesn't rely on Claude remembering"
# mechanism: memory is loaded automatically, every time.
#
# Input: JSON on stdin (session_id, cwd, source, ...). We don't need it, but read it
# so the pipe doesn't break.
#
# Keep this FAST (it runs on every session). Output is plain text/markdown.

cat >/dev/null  # consume stdin

MEM=".ai-memory"

# If there's no memory here, gently suggest initializing and exit.
if [ ! -d "$MEM" ]; then
  echo "ℹ️ No .ai-memory/ found in this project. Run /e4:init to set up persistent memory."
  exit 0
fi

# Tier-3: active-tasks.md is a GENERATED view of tasks/* (gitignored). Refresh it from the task
# files before reading (best-effort, fast). No-op if the generator isn't installed.
PYDIR="${AI_MEMORY_PYDIR:-$HOME/AI-Vault/automation/python}"
[ -f "$PYDIR/generate_status.py" ] && python3 "$PYDIR/generate_status.py" --project "$PWD" >/dev/null 2>&1

echo "## 🧠 Project Memory (auto-loaded) — e4"
echo "*forged by Omar Araby & contributors*"
echo
echo "Read order if you need detail: onboarding → context → architecture → active-tasks → progress → decisions → sessions/ (newest file)."
echo

# --- Current state from progress.md ---
if [ -f "$MEM/progress.md" ]; then
  echo "### Status"
  grep -E -i "Phase:|completion|Last updated:|Health:" "$MEM/progress.md" | head -4 \
    | sed -E 's/\*\*//g; s/^[[:space:]]*[-*][[:space:]]*//; s/^/- /'
  echo
fi

# --- Current task from active-tasks.md ---
if [ -f "$MEM/active-tasks.md" ]; then
  echo "### Current focus"
  # print the Current Task section up to the next header (tolerant of header style:
  # "## 🎯 Current Task" or "## Current task"; stops at the next "## " header)
  awk '/^##[[:space:]].*[Cc]urrent/{flag=1; next} /^##[[:space:]]/{flag=0} flag' "$MEM/active-tasks.md" | sed '/^[[:space:]]*$/d' | head -8
  echo
fi

# --- Most recent session (the resume anchor) ---
# Tier-2: sessions live as one file per session in sessions/ (timestamp-prefixed names
# sort chronologically, so the newest is the lexically-last). Falls back to the legacy
# single session-log.md for repos not yet migrated.
LATEST=""
if [ -d "$MEM/sessions" ]; then
  LATEST="$(ls -1 "$MEM/sessions"/*.md 2>/dev/null | LC_ALL=C sort | tail -1)"
fi
if [ -n "$LATEST" ] && [ -f "$LATEST" ]; then
  echo "### Last session"
  echo "_$(basename "$LATEST")_"
  head -40 "$LATEST"
  echo
  if ! grep -qi "TL;DR" "$LATEST"; then
    echo "> ⚠️ The latest session file looks incomplete (no TL;DR). The previous session may have ended without /e4:end-session. Reconstruct what happened from git log and write a proper session file."
    echo
  fi
elif [ -f "$MEM/session-log.md" ]; then
  # Legacy single-file fallback (pre-Tier-2 repos)
  echo "### Last session"
  awk '/^## Session/{c++} c==1{print} c==2{exit}' "$MEM/session-log.md" | head -30
  echo
  LAST_BLOCK="$(awk '/^## Session/{c++} c==1{print} c==2{exit}' "$MEM/session-log.md")"
  if ! echo "$LAST_BLOCK" | grep -qi "TL;DR"; then
    echo "> ⚠️ The last session-log entry looks incomplete (no TL;DR). The previous session may have ended without /e4:end-session. Reconstruct what happened from git log and offer to write a proper entry."
    echo
  fi
fi

echo "When this session involves real work, run /e4:end-session before finishing so the next session can resume."
exit 0
