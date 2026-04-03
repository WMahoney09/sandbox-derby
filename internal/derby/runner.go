package derby

import (
	"fmt"
	"sync"
)

// Runner orchestrates a derby — launching sandboxes and collecting results.
type Runner struct {
	config *Config
}

// NewRunner creates a derby runner from a config.
func NewRunner(cfg *Config) *Runner {
	return &Runner{config: cfg}
}

// Run executes the derby: launches all sandboxes concurrently (bounded by
// the concurrency limit) and collects results.
func (r *Runner) Run() ([]SandboxResult, error) {
	// Check that image exists
	if err := CheckImage(r.config.Image); err != nil {
		return nil, fmt.Errorf("image %s not found — build it first: %w", r.config.Image, err)
	}

	// Build sandbox specs with sequential IDs
	sandboxID := 1
	var specs []SandboxSpec
	for _, entry := range r.config.Entries {
		for replica := 1; replica <= entry.Replicas; replica++ {
			specs = append(specs, SandboxSpec{
				ID:              sandboxID,
				Name:            r.config.Name,
				EntryName:       entry.Name,
				ReplicaNum:      replica,
				Image:           r.config.Image,
				LoadoutPath:     entry.Loadout,
				CoursePath:      entry.Course,
				RepoURL:         r.config.Workspace.Repo,
				EnvFile:         r.config.EnvFile,
				SkipPermissions: entry.SkipPermissions,
				Resources:       entry.Resources,
			})
			sandboxID++
		}
	}

	total := len(specs)
	fmt.Printf("Derby '%s': launching %d sandboxes (concurrency: %d)\n",
		r.config.Name, total, r.config.Concurrency)

	// Run sandboxes with bounded concurrency
	results := make([]SandboxResult, total)
	sem := make(chan struct{}, r.config.Concurrency)
	var wg sync.WaitGroup

	for i, spec := range specs {
		wg.Add(1)
		sem <- struct{}{} // acquire semaphore

		go func(idx int, s SandboxSpec) {
			defer wg.Done()
			defer func() { <-sem }() // release semaphore

			fmt.Printf("  Sandbox #%d: starting (%s, replica %d)\n",
				s.ID, s.EntryName, s.ReplicaNum)

			result := RunSandbox(s, "")
			results[idx] = result

			if result.Error != nil {
				fmt.Printf("  Sandbox #%d: failed — %v\n", s.ID, result.Error)
			} else {
				fmt.Printf("  Sandbox #%d: complete — exit %d, %s\n",
					s.ID, result.ExitCode, result.Duration.Round(100*1e6))
			}
		}(i, spec)
	}

	wg.Wait()
	fmt.Printf("Derby '%s': all sandboxes complete.\n", r.config.Name)

	return results, nil
}

// CheckImage verifies that a Docker image exists locally.
func CheckImage(image string) error {
	cmd := newCommand("docker", "image", "inspect", image)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}
