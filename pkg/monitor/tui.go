package monitor

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"tole/pkg/ui"
	"tole/pkg/version"
)

type tickMsg time.Time

type Model struct {
	monitor *Monitor
	stats   SystemStats
	width   int
	height  int
}

func NewModel() Model {
	m := NewMonitor()
	return Model{
		monitor: m,
		stats:   m.GetStats(),
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		m.stats = m.monitor.GetStats()
		return m, tickCmd()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	return m, nil
}

func renderBar(percent float64, width int, filledColor lipgloss.Color) string {
	if width <= 0 {
		width = 36
	}
	filled := int((percent / 100.0) * float64(width))
	if filled > width {
		filled = width
	}
	empty := width - filled

	filledStr := strings.Repeat("█", filled)
	emptyStr := strings.Repeat("░", empty)

	return lipgloss.NewStyle().Foreground(filledColor).Render(filledStr) +
		lipgloss.NewStyle().Foreground(lipgloss.Color("#334155")).Render(emptyStr)
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	mins := int(d.Minutes()) % 60
	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, mins)
	}
	return fmt.Sprintf("%dh %dm", hours, mins)
}

func (m Model) View() string {
	var b strings.Builder

	// Styles: single emerald accent, slate borders
	panelStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#334155")).
		Padding(0, 1).
		Width(68).
		MarginBottom(1)

	headerTitle := ui.HighlightStyle.Render("  tole / live monitor")

	devSubtitle := ui.MutedStyle.Render(
		fmt.Sprintf("  %s  •  v%s  •  polling 1000ms", m.stats.Hostname, version.CurrentVersion))

	b.WriteString("\n" + headerTitle + "\n" + devSubtitle + "\n\n")

	// Context pills: single muted style, one live accent
	hostPill := ui.BadgeMuted.Render(m.stats.Hostname)
	osPill := ui.BadgeMuted.Render(fmt.Sprintf("win %s", m.stats.Arch))
	corePill := ui.BadgeMuted.Render(fmt.Sprintf("%d cores", m.stats.NumCPU))
	uptimePill := ui.BadgeMuted.Render(formatDuration(m.stats.Uptime))
	livePill := ui.BadgeSuccess.Render("live")

	b.WriteString(fmt.Sprintf("  %s %s %s %s %s\n\n",
		hostPill, osPill, corePill, uptimePill, livePill))

	barWidth := 38

	// 1. CPU PANEL
	cpuColor := ui.ColorSuccess
	cpuStatusBadge := ui.BadgeSuccess.Render("OPTIMAL")
	if m.stats.CPUPercent > 80 {
		cpuColor = ui.ColorDanger
		cpuStatusBadge = ui.BadgeDanger.Render("HIGH LOAD")
	} else if m.stats.CPUPercent > 50 {
		cpuColor = ui.ColorWarning
		cpuStatusBadge = ui.BadgeWarning.Render("ELEVATED")
	}

	cpuContent := fmt.Sprintf(
		"  %s  %s\n"+
			"  [%s]  %5.1f%%\n\n"+
			"  %-20s : %s\n"+
			"  %-20s : %s",
		ui.MutedStyle.Render("CPU"),
		cpuStatusBadge,
		renderBar(m.stats.CPUPercent, barWidth, cpuColor),
		m.stats.CPUPercent,
		ui.MutedStyle.Render("Hardware Cores"),
		fmt.Sprintf("%d physical/logical cores", m.stats.NumCPU),
		ui.MutedStyle.Render("Active Goroutines"),
		fmt.Sprintf("%d runtime tasks", m.stats.NumGoroutine),
	)
	b.WriteString(panelStyle.Render(cpuContent) + "\n")

	// 2. RAM PANEL
	ramColor := ui.ColorSuccess
	ramStatusBadge := ui.BadgeSuccess.Render("HEALTHY")
	if m.stats.RAMPercent > 85 {
		ramColor = ui.ColorDanger
		ramStatusBadge = ui.BadgeDanger.Render("CRITICAL")
	} else if m.stats.RAMPercent > 65 {
		ramColor = ui.ColorWarning
		ramStatusBadge = ui.BadgeWarning.Render("MODERATE")
	}

	usedRAMStr := ui.FormatBytes(int64(m.stats.UsedRAM))
	totalRAMStr := ui.FormatBytes(int64(m.stats.TotalRAM))
	freeRAMStr := ui.FormatBytes(int64(m.stats.FreeRAM))

	ramContent := fmt.Sprintf(
		"  %s  %s\n"+
			"  [%s]  %5.1f%%\n\n"+
			"  %-20s : %s / %s\n"+
			"  %-20s : %s",
		ui.MutedStyle.Render("Memory"),
		ramStatusBadge,
		renderBar(m.stats.RAMPercent, barWidth, ramColor),
		m.stats.RAMPercent,
		ui.MutedStyle.Render("Memory In Use"),
		usedRAMStr, totalRAMStr,
		ui.MutedStyle.Render("Available Standby"),
		ui.SuccessStyle.Render(freeRAMStr),
	)
	b.WriteString(panelStyle.Render(ramContent) + "\n")

	// 3. STORAGE PANEL (DRIVE C:)
	var diskPercent float64
	if m.stats.TotalDisk > 0 {
		diskPercent = (float64(m.stats.UsedDisk) / float64(m.stats.TotalDisk)) * 100.0
	}
	usedDiskStr := ui.FormatBytes(int64(m.stats.UsedDisk))
	totalDiskStr := ui.FormatBytes(int64(m.stats.TotalDisk))
	freeDiskStr := ui.FormatBytes(int64(m.stats.FreeDisk))

	diskBadge := ui.BadgeSuccess.Render("PLENTY")
	if diskPercent > 90 {
		diskBadge = ui.BadgeDanger.Render("CRITICAL")
	} else if diskPercent > 75 {
		diskBadge = ui.BadgeWarning.Render("LOW SPACE")
	}

	diskContent := fmt.Sprintf(
		"  %s  %s\n"+
			"  [%s]  %5.1f%%\n\n"+
			"  %-20s : %s / %s\n"+
			"  %-20s : %s",
		ui.MutedStyle.Render("Storage C:"),
		diskBadge,
		renderBar(diskPercent, barWidth, ui.ColorPrimary),
		diskPercent,
		ui.MutedStyle.Render("Storage In Use"),
		usedDiskStr, totalDiskStr,
		ui.MutedStyle.Render("Free Space"),
		ui.HighlightStyle.Render(freeDiskStr),
	)
	b.WriteString(panelStyle.Render(diskContent) + "\n")

	// Footer Help
	b.WriteString(ui.MutedStyle.Render("  q / esc exit") + "\n")

	return b.String()
}
