package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"tole/pkg/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the Tole version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("tole v%s (built %s) by %s\n",
			version.CurrentVersion,
			version.BuildDate,
			version.Author,
		)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
