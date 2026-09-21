package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"tole/pkg/ui"
)

var rootCmd = &cobra.Command{
	Use:   "tole",
	Short: "Tole - Windows System Maintenance & Optimization Toolkit",
	Long: `Tole is a fast, native Windows CLI toolkit designed to clean,
uninstall, optimize, and analyze your PC with the same philosophy and elegance as Mole for Mac.`,
	Run: runInteractiveMenu,
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(purgeCmd)
}

func runInteractiveMenu(cmd *cobra.Command, args []string) {
	for {
		// Clean screen before presenting menu
		fmt.Print("\033[H\033[2J")
		ui.ShowBanner()

		var choice string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select action").
					Options(
						huh.NewOption("Maintain / Deep Clean — cache, temp, browser", "clean"),
						huh.NewOption("Maintain / System Optimize — DNS, RAM, icon cache", "optimize"),
						huh.NewOption("Maintain / Smart Uninstall — apps + leftovers", "uninstall"),
						huh.NewOption("Reclaim / Project Purge — node_modules, target, .venv", "purge"),
						huh.NewOption("Reclaim / Disk Analyzer — visual explorer", "analyze"),
						huh.NewOption("Reclaim / Installer Cleaner — .msi, .exe, .iso", "installer"),
						huh.NewOption("System / Live Monitor — CPU, RAM, disk", "status"),
						huh.NewOption("System / Update Tole", "update"),
						huh.NewOption("System / Exit", "exit"),
					).
					Value(&choice),
			),
		)

		err := form.Run()
		if err != nil || choice == "exit" || choice == "" {
			fmt.Println(ui.MutedStyle.Render("  Sampai jumpa! Terima kasih telah menggunakan Tole."))
			fmt.Println()
			return
		}

		// Execute chosen action
		switch choice {
		case "clean":
			fmt.Print("\033[H\033[2J")
			cleanDryRun = true
			runClean(cleanCmd, nil)
			var proceed bool
			huh.NewConfirm().
				Title("Hasil di atas adalah PREVIEW. Jalankan penghapusan permanen?").
				Description("File dihapus PERMANEN, tidak masuk Recycle Bin.").
				Affirmative("Ya, hapus permanen").
				Negative("Batal").
				Value(&proceed).
				Run()
			if proceed {
				cleanDryRun = false
				fmt.Print("\033[H\033[2J")
				runClean(cleanCmd, nil)
			}
			ui.WaitEnter("Tekan [Enter] untuk kembali ke menu utama...")

		case "purge":
			fmt.Print("\033[H\033[2J")
			prevDry, prevYes := purgeDryRun, purgeYes
			purgeDryRun = true
			runPurge(purgeCmd, nil)
			var proceed bool
			huh.NewConfirm().
				Title("Hasil di atas adalah PREVIEW. Lanjut hapus artifact?").
				Affirmative("Ya, hapus").
				Negative("Batal").
				Value(&proceed).
				Run()
			if proceed {
				purgeDryRun, purgeYes = false, true
				fmt.Print("\033[H\033[2J")
				runPurge(purgeCmd, nil)
			}
			purgeDryRun, purgeYes = prevDry, prevYes
			ui.WaitEnter("Tekan [Enter] untuk kembali ke menu utama...")

		case "analyze":
			analyzeCmd.Run(analyzeCmd, nil)
			// analyze exits Bubble Tea alt screen automatically

		case "optimize":
			fmt.Print("\033[H\033[2J")
			prevOptDry := optimizeDryRun
			optimizeDryRun = true
			runOptimize(optimizeCmd, nil)
			var proceedOpt bool
			huh.NewConfirm().
				Title("Hasil di atas adalah PREVIEW. Jalankan optimasi sistem?").
				Affirmative("Ya, jalankan").
				Negative("Batal").
				Value(&proceedOpt).
				Run()
			if proceedOpt {
				optimizeDryRun = false
				fmt.Print("\033[H\033[2J")
				runOptimize(optimizeCmd, nil)
			}
			optimizeDryRun = prevOptDry
			ui.WaitEnter("Tekan [Enter] untuk kembali ke menu utama...")

		case "uninstall":
			fmt.Print("\033[H\033[2J")
			runUninstall(uninstallCmd, nil)
			ui.WaitEnter("Tekan [Enter] untuk kembali ke menu utama...")

		case "status":
			statusCmd.Run(statusCmd, nil)
			// status exits Bubble Tea alt screen automatically

		case "installer":
			fmt.Print("\033[H\033[2J")
			runInstaller(installerCmd, nil)
			ui.WaitEnter("Tekan [Enter] untuk kembali ke menu utama...")

		case "update":
			fmt.Print("\033[H\033[2J")
			runUpdate(updateCmd, nil)
			ui.WaitEnter("Tekan [Enter] untuk kembali ke menu utama...")
		}
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
