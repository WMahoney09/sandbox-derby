package derby

import (
	"fmt"
	"os"
	"path/filepath"
)

// CoastConfig holds the parameters for an autonomous coast run.
type CoastConfig struct {
	Image           string
	Loadout         string
	Course          string
	Repo            string
	EnvFile         string
	SkipPermissions bool
}

// Coast runs an autonomous sandbox session. It launches a single container
// that clones the target repo, executes the course via Claude, and exits.
// The container is removed automatically via --rm.
func Coast(cfg CoastConfig) error {
	if err := CheckImage(cfg.Image); err != nil {
		return fmt.Errorf("image %q not found — run: docker compose build\n  %w", cfg.Image, err)
	}

	absCourse, err := filepath.Abs(cfg.Course)
	if err != nil {
		return fmt.Errorf("resolving course path: %w", err)
	}
	absEnvFile, err := filepath.Abs(cfg.EnvFile)
	if err != nil {
		return fmt.Errorf("resolving env file path: %w", err)
	}

	args := []string{
		"run", "--rm",
		"--env-file", absEnvFile,
		"-e", fmt.Sprintf("TARGET_REPO=%s", cfg.Repo),
	}

	if cfg.SkipPermissions {
		args = append(args, "-e", "SKIP_PERMISSIONS=true")
	}

	if isRemoteLoadout(cfg.Loadout) {
		args = append(args, "-e", fmt.Sprintf("LOADOUT_REPO=%s", cfg.Loadout))
	} else {
		absLoadout, err := filepath.Abs(cfg.Loadout)
		if err != nil {
			return fmt.Errorf("resolving loadout path: %w", err)
		}
		args = append(args, "-v", fmt.Sprintf("%s:/home/agent/loadout:ro", absLoadout))
	}

	args = append(args,
		"-v", fmt.Sprintf("%s:/home/agent/course/course.md:ro", absCourse),
		"--cpus", "2",
		"--memory", "4g",
		"--pids-limit", "256",
		cfg.Image,
		"./entrypoint-coast.sh",
	)

	cmd := newCommand("docker", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("coast run failed: %w", err)
	}

	return nil
}
