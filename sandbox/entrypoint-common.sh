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

# Load loadout into runtime path (~/.claude/)
# Two sources: a git URL (LOADOUT_REPO env var) or a volume mount at the staging path.
# LOADOUT_REPO takes precedence if set.
LOADOUT_STAGING="/home/agent/loadout"
CLAUDE_DIR="/home/agent/.claude"

if [ -n "${LOADOUT_REPO:-}" ]; then
    mkdir -p "${CLAUDE_DIR}"
    echo "Cloning loadout from ${LOADOUT_REPO}..."
    git clone --depth 1 "${LOADOUT_REPO}" /tmp/loadout-clone
    cp -r /tmp/loadout-clone/. "${CLAUDE_DIR}/"
    rm -rf /tmp/loadout-clone
    export LOADOUT_STATUS="loaded (remote: ${LOADOUT_REPO})"
    echo "Loadout cloned into ${CLAUDE_DIR}"
elif [ -d "${LOADOUT_STAGING}" ] && [ "$(ls -A "${LOADOUT_STAGING}" 2>/dev/null)" ]; then
    mkdir -p "${CLAUDE_DIR}"
    cp -r "${LOADOUT_STAGING}/." "${CLAUDE_DIR}/"
    export LOADOUT_STATUS="loaded"
    echo "Loadout copied from staging to ${CLAUDE_DIR}"
else
    export LOADOUT_STATUS="bare"
    echo "No loadout found — running bare"
fi
