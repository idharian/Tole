package uninstaller

import (
	"os"
	"path/filepath"
	"strings"
	"tole/pkg/winapi"
)

type LeftoverItem struct {
	Path        string
	Name        string
	Location    string
	Size        int64
	FileCount   int64
}

// FindLeftovers searches AppData, ProgramData, and common locations for folders matching appName.
func FindLeftovers(appName string) []*LeftoverItem {
	if appName == "" {
		return nil
	}

	searchDirs := []struct {
		Label string
		Path  string
	}{
		{"LocalAppData", winapi.GetLocalAppDataDir()},
		{"RoamingAppData", winapi.GetRoamingAppDataDir()},
		{"ProgramData", winapi.GetProgramDataDir()},
	}

	var results []*LeftoverItem
	normalizedName := strings.ToLower(strings.TrimSpace(appName))

	// Significant words in the app name (skip generic terms)
	genericWords := map[string]bool{
		"setup": true, "install": true, "installer": true, "client": true,
		"desktop": true, "tools": true, "tool": true, "windows": true,
		"free": true, "pro": true, "suite": true, "version": true, "software": true,
		"app": true, "program": true, "manager": true, "studio": true, "code": true,
		"system": true, "update": true, "helper": true, "launcher": true, "player": true,
	}
	var sigWords []string
	for _, w := range strings.Fields(normalizedName) {
		if len(w) >= 3 && !genericWords[w] {
			sigWords = append(sigWords, w)
		}
	}
	if len(sigWords) == 0 {
		return nil
	}

	for _, s := range searchDirs {
		entries, err := os.ReadDir(s.Path)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			entryLower := strings.ToLower(entry.Name())

			// Skip standard Windows/Microsoft/Core directories
			skipFolders := map[string]bool{
				"microsoft":        true,
				"windows":          true,
				"temp":             true,
				"packages":         true,
				"system":           true,
				"system32":         true,
				"common files":     true,
				"appdata":          true,
				"users":            true,
				"default":          true,
				"all users":        true,
				"winstore":         true,
				"microsoft shared": true,
				"directx":          true,
				"windows defender": true,
				"windowsapps":      true,
				"crashdumps":       true,
				"wer":              true,
			}
			if skipFolders[entryLower] {
				continue
			}

			// Strict match: exact name, or every significant word present.
			matched := false
			if entryLower == normalizedName || strings.HasPrefix(entryLower, normalizedName+"-") || strings.HasPrefix(entryLower, normalizedName+"_") {
				matched = true
			} else {
				allPresent := true
				for _, w := range sigWords {
					if !strings.Contains(entryLower, w) {
						allPresent = false
						break
					}
				}
				matched = allPresent
			}

			if matched {
				fullPath := filepath.Join(s.Path, entry.Name())
				var dirSize int64
				var fileCount int64

				filepath.Walk(fullPath, func(_ string, info os.FileInfo, err error) error {
					if err == nil && !info.IsDir() {
						dirSize += info.Size()
						fileCount++
					}
					return nil
				})

				results = append(results, &LeftoverItem{
					Path:      fullPath,
					Name:      entry.Name(),
					Location:  s.Label,
					Size:      dirSize,
					FileCount: fileCount,
				})
			}
		}
	}

	return results
}
