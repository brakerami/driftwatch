// Command driftwatch is a lightweight daemon that detects configuration drift
// between running containers and their source manifests.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/inspector"
	"github.com/yourorg/driftwatch/internal/manifest"
	"github.com/yourorg/driftwatch/internal/reporter"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		manifestPath string
		formatStr    string
		outputPath   string
		watch        bool
	)

	flag.StringVar(&manifestPath, "manifest", "", "path to the manifest YAML file (required)")
	flag.StringVar(&formatStr, "format", "text", "output format: text or json")
	flag.StringVar(&outputPath, "output", "", "write report to file instead of stdout")
	flag.BoolVar(&watch, "watch", false, "continuously watch for drift (not yet implemented)")
	flag.Parse()

	if manifestPath == "" {
		flag.Usage()
		return fmt.Errorf("--manifest is required")
	}

	// Parse the desired output format.
	fmt_, err := reporter.ParseFormat(formatStr)
	if err != nil {
		return fmt.Errorf("invalid format %q: %w", formatStr, err)
	}

	// Load the manifest describing expected container state.
	loader := manifest.NewLoader()
	specs, err := loader.Load(manifestPath)
	if err != nil {
		return fmt.Errorf("loading manifest %q: %w", manifestPath, err)
	}

	// Connect to the Docker daemon.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	insp, err := inspector.New(ctx)
	if err != nil {
		return fmt.Errorf("creating inspector: %w", err)
	}

	// Run drift detection across all specified containers.
	detector := drift.New(insp)

	var allFindings []drift.Finding
	for _, spec := range specs {
		findings, err := detector.Detect(ctx, spec)
		if err != nil {
			// Log non-fatal per-container errors and continue.
			fmt.Fprintf(os.Stderr, "warn: detecting drift for %q: %v\n", spec.Name, err)
			continue
		}
		allFindings = append(allFindings, findings...)
	}

	// Set up the report writer.
	var rOpts []reporter.Option
	if outputPath != "" {
		f, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("opening output file %q: %w", outputPath, err)
		}
		defer f.Close()
		rOpts = append(rOpts, reporter.WithWriter(f))
	}

	rep := reporter.New(fmt_, rOpts...)
	if err := rep.Write(allFindings); err != nil {
		return fmt.Errorf("writing report: %w", err)
	}

	// Exit with a non-zero status when drift is detected so CI pipelines
	// can treat the result as a failure.
	if len(allFindings) > 0 {
		os.Exit(2)
	}

	return nil
}
