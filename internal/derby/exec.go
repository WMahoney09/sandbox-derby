package derby

import (
	"os/exec"
	"strings"
)

// newCommand creates an exec.Cmd. Extracted for future testability.
func newCommand(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}

// isRemoteLoadout returns true if the loadout string is a git URL
// rather than a local path.
func isRemoteLoadout(loadout string) bool {
	return strings.HasPrefix(loadout, "https://") ||
		strings.HasPrefix(loadout, "git@") ||
		strings.HasPrefix(loadout, "http://")
}
