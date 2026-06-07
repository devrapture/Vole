package cleaner

import (
	"fmt"
	"os"

	"github.com/devrapture/vole/internal/scanner"
	"github.com/fatih/color"
)

var (
	dryStyle = color.New(color.FgCyan)
	delStyle = color.New(color.FgGreen)
	errStyle = color.New(color.FgRed, color.Bold)
)

type Result struct {
	Deleted         []string
	Skipped         []string
	Errors          []string
	SpaceSavedBytes int64
}

type Options struct {
	DryRun  bool
	Verbose bool
}

func Clean(scanResult *scanner.ScanResult, opts Options) (*Result, error) {
	result := &Result{}

	for _, asset := range scanResult.UnusedList() {
		if opts.DryRun {
			if opts.Verbose {
				fmt.Println(dryStyle.Sprintf("vole [dry-run] would delete: %s", asset.RelPath))
			}
			result.Skipped = append(result.Skipped, asset.AbsPath)
			continue
		}

		if err := os.Remove(asset.AbsPath); err != nil {
			msg := fmt.Sprintf("failed to delete %s:%v", asset.RelPath, err)
			fmt.Fprintln(os.Stderr, errStyle.Sprint("vole error: ")+msg)
			result.Errors = append(result.Errors, msg)
			continue
		}

		if opts.Verbose {
			fmt.Println(delStyle.Sprintf("vole deleted: %s", asset.RelPath))
		}

		result.SpaceSavedBytes += asset.SizeBytes
		result.Deleted = append(result.Deleted, asset.AbsPath)
	}

	return result, nil
}
