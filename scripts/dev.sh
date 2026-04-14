#!/usr/bin/env bash
# Local development runner - simulates Alfred's Script Filter execution.
#
# Usage:
#   ./scripts/dev.sh "search foo"
#   ./scripts/dev.sh "config"
#   ./scripts/dev.sh ""
#
# Output is pretty-printed JSON if `jq` is available.
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENTRY="$REPO_ROOT/workflow/scripts/entry.py"

QUERY="${1:-}"

# ---------------------------------------------------------------------------
# Read use_uv from the config builder default (info.plist variables).
# Alfred sets this as an env var at runtime; dev.sh simulates that here.
# ---------------------------------------------------------------------------
use_uv=$(python3 -c "
import plistlib, pathlib
d = plistlib.loads(pathlib.Path('$REPO_ROOT/workflow/info.plist').read_bytes())
entry = next((e for e in d.get('userconfigurationconfig', []) if e['variable'] == 'use_uv'), {})
print(1 if entry.get('config', {}).get('default', True) else 0)
" 2>/dev/null || echo 1)
export use_uv

if [[ "$use_uv" == "1" ]] && command -v uv &>/dev/null; then
  PYTHON_CMD="uv run python"
else
  PYTHON_CMD="python3"
fi

# Make src/ importable without installing the package
# In the packaged workflow, entry.py adds workflow_root/src/ to sys.path.
# During development, workflow/src/ does not exist – the real source is at
# repo_root/src/, so we set PYTHONPATH explicitly.
export PYTHONPATH="$REPO_ROOT/src${PYTHONPATH:+:$PYTHONPATH}"

# Set environment variables that Alfred normally provides
export alfred_workflow_bundleid="${alfred_workflow_bundleid:-com.example.dev}"
export alfred_workflow_cache="${alfred_workflow_cache:-/tmp/alfred-dev-cache}"
export alfred_workflow_data="${alfred_workflow_data:-/tmp/alfred-dev-data}"
export alfred_workflow_version="${alfred_workflow_version:-0.0.0-dev}"
export alfred_workflow_uid="${alfred_workflow_uid:-com.example.dev}"
export alfred_workflow_name="${alfred_workflow_name:-Workflow Dev}"

mkdir -p "$alfred_workflow_cache" "$alfred_workflow_data"

echo "─────────────────────────────────────"
echo "  Alfred Script Filter Simulator"
echo "  Query: \"$QUERY\""
echo "─────────────────────────────────────"

if command -v jq &>/dev/null; then
  $PYTHON_CMD "$ENTRY" "$QUERY" | jq .
else
  $PYTHON_CMD "$ENTRY" "$QUERY"
fi
