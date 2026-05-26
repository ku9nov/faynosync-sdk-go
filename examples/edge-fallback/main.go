package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	faynosync "github.com/ku9nov/faynosync-sdk-go"
)

func main() {
	client := faynosync.NewClient(faynosync.Config{
		BaseURL: "http://localhost:9000",
		EdgeURL: "http://faynosync-cdn-edge.web.garage.localhost:3902",
	})

	resp, err := client.CheckForUpdates(context.Background(), faynosync.CheckOptions{
		Owner:    "admin",
		AppName:  "test",
		Version:  "0.0.0.1",
		Channel:  "nightly",
		Platform: "darwin",
		Arch:     "arm64",
	})
	if err != nil {
		if errors.Is(err, faynosync.ErrRequestFailed) {
			log.Fatalf("update check request failed: %v", err)
		}
		log.Fatalf("invalid update check options: %v", err)
	}

	printUpdateResponse(resp)
}

func printUpdateResponse(resp *faynosync.UpdateResponse) {
	const pad = "  "

	fmt.Println("-- Update check response --")
	fmt.Printf("%supdate_available:         %t\n", pad, resp.UpdateAvailable)
	fmt.Printf("%scritical:                 %t\n", pad, resp.Critical)
	fmt.Printf("%sis_intermediate_required: %t\n", pad, resp.IsIntermediateRequired)
	fmt.Printf("%spossible_rollback:        %t\n", pad, resp.PossibleRollback)
	fmt.Printf("%ssource:                   %s\n", pad, formatUpdateSource(resp.Source))

	if resp.UpdateURL != "" {
		fmt.Printf("%supdate_url:               %s\n", pad, resp.UpdateURL)
	} else {
		fmt.Printf("%supdate_url:               (empty)\n", pad)
	}

	if resp.Changelog != "" {
		fmt.Printf("%schangelog:\n", pad)
		for _, line := range strings.Split(strings.TrimRight(resp.Changelog, "\n"), "\n") {
			fmt.Printf("%s  %s\n", pad, line)
		}
	} else {
		fmt.Printf("%schangelog:               (empty)\n", pad)
	}

	if len(resp.PackageURLs) == 0 {
		fmt.Printf("%spackage_urls:           (none)\n", pad)
		return
	}

	fmt.Printf("%spackage_urls:\n", pad)
	for _, pkg := range resp.PackageURLs {
		fmt.Printf("%s  %-12s %s\n", pad, pkg.Package+":", pkg.URL)
	}
}

func formatUpdateSource(source faynosync.UpdateSource) string {
	switch source {
	case faynosync.SourceEdge:
		return "edge"
	case faynosync.SourceAPI:
		return "api"
	default:
		return "unknown"
	}
}
