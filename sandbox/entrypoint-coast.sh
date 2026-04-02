#!/bin/bash
set -euo pipefail

# ---------------------------------------------------------------------------
# Sandbox Derby — Coast Entrypoint
#
# Autonomous mode: clone workspace, execute course, exit.
# ---------------------------------------------------------------------------

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/entrypoint-common.sh"

# Validate coast-specific requirements
: "${TARGET_REPO:?TARGET_REPO is required for coast mode}"

WORKSPACE="/home/agent/workspace"
COURSE_STAGING="/home/agent/course"

# Clone the target repo into the workspace
echo "Cloning ${TARGET_REPO} into ${WORKSPACE}..."
git clone "${TARGET_REPO}" "${WORKSPACE}"

# Copy course from staging into workspace (after clone so the workspace exists)
COURSE_FILE=""
if [ -d "${COURSE_STAGING}" ]; then
    COURSE_FILE=$(find "${COURSE_STAGING}" -name "*.md" -type f | head -1)
    if [ -n "${COURSE_FILE}" ]; then
        COURSE_BASENAME=$(basename "${COURSE_FILE}")
        cp "${COURSE_FILE}" "${WORKSPACE}/${COURSE_BASENAME}"
        echo "Course copied to ${WORKSPACE}/${COURSE_BASENAME}"
    fi
fi

if [ -z "${COURSE_FILE}" ]; then
    echo "ERROR: No course file found at ${COURSE_STAGING}/"
    exit 1
fi

# Read course content and execute via Claude
COURSE_CONTENT=$(cat "${WORKSPACE}/${COURSE_BASENAME}")

echo "============================================"
echo "  Sandbox Derby — Coast Mode"
echo "============================================"
echo "  Workspace: ${TARGET_REPO}"
echo "  Course:    ${COURSE_BASENAME}"
echo "  Loadout:   ${LOADOUT_STATUS:-bare}"
echo "============================================"
echo ""

cd "${WORKSPACE}"

claude -p "You are working in a git repository cloned from ${TARGET_REPO}. Follow the course instructions below. Commit your work as you go.

--- COURSE ---
${COURSE_CONTENT}
--- END COURSE ---"

echo ""
echo "============================================"
echo "  Workspace Summary"
echo "============================================"
git -C /home/agent/workspace log --oneline 2>/dev/null || echo "(no commits found)"
echo ""
echo "============================================"
echo "  Coast complete."
echo "============================================"
