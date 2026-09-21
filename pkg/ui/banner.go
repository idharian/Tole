package ui

import (
	"fmt"
	"runtime"
	"strings"
	"tole/pkg/version"
	"tole/pkg/winapi"
)

// Logo is the TOLE wordmark: heavy block letters, emerald on dark.
const Logo = `
  ████████╗ ██████╗ ██╗     ███████╗
  ╚══██╔══╝██╔═══██╗██║     ██╔════╝
     ██║   ██║   ██║██║     █████╗
     ██║   ██║   ██║██║     ██╔══╝
     ██║   ╚██████╔╝███████╗███████╗
     ╚═╝    ╚═════╝ ╚══════╝╚══════╝`

// ShowBanner prints the TOLE brand header.
// Design read: CLI brand mark for Windows power users, dark-tech premium
// language, Lipgloss system. Dials 6/2/4: bold wordmark block, static,
// compact. Single emerald accent, slate neutrals, no emoji.
func ShowBanner() {
	// Wordmark block — the brand moment.
	fmt.Println(HighlightStyle.Render(Logo))

	// Tagline + version on one line, author muted.
	tagLine := fmt.Sprintf("  %s  %s  %s",
		MutedStyle.Render("windows maintenance toolkit"),
		BadgeMuted.Render("v"+version.CurrentVersion),
		MutedStyle.Render("by agushariyanto"))
	fmt.Println(tagLine)

	// Context line: platform · privilege · disk + mini usage bar.
	priv := "standard"
	if winapi.IsAdmin() {
		priv = "admin"
	}

	disk := ""
	bar := ""
	if free, total, err := winapi.GetDiskSpace("C:\\"); err == nil && total > 0 {
		usedPct := 100.0 * (1 - float64(free)/float64(total))
		disk = fmt.Sprintf("C: %s free", FormatBytes(int64(free)))
		bar = " " + miniBar(usedPct, 10)
	}

	parts := []string{
		fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		priv,
		disk,
	}
	context := MutedStyle.Render(strings.Join(parts, "  ·  ") + bar)
	if winapi.IsAdmin() {
		context = MutedStyle.Render(parts[0]+"  ·  ") +
			HighlightStyle.Render("admin") +
			MutedStyle.Render("  ·  "+disk+bar)
	}
	fmt.Printf("  %s\n", context)
	fmt.Println()
}

// miniBar renders a compact usage bar: filled emerald, empty slate.
func miniBar(pct float64, width int) string {
	if width <= 0 {
		width = 10
	}
	filled := int(pct / 100 * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	return HighlightStyle.Render(strings.Repeat("━", filled)) +
		MutedStyle.Render(strings.Repeat("━", width-filled))
}
