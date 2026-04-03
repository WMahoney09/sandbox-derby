package derby

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// DriveConfig holds the parameters for an interactive drive session.
type DriveConfig struct {
	Image   string
	Loadout string
	EnvFile string
}

// Drive starts an interactive sandbox session. It launches a detached
// container, execs bash into it for interactive use, then cleans up
// the container when the session ends.
func Drive(cfg DriveConfig) error {
	if err := CheckImage(cfg.Image); err != nil {
		return fmt.Errorf("image %q not found — run: docker compose build\n  %w", cfg.Image, err)
	}

	absLoadout, err := filepath.Abs(cfg.Loadout)
	if err != nil {
		return fmt.Errorf("resolving loadout path: %w", err)
	}
	absEnvFile, err := filepath.Abs(cfg.EnvFile)
	if err != nil {
		return fmt.Errorf("resolving env file path: %w", err)
	}

	// Unique container name to allow multiple concurrent drive sessions
	containerName := fmt.Sprintf("derby-drive-%d", time.Now().UnixNano())

	// Start the container in detached mode
	runArgs := []string{
		"run", "-d",
		"--name", containerName,
		"--env-file", absEnvFile,
		"-v", fmt.Sprintf("%s:/home/agent/loadout:ro", absLoadout),
		"--cpus", "2",
		"--memory", "4g",
		"--pids-limit", "256",
		cfg.Image,
		"./entrypoint-drive.sh",
	}

	startCmd := newCommand("docker", runArgs...)
	startCmd.Stderr = os.Stderr
	if out, err := startCmd.Output(); err != nil {
		return fmt.Errorf("starting container: %w", err)
	} else {
		_ = out // container ID, not needed
	}

	// Always clean up the container when we're done
	defer func() {
		stopCmd := newCommand("docker", "rm", "-f", containerName)
		stopCmd.Stdout = nil
		stopCmd.Stderr = nil
		_ = stopCmd.Run()
	}()

	fmt.Printf("Drive session starting (container: %s)...\n", containerName)

	// Exec into the container interactively
	execCmd := newCommand("docker", "exec", "-it", containerName, "bash")
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	if err := execCmd.Run(); err != nil {
		// The user ended the session — any exit from the interactive shell
		// (exit, Ctrl+D, Ctrl+C) is a normal outcome, not an error.
		if _, ok := err.(*exec.ExitError); ok {
			fmt.Println("Drive session ended. Container removed.")
			return nil
		}
		// A non-ExitError means something else went wrong (e.g., docker not found)
		return fmt.Errorf("exec session failed: %w", err)
	}

	fmt.Println("Drive session ended. Container removed.")
	return nil
}
