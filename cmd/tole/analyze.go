package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"tole/pkg/analyzer"
	"tole/pkg/ui"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze [path]",
	Short: "Interactive TUI disk space explorer",
	Long: `Analyze disk space usage with an interactive terminal tree visualizer.
Provides live drill-down navigation, proportional visual bar charts,
and the ability to move heavy files and folders to the Windows Recycle Bin safely.`,
	Run: runAnalyze,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}

func runAnalyze(cmd *cobra.Command, args []string) {
	targetPath := "."
	if len(args) > 0 {
		targetPath = args[0]
	}

	info, err := os.Stat(targetPath)
	if err != nil {
		fmt.Printf("%s Path does not exist: %s\n", ui.WarningStyle.Render("✗"), targetPath)
		return
	}
	if !info.IsDir() {
		fmt.Printf("%s Path is not a directory: %s\n", ui.WarningStyle.Render("✗"), targetPath)
		return
	}

	m := analyzer.NewModel(targetPath)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running analyzer: %v\n", err)
	}
}
