package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/WMahoney09/sandbox-derby/internal/derby"
)

const usage = `Usage: derby <command> [options]

Commands:
  drive   Start an interactive sandbox session
  coast   Run an autonomous sandbox with a course
  run     Launch a derby from a config file

Run 'derby <command> -help' for command-specific options.`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "drive":
		cmdDrive(os.Args[2:])
	case "coast":
		cmdCoast(os.Args[2:])
	case "run":
		cmdRun(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n%s\n", os.Args[1], usage)
		os.Exit(1)
	}
}

func cmdDrive(args []string) {
	fs := flag.NewFlagSet("drive", flag.ExitOnError)
	loadout := fs.String("loadout", "./loadouts/bare", "Path to loadout directory")
	image := fs.String("image", "sandbox-derby:latest", "Docker image to use")
	envFile := fs.String("env-file", ".env", "Path to environment file")
	fs.Parse(args)

	cfg := derby.DriveConfig{
		Image:   *image,
		Loadout: *loadout,
		EnvFile: *envFile,
	}

	if err := derby.Drive(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func cmdCoast(args []string) {
	fs := flag.NewFlagSet("coast", flag.ExitOnError)
	loadout := fs.String("loadout", "./loadouts/bare", "Path to loadout directory")
	course := fs.String("course", "", "Path to course file (required)")
	repo := fs.String("repo", "", "Target repository URL (required)")
	skipPermissions := fs.Bool("skip-permissions", false, "Skip Claude permission prompts (dangerous)")
	image := fs.String("image", "sandbox-derby:latest", "Docker image to use")
	envFile := fs.String("env-file", ".env", "Path to environment file")
	fs.Parse(args)

	if *course == "" {
		fmt.Fprintln(os.Stderr, "Error: --course is required")
		fs.Usage()
		os.Exit(1)
	}
	if *repo == "" {
		fmt.Fprintln(os.Stderr, "Error: --repo is required")
		fs.Usage()
		os.Exit(1)
	}

	cfg := derby.CoastConfig{
		Image:           *image,
		Loadout:         *loadout,
		Course:          *course,
		Repo:            *repo,
		EnvFile:         *envFile,
		SkipPermissions: *skipPermissions,
	}

	if err := derby.Coast(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func cmdRun(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: derby run <config.yaml>")
		os.Exit(1)
	}

	configPath := args[0]

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
