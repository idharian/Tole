package cleaner

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"tole/pkg/winapi"
)

type Category string

const (
	CategoryUserEssentials Category = "User Essentials"
	CategorySystem         Category = "System"
	CategoryBrowsers       Category = "Browsers"
	CategoryDeveloper      Category = "Developer Tools"
)

type TargetType int

const (
	TypeDirectory TargetType = iota
	TypeRecycleBin
)

type Target struct {
	ID           string
	Name         string
	Category     Category
	Type         TargetType
	Paths        []string
	AdminOnly    bool
	Size         int64
	Count        int64
	Scanned      bool
	CleanedSize  int64
	CleanedCount int64
	SkippedCount int64
	ErrorMsg     string
}

// GetDefaultTargets returns all scan targets configured for Windows.
func GetDefaultTargets() []*Target {
	localAppData := winapi.GetLocalAppDataDir()
	appData := winapi.GetRoamingAppDataDir()
	userProfile := winapi.GetUserProfileDir()
	isAdmin := winapi.IsAdmin()

	targets := []*Target{
		// --- USER ESSENTIALS ---
		{
			ID:       "user_temp",
			Name:     "User Temp Files",
			Category: CategoryUserEssentials,
			Type:     TypeDirectory,
			Paths:    []string{winapi.GetUserTempDir()},
		},
		{
			ID:       "crash_dumps",
			Name:     "Crash Dumps & Diagnostics",
			Category: CategoryUserEssentials,
			Type:     TypeDirectory,
			Paths: []string{
				filepath.Join(localAppData, "CrashDumps"),
				filepath.Join(localAppData, "Microsoft", "Windows", "WER"),
			},
		},
		{
			ID:       "recycle_bin",
			Name:     "Recycle Bin",
			Category: CategoryUserEssentials,
			Type:     TypeRecycleBin,
		},

		// --- BROWSERS ---
		{
			ID:       "chrome_cache",
			Name:     "Google Chrome Cache",
			Category: CategoryBrowsers,
			Type:     TypeDirectory,
			Paths: []string{
				filepath.Join(localAppData, "Google", "Chrome", "User Data", "Default", "Cache"),
				filepath.Join(localAppData, "Google", "Chrome", "User Data", "Default", "Code Cache"),
				filepath.Join(localAppData, "Google", "Chrome", "User Data", "Default", "GPUCache"),
			},
		},
		{
			ID:       "edge_cache",
			Name:     "Microsoft Edge Cache",
			Category: CategoryBrowsers,
			Type:     TypeDirectory,
			Paths: []string{
				filepath.Join(localAppData, "Microsoft", "Edge", "User Data", "Default", "Cache"),
				filepath.Join(localAppData, "Microsoft", "Edge", "User Data", "Default", "Code Cache"),
				filepath.Join(localAppData, "Microsoft", "Edge", "User Data", "Default", "GPUCache"),
			},
		},
		{
			ID:       "brave_cache",
			Name:     "Brave Browser Cache",
			Category: CategoryBrowsers,
			Type:     TypeDirectory,
			Paths: []string{
				filepath.Join(localAppData, "BraveSoftware", "Brave-Browser", "User Data", "Default", "Cache"),
				filepath.Join(localAppData, "BraveSoftware", "Brave-Browser", "User Data", "Default", "Code Cache"),
			},
		},
		{
			ID:       "firefox_cache",
			Name:     "Mozilla Firefox Cache",
			Category: CategoryBrowsers,
			Type:     TypeDirectory,
			Paths: []string{
				filepath.Join(localAppData, "Mozilla", "Firefox", "Profiles"),
			},
		},

		// --- DEVELOPER TOOLS ---
		{
			ID:       "npm_cache",
			Name:     "npm Cache",
			Category: CategoryDeveloper,
			Type:     TypeDirectory,
			Paths: []string{
				filepath.Join(localAppData, "npm-cache"),
				filepath.Join(appData, "npm-cache"),
			},
		},
		{
			ID:       "pnpm_cache",
			Name:     "pnpm Cache",
			Category: CategoryDeveloper,
			Type:     TypeDirectory,
			Paths: []string{
				filepath.Join(localAppData, "pnpm", "store"),
			},
		},
		{
			ID:       "yarn_cache",
			Name:     "Yarn Cache",
			Category: CategoryDeveloper,
			Type:     TypeDirectory,
			Paths: []string{
				filepath.Join(localAppData, "Yarn", "Cache"),
			},
		},
		{
			ID:       "pip_cache",
			Name:     "Python pip Cache",
			Category: CategoryDeveloper,
			Type:     TypeDirectory,
			Paths: []string{
				filepath.Join(localAppData, "pip", "cache"),
			},
		},
		{
			ID:       "go_cache",
			Name:     "Go Build Cache",
			Category: CategoryDeveloper,
			Type:     TypeDirectory,
			Paths: []string{
				filepath.Join(localAppData, "go-build"),
			},
		},
		{
			ID:       "cargo_cache",
			Name:     "Rust Cargo Cache",
			Category: CategoryDeveloper,
			Type:     TypeDirectory,
			Paths: []string{
				filepath.Join(userProfile, ".cargo", "registry", "cache"),
			},
		},
		{
			ID:       "nuget_cache",
			Name:     "NuGet HTTP Cache",
			Category: CategoryDeveloper,
			Type:     TypeDirectory,
			Paths: []string{
				filepath.Join(localAppData, "NuGet", "v3-cache"),
			},
		},
	}

	// --- SYSTEM TARGETS (Administrator only) ---
	if isAdmin {
		targets = append(targets,
			&Target{
				ID:        "sys_temp",
				Name:      "Windows System Temp",
				Category:  CategorySystem,
				Type:      TypeDirectory,
				AdminOnly: true,
				Paths:     []string{winapi.GetSystemTempDir()},
			},
			&Target{
				ID:        "windows_update",
				Name:      "Windows Update Download Cache",
				Category:  CategorySystem,
				Type:      TypeDirectory,
				AdminOnly: true,
				Paths:     []string{winapi.GetSoftwareDistributionDownloadDir()},
			},
			&Target{
				ID:        "prefetch",
				Name:      "Windows Prefetch",
				Category:  CategorySystem,
				Type:      TypeDirectory,
				AdminOnly: true,
				Paths:     []string{winapi.GetWindowsPrefetchDir()},
			},
		)
	}

	return targets
}

// walkDirSize computes size + file count in one WalkDir pass,
// skipping reparse points (junctions/symlinks) to avoid cycles.
func walkDirSize(root string) (int64, int64) {
	var size, count int64
	_ = filepath.WalkDir(root, func(_ string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !d.IsDir() {
			if info, ierr := d.Info(); ierr == nil {
				size += info.Size()
				count++
			}
		}
		return nil
	})
	return size, count
}

// cleanDirContents deletes all children of dir in a single WalkDir pass,
// summing sizes of files actually removed. Returns freed bytes, removed
// files, skipped files. Empty dirs are removed bottom-up afterwards.
func cleanDirContents(dir string) (int64, int64, int64) {
	var freed, removed, skipped int64
	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if path == dir {
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			if os.Remove(path) == nil {
				removed++
			} else {
				skipped++
			}
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !d.IsDir() {
			var sz int64
			if info, ierr := d.Info(); ierr == nil {
				sz = info.Size()
			}
			if os.Remove(path) == nil {
				freed += sz
				removed++
			} else {
				skipped++
			}
		}
		return nil
	})
	removeEmptyDirs(dir)
	return freed, removed, skipped
}

// removeEmptyDirs removes empty subdirs deepest-first, best effort.
func removeEmptyDirs(root string) {
	var dirs []string
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err == nil && d.IsDir() && path != root {
			dirs = append(dirs, path)
		}
		return nil
	})
	for i := len(dirs) - 1; i >= 0; i-- {
		_ = os.Remove(dirs[i]) // fails if non-empty, which is fine
	}
}
func ScanTarget(t *Target) {
	if t.Type == TypeRecycleBin {
		size, count, err := winapi.GetRecycleBinInfo()
		if err == nil {
			t.Size = size
			t.Count = count
		}
		t.Scanned = true
		return
	}

	var totalSize, totalCount int64
	for _, p := range t.Paths {
		if _, err := os.Stat(p); err != nil {
			continue
		}
		if t.ID == "firefox_cache" {
			_ = filepath.WalkDir(p, func(path string, d os.DirEntry, err error) error {
				if err != nil {
					return nil
				}
				if d.IsDir() && strings.EqualFold(d.Name(), "cache2") {
					size, count := walkDirSize(path)
					totalSize += size
					totalCount += count
					return filepath.SkipDir
				}
				return nil
			})
			continue
		}
		size, count := walkDirSize(p)
		totalSize += size
		totalCount += count
	}
	t.Size = totalSize
	t.Count = totalCount
	t.Scanned = true
}

// ScanAll scans targets with bounded concurrency to avoid disk thrashing.
func ScanAll(targets []*Target) {
	workers := runtime.GOMAXPROCS(0)
	if workers > 6 {
		workers = 6
	}
	if workers < 2 {
		workers = 2
	}
	jobs := make(chan *Target)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for target := range jobs {
				ScanTarget(target)
			}
		}()
	}
	for _, target := range targets {
		jobs <- target
	}
	close(jobs)
	wg.Wait()
}

// CleanTarget cleans items from the target and records reclaimed bytes.
func CleanTarget(t *Target, dryRun bool) {
	if dryRun {
		t.CleanedSize = t.Size
		t.CleanedCount = t.Count
		return
	}

	if t.Type == TypeRecycleBin {
		if t.Count > 0 {
			err := winapi.EmptyRecycleBin(true)
			if err == nil {
				t.CleanedSize = t.Size
				t.CleanedCount = t.Count
			} else {
				t.ErrorMsg = err.Error()
			}
		}
		return
	}

	var skippedItems int64
	for _, p := range t.Paths {
		if t.ID == "firefox_cache" {
			_ = filepath.WalkDir(p, func(path string, d os.DirEntry, err error) error {
				if err != nil {
					return nil
				}
				if d.IsDir() && strings.EqualFold(d.Name(), "cache2") {
					entries, rerr := os.ReadDir(path)
					if rerr == nil {
						for _, entry := range entries {
							if os.RemoveAll(filepath.Join(path, entry.Name())) != nil {
								skippedItems++
							}
						}
					}
					return filepath.SkipDir
				}
				return nil
			})
			continue
		}

		entries, err := os.ReadDir(p)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if os.RemoveAll(filepath.Join(p, entry.Name())) != nil {
				skippedItems++
			}
		}
	}

	// Scan already measured the payload. Avoid walking every item again.
	if skippedItems == 0 {
		t.CleanedSize = t.Size
		t.CleanedCount = t.Count
	}
	t.SkippedCount = skippedItems
}
