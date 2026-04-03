package derby

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents a derby configuration file.
type Config struct {
	Name        string    `yaml:"name"`
	Image       string    `yaml:"image"`
	EnvFile     string    `yaml:"env_file"`
	Concurrency int       `yaml:"concurrency"`
	Workspace   Workspace `yaml:"workspace"`
	Entries     []Entry   `yaml:"entries"`
}

// Workspace defines the target repository for all sandboxes.
type Workspace struct {
	Repo string `yaml:"repo"`
}

// Entry defines a single derby entry — a loadout/course combination
// that may be replicated multiple times.
type Entry struct {
	Name      string    `yaml:"name"`
	Loadout   string    `yaml:"loadout"`
	Course    string    `yaml:"course"`
	Replicas  int       `yaml:"replicas"`
	Resources Resources `yaml:"resources"`
}

// Resources defines container resource limits.
type Resources struct {
	CPUs   string `yaml:"cpus"`
	Memory string `yaml:"memory"`
}

// LoadConfig reads and validates a derby configuration file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	cfg.applyDefaults()

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	if c.Workspace.Repo == "" {
		return fmt.Errorf("workspace.repo is required")
	}
	if len(c.Entries) == 0 {
		return fmt.Errorf("at least one entry is required")
	}
	for i, e := range c.Entries {
		if e.Name == "" {
			return fmt.Errorf("entry %d: name is required", i)
		}
		if e.Loadout == "" {
			return fmt.Errorf("entry %d (%s): loadout is required", i, e.Name)
		}
		if e.Course == "" {
			return fmt.Errorf("entry %d (%s): course is required", i, e.Name)
		}
	}
	return nil
}

func (c *Config) applyDefaults() {
	if c.Image == "" {
		c.Image = "sandbox-derby:latest"
	}
	if c.EnvFile == "" {
		c.EnvFile = ".env"
	}
	for i := range c.Entries {
		if c.Entries[i].Replicas <= 0 {
			c.Entries[i].Replicas = 1
		}
		if c.Entries[i].Resources.CPUs == "" {
			c.Entries[i].Resources.CPUs = "2"
		}
		if c.Entries[i].Resources.Memory == "" {
			c.Entries[i].Resources.Memory = "4g"
		}
	}
	if c.Concurrency <= 0 {
		total := 0
		for _, e := range c.Entries {
			total += e.Replicas
		}
		c.Concurrency = total
	}
}

// TotalSandboxes returns the total number of sandboxes the derby will run.
func (c *Config) TotalSandboxes() int {
	total := 0
	for _, e := range c.Entries {
		total += e.Replicas
	}
	return total
}
