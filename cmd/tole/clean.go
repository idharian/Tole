package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"tole/pkg/cleaner"
	"tole/pkg/ui"
	"tole/pkg/winapi"
)

var (
	cleanDryRun bool
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Deep cleanup for Windows caches, logs, temp files, and recycle bin",
	Long: `Clean Windows system caches, user temp directories, browser caches,
developer caches (npm, pnpm, yarn, cargo, pip, go), and Recycle Bin.
Use --dry-run to preview what will be cleaned without deleting anything.`,
	Run: runClean,
}

func init() {
	cleanCmd.Flags().BoolVarP(&cleanDryRun, "dry-run", "n", false, "Preview cleanup plan without deleting files")
}

func runClean(cmd *cobra.Command, args []string) {
	ui.ShowBanner()

	fmt.Println(ui.SectionHeader("clean", cleanDryRun))
	if !cleanDryRun {
		fmt.Printf("  %s\n\n", ui.WarningStyle.Render("Files are deleted permanently, not moved to Recycle Bin."))
	} else {
		fmt.Println()
	}

	targets := cleaner.GetDefaultTargets()

	sp := ui.StartSpinner("Scanning caches...")
	cleaner.ScanAll(targets)
	sp.StopSuccess("Scan complete")
	fmt.Println()

	// Group by category
	categories := []cleaner.Category{
		cleaner.CategoryUserEssentials,
		cleaner.CategorySystem,
		cleaner.CategoryBrowsers,
		cleaner.CategoryDeveloper,
	}

	var totalCleanedSize int64
	var totalCleanedCount int64
	var activeCategories int

	for _, cat := range categories {
		var catTargets []*cleaner.Target
		for _, t := range targets {
			if t.Category == cat && (t.Size > 0 || t.Count > 0) {
				catTargets = append(catTargets, t)
			}
		}

		if len(catTargets) == 0 {
			continue
		}

		activeCategories++
		fmt.Printf("  %s\n",
			ui.MutedStyle.Render(string(cat)),
		)

		for i, t := range catTargets {
			cleaner.CleanTarget(t, cleanDryRun)

			sizeStr := ui.FormatBytes(t.CleanedSize)
			branch := "├──"
			if i == len(catTargets)-1 {
				branch = "└──"
			}

			statusIcon := ui.SuccessStyle.Render("✓")
			if t.ErrorMsg != "" {
				statusIcon = ui.WarningStyle.Render("!")
			}

			if cleanDryRun {
				fmt.Printf("    %s %s %-32s  %s  %s\n",
					ui.MutedStyle.Render(branch),
					statusIcon,
					t.Name,
					ui.MutedStyle.Render(fmt.Sprintf("%6d items", t.CleanedCount)),
					ui.HighlightStyle.Render(sizeStr),
				)
			} else {
				detail := fmt.Sprintf("%6d items", t.CleanedCount)
				if t.SkippedCount > 0 {
					detail += fmt.Sprintf(" (%d skipped)", t.SkippedCount)
				}
				fmt.Printf("    %s %s %-32s  %s  %s\n",
					ui.MutedStyle.Render(branch),
					statusIcon,
					t.Name,
					ui.MutedStyle.Render(detail),
					ui.SuccessStyle.Render(sizeStr),
				)
			}

			totalCleanedSize += t.CleanedSize
			totalCleanedCount += t.CleanedCount
		}
		fmt.Println()
	}

	freeAfter, _, _ := winapi.GetDiskSpace("C:\\")

	title := "cleanup complete"
	if cleanDryRun {
		title = "dry-run estimate"
	}
	fmt.Println(ui.Summary(title, [][2]string{
		{"reclaimed", ui.FormatBytes(totalCleanedSize)},
		{"items", fmt.Sprintf("%d", totalCleanedCount)},
		{"categories", fmt.Sprintf("%d", activeCategories)},
		{"C: free", ui.FormatBytes(int64(freeAfter))},
	}, cleanDryRun))
	fmt.Println()
	time.Sleep(200 * time.Millisecond)
}
