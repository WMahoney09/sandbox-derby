#!/bin/bash
set -euo pipefail

# ---------------------------------------------------------------------------
# Sandbox Derby — Common Entrypoint
#
# Shared setup sourced by both drive and coast entrypoints.
# Handles: env validation, git identity, loadout copy-in.
# ---------------------------------------------------------------------------

# Validate required environment variables
: "${ANTHROPIC_API_KEY:?ANTHROPIC_API_KEY is required}"

# Configure git identity
git config --global user.name "${GIT_USER_NAME:-Sandbox Derby Agent}"
git config --global user.email "${GIT_USER_EMAIL:-sandbox-derby[bot]@noreply.github.com}"

# Copy loadout from staging to runtime path
# Loadout is mounted read-only at /home/agent/loadout/ and copied to
# /home/agent/.claude/ so Claude Code can write freely at runtime.
LOADOUT_STAGING="/home/agent/loadout"
CLAUDE_DIR="/home/agent/.claude"

if [ -d "${LOADOUT_STAGING}" ] && [ "$(ls -A "${LOADOUT_STAGING}" 2>/dev/null)" ]; then
    mkdir -p "${CLAUDE_DIR}"
    cp -r "${LOADOUT_STAGING}/." "${CLAUDE_DIR}/"
    export LOADOUT_STATUS="loaded"
    echo "Loadout copied from staging to ${CLAUDE_DIR}"
else
    export LOADOUT_STATUS="bare"
    echo "No loadout found at ${LOADOUT_STAGING} — running bare"
fi
