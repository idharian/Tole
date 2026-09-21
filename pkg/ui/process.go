package ui

import (
	"fmt"
	"strings"
)

// Process display rhythm: one header shape, one summary shape.
// [slug] names the task, dryRun toggles the mode badge. No emojis,
// single emerald accent, muted labels, bright values.

// ModeBadge renders the run-mode pill.
func ModeBadge(dryRun bool) string {
	if dryRun {
		return BadgeWarning.Render("dry-run")
	}
	return BadgeSuccess.Render("live")
}

// SectionHeader prints "tole / <slug>" plus the mode badge.
func SectionHeader(slug string, dryRun bool) string {
	return fmt.Sprintf("  %s  %s",
		HighlightStyle.Render("tole / "+slug),
		ModeBadge(dryRun),
	)
}

// Summary renders a titled key/value box. Values render bright,
// labels muted. Warnings pass warn=true to tint the title amber.
func Summary(title string, rows [][2]string, warn bool) string {
	titleStyle := HighlightStyle
	if warn {
		titleStyle = WarningStyle
	}
	var b strings.Builder
	b.WriteString("  " + titleStyle.Render(title) + "\n")
	for _, r := range rows {
		b.WriteString(fmt.Sprintf("  %-24s : %s\n",
			MutedStyle.Render(r[0]),
			HighlightStyle.Render(r[1]),
		))
	}
	return SummaryBox.Render(strings.TrimSuffix(b.String(), "\n"))
}

// StepTag renders a "[i/n]" counter in muted style.
func StepTag(i, n int) string {
	return MutedStyle.Render(fmt.Sprintf("[%d/%d]", i, n))
}
