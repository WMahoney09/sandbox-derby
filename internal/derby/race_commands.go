package derby

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// RaceSetup loads the config, validates it, builds the lineup, and saves
// the initial race state with status "setup".
func RaceSetup(configPath string) error {
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	if err := CheckImage(cfg.Image); err != nil {
		return fmt.Errorf("image %s not found — build it first: %w", cfg.Image, err)
	}

	// Create output directory
	timestamp := time.Now().Format("20060102-150405")
	dirName := fmt.Sprintf("%s-%s", cfg.Name, timestamp)
	outputDir := filepath.Join("derby-results", dirName)

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	// Resolve the config path to absolute so it can be loaded from any cwd later
	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		return fmt.Errorf("resolving config path: %w", err)
	}

	// Build the lineup
	sandboxID := 1
	var sandboxes []RaceSandbox
	for _, entry := range cfg.Entries {
		for replica := 1; replica <= entry.Replicas; replica++ {
			containerName := fmt.Sprintf("derby-%s-%d", cfg.Name, sandboxID)
			sandboxes = append(sandboxes, RaceSandbox{
				ID:              sandboxID,
				EntryName:       entry.Name,
				ReplicaNum:      replica,
				ContainerName:   containerName,
				Loadout:         entry.Loadout,
				Course:          entry.Course,
				SkipPermissions: entry.SkipPermissions,
				Status:          "pending",
			})
			sandboxID++
		}
	}

	state := &RaceState{
		Name:      cfg.Name,
		Config:    absConfigPath,
		Status:    "setup",
		CreatedAt: time.Now(),
		OutputDir: outputDir,
		Sandboxes: sandboxes,
	}

	if err := SaveRaceState(state); err != nil {
		return fmt.Errorf("saving race state: %w", err)
	}

	// Print the lineup card
	fmt.Printf("Race: %s\n", cfg.Name)
	fmt.Printf("Sandboxes: %d\n\n", len(sandboxes))
	for _, sb := range sandboxes {
		perm := "safe"
		if sb.SkipPermissions {
			perm = "skip"
		}
		fmt.Printf("  #%-3d %-18s loadout: %-24s permissions: %s\n",
			sb.ID, sb.EntryName, sb.Loadout, perm)
	}
	fmt.Printf("\nState saved to: %s/race.yaml\n", outputDir)

	return nil
}

// RaceStart loads race state, launches all containers synchronously
// (with bounded concurrency), collects results, generates a report,
// and updates the state to "concluded".
func RaceStart(outputDir string) error {
	if outputDir == "" {
		dir, err := findRaceByStatus("setup")
		if err != nil {
			return fmt.Errorf("finding race in setup state: %w", err)
		}
		outputDir = dir
	}

	state, err := LoadRaceState(outputDir)
	if err != nil {
		return err
	}

	if state.Status != "setup" {
		return fmt.Errorf("race status is %q, expected \"setup\"", state.Status)
	}

	// Load the original config to get shared settings
	cfg, err := LoadConfig(state.Config)
	if err != nil {
		return fmt.Errorf("loading original config: %w", err)
	}

	// Mark as started
	now := time.Now()
	state.Status = "started"
	state.StartedAt = &now
	if err := SaveRaceState(state); err != nil {
		return fmt.Errorf("saving started state: %w", err)
	}

	total := len(state.Sandboxes)
	fmt.Printf("Race '%s': launching %d sandboxes (concurrency: %d)\n",
		state.Name, total, cfg.Concurrency)

	// Build specs from race sandboxes
	specs := make([]SandboxSpec, total)
	for i, sb := range state.Sandboxes {
		// Find the matching entry for resource limits
		var resources Resources
		for _, entry := range cfg.Entries {
			if entry.Name == sb.EntryName {
				resources = entry.Resources
				break
			}
		}

		specs[i] = SandboxSpec{
			ID:              sb.ID,
			Name:            state.Name,
			EntryName:       sb.EntryName,
			ReplicaNum:      sb.ReplicaNum,
			Image:           cfg.Image,
			LoadoutPath:     sb.Loadout,
			CoursePath:      sb.Course,
			RepoURL:         cfg.Workspace.Repo,
			EnvFile:         cfg.EnvFile,
			SkipPermissions: sb.SkipPermissions,
			Resources:       resources,
		}
	}

	// Run sandboxes with bounded concurrency (same pattern as runner.go)
	results := make([]SandboxResult, total)
	sem := make(chan struct{}, cfg.Concurrency)
	var wg sync.WaitGroup

	for i, spec := range specs {
		wg.Add(1)
		sem <- struct{}{} // acquire semaphore

		go func(idx int, s SandboxSpec) {
			defer wg.Done()
			defer func() { <-sem }() // release semaphore

			// Mark sandbox as running
			state.Sandboxes[idx].Status = "running"

			fmt.Printf("  Sandbox #%d: starting (%s, replica %d)\n",
				s.ID, s.EntryName, s.ReplicaNum)

			result := RunSandbox(s, state.OutputDir)
			results[idx] = result

			// Update sandbox state
			state.Sandboxes[idx].ExitCode = result.ExitCode
			state.Sandboxes[idx].Duration = result.Duration.Round(time.Second).String()

			if result.Error != nil {
				state.Sandboxes[idx].Status = "dnf"
				fmt.Printf("  Sandbox #%d: failed — %v\n", s.ID, result.Error)
			} else {
				state.Sandboxes[idx].Status = "finished"
				fmt.Printf("  Sandbox #%d: complete — exit %d, %s\n",
					s.ID, result.ExitCode, result.Duration.Round(100*1e6))
			}
		}(i, spec)
	}

	wg.Wait()

	// Generate and write the report
	report, err := GenerateReport(cfg, results)
	if err != nil {
		return fmt.Errorf("generating report: %w", err)
	}

	reportPath := filepath.Join(state.OutputDir, "report.md")
	if err := os.WriteFile(reportPath, []byte(report), 0o644); err != nil {
		return fmt.Errorf("writing report: %w", err)
	}

	// Update state to concluded
	endTime := time.Now()
	state.Status = "concluded"
	state.EndedAt = &endTime
	if err := SaveRaceState(state); err != nil {
		return fmt.Errorf("saving concluded state: %w", err)
	}

	fmt.Printf("Race '%s': all sandboxes complete.\n", state.Name)
	fmt.Printf("Report written to: %s\n", reportPath)

	return nil
}

// RaceStatus prints the current state of each sandbox in the race.
// If the race is "started", it checks live container status.
func RaceStatus(outputDir string) error {
	if outputDir == "" {
		dir, err := FindLatestRace()
		if err != nil {
			return err
		}
		outputDir = dir
	}

	state, err := LoadRaceState(outputDir)
	if err != nil {
		return err
	}

	// If race is started, check live container status
	if state.Status == "started" {
		for i, sb := range state.Sandboxes {
			if sb.Status != "running" {
				continue
			}
			cmd := newCommand("docker", "inspect", "--format", "{{.State.Status}}", sb.ContainerName)
			out, err := cmd.Output()
			if err != nil {
				state.Sandboxes[i].Status = "unknown"
				continue
			}
			dockerStatus := strings.TrimSpace(string(out))
			if dockerStatus == "exited" {
				state.Sandboxes[i].Status = "finished"
			}
		}
	}

	fmt.Printf("Race: %s\n", state.Name)
	fmt.Printf("Status: %s\n", state.Status)
	fmt.Printf("Created: %s\n", state.CreatedAt.Format("2006-01-02 15:04:05"))
	if state.StartedAt != nil {
		fmt.Printf("Started: %s\n", state.StartedAt.Format("2006-01-02 15:04:05"))
	}
	if state.EndedAt != nil {
		fmt.Printf("Ended: %s\n", state.EndedAt.Format("2006-01-02 15:04:05"))
	}
	fmt.Printf("Output: %s\n\n", state.OutputDir)

	fmt.Printf("  %-4s %-18s %-8s %-10s %-10s\n", "ID", "Entry", "Replica", "Status", "Duration")
	fmt.Printf("  %-4s %-18s %-8s %-10s %-10s\n", "---", "-----", "-------", "------", "--------")
	for _, sb := range state.Sandboxes {
		dur := sb.Duration
		if dur == "" {
			dur = "-"
		}
		fmt.Printf("  #%-3d %-18s %-8d %-10s %-10s\n",
			sb.ID, sb.EntryName, sb.ReplicaNum, sb.Status, dur)
	}

	return nil
}

// RaceConclude stops running containers, extracts artifacts, generates
// a report, and marks the race as concluded.
func RaceConclude(outputDir string) error {
	if outputDir == "" {
		dir, err := findActiveRace()
		if err != nil {
			return err
		}
		outputDir = dir
	}

	state, err := LoadRaceState(outputDir)
	if err != nil {
		return err
	}

	if state.Status != "started" && state.Status != "setup" {
		return fmt.Errorf("race status is %q, expected \"started\" or \"setup\"", state.Status)
	}

	cfg, err := LoadConfig(state.Config)
	if err != nil {
		return fmt.Errorf("loading original config: %w", err)
	}

	fmt.Printf("Concluding race '%s'...\n", state.Name)

	// Stop any running containers and extract artifacts
	var results []SandboxResult
	for i := range state.Sandboxes {
		sb := &state.Sandboxes[i]

		if sb.Status == "running" {
			fmt.Printf("  Stopping sandbox #%d (%s)...\n", sb.ID, sb.ContainerName)
			stopCmd := newCommand("docker", "stop", sb.ContainerName)
			stopCmd.Stdout = nil
			stopCmd.Stderr = nil
			_ = stopCmd.Run()
			sb.Status = "dnf"
		}

		// Build a result for each sandbox
		var resources Resources
		for _, entry := range cfg.Entries {
			if entry.Name == sb.EntryName {
				resources = entry.Resources
				break
			}
		}

		result := SandboxResult{
			Spec: SandboxSpec{
				ID:              sb.ID,
				Name:            state.Name,
				EntryName:       sb.EntryName,
				ReplicaNum:      sb.ReplicaNum,
				Image:           cfg.Image,
				LoadoutPath:     sb.Loadout,
				CoursePath:      sb.Course,
				RepoURL:         cfg.Workspace.Repo,
				EnvFile:         cfg.EnvFile,
				SkipPermissions: sb.SkipPermissions,
				Resources:       resources,
			},
			ExitCode: sb.ExitCode,
		}

		if sb.Status == "dnf" {
			result.Error = fmt.Errorf("did not finish (stopped)")
		}

		// Extract artifacts from containers that ran
		if sb.Status == "finished" || sb.Status == "dnf" {
			result.GitLog, result.FileList = extractArtifactsNoCleanup(sb.ContainerName)
			cleanupContainer(sb.ContainerName)
		}

		results = append(results, result)
	}

	// Generate and write the report
	report, err := GenerateReport(cfg, results)
	if err != nil {
		return fmt.Errorf("generating report: %w", err)
	}

	reportPath := filepath.Join(state.OutputDir, "report.md")
	if err := os.WriteFile(reportPath, []byte(report), 0o644); err != nil {
		return fmt.Errorf("writing report: %w", err)
	}

	// Update state to concluded
	endTime := time.Now()
	state.Status = "concluded"
	state.EndedAt = &endTime
	if err := SaveRaceState(state); err != nil {
		return fmt.Errorf("saving concluded state: %w", err)
	}

	fmt.Printf("Race '%s': concluded.\n", state.Name)
	fmt.Printf("Report written to: %s\n", reportPath)

	return nil
}

// RaceResults prints the path to the report for a concluded race.
func RaceResults(outputDir string) error {
	if outputDir == "" {
		dir, err := FindLatestRace()
		if err != nil {
			return err
		}
		outputDir = dir
	}

	state, err := LoadRaceState(outputDir)
	if err != nil {
		return err
	}

	reportPath := filepath.Join(outputDir, "report.md")
	if _, err := os.Stat(reportPath); os.IsNotExist(err) {
		return fmt.Errorf("no report found at %s (race status: %s)", reportPath, state.Status)
	}

	fmt.Printf("Race: %s\n", state.Name)
	fmt.Printf("Status: %s\n", state.Status)
	fmt.Printf("Report: %s\n", reportPath)

	// Print a quick summary
	finished := 0
	dnf := 0
	for _, sb := range state.Sandboxes {
		switch sb.Status {
		case "finished":
			finished++
		case "dnf":
			dnf++
		}
	}
	fmt.Printf("\nSandboxes: %d total, %d finished, %d dnf\n",
		len(state.Sandboxes), finished, dnf)

	return nil
}

// extractArtifactsNoCleanup copies the workspace out of a stopped container
// and extracts git log and file list, but does NOT remove the container.
// The caller is responsible for cleanup.
func extractArtifactsNoCleanup(containerName string) (gitLog string, fileList string) {
	tmpDir, err := os.MkdirTemp("", "derby-extract-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: failed to create temp dir for artifact extraction: %v\n", err)
		return "", ""
	}
	defer os.RemoveAll(tmpDir)

	workspaceDst := filepath.Join(tmpDir, "workspace")

	cpCmd := newCommand("docker", "cp", containerName+":/home/agent/workspace", workspaceDst)
	if out, err := cpCmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: docker cp failed: %v\n%s\n", err, string(out))
		return "", ""
	}

	gitCmd := newCommand("git", "-C", workspaceDst, "log", "--oneline")
	if out, err := gitCmd.Output(); err == nil {
		gitLog = strings.TrimSpace(string(out))
	} else {
		fmt.Fprintf(os.Stderr, "warning: git log extraction failed: %v\n", err)
	}

	var files []string
	filepath.Walk(workspaceDst, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
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

// findRaceByStatus finds the most recent race with the given status.
func findRaceByStatus(status string) (string, error) {
	resultsDir := "derby-results"
	entries, err := os.ReadDir(resultsDir)
	if err != nil {
		return "", fmt.Errorf("reading derby-results directory: %w", err)
	}

	type raceDir struct {
		path    string
		modTime time.Time
	}

	var matches []raceDir
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dir := filepath.Join(resultsDir, entry.Name())
		state, err := LoadRaceState(dir)
		if err != nil {
			continue
		}
		if state.Status == status {
			info, _ := os.Stat(filepath.Join(dir, "race.yaml"))
			if info != nil {
				matches = append(matches, raceDir{path: dir, modTime: info.ModTime()})
			}
		}
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no race with status %q found in derby-results/", status)
	}

	// Sort by most recent
	for i := 0; i < len(matches)-1; i++ {
		for j := i + 1; j < len(matches); j++ {
			if matches[j].modTime.After(matches[i].modTime) {
				matches[i], matches[j] = matches[j], matches[i]
			}
		}
	}

	return matches[0].path, nil
}

// findActiveRace finds the most recent race that is either "started" or "setup".
func findActiveRace() (string, error) {
	resultsDir := "derby-results"
	entries, err := os.ReadDir(resultsDir)
	if err != nil {
		return "", fmt.Errorf("reading derby-results directory: %w", err)
	}

	type raceDir struct {
		path    string
		modTime time.Time
	}

	var matches []raceDir
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dir := filepath.Join(resultsDir, entry.Name())
		state, err := LoadRaceState(dir)
		if err != nil {
			continue
		}
		if state.Status == "started" || state.Status == "setup" {
			info, _ := os.Stat(filepath.Join(dir, "race.yaml"))
			if info != nil {
				matches = append(matches, raceDir{path: dir, modTime: info.ModTime()})
			}
		}
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no active race found in derby-results/")
	}

	for i := 0; i < len(matches)-1; i++ {
		for j := i + 1; j < len(matches); j++ {
			if matches[j].modTime.After(matches[i].modTime) {
				matches[i], matches[j] = matches[j], matches[i]
			}
		}
	}

	return matches[0].path, nil
}
