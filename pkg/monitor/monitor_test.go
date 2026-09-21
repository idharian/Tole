package monitor

import (
	"testing"
)

func TestMonitorStats(t *testing.T) {
	m := NewMonitor()
	stats := m.GetStats()

	if stats.TotalRAM == 0 {
		t.Errorf("Expected TotalRAM > 0, got 0")
	}

	if stats.TotalDisk == 0 {
		t.Errorf("Expected TotalDisk > 0, got 0")
	}

	if stats.NumCPU <= 0 {
		t.Errorf("Expected NumCPU > 0, got %d", stats.NumCPU)
	}

	t.Logf("Stats: RAM %d MB (%.1f%%), CPU: %.1f%%, Cores: %d",
		stats.TotalRAM/(1024*1024), stats.RAMPercent, stats.CPUPercent, stats.NumCPU)
}
