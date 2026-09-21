package ui

import (
	"bufio"
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Locked palette: single emerald accent + slate neutrals.
	ColorPrimary   = lipgloss.Color("#34D399") // Emerald (sole accent)
	ColorSecondary = lipgloss.Color("#CBD5E1") // Slate light (secondary text)
	ColorAccent    = lipgloss.Color("#34D399") // Alias: same accent, no purple
	ColorSuccess   = lipgloss.Color("#34D399") // Emerald green
	ColorWarning   = lipgloss.Color("#FBBF24") // Amber (status only)
	ColorDanger    = lipgloss.Color("#F87171") // Rose red (status only)
	ColorMuted     = lipgloss.Color("#94A3B8") // Slate gray
	ColorDark      = lipgloss.Color("#1E293B") // Dark slate
	ColorBorder    = lipgloss.Color("#334155") // Subtle border slate

	// Text Styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	CategoryStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSecondary)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess)

	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorWarning)

	MutedStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	HighlightStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true)

	// Modern Pill Badges
	BadgePrimary = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#0F172A")).
			Background(ColorPrimary).
			Padding(0, 1).
			Bold(true)

	BadgeSuccess = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#0F172A")).
			Background(ColorSuccess).
			Padding(0, 1).
			Bold(true)

	BadgeWarning = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#0F172A")).
			Background(ColorWarning).
			Padding(0, 1).
			Bold(true)

	BadgeDanger = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorDanger).
			Padding(0, 1).
			Bold(true)

	BadgeMuted = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CBD5E1")).
			Background(lipgloss.Color("#334155")).
			Padding(0, 1)

	BadgeAdmin = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#0F172A")).
			Background(ColorSuccess).
			Padding(0, 1).
			Bold(true)

	// Rounded Card & Box Containers
	SummaryBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 2).
			MarginTop(1)

	CardBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(1, 2).
			MarginTop(1)
)

// WaitEnter prompts the user to press Enter before continuing back to the menu.
func WaitEnter(prompt string) {
	if prompt == "" {
		prompt = "Tekan [Enter] untuk kembali ke menu utama..."
	}
	fmt.Println()
	promptStyled := lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Render("  ↵ " + prompt)
	fmt.Print(promptStyled)
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
	fmt.Println()
}
