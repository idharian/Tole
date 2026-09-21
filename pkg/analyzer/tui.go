package analyzer

import (
	"fmt"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"tole/pkg/ui"
	"tole/pkg/winapi"
)

type scanResultMsg struct {
	items     []*FileItem
	totalSize int64
	err       error
}

type Model struct {
	currentDir    string
	history       []string
	items         []*FileItem
	cursor        int
	totalSize     int64
	loading       bool
	confirmDelete bool
	statusMsg     string
	width         int
	height        int
}

func NewModel(initialDir string) Model {
	abs, err := filepath.Abs(initialDir)
	if err != nil {
		abs = initialDir
	}
	return Model{
		currentDir: abs,
		loading:    true,
		cursor:     0,
	}
}

func scanCmd(dir string) tea.Cmd {
	return func() tea.Msg {
		items, total, err := ScanDirectory(dir)
		return scanResultMsg{items: items, totalSize: total, err: err}
	}
}

func (m Model) Init() tea.Cmd {
	return scanCmd(m.currentDir)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case scanResultMsg:
		m.loading = false
		if msg.err != nil {
			m.statusMsg = fmt.Sprintf("Error scanning: %v", msg.err)
		} else {
			m.items = msg.items
			m.totalSize = msg.totalSize
			m.cursor = 0
			m.statusMsg = ""
		}
		return m, nil

	case tea.KeyMsg:
		if m.confirmDelete {
			switch msg.String() {
			case "y", "Y":
				if len(m.items) > 0 && m.cursor < len(m.items) {
					target := m.items[m.cursor]
					err := winapi.MoveToRecycleBin(target.Path)
					if err != nil {
						m.statusMsg = fmt.Sprintf("Error moving to Recycle Bin: %v", err)
					} else {
						m.statusMsg = fmt.Sprintf("Moved '%s' to Recycle Bin", target.Name)
						m.loading = true
						m.confirmDelete = false
						return m, scanCmd(m.currentDir)
					}
				}
				m.confirmDelete = false
				return m, nil
			case "n", "N", "esc":
				m.confirmDelete = false
				m.statusMsg = "Deletion cancelled"
				return m, nil
			default:
				return m, nil
			}
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}

		case "enter", "right", "l":
			if len(m.items) > 0 && m.cursor < len(m.items) {
				selected := m.items[m.cursor]
				if selected.Type == ItemDir {
					m.history = append(m.history, m.currentDir)
					m.currentDir = selected.Path
					m.loading = true
					return m, scanCmd(m.currentDir)
				}
			}

		case "esc", "backspace", "left", "h":
			if len(m.history) > 0 {
				parent := m.history[len(m.history)-1]
				m.history = m.history[:len(m.history)-1]
				m.currentDir = parent
				m.loading = true
				return m, scanCmd(m.currentDir)
			} else {
				// Navigate up one folder if possible
				parent := filepath.Dir(m.currentDir)
				if parent != m.currentDir {
					m.currentDir = parent
					m.loading = true
					return m, scanCmd(m.currentDir)
				}
			}

		case "d":
			if len(m.items) > 0 && m.cursor < len(m.items) {
				m.confirmDelete = true
			}
		}
	}

	return m, nil
}

func renderBar(percentage float64, barWidth int) string {
	if barWidth <= 0 {
		barWidth = 15
	}
	filledLen := int((percentage / 100.0) * float64(barWidth))
	if filledLen > barWidth {
		filledLen = barWidth
	}
	emptyLen := barWidth - filledLen

	filled := strings.Repeat("█", filledLen)
	empty := strings.Repeat("░", emptyLen)

	return ui.HighlightStyle.Render(filled) + ui.MutedStyle.Render(empty)
}

func (m Model) View() string {
	var b strings.Builder

	// Header: slim title + context, matching banner rhythm
	b.WriteString(ui.HighlightStyle.Render("  tole / disk analyzer\n"))
	b.WriteString(ui.MutedStyle.Render(fmt.Sprintf("  %s  •  %s total, %d items\n\n", m.currentDir, ui.FormatBytes(m.totalSize), len(m.items))))

	if m.loading {
		b.WriteString(ui.MutedStyle.Render("  Scanning directory... please wait.\n\n"))
		return b.String()
	}

	if len(m.items) == 0 {
		b.WriteString(ui.MutedStyle.Render("  (Directory is empty)\n\n"))
	} else {
		// Calculate visible window
		maxRows := 14
		start := 0
		if m.cursor >= maxRows {
			start = m.cursor - maxRows + 1
		}
		end := start + maxRows
		if end > len(m.items) {
			end = len(m.items)
		}

		for i := start; i < end; i++ {
			item := m.items[i]
			isCursor := i == m.cursor

			cursorPrefix := "  "
			if isCursor {
				cursorPrefix = ui.HighlightStyle.Render("➤ ")
			}

			icon := "dir "
			if item.Type == ItemFile {
				icon = "file"
			}

			nameStyle := lipgloss.NewStyle()
			if isCursor {
				nameStyle = ui.HighlightStyle
			} else if item.Type == ItemDir {
				nameStyle = ui.CategoryStyle
			}

			name := item.Name
			if len(name) > 30 {
				name = name[:27] + "..."
			}

			sizeStr := ui.FormatBytes(item.Size)
			bar := renderBar(item.Percentage, 12)

			line := fmt.Sprintf("%s%s%-30s  %10s  %s %5.1f%%",
				cursorPrefix,
				icon,
				nameStyle.Render(name),
				sizeStr,
				bar,
				item.Percentage,
			)

			if isCursor {
				line = lipgloss.NewStyle().Background(lipgloss.Color("#1E293B")).Render(line)
			}
			b.WriteString(line + "\n")
		}
		b.WriteString("\n")
	}

	// Status / Confirm bar
	if m.confirmDelete {
		selectedName := m.items[m.cursor].Name
		confirmBox := ui.WarningStyle.Render(fmt.Sprintf("  Move '%s' to Recycle Bin? [y] Yes / [n] Cancel", selectedName))
		b.WriteString(confirmBox + "\n\n")
	} else if m.statusMsg != "" {
		b.WriteString(ui.MutedStyle.Render("  "+m.statusMsg) + "\n\n")
	}

	// Keybindings footer
	help := "  ↑/↓ navigate • enter open • esc parent • d recycle • q exit"
	b.WriteString(ui.MutedStyle.Render(help) + "\n")

	return b.String()
}
