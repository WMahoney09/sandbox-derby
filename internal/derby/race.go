package derby

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"gopkg.in/yaml.v3"
)

// RaceState represents the persisted state of a race across its lifecycle.
type RaceState struct {
	Name      string        `yaml:"name"`
	Config    string        `yaml:"config"`
	Status    string        `yaml:"status"`
	CreatedAt time.Time     `yaml:"created_at"`
	StartedAt *time.Time    `yaml:"started_at,omitempty"`
	EndedAt   *time.Time    `yaml:"ended_at,omitempty"`
	OutputDir string        `yaml:"output_dir"`
	Sandboxes []RaceSandbox `yaml:"sandboxes"`
}

// RaceSandbox tracks the state of a single sandbox within a race.
type RaceSandbox struct {
	ID              int    `yaml:"id"`
	EntryName       string `yaml:"entry_name"`
	ReplicaNum      int    `yaml:"replica_num"`
	ContainerName   string `yaml:"container_name"`
	Loadout         string `yaml:"loadout"`
	Course          string `yaml:"course"`
	SkipPermissions bool   `yaml:"skip_permissions"`
	Status          string `yaml:"status"`
	ExitCode        int    `yaml:"exit_code,omitempty"`
	Duration        string `yaml:"duration,omitempty"`
}

// SaveRaceState writes the race state to <outputDir>/race.yaml.
func SaveRaceState(state *RaceState) error {
	data, err := yaml.Marshal(state)
	if err != nil {
		return fmt.Errorf("marshalling race state: %w", err)
	}

	path := filepath.Join(state.OutputDir, "race.yaml")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing race state: %w", err)
	}

	return nil
}

// LoadRaceState reads race state from <outputDir>/race.yaml.
func LoadRaceState(outputDir string) (*RaceState, error) {
	path := filepath.Join(outputDir, "race.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading race state: %w", err)
	}

	var state RaceState
	if err := yaml.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("parsing race state: %w", err)
	}

	return &state, nil
}

// FindLatestRace scans derby-results/ for the most recently modified
// race.yaml and returns its output directory path.
func FindLatestRace() (string, error) {
	resultsDir := "derby-results"
	entries, err := os.ReadDir(resultsDir)
	if err != nil {
		return "", fmt.Errorf("reading derby-results directory: %w", err)
	}

	type raceDir struct {
		path    string
		modTime time.Time
	}

	var races []raceDir
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dir := filepath.Join(resultsDir, entry.Name())
		racePath := filepath.Join(dir, "race.yaml")
		info, err := os.Stat(racePath)
		if err != nil {
			continue // no race.yaml in this directory
		}
		races = append(races, raceDir{path: dir, modTime: info.ModTime()})
	}

	if len(races) == 0 {
		return "", fmt.Errorf("no races found in derby-results/")
	}

	sort.Slice(races, func(i, j int) bool {
		return races[i].modTime.After(races[j].modTime)
	})

	return races[0].path, nil
}
