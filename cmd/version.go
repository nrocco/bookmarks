package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version and build information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Bookmarks:\n")
		fmt.Printf("  Version:    %s\n", version)
		fmt.Printf("  Commit:     %s\n", commit)
		fmt.Printf("  Build Date: %s\n", buildDate)
		fmt.Printf("  Platform:   %s (%s)\n", runtime.GOOS, runtime.GOARCH)
		fmt.Printf("  Build Info: %s (%s)\n", runtime.Version(), runtime.Compiler)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
