package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"tole/pkg/purge"
	"tole/pkg/ui"
)

var (
	purgePath   string
	purgeDryRun bool
	purgeYes    bool
	purgeDepth  int
)

var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Scan and purge heavy project build artifacts (node_modules, target, .venv, etc.)",
	Long: `Recursively search directories for dev build artifacts such as
node_modules, Rust target folders, .NET bin/obj, .gradle caches, and Python .venv.
Allows you to reclaim tens of gigabytes from old or unused projects safely.`,
	Run: runPurge,
}

func init() {
	purgeCmd.Flags().StringVarP(&purgePath, "path", "p", ".", "Directory path to scan")
	purgeCmd.Flags().BoolVarP(&purgeDryRun, "dry-run", "n", false, "Preview artifacts without deleting")
	purgeCmd.Flags().BoolVarP(&purgeYes, "yes", "y", false, "Confirm deletion without interactive prompt")
	purgeCmd.Flags().IntVarP(&purgeDepth, "depth", "d", 6, "Maximum directory scan depth")
}

func runPurge(cmd *cobra.Command, args []string) {
	ui.ShowBanner()

	fmt.Println(ui.SectionHeader("purge", purgeDryRun))
	fmt.Println()

	fmt.Printf("  %s %s\n\n",
		ui.MutedStyle.Render("target:"),
		ui.HighlightStyle.Render(purgePath),
	)

	sp := ui.StartSpinner("Scanning projects...")
	artifacts, err := purge.ScanProjects(purgePath, purgeDepth)
	if err != nil {
		sp.StopError(fmt.Sprintf("Scan failed: %v", err))
		return
	}
	sp.StopSuccess("Scan complete")
	fmt.Println()

	if len(artifacts) == 0 {
		fmt.Println(ui.MutedStyle.Render("  No build artifacts found."))
		fmt.Println()
		return
	}

	var totalSize int64
	var totalFiles int64

	fmt.Println(ui.MutedStyle.Render("  artifacts:"))
	for idx, a := range artifacts {
		sizeStr := ui.FormatBytes(a.Size)
		fmt.Printf("  %s %-28s %10s  (%d files)\n       %s\n",
			ui.StepTag(idx+1, len(artifacts)),
			a.Name,
			ui.HighlightStyle.Render(sizeStr),
			a.ItemCount,
			ui.MutedStyle.Render(a.Path),
		)
		totalSize += a.Size
		totalFiles += a.ItemCount
	}
	fmt.Println()

	fmt.Println(ui.Summary("scan result", [][2]string{
		{"folders", fmt.Sprintf("%d", len(artifacts))},
		{"files", fmt.Sprintf("%d", totalFiles)},
		{"reclaimable", ui.FormatBytes(totalSize)},
	}, false))
	fmt.Println()

	if purgeDryRun {
		fmt.Println(ui.MutedStyle.Render("  dry-run: nothing deleted."))
		fmt.Println()
		return
	}

	if !purgeYes {
		var confirm bool
		err := huh.NewConfirm().
			Title("Yakin ingin menghapus seluruh build artifact di atas?").
			Description("Folder seperti node_modules dan target bisa di-generate ulang nanti.").
			Affirmative("Ya, bersihkan").
			Negative("Batal").
			Value(&confirm).
			Run()

		if err != nil || !confirm {
			fmt.Println(ui.MutedStyle.Render("  Dibatalkan oleh pengguna."))
			fmt.Println()
			return
		}
	}

	fmt.Println()
	spDel := ui.StartSpinner("Deleting artifacts...")
	reclaimed, errs := purge.PurgeAll(artifacts, false)
	spDel.StopSuccess("Delete complete")

	fmt.Printf("\n  %s %s",
		ui.HighlightStyle.Render("reclaimed"),
		ui.HighlightStyle.Render(ui.FormatBytes(reclaimed)),
	)
	if errs > 0 {
		fmt.Printf("  %s", ui.MutedStyle.Render(fmt.Sprintf("(%d locked, skipped)", errs)))
	}
	fmt.Println()
	_ = os.Stdout.Sync()
}
