package main

import (
	"fmt"
	"os"

	"github.com/WMahoney09/sandbox-derby/internal/derby"
)

func main() {
	if len(os.Args) < 3 || os.Args[1] != "run" {
		fmt.Fprintf(os.Stderr, "Usage: derby run <config.yaml>\n")
		os.Exit(1)
	}

	configPath := os.Args[2]

	cfg, err := derby.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	runner := derby.NewRunner(cfg)
	results, err := runner.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Derby run failed: %v\n", err)
		os.Exit(1)
	}

	report, err := derby.GenerateReport(cfg, results)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating report: %v\n", err)
		os.Exit(1)
	}

	outputPath, err := derby.WriteReport(cfg.Name, report)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing report: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Derby complete. Report written to: %s\n", outputPath)
}
