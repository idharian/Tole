package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"tole/pkg/ui"
	"tole/pkg/version"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Tole to the latest version",
	Long: `Update Tole to the latest version.
If installed from source, pulls the latest git commits, rebuilds the binary,
and refreshes the installation in %LOCALAPPDATA%\Tole.`,
	Run: runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// If destination file is currently running, Windows locks it.
	// Trick: In Windows, you can RENAME a running executable!
	// So we rename the old dst to .old, write new dst, then remove .old.
	oldDst := dst + ".old"
	_ = os.Remove(oldDst)
	_ = os.Rename(dst, oldDst)

	out, err := os.Create(dst)
	if err != nil {
		// Restore if failed
		_ = os.Rename(oldDst, dst)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err == nil {
		_ = os.Remove(oldDst)
	}
	return err
}

func runUpdate(cmd *cobra.Command, args []string) {
	ui.ShowBanner()

	fmt.Println(ui.SectionHeader("update", false))
	fmt.Println(ui.MutedStyle.Render(fmt.Sprintf("  current: v%s", version.CurrentVersion)))
	fmt.Println()

	// Check where this executable is running from
	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("  %s Could not determine executable path: %v\n", ui.WarningStyle.Render("x"), err)
		return
	}

	installDir := filepath.Join(os.Getenv("LOCALAPPDATA"), "Tole")
	targetExe := filepath.Join(installDir, "tole.exe")

	fmt.Printf("  %s %s\n", ui.MutedStyle.Render("running:"), ui.MutedStyle.Render(execPath))
	fmt.Printf("  %s %s\n\n", ui.MutedStyle.Render("target: "), ui.MutedStyle.Render(targetExe))

	// Look for source directory (e.g. current directory or git repo)
	cwd, _ := os.Getwd()
	sourceDir := ""

	// Check if cwd has go.mod with module tole
	if isToleSource(cwd) {
		sourceDir = cwd
	} else if parent := filepath.Dir(execPath); isToleSource(parent) {
		sourceDir = parent
	}

	if sourceDir != "" {
		fmt.Println(ui.MutedStyle.Render("  source channel"))
		fmt.Printf("  %s\n\n", ui.MutedStyle.Render(sourceDir))

		// Check if git is available in sourceDir
		gitDir := filepath.Join(sourceDir, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			spGit := ui.StartSpinner("git pull...")
			gitCmd := exec.Command("git", "-C", sourceDir, "pull")
			gitOut, gitErr := gitCmd.CombinedOutput()
			if gitErr != nil {
				spGit.StopWarning(strings.TrimSpace(string(gitOut)))
			} else {
				spGit.StopSuccess(strings.TrimSpace(string(gitOut)))
			}
		}

		// Recompile tole.exe
		spBuild := ui.StartSpinner("Building tole.exe...")
		tempBuiltExe := filepath.Join(sourceDir, "tole_new.exe")
		buildCmd := exec.Command("go", "build", "-o", tempBuiltExe, "./cmd/tole")
		buildCmd.Dir = sourceDir
		buildOut, buildErr := buildCmd.CombinedOutput()
		if buildErr != nil {
			spBuild.StopError("Build failed")
			fmt.Printf("      %s\n\n", ui.WarningStyle.Render(string(buildOut)))
			return
		}
		spBuild.StopSuccess("Build complete")

		// Copy to local install directory
		spInstall := ui.StartSpinner("Installing...")
		if err := copyFile(tempBuiltExe, targetExe); err != nil {
			spInstall.StopError(fmt.Sprintf("Copy failed: %v", err))
			_ = os.Remove(tempBuiltExe)
			return
		}

		// Also update the source directory's tole.exe
		srcMainExe := filepath.Join(sourceDir, "tole.exe")
		_ = copyFile(tempBuiltExe, srcMainExe)
		_ = os.Remove(tempBuiltExe)

		spInstall.StopSuccess(fmt.Sprintf("Installed to %s", targetExe))
		fmt.Println()

		fmt.Println(ui.Summary("update complete", [][2]string{
			{"target", targetExe},
		}, false))
		fmt.Println()
		return
	}

	// Standalone binary update advice
	fmt.Println(ui.MutedStyle.Render(fmt.Sprintf("  tole v%s — latest local build", version.CurrentVersion)))
	fmt.Println(ui.MutedStyle.Render("  standalone install: download the latest release or git pull in the source repo."))
	fmt.Println()
}

func isToleSource(dir string) bool {
	modFile := filepath.Join(dir, "go.mod")
	content, err := os.ReadFile(modFile)
	if err != nil {
		return false
	}
	return strings.Contains(string(content), "module tole")
}
