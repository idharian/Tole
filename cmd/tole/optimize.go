package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"tole/pkg/optimizer"
	"tole/pkg/ui"
	"tole/pkg/winapi"
)

var (
	optimizeDryRun bool
)

var optimizeCmd = &cobra.Command{
	Use:   "optimize",
	Short: "System optimization and maintenance (DNS flush, RAM trim, Icon cache)",
	Long: `Optimize Windows system services, trim memory working sets,
flush DNS cache, and rebuild icon/thumbnail databases.
Run with Administrator privileges to also trigger DISM Component Store cleanup and SSD TRIM.`,
	Run: runOptimize,
}

func init() {
	optimizeCmd.Flags().BoolVarP(&optimizeDryRun, "dry-run", "n", false, "Preview optimization actions without executing")
	rootCmd.AddCommand(optimizeCmd)
}

func runOptimize(cmd *cobra.Command, args []string) {
	ui.ShowBanner()

	fmt.Println(ui.SectionHeader("optimize", optimizeDryRun))
	fmt.Println()

	isAdmin := winapi.IsAdmin()
	tasks := optimizer.GetTasks()

	var successCount int
	var skippedCount int

	for idx, t := range tasks {
		tag := ui.StepTag(idx+1, len(tasks))

		if t.AdminOnly && !isAdmin {
			fmt.Printf("  %s %s %s\n      %s\n\n",
				tag,
				ui.BadgeMuted.Render("admin"),
				ui.MutedStyle.Render(t.Name),
				ui.MutedStyle.Render("skipped — rerun as administrator"),
			)
			skippedCount++
			continue
		}

		sp := ui.StartSpinner(fmt.Sprintf("%s %s...", tag, t.Name))
		result, err := t.Action(optimizeDryRun)
		if err != nil {
			sp.StopError(fmt.Sprintf("%s %s", tag, t.Name))
			fmt.Printf("      %s\n\n", ui.WarningStyle.Render(fmt.Sprintf("%v", err)))
		} else {
			sp.StopSuccess(fmt.Sprintf("%s %s", tag, t.Name))
			fmt.Printf("      %s\n\n", ui.MutedStyle.Render(result))
			successCount++
		}
	}

	fmt.Println(ui.Summary("optimization summary", [][2]string{
		{"executed", fmt.Sprintf("%d", successCount)},
		{"skipped", fmt.Sprintf("%d", skippedCount)},
	}, false))
	if !isAdmin {
		fmt.Println(ui.MutedStyle.Render("  tip: rerun as administrator to unlock DISM + SSD TRIM"))
	}
	fmt.Println()
}
