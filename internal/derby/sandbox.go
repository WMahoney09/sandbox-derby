package derby

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// SandboxSpec defines a single sandbox to run.
type SandboxSpec struct {
	Name        string
	EntryName   string
	ReplicaNum  int
	Image       string
	LoadoutPath string
	CoursePath  string
	RepoURL     string
	Resources   Resources
}

// SandboxResult captures the outcome of a sandbox run.
type SandboxResult struct {
	Spec      SandboxSpec
	ExitCode  int
	Duration  time.Duration
	GitLog    string
	FileList  string
	Stdout    string
	Stderr    string
	Error     error
}

// RunSandbox runs a single sandbox container and collects results.
func RunSandbox(spec SandboxSpec, outputDir string) SandboxResult {
	start := time.Now()

	containerName := fmt.Sprintf("derby-%s-%s-%d", spec.Name, spec.EntryName, spec.ReplicaNum)

	// Resolve absolute paths for volume mounts
	absLoadout, err := filepath.Abs(spec.LoadoutPath)
	if err != nil {
		return SandboxResult{Spec: spec, Error: fmt.Errorf("resolving loadout path: %w", err)}
	}
	absCourse, err := filepath.Abs(spec.CoursePath)
	if err != nil {
		return SandboxResult{Spec: spec, Error: fmt.Errorf("resolving course path: %w", err)}
	}

	// Build docker run arguments
	args := []string{
		"run",
		"--rm",
		"--name", containerName,
		"-e", fmt.Sprintf("ANTHROPIC_API_KEY=%s", os.Getenv("ANTHROPIC_API_KEY")),
		"-e", fmt.Sprintf("GITHUB_TOKEN=%s", os.Getenv("GITHUB_TOKEN")),
		"-e", fmt.Sprintf("TARGET_REPO=%s", spec.RepoURL),
		"-v", fmt.Sprintf("%s:/home/agent/loadout:ro", absLoadout),
		"-v", fmt.Sprintf("%s:/home/agent/course/course.md:ro", absCourse),
		"--cpus", spec.Resources.CPUs,
		"--memory", spec.Resources.Memory,
		"--pids-limit", "256",
		spec.Image,
		"./entrypoint-coast.sh",
	}

	var stdout, stderr bytes.Buffer
	cmd := exec.Command("docker", args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	result := SandboxResult{
		Spec:     spec,
		Duration: time.Since(start),
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.Error = err
			return result
		}
	}

	// Extract git log and file list from a fresh container with the same setup
	// For PoC, we parse what we can from stdout. A more robust approach would
	// use named volumes and docker cp.
	result.GitLog = extractSection(stdout.String(), "git log")
	result.FileList = extractSection(stdout.String(), "file list")

	return result
}

// extractSection is a placeholder for parsing structured output from sandboxes.
// For PoC, the coast entrypoint doesn't produce structured sections, so this
// returns empty. The report uses stdout/stderr directly.
func extractSection(output string, section string) string {
	return ""
}
