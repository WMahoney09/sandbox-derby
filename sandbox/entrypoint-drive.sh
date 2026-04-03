#!/bin/bash
set -euo pipefail

# ---------------------------------------------------------------------------
# Sandbox Derby — Drive Entrypoint
#
# Interactive mode: set up the sandbox and keep it alive for docker exec.
# ---------------------------------------------------------------------------

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/entrypoint-common.sh"

SANDBOX_LABEL="Sandbox"
if [ -n "${SANDBOX_ID:-}" ]; then
    SANDBOX_LABEL="Sandbox #${SANDBOX_ID}"
fi

echo "============================================"
echo "  ${SANDBOX_LABEL} — Drive Mode"
echo "============================================"
echo ""
echo "  Available tools:"
echo "    claude   — Claude Code CLI"
echo "    git      — version control"
echo "    gh       — GitHub CLI"
echo "    curl     — HTTP client"
echo "    python3  — Python runtime"
echo "    node     — Node.js runtime"
echo ""
echo "  Loadout: ${LOADOUT_STATUS:-bare}"
echo ""
echo "  Attach with: docker exec -it <container> bash"
echo "============================================"

# Keep container alive for interactive use
exec tail -f /dev/null
