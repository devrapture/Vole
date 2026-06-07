package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/devrapture/vole/internal/cleaner"
	"github.com/devrapture/vole/internal/config"
	"github.com/devrapture/vole/internal/scanner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	flagProject  string
	flagAssetDir string
	flagIgnore   []string
	flagVerbose  bool
	flagNoPrompt bool

	sepStyle    = color.New(color.Bold, color.FgCyan)
	labelStyle  = color.New(color.Faint)
	usedStyle   = color.New(color.FgGreen, color.Bold)
	unusedStyle = color.New(color.FgRed, color.Bold)
	okStyle     = color.New(color.FgGreen)
	warnStyle   = color.New(color.FgYellow)
	errStyle    = color.New(color.FgRed, color.Bold)
	headerStyle = color.New(color.Bold)
)

var rootCmd = &cobra.Command{
	Use:   "vole",
	Short: "Find and remove unused image assets in a React project",
	Long:  "vole is a developer tool that scans your React/TypeScript project, identifies image files inside a chosen assets directory that are never referenced in your source code, and lets you delete them safely",
	RunE:  runRoot,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagProject, "project", ".", "React project root")
	rootCmd.PersistentFlags().StringVar(&flagAssetDir, "assets", "src/assets", "Image assets sub-directory")
	rootCmd.PersistentFlags().StringSliceVar(&flagIgnore, "ignore", nil, "Extra directory names to ignore")
	rootCmd.PersistentFlags().BoolVar(&flagVerbose, "verbose", false, "Log every file vole reads")
	rootCmd.PersistentFlags().BoolVar(&flagNoPrompt, "yes", false, "Delete without prompt")
}

func runRoot(cmd *cobra.Command, args []string) error {
	opts, err := scannerOpts(cmd)
	if err != nil {
		return err
	}

	result, err := scanner.NewScanner(*opts).Scan()
	if err != nil {
		return err
	}

	printReport(result)

	if result.UnusedAssets == 0 {
		return nil
	}

	doDelete, err := askDelete(result.UnusedList(), flagNoPrompt)
	if err != nil {
		return err
	}

	if !doDelete {
		fmt.Println()
		fmt.Println("  " + labelStyle.Sprint("No files deleted.") + "  Run  vole --yes  to skip this prompt.")
		fmt.Println()
		return nil
	}

	cleanResult, err := cleaner.Clean(result, cleaner.Options{
		Verbose: flagVerbose,
	})
	if err != nil {
		return err
	}

	printDeleteSummary(cleanResult)

	return nil
}

func scannerOpts(cmd *cobra.Command) (*scanner.Options, error) {
	cfg, err := config.Load(flagProject)
	if err != nil {
		return nil, err
	}
	assetsDirs := cfg.Assets

	if len(assetsDirs) == 0 || cmd.Flags().Changed("assetsDir") {
		assetsDirs = []string{}
	}
	ignoreDirs := cfg.Ignore

	if cmd.Flags().Changed("ignore") {
		ignoreDirs = flagIgnore
	}
	return &scanner.Options{
		ProjectPath: flagProject,
		AssetsDirs:  assetsDirs,
		IgnoreDirs:  ignoreDirs,
		Verbose:     flagVerbose,
	}, nil
}

func printReport(result *scanner.ScanResult) {
	sep := sepStyle.Sprint(strings.Repeat("─", 55))
	fmt.Println(sep)
	fmt.Printf("  %s %s\n", labelStyle.Sprintf("%-28s", "Project"), color.New(color.Bold).Sprint(result.ProjectPath))
	fmt.Printf("  %s %s\n", labelStyle.Sprintf("%-28s", "Assets directories"), strings.Join(result.AssetsDirs, ", "))
	fmt.Println(sep)
	fmt.Printf("  %s %s\n", labelStyle.Sprintf("%-28s", "Total assets"), color.New(color.Bold).Sprintf("%d", result.TotalAssets))
	fmt.Printf("  %s %s\n", labelStyle.Sprintf("%-28s", "Used"), usedStyle.Sprintf("%d", result.UsedAssets))
	fmt.Printf("  %s %s\n", labelStyle.Sprintf("%-28s", "Unused"), unusedStyle.Sprintf("%d", result.UnusedAssets))
	fmt.Println(sep)

	if result.UnusedAssets == 0 {
		fmt.Println()
		fmt.Println("  " + okStyle.Sprint("✓") + "  No unused assets — your project is clean!")
		fmt.Println()
		return
	}

	fmt.Println()
	fmt.Println("  " + headerStyle.Sprint("Unused assets:"))
	fmt.Println()
	for _, a := range result.UnusedList() {
		fmt.Printf("    %s  %s\n", errStyle.Sprint("✗"), a.RelPath)
	}
	fmt.Println()
}

func askDelete(unused []*scanner.ImageAsset, noPrompt bool) (bool, error) {
	if noPrompt {
		return true, nil
	}
	fmt.Printf("   Delete %d unused file(s)? [y/N]", len(unused))
	reader := bufio.NewReader(os.Stdin)
	answer, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("reading answer: %w", err)
	}

	answer = strings.TrimSpace(strings.ToLower(answer))

	return answer == "y" || answer == "yes", nil
}

func printDeleteSummary(cr *cleaner.Result) {
	fmt.Println()
	if len(cr.Errors) > 0 {
		fmt.Printf("  %s  %s\n", warnStyle.Sprint("⚠"), warnStyle.Sprintf("Finished with %d error(s) — check stderr.", len(cr.Errors)))
	} else {
		fmt.Printf("  %s  %s\n", okStyle.Sprint("✓"), okStyle.Sprintf("Deleted %d file(s).", len(cr.Deleted)))
		fmt.Printf("  %s %s\n", okStyle.Sprint("✓"), okStyle.Sprintf("Space saved: %s", formatBytes(cr.SpaceSavedBytes)))
	}
	fmt.Println()
}

func formatBytes(bytes int64) string {
	const unit = 1024

	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0

	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"KB", "MB", "GB", "TB"}
	if exp >= len(units) {
		exp = len(units) - 1
		div = int64(unit)
		for i := 0; i < exp; i++ {
			div *= unit
		}
	}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}
