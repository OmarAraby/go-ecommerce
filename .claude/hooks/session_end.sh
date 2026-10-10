#!/usr/bin/env bash
# session_end.sh — Claude Code SessionEnd hook.
# Runs automatically when the session closes (reasons: clear|resume|logout|prompt_input_exit).
#
# IMPORTANT LIMITATION: the session is ending, so Claude will NOT generate more text.
# This script therefore does only MECHANICAL work — it cannot write an intelligent
# summary (that's what /e4:end-session is for, while Claude is still active).
#
# What it guarantees:
#   1. A breadcrumb FILE is written to .ai-memory/sessions/ IF no real session was logged
#      today, so a session is never completely unrecorded — even if the user closed the
#      terminal without running /e4:end-session. (One file per session → never a merge conflict.)
#   2. Completed tasks are archived and the daily summary is refreshed (best-effort).
#
# Heavy/slow work is backgrounded so the hook returns immediately (avoids "Hook cancelled").
#
# Input: JSON on stdin (session_id, reason, cwd, transcript_path, ...).

INPUT="$(cat)"                       # consume stdin
MEM=".ai-memory"
TODAY="$(date +%F)"
NOW="$(date '+%Y-%m-%d %H:%M')"

# Try to pull the reason out of the JSON without requiring jq.
REASON="$(printf '%s' "$INPUT" | sed -n 's/.*"reason"[: ]*"\([^"]*\)".*/\1/p')"
[ -z "$REASON" ] && REASON="unknown"

[ ! -d "$MEM" ] && exit 0

# Tier-2: sessions are one file per session in sessions/. Decide whether a REAL session
# was already logged today; if not, drop a breadcrumb FILE (never edit a shared file → no conflicts).
BRANCH="$(git rev-parse --abbrev-ref HEAD 2>/dev/null)"
[ -z "$BRANCH" ] && BRANCH="nogit"
BRANCH_SAN="$(printf '%s' "$BRANCH" | tr -c 'A-Za-z0-9._-' '-')"

logged=false
# A real entry for today = a sessions/ file named with today's date that contains a TL;DR.
if [ -d "$MEM/sessions" ]; then
  for f in "$MEM/sessions/$TODAY"*.md; do
    [ -f "$f" ] || continue
    grep -qi "TL;DR" "$f" && { logged=true; break; }
  done
fi
# Legacy fallback: a same-day TL;DR entry in the old single session-log.md also counts.
if [ "$logged" = false ] && [ -f "$MEM/session-log.md" ]; then
  LAST_BLOCK="$(awk '/^## Session/{c++} c==1{print} c==2{exit}' "$MEM/session-log.md")"
  echo "$LAST_BLOCK" | grep -q "$TODAY" && echo "$LAST_BLOCK" | grep -qi "TL;DR" && logged=true
fi

if [ "$logged" = false ]; then
  mkdir -p "$MEM/sessions"
  STAMP="$(date '+%Y-%m-%d-%H%M%S')"
  DIRTY="$(git status --short 2>/dev/null | wc -l | tr -d ' ')"
  CRUMB="$MEM/sessions/${STAMP}-${BRANCH_SAN}-auto.md"
  cat > "$CRUMB" <<EOF
---
type: session
auto: true
date: $TODAY
branch: $BRANCH
---
## Session — $TODAY (auto-recorded by SessionEnd hook)
**Time**: ended $NOW
**Branch**: $BRANCH
**TL;DR**: ⚠️ Session ended via '$REASON' without /e4:end-session. This is an auto-breadcrumb — next session should reconstruct details from git log and write a proper session file.
**Uncommitted files at end**: $DIRTY
**Next actions**:
1. Review git log/diff for this session and backfill a proper session file.
EOF
fi

# --- Best-effort maintenance, backgrounded so we exit fast ---
VAULT="${AI_VAULT:-$HOME/AI-Vault}"
PROJECTS_ROOT="${AI_PROJECTS_ROOT:-$(dirname "$PWD")}"
PYDIR="${AI_MEMORY_PYDIR:-$HOME/AI-Vault/automation/python}"
(
  if [ -f "$PYDIR/archive_completed_tasks.py" ]; then
    python3 "$PYDIR/archive_completed_tasks.py" --project "$PWD" --vault "$VAULT" --infer
  fi
  if [ -f "$PYDIR/generate_daily_summary.py" ]; then
    python3 "$PYDIR/generate_daily_summary.py" --vault "$VAULT" --projects-root "$PROJECTS_ROOT"
  fi
) >/dev/null 2>&1 &

exit 0
