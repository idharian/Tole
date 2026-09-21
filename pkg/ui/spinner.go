package ui

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var SpinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

type Spinner struct {
	message string
	active  bool
	mu      sync.Mutex
	stop    chan struct{}
}

// StartSpinner starts a smooth animated CLI spinner with a message.
func StartSpinner(message string) *Spinner {
	s := &Spinner{
		message: message,
		active:  true,
		stop:    make(chan struct{}),
	}

	go func() {
		idx := 0
		spinnerStyle := lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)
		msgStyle := lipgloss.NewStyle().Foreground(ColorMuted)

		for {
			select {
			case <-s.stop:
				return
			default:
				s.mu.Lock()
				msg := s.message
				s.mu.Unlock()

				frame := SpinnerFrames[idx%len(SpinnerFrames)]
				line := fmt.Sprintf("  %s %s", spinnerStyle.Render(frame), msgStyle.Render(msg))
				// Pad with spaces to clear trailing characters if message shortens
				if len(line) < 70 {
					line += strings.Repeat(" ", 70-len(line))
				}
				fmt.Printf("\r%s", line)
				idx++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()

	return s
}

// UpdateMessage updates the active text next to the spinner.
func (s *Spinner) UpdateMessage(msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.message = msg
}

// Stop stops the spinner and writes a clean final status line.
func (s *Spinner) Stop(statusIcon, completionMessage string) {
	s.mu.Lock()
	if !s.active {
		s.mu.Unlock()
		return
	}
	s.active = false
	close(s.stop)
	s.mu.Unlock()

	// Wipe line with spaces and print final status
	blank := strings.Repeat(" ", 75)
	fmt.Printf("\r%s\r  %s %s\n", blank, statusIcon, completionMessage)
}

func (s *Spinner) StopSuccess(message string) {
	s.Stop(SuccessStyle.Render("✓"), HighlightStyle.Render(message))
}

func (s *Spinner) StopWarning(message string) {
	s.Stop(WarningStyle.Render("!"), WarningStyle.Render(message))
}

func (s *Spinner) StopError(message string) {
	s.Stop(lipgloss.NewStyle().Foreground(ColorDanger).Render("✗"), WarningStyle.Render(message))
}
