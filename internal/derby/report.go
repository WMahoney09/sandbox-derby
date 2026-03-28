package derby

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GenerateReport produces a markdown report from derby results.
func GenerateReport(cfg *Config, results []SandboxResult) (string, error) {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("# Derby Report: %s\n\n", cfg.Name))
	b.WriteString(fmt.Sprintf("**Date:** %s\n\n", time.Now().Format("2006-01-02 15:04:05")))
	b.WriteString(fmt.Sprintf("**Image:** %s\n\n", cfg.Image))
	b.WriteString(fmt.Sprintf("**Workspace:** %s\n\n", cfg.Workspace.Repo))
	b.WriteString(fmt.Sprintf("**Total sandboxes:** %d\n\n", len(results)))

	// Configuration summary
	b.WriteString("## Configuration\n\n")
	b.WriteString("| Entry | Loadout | Course | Replicas | CPUs | Memory |\n")
	b.WriteString("|-------|---------|--------|----------|------|--------|\n")
	for _, e := range cfg.Entries {
		b.WriteString(fmt.Sprintf("| %s | %s | %s | %d | %s | %s |\n",
			e.Name, e.Loadout, e.Course, e.Replicas, e.Resources.CPUs, e.Resources.Memory))
	}
	b.WriteString("\n")

	// Results summary
	b.WriteString("## Results Summary\n\n")
	b.WriteString("| Entry | Replica | Exit Code | Duration | Status |\n")
	b.WriteString("|-------|---------|-----------|----------|--------|\n")
	for _, r := range results {
		status := "OK"
		if r.Error != nil {
			status = fmt.Sprintf("Error: %v", r.Error)
		} else if r.ExitCode != 0 {
			status = fmt.Sprintf("Failed (exit %d)", r.ExitCode)
		}
		b.WriteString(fmt.Sprintf("| %s | %d | %d | %s | %s |\n",
			r.Spec.EntryName, r.Spec.ReplicaNum, r.ExitCode,
			r.Duration.Round(time.Second), status))
	}
	b.WriteString("\n")

	// Per-entry comparison
	b.WriteString("## Per-Entry Details\n\n")

	// Group results by entry name
	grouped := make(map[string][]SandboxResult)
	for _, r := range results {
		grouped[r.Spec.EntryName] = append(grouped[r.Spec.EntryName], r)
	}

	for _, entry := range cfg.Entries {
		entryResults := grouped[entry.Name]
		b.WriteString(fmt.Sprintf("### %s\n\n", entry.Name))
		b.WriteString(fmt.Sprintf("**Loadout:** %s\n\n", entry.Loadout))
		b.WriteString(fmt.Sprintf("**Course:** %s\n\n", entry.Course))

		for _, r := range entryResults {
			b.WriteString(fmt.Sprintf("#### Replica %d\n\n", r.Spec.ReplicaNum))
			b.WriteString(fmt.Sprintf("- **Exit code:** %d\n", r.ExitCode))
			b.WriteString(fmt.Sprintf("- **Duration:** %s\n", r.Duration.Round(time.Second)))

			if r.Error != nil {
				b.WriteString(fmt.Sprintf("- **Error:** %v\n", r.Error))
			}

			if r.Stdout != "" {
				b.WriteString("\n<details><summary>Stdout</summary>\n\n```\n")
				b.WriteString(truncate(r.Stdout, 5000))
				b.WriteString("\n```\n\n</details>\n")
			}
			if r.Stderr != "" {
				b.WriteString("\n<details><summary>Stderr</summary>\n\n```\n")
				b.WriteString(truncate(r.Stderr, 5000))
				b.WriteString("\n```\n\n</details>\n")
			}
			b.WriteString("\n")
		}
	}

	return b.String(), nil
}

// WriteReport saves the report to the derby-results directory.
func WriteReport(derbyName string, report string) (string, error) {
	timestamp := time.Now().Format("20060102-150405")
	dirName := fmt.Sprintf("%s-%s", derbyName, timestamp)
	outputDir := filepath.Join("derby-results", dirName)

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return "", fmt.Errorf("creating output directory: %w", err)
	}

	reportPath := filepath.Join(outputDir, "report.md")
	if err := os.WriteFile(reportPath, []byte(report), 0o644); err != nil {
		return "", fmt.Errorf("writing report: %w", err)
	}

	return reportPath, nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "\n... (truncated)"
}
