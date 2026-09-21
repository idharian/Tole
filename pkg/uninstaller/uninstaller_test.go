package uninstaller

import (
	"testing"
)

func TestGetInstalledApps(t *testing.T) {
	apps, err := GetInstalledApps()
	if err != nil {
		t.Fatalf("Failed to get installed apps: %v", err)
	}

	if len(apps) == 0 {
		t.Log("Warning: No installed apps found (unusual on standard Windows)")
	} else {
		t.Logf("Successfully found %d installed apps in registry", len(apps))
	}
}

func TestFindLeftovers(t *testing.T) {
	// Search for a non-existent app, should return 0 items safely
	leftovers := FindLeftovers("NonExistentAppNameXYZ12345")
	if len(leftovers) != 0 {
		t.Errorf("Expected 0 leftovers for dummy app, got %d", len(leftovers))
	}
}
