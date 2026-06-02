package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vole",
	Short: "Find and remove unused image assets in a React project",
	Long:  "vole is a developer tool that scans your React/TypeScript project, identifies image files inside a chosen assets directory that are never referenced in your source code, and lets you delete them safely",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().Bool("verbose", false, "Print each file as it is processed")
}
