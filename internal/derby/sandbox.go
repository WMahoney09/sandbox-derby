package derby

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// SandboxSpec defines a single sandbox to run.
type SandboxSpec struct {
	ID              int
	Name            string
	EntryName       string
	ReplicaNum      int
	Image           string
	LoadoutPath     string
	CoursePath      string
	RepoURL         string
	EnvFile         string
	SkipPermissions bool
	Resources       Resources
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

	containerName := fmt.Sprintf("derby-%s-%d", spec.Name, spec.ID)

	// Resolve absolute paths
	absCourse, err := filepath.Abs(spec.CoursePath)
	if err != nil {
		return SandboxResult{Spec: spec, Error: fmt.Errorf("resolving course path: %w", err)}
	}
	absEnvFile, err := filepath.Abs(spec.EnvFile)
	if err != nil {
		return SandboxResult{Spec: spec, Error: fmt.Errorf("resolving env file path: %w", err)}
	}

	// Build docker run arguments (no --rm: we need the container to stick around
	// so we can docker-cp the workspace out before cleaning up)
	args := []string{
		"run",
		"--name", containerName,
		"--env-file", absEnvFile,
		"-e", fmt.Sprintf("SANDBOX_ID=%d", spec.ID),
		"-e", fmt.Sprintf("TARGET_REPO=%s", spec.RepoURL),
	}

	if isRemoteLoadout(spec.LoadoutPath) {
		args = append(args, "-e", fmt.Sprintf("LOADOUT_REPO=%s", spec.LoadoutPath))
	} else {
		absLoadout, err := filepath.Abs(spec.LoadoutPath)
		if err != nil {
			return SandboxResult{Spec: spec, Error: fmt.Errorf("resolving loadout path: %w", err)}
		}
		args = append(args, "-v", fmt.Sprintf("%s:/home/agent/loadout:ro", absLoadout))
	}

	if spec.SkipPermissions {
		args = append(args, "-e", "SKIP_PERMISSIONS=true")
	}

	args = append(args,
		"-v", fmt.Sprintf("%s:/home/agent/course/course.md:ro", absCourse),
		"--cpus", spec.Resources.CPUs,
		"--memory", spec.Resources.Memory,
		"--pids-limit", "256",
		spec.Image,
		"./entrypoint-coast.sh",
	)

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
			// Still attempt cleanup even on hard errors
			cleanupContainer(containerName)
			return result
		}
	}

	// Extract artifacts from the stopped container, then remove it
	result.GitLog, result.FileList = extractArtifacts(containerName)

	return result
}

// extractArtifacts copies the workspace out of a stopped container, extracts
// the git log and file list, then removes the container and temp directory.
// If any extraction step fails it returns what it managed to collect.
func extractArtifacts(containerName string) (gitLog string, fileList string) {
	// Always clean up the container when we're done
	defer cleanupContainer(containerName)

	// Create a temp directory to hold the copied workspace
	tmpDir, err := os.MkdirTemp("", "derby-extract-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: failed to create temp dir for artifact extraction: %v\n", err)
		return "", ""
	}
	defer os.RemoveAll(tmpDir)

	workspaceDst := filepath.Join(tmpDir, "workspace")

	// docker cp works on stopped containers
	cpCmd := newCommand("docker", "cp", containerName+":/home/agent/workspace", workspaceDst)
	if out, err := cpCmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: docker cp failed: %v\n%s\n", err, string(out))
		return "", ""
	}

	// Extract git log
	gitCmd := newCommand("git", "-C", workspaceDst, "log", "--oneline")
	if out, err := gitCmd.Output(); err == nil {
		gitLog = strings.TrimSpace(string(out))
	} else {
		fmt.Fprintf(os.Stderr, "warning: git log extraction failed: %v\n", err)
	}

	// Extract file list via filepath.Walk, excluding .git directory
	var files []string
	filepath.Walk(workspaceDst, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip files we can't stat
		}
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			rel, relErr := filepath.Rel(workspaceDst, path)
			if relErr == nil {
				files = append(files, rel)
			}
		}
		return nil
	})
	fileList = strings.Join(files, "\n")

	return gitLog, fileList
}

// cleanupContainer removes a container by name, ignoring errors (the container
// may already be gone if docker run was interrupted).
func cleanupContainer(containerName string) {
	rmCmd := newCommand("docker", "rm", "-f", containerName)
	if out, err := rmCmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: container cleanup failed for %s: %v\n%s\n", containerName, err, string(out))
	}
}
