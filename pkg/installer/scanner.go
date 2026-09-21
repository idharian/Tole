package installer

import (
	"os"
	"path/filepath"
	"strings"
	"tole/pkg/winapi"
)

type InstallerFile struct {
	Path      string
	Name      string
	Extension string
	Size      int64
}

// ScanInstallers scans user Downloads and Desktop directories for installer packages.
func ScanInstallers() ([]*InstallerFile, int64, error) {
	userProfile := winapi.GetUserProfileDir()
	dirsToScan := []string{
		filepath.Join(userProfile, "Downloads"),
		filepath.Join(userProfile, "Desktop"),
	}

	var results []*InstallerFile
	var totalSize int64

	// scanDir lists installer candidates in dir; recurse indicates subdir depth left.
	var scanDir func(dir string, depth int)
	scanDir = func(dir string, depth int) {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return
		}

		for _, entry := range entries {
			name := entry.Name()
			full := filepath.Join(dir, name)

			if entry.IsDir() {
				if depth > 0 && !strings.HasPrefix(name, ".") {
					scanDir(full, depth-1)
				}
				continue
			}

			ext := strings.ToLower(filepath.Ext(name))

			isInstaller := false
			if ext == ".msi" || ext == ".msix" || ext == ".appx" || ext == ".iso" || ext == ".img" {
				isInstaller = true
			} else if ext == ".exe" {
				lowerName := strings.ToLower(name)
				if strings.Contains(lowerName, "setup") ||
					strings.Contains(lowerName, "install") ||
					strings.Contains(lowerName, "update") ||
					strings.Contains(lowerName, "patch") {
					isInstaller = true
				}
			}

			if isInstaller {
				info, err := entry.Info()
				if err != nil {
					continue
				}

				sz := info.Size()
				// Skip tiny stubs (<1MB) to reduce noise
				if sz < 1024*1024 {
					continue
				}
				results = append(results, &InstallerFile{
					Path:      full,
					Name:      name,
					Extension: ext,
					Size:      sz,
				})
				totalSize += sz
			}
		}
	}

	for _, dir := range dirsToScan {
		scanDir(dir, 1)
	}

	return results, totalSize, nil
}
