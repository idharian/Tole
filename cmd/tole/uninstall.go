package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"tole/pkg/ui"
	"tole/pkg/uninstaller"
	"tole/pkg/winapi"
)

var (
	uninstallList   bool
	uninstallSearch string
	uninstallDryRun bool
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Smart application uninstaller with leftover cleanup",
	Long: `List and uninstall Windows applications from the Registry.
After running the vendor uninstaller, Tole performs a smart sweep
of AppData and ProgramData to eliminate leftover remnants and folders.`,
	Run: runUninstall,
}

func init() {
	uninstallCmd.Flags().BoolVarP(&uninstallList, "list", "l", false, "List all installed applications")
	uninstallCmd.Flags().StringVarP(&uninstallSearch, "search", "s", "", "Search for an installed application by name")
	uninstallCmd.Flags().BoolVarP(&uninstallDryRun, "dry-run", "n", false, "Preview without running uninstaller or deleting files")
	rootCmd.AddCommand(uninstallCmd)
}

func runUninstall(cmd *cobra.Command, args []string) {
	ui.ShowBanner()

	fmt.Println(ui.SectionHeader("uninstall", uninstallDryRun))
	fmt.Println()

	sp := ui.StartSpinner("Reading installed apps...")
	apps, err := uninstaller.GetInstalledApps()
	if err != nil {
		sp.StopError(fmt.Sprintf("Registry read failed: %v", err))
		return
	}
	sp.StopSuccess(fmt.Sprintf("Found %d apps", len(apps)))
	fmt.Println()

	// Filter by search query if provided
	if uninstallSearch != "" {
		var filtered []*uninstaller.InstalledApp
		query := strings.ToLower(uninstallSearch)
		for _, a := range apps {
			if strings.Contains(strings.ToLower(a.DisplayName), query) ||
				strings.Contains(strings.ToLower(a.Publisher), query) {
				filtered = append(filtered, a)
			}
		}
		apps = filtered
	}

	if len(apps) == 0 {
		fmt.Println(ui.MutedStyle.Render("  No matching apps."))
		fmt.Println()
		return
	}

	if uninstallList {
		fmt.Printf("  %-45s %-15s %s\n", "name", "version", "publisher")
		fmt.Println("  " + ui.MutedStyle.Render(strings.Repeat("-", 80)))
		for _, a := range apps {
			ver := a.DisplayVersion
			if ver == "" {
				ver = "-"
			}
			pub := a.Publisher
			if pub == "" {
				pub = "-"
			}
			name := a.DisplayName
			if len(name) > 43 {
				name = name[:40] + "..."
			}
			fmt.Printf("  %-45s %-15s %s\n",
				ui.CategoryStyle.Render(name),
				ui.MutedStyle.Render(ver),
				ui.MutedStyle.Render(pub),
			)
		}
		fmt.Println()
		return
	}

	// Interactive selection if not just listing.
	// Apps sorted by name; filter via --search when list exceeds display limit.
	displayApps := apps
	if len(apps) > 60 {
		fmt.Println(ui.MutedStyle.Render(fmt.Sprintf("  %d apps found, showing first 60 (use --search to filter).", len(apps))))
		fmt.Println()
		displayApps = apps[:60]
	}
	options := make([]huh.Option[string], 0, len(displayApps))
	for idx, a := range displayApps {
		label := a.DisplayName
		if a.DisplayVersion != "" {
			label += " (" + a.DisplayVersion + ")"
		}
		options = append(options, huh.NewOption(label, fmt.Sprintf("%d", idx)))
		if idx >= 60 { // Limit select display to top 60 for terminal performance
			break
		}
	}

	var selectedIndexStr string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select app to uninstall:").
				Options(options...).
				Value(&selectedIndexStr),
		),
	)

	err = form.Run()
	if err != nil || selectedIndexStr == "" {
		fmt.Println(ui.MutedStyle.Render("  Cancelled."))
		fmt.Println()
		return
	}

	var selectedIdx int
	selectedIdx, err = strconv.Atoi(selectedIndexStr)
	if err != nil || selectedIdx < 0 || selectedIdx >= len(apps) {
		fmt.Println(ui.WarningStyle.Render("  Invalid selection."))
		fmt.Println()
		return
	}
	targetApp := apps[selectedIdx]

	fmt.Println()
	fmt.Println("  " + ui.MutedStyle.Render("selected:") + " " + ui.HighlightStyle.Render(targetApp.DisplayName))
	if targetApp.Publisher != "" {
		fmt.Println("  " + ui.MutedStyle.Render("publisher:") + " " + targetApp.Publisher)
	}
	if targetApp.UninstallString != "" {
		fmt.Println("  " + ui.MutedStyle.Render("command:  ") + " " + ui.MutedStyle.Render(targetApp.UninstallString))
	}
	fmt.Println()

	if uninstallDryRun {
		fmt.Println(ui.MutedStyle.Render("  dry-run: leftover scan only, uninstaller not executed."))
	} else {
		var confirm bool
		huh.NewConfirm().
			Title(fmt.Sprintf("Jalankan uninstaller resmi untuk '%s'?", targetApp.DisplayName)).
			Affirmative("Ya, jalankan").
			Negative("Batal").
			Value(&confirm).
			Run()

		if !confirm {
			fmt.Println(ui.MutedStyle.Render("  Dibatalkan."))
			fmt.Println()
			return
		}

		if targetApp.UninstallString != "" {
			fmt.Println(ui.HighlightStyle.Render("  Menjalankan uninstaller..."))
			cmdExec := exec.Command("cmd", "/C", targetApp.UninstallString)
			cmdExec.Stdout = os.Stdout
			cmdExec.Stderr = os.Stderr
			_ = cmdExec.Run()
		}
	}

	// Leftovers Sweep
	fmt.Println()
	spLeft := ui.StartSpinner("Scanning leftovers...")
	leftovers := uninstaller.FindLeftovers(targetApp.DisplayName)
	spLeft.StopSuccess("Leftover scan complete")
	fmt.Println()

	if len(leftovers) == 0 {
		fmt.Println(ui.MutedStyle.Render("  No leftover folders."))
		fmt.Println()
		return
	}

	var totalLeftoverBytes int64
	for _, l := range leftovers {
		totalLeftoverBytes += l.Size
		fmt.Printf("  %s %-28s %10s  (%d files)\n     %s\n",
			ui.StepTag(1, 1),
			l.Name,
			ui.HighlightStyle.Render(ui.FormatBytes(l.Size)),
			l.FileCount,
			ui.MutedStyle.Render(l.Path),
		)
	}
	fmt.Println()

	if uninstallDryRun {
		fmt.Printf("  %s\n\n", ui.MutedStyle.Render(fmt.Sprintf("dry-run: %s across %d locations.", ui.FormatBytes(totalLeftoverBytes), len(leftovers))))
		return
	}

	var deleteLeftovers bool
	huh.NewConfirm().
		Title(fmt.Sprintf("Pindahkan %d folder sisa (%s) di bawah ini ke Windows Recycle Bin?", len(leftovers), ui.FormatBytes(totalLeftoverBytes))).
		Description("Periksa daftar full path di atas sebelum menyetujui.").
		Affirmative("Ya, pindahkan").
		Negative("Jangan hapus").
		Value(&deleteLeftovers).
		Run()

	if deleteLeftovers {
		for _, l := range leftovers {
			_ = winapi.MoveToRecycleBin(l.Path)
		}
		fmt.Println(ui.HighlightStyle.Render("  done — leftovers moved to Recycle Bin."))
	}
	fmt.Println()
}
