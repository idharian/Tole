package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"tole/pkg/largefiles"
	"tole/pkg/ui"
	"tole/pkg/winapi"
)

var (
	largeDryRun  bool
	largeMinSize string
	largeTop     int
)

var largeCmd = &cobra.Command{
	Use:   "large [path]",
	Short: "Find the largest files under a directory",
	Long: `Scan a directory tree and list the biggest files, so you can reclaim
space in one look instead of drilling through folders. Selected files
are moved to the Recycle Bin, nothing is deleted permanently.`,
	Run: runLarge,
}

func init() {
	largeCmd.Flags().BoolVarP(&largeDryRun, "dry-run", "n", false, "Preview found files without moving anything")
	largeCmd.Flags().StringVar(&largeMinSize, "min", "100MB", "Minimum file size (e.g. 50MB, 1GB)")
	largeCmd.Flags().IntVar(&largeTop, "top", 25, "Show only the N largest files")
	rootCmd.AddCommand(largeCmd)
}

func runLarge(cmd *cobra.Command, args []string) {
	ui.ShowBanner()

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

	minSize := int64(100 * 1024 * 1024)
	if parsed, perr := parseSize(largeMinSize); perr == nil {
		minSize = parsed
	} else {
		fmt.Printf("%s Invalid --min value %q: %v\n", ui.WarningStyle.Render("✗"), largeMinSize, perr)
		return
	}

	fmt.Println(ui.SectionHeader("large", largeDryRun))
	fmt.Println(ui.MutedStyle.Render(fmt.Sprintf("  path: %s  ·  min: %s  ·  top: %d",
		targetPath, ui.FormatBytes(minSize), largeTop)))
	fmt.Println()

	sp := ui.StartSpinner("Scanning for large files...")
	files, scannedBytes, scanErr := largefiles.FindLargeFiles(targetPath, minSize, largeTop)
	if scanErr != nil {
		sp.StopError(fmt.Sprintf("Scan failed: %v", scanErr))
		return
	}
	sp.StopSuccess("Scan complete")

	if len(files) == 0 {
		fmt.Println(ui.MutedStyle.Render("  No files above the size threshold."))
		fmt.Println()
		return
	}

	fmt.Println()
	for idx, f := range files {
		fmt.Printf("  %s %-44s %s\n       %s\n",
			ui.StepTag(idx+1, len(files)),
			f.Name,
			ui.HighlightStyle.Render(ui.FormatBytes(f.Size)),
			ui.MutedStyle.Render(f.Path),
		)
	}
	fmt.Println()

	fmt.Println(ui.Summary("scan result", [][2]string{
		{"files shown", fmt.Sprintf("%d", len(files))},
		{"size shown", ui.FormatBytes(sumSizes(files))},
		{"all files scanned", ui.FormatBytes(scannedBytes)},
	}, false))
	fmt.Println()

	if largeDryRun {
		fmt.Println(ui.MutedStyle.Render("  dry-run: nothing moved."))
		fmt.Println()
		return
	}

	var confirm bool
	huh.NewConfirm().
		Title(fmt.Sprintf("Pindahkan seluruh %d file (%s) ke Recycle Bin?",
			len(files), ui.FormatBytes(sumSizes(files)))).
		Value(&confirm).
		Run()

	if !confirm {
		fmt.Println(ui.MutedStyle.Render("  Dibatalkan oleh pengguna."))
		fmt.Println()
		return
	}

	var movedCount int
	var movedBytes int64
	for _, f := range files {
		if err := winapi.MoveToRecycleBin(f.Path); err == nil {
			movedCount++
			movedBytes += f.Size
		}
	}

	fmt.Printf("\n  %s %d files moved to Recycle Bin (%s).\n\n",
		ui.HighlightStyle.Render("done"),
		movedCount,
		ui.FormatBytes(movedBytes),
	)
}

func sumSizes(files []*largefiles.LargeFile) int64 {
	var total int64
	for _, f := range files {
		total += f.Size
	}
	return total
}

// parseSize accepts forms like "512KB", "100MB", "2GB", or a plain byte count.
func parseSize(s string) (int64, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	if s == "" {
		return 0, fmt.Errorf("empty size")
	}
	idx := 0
	for idx < len(s) && ((s[idx] >= '0' && s[idx] <= '9') || s[idx] == '.') {
		idx++
	}
	if idx == len(s) {
		// Plain number: bytes
		return strconv.ParseInt(s, 10, 64)
	}
	unit := strings.TrimSpace(s[idx:])
	num, err := strconv.ParseFloat(s[:idx], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number in %q", s)
	}
	var mult float64
	switch unit {
	case "B":
		mult = 1
	case "KB":
		mult = 1024
	case "MB":
		mult = 1024 * 1024
	case "GB":
		mult = 1024 * 1024 * 1024
	case "TB":
		mult = 1024 * 1024 * 1024 * 1024
	default:
		return 0, fmt.Errorf("unknown unit %q (use B, KB, MB, GB, TB)", unit)
	}
	return int64(num * mult), nil
}
