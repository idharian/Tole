package monitor

import (
	"strings"
	"testing"
)

func TestModelView(t *testing.T) {
	m := NewModel()
	rendered := m.View()

	if !strings.Contains(rendered, "tole / live monitor") {
		t.Errorf("View missing header title")
	}

	if !strings.Contains(rendered, "CPU") {
		t.Errorf("View missing CPU panel")
	}

	if !strings.Contains(rendered, "Memory") {
		t.Errorf("View missing RAM panel")
	}

	if !strings.Contains(rendered, "Storage C:") {
		t.Errorf("View missing Storage panel")
	}
}
