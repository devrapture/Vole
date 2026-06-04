package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/devrapture/vole/internal/cleaner"
	"github.com/devrapture/vole/internal/scanner"
	"github.com/spf13/cobra"
)

var (
	flagProject  string
	flagAssetDir string
	flagIgnore   []string
	flagVerbose  bool
	flagNoPrompt bool
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
	opts := scannerOpts()

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
		fmt.Println("  No files deleted.  Run  vole clean --yes  to skip this prompt.")
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

func scannerOpts() *scanner.Options {
	return &scanner.Options{
		ProjectPath: flagProject,
		AssetsDir:   flagAssetDir,
		IgnoreDirs:  flagIgnore,
		Verbose:     flagVerbose,
	}
}

func printReport(result *scanner.ScanResult) {
	sep := strings.Repeat("─", 55)
	fmt.Println(sep)
	fmt.Printf("  %-28s %s\n", "Project", result.ProjectPath)
	fmt.Printf("  %-28s %s\n", "Assets directory", result.AssetsDir)
	fmt.Println(sep)
	fmt.Printf("  %-28s %d\n", "Total assets", result.TotalAssets)
	fmt.Printf("  %-28s %d\n", "Used", result.UsedAssets)
	fmt.Printf("  %-28s %d\n", "Unused", result.UnusedAssets)
	fmt.Println(sep)

	if result.UnusedAssets == 0 {
		fmt.Println()
		fmt.Println("  ✓  No unused assets — your project is clean!")
		fmt.Println()
		return
	}

	fmt.Println()
	fmt.Println("  Unused assets:")
	fmt.Println()
	for _, a := range result.UnusedList() {
		fmt.Printf("    ✗  %s\n", a.RelPath)
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
		fmt.Printf("  ⚠  Finished with %d error(s) — check stderr.\n", len(cr.Errors))
	} else {
		fmt.Printf("  ✓  Deleted %d file(s).\n", len(cr.Deleted))
	}
	fmt.Println()
}
