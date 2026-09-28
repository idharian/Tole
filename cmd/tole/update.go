package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"tole/pkg/ui"
	"tole/pkg/version"
)

const (
	githubReleaseAPI  = "https://api.github.com/repos/idharian/Tole/releases/latest"
	releaseAssetName  = "tole-windows-amd64.exe"
	checksumAssetName = "SHA256SUMS"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Tole from the latest GitHub release",
	Long:  "Checks GitHub Releases, downloads the matching Windows binary, verifies its SHA256 checksum, and replaces the installed binary.",
	Run:   runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) {
	ui.ShowBanner()
	fmt.Println(ui.SectionHeader("update", false))
	fmt.Println(ui.MutedStyle.Render(fmt.Sprintf("  current: v%s", version.CurrentVersion)))
	fmt.Println()

	installDir := filepath.Join(os.Getenv("LOCALAPPDATA"), "Tole")
	targetExe := filepath.Join(installDir, "tole.exe")
	if err := os.MkdirAll(installDir, 0o755); err != nil {
		fmt.Printf("  %s %v\n", ui.WarningStyle.Render("update failed:"), err)
		return
	}

	sp := ui.StartSpinner("Checking GitHub Releases...")
	release, err := fetchLatestRelease()
	if err != nil {
		sp.StopError("GitHub check failed")
		fmt.Printf("      %s\n\n", ui.WarningStyle.Render(err.Error()))
		return
	}
	sp.StopSuccess("GitHub release found: " + release.TagName)

	latest := strings.TrimPrefix(strings.TrimSpace(release.TagName), "v")
	if compareVersions(latest, version.CurrentVersion) <= 0 {
		fmt.Println(ui.MutedStyle.Render("  Tole is already up to date."))
		fmt.Println()
		return
	}

	assetURL, checksumURL := releaseAssets(release)
	if assetURL == "" {
		fmt.Printf("  %s release has no %s asset\n\n", ui.WarningStyle.Render("update failed:"), release.TagName)
		return
	}

	tempExe := filepath.Join(installDir, "tole-update.exe")
	spDownload := ui.StartSpinner("Downloading " + releaseAssetName + "...")
	if err := downloadFile(assetURL, tempExe); err != nil {
		spDownload.StopError("Download failed")
		_ = os.Remove(tempExe)
		fmt.Printf("      %s\n\n", ui.WarningStyle.Render(err.Error()))
		return
	}
	spDownload.StopSuccess("Download complete")

	if checksumURL == "" {
		_ = os.Remove(tempExe)
		fmt.Printf("  %s checksum asset is missing\n\n", ui.WarningStyle.Render("update failed:"))
		return
	}
	spVerify := ui.StartSpinner("Verifying SHA256...")
	checksumText, err := fetchText(checksumURL)
	if err != nil {
		spVerify.StopError("Checksum download failed")
		_ = os.Remove(tempExe)
		fmt.Printf("      %s\n\n", ui.WarningStyle.Render(err.Error()))
		return
	}
	expected, err := checksumForAsset(checksumText, releaseAssetName)
	if err != nil {
		spVerify.StopError("Checksum entry missing")
		_ = os.Remove(tempExe)
		fmt.Printf("      %s\n\n", ui.WarningStyle.Render(err.Error()))
		return
	}
	actual, err := fileSHA256(tempExe)
	if err != nil || !strings.EqualFold(expected, actual) {
		spVerify.StopError("Checksum mismatch")
		_ = os.Remove(tempExe)
		fmt.Printf("      expected %s, got %s\n\n", expected, actual)
		return
	}
	spVerify.StopSuccess("SHA256 verified")

	if err := scheduleReplacement(tempExe, targetExe); err != nil {
		_ = os.Remove(tempExe)
		fmt.Printf("  %s %v\n\n", ui.WarningStyle.Render("update failed:"), err)
		return
	}
	fmt.Println(ui.Summary("update scheduled", [][2]string{
		{"from", "GitHub Releases"},
		{"version", "v" + latest},
		{"target", targetExe},
	}, false))
	fmt.Println()
}

func fetchLatestRelease() (githubRelease, error) {
	var release githubRelease
	body, err := fetchText(githubReleaseAPI)
	if err != nil {
		return release, err
	}
	if err := json.Unmarshal([]byte(body), &release); err != nil {
		return release, err
	}
	if release.TagName == "" {
		return release, fmt.Errorf("GitHub returned no release tag")
	}
	return release, nil
}

func releaseAssets(release githubRelease) (string, string) {
	var binaryURL, checksumURL string
	for _, asset := range release.Assets {
		switch asset.Name {
		case releaseAssetName:
			binaryURL = asset.BrowserDownloadURL
		case checksumAssetName, checksumAssetName + ".txt":
			checksumURL = asset.BrowserDownloadURL
		}
	}
	return binaryURL, checksumURL
}

func fetchText(url string) (string, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "Tole-Updater/"+version.CurrentVersion)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %s", resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	return string(body), err
}

func downloadFile(url, path string) error {
	client := &http.Client{Timeout: 5 * time.Minute}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Tole-Updater/"+version.CurrentVersion)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %s", resp.Status)
	}
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, io.LimitReader(resp.Body, 250<<20))
	return err
}

func checksumForAsset(text, asset string) (string, error) {
	for _, line := range strings.Split(text, "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 && strings.EqualFold(filepath.Base(fields[len(fields)-1]), asset) {
			checksum := strings.TrimSpace(fields[0])
			if len(checksum) == sha256.Size*2 {
				return checksum, nil
			}
		}
	}
	return "", fmt.Errorf("no SHA256 entry for %s", asset)
}

func fileSHA256(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func scheduleReplacement(source, target string) error {
	script := `param([string]$Source, [string]$Target, [string]$Script)
Start-Sleep -Seconds 1
$old = "$Target.old"
Remove-Item -LiteralPath $old -Force -ErrorAction SilentlyContinue
Move-Item -LiteralPath $Target -Destination $old -Force -ErrorAction SilentlyContinue
Move-Item -LiteralPath $Source -Destination $Target -Force
$null = & $Target version 2>$null
if ($LASTEXITCODE -eq 0) {
    Remove-Item -LiteralPath $old -Force -ErrorAction SilentlyContinue
} else {
    Copy-Item -LiteralPath $old -Destination $Target -Force -ErrorAction SilentlyContinue
}
Start-Process -FilePath $Target
Remove-Item -LiteralPath $Script -Force -ErrorAction SilentlyContinue`

	scriptPath := filepath.Join(os.TempDir(), "tole-updater.ps1")
	if err := os.WriteFile(scriptPath, []byte(script), 0o600); err != nil {
		return err
	}
	process := exec.Command("powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", scriptPath, source, target, scriptPath)
	process.Stdout = nil
	process.Stderr = nil
	return process.Start()
}

func compareVersions(a, b string) int {
	parse := func(v string) []int {
		v = strings.TrimPrefix(strings.TrimSpace(v), "v")
		parts := strings.Split(v, ".")
		out := make([]int, 3)
		for i := 0; i < len(parts) && i < 3; i++ {
			out[i], _ = strconv.Atoi(parts[i])
		}
		return out
	}
	x, y := parse(a), parse(b)
	for i := range x {
		if x[i] < y[i] {
			return -1
		}
		if x[i] > y[i] {
			return 1
		}
	}
	return 0
}
