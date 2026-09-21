package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"tole/pkg/monitor"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Live system health dashboard (CPU, Memory, Disk)",
	Long:  "Monitor real-time system performance counters and resource utilization with a live terminal dashboard.",
	Run: runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) {
	m := monitor.NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running status monitor: %v\n", err)
	}
}
