package main

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"tole/pkg/installer"
	"tole/pkg/ui"
	"tole/pkg/winapi"
)

var (
	installerDryRun bool
)

var installerCmd = &cobra.Command{
	Use:   "installer",
	Short: "Find and clean leftover installer files (.msi, .exe, .iso)",
	Long:  "Scan your Downloads folder and desktop for forgotten setup packages, .msi files, and disk images.",
	Run: runInstaller,
}

func init() {
	installerCmd.Flags().BoolVarP(&installerDryRun, "dry-run", "n", false, "Preview found installer packages without deleting")
	rootCmd.AddCommand(installerCmd)
}

func runInstaller(cmd *cobra.Command, args []string) {
	ui.ShowBanner()

	fmt.Println(ui.SectionHeader("installer", installerDryRun))
	fmt.Println()

	sp := ui.StartSpinner("Scanning Downloads + Desktop...")
	files, totalSize, err := installer.ScanInstallers()
	if err != nil {
		sp.StopError(fmt.Sprintf("Scan failed: %v", err))
		return
	}
	sp.StopSuccess("Scan complete")

	if len(files) == 0 {
		fmt.Println(ui.MutedStyle.Render("  No stale installers found."))
		fmt.Println()
		return
	}

	fmt.Println()
	for idx, f := range files {
		sizeStr := ui.FormatBytes(f.Size)
		fmt.Printf("  %s %-34s %s\n       %s\n",
			ui.StepTag(idx+1, len(files)),
			f.Name,
			ui.HighlightStyle.Render(sizeStr),
			ui.MutedStyle.Render(f.Path),
		)
	}
	fmt.Println()

	fmt.Println(ui.Summary("scan result", [][2]string{
		{"packages", fmt.Sprintf("%d", len(files))},
		{"size", ui.FormatBytes(totalSize)},
	}, false))
	fmt.Println()

	if installerDryRun {
		fmt.Println(ui.MutedStyle.Render("  dry-run: nothing moved."))
		fmt.Println()
		return
	}

	var confirm bool
	huh.NewConfirm().
		Title(fmt.Sprintf("Pindahkan seluruh %d file installer (%s) ke Recycle Bin?", len(files), ui.FormatBytes(totalSize))).
		Value(&confirm).
		Run()

	if !confirm {
		fmt.Println(ui.MutedStyle.Render("  Dibatalkan oleh pengguna."))
		fmt.Println()
		return
	}

	var movedCount int
	for _, f := range files {
		if err := winapi.MoveToRecycleBin(f.Path); err == nil {
			movedCount++
		}
	}

	fmt.Printf("\n  %s %d packages moved to Recycle Bin.\n\n",
		ui.HighlightStyle.Render("done"),
		movedCount,
	)
}
