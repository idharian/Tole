package purge

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
)

// ArtifactFolder names that are candidates for purging
var DefaultArtifactDirs = map[string]string{
	"node_modules": "Node.js Dependencies",
	"target":       "Rust Cargo Target",
	"bin":          ".NET / Binary Build",
	"obj":          ".NET Intermediate Build",
	".gradle":      "Gradle Build Cache",
	".venv":        "Python Virtual Environment",
	"venv":         "Python Virtual Environment",
	"__pycache__":  "Python Bytecode Cache",
	"dist":         "Frontend Build Output",
	".next":        "Next.js Build Cache",
	".nuxt":        "Nuxt.js Build Cache",
}

type FoundArtifact struct {
	Path        string
	Name        string
	Description string
	Size        int64
	ItemCount   int64
}

// CalculateDirSize recursively computes the size and file count of a directory.
func CalculateDirSize(path string) (int64, int64) {
	var totalSize, totalCount int64
	_ = filepath.WalkDir(path, func(_ string, d os.DirEntry, err error) error {
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
				totalSize += info.Size()
				totalCount++
			}
		}
		return nil
	})
	return totalSize, totalCount
}

func isProtectedSystemPath(p string) bool {
	pLower := strings.ToLower(p)
	systemRoots := []string{
		`c:\windows`,
		`c:\program files`,
		`c:\program files (x86)`,
		`c:\programdata`,
		`c:\$recycle.bin`,
		`c:\system volume information`,
		`c:\recovery`,
		`c:\boot`,
	}
	for _, sys := range systemRoots {
		if strings.HasPrefix(pLower, sys) {
			return true
		}
	}
	return false
}

func isDevProjectDirectory(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	projectMarkers := map[string]bool{
		"package.json":     true,
		"cargo.toml":       true,
		"go.mod":           true,
		"pom.xml":          true,
		"build.gradle":     true,
		"build.gradle.kts": true,
		"requirements.txt": true,
		"pyproject.toml":   true,
		"makefile":         true,
		"cmakelists.txt":   true,
	}
	for _, e := range entries {
		name := strings.ToLower(e.Name())
		if projectMarkers[name] {
			return true
		}
		if strings.HasSuffix(name, ".csproj") || strings.HasSuffix(name, ".sln") || strings.HasSuffix(name, ".fsproj") {
			return true
		}
	}
	return false
}

// ScanProjects scans rootDir for build artifact folders.
func ScanProjects(rootDir string, maxDepth int) ([]*FoundArtifact, error) {
	if rootDir == "" {
		rootDir = "."
	}

	absRoot, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, err
	}

	var results []*FoundArtifact
	var mu sync.Mutex
	var wg sync.WaitGroup
	workers := runtime.GOMAXPROCS(0)
	if workers > 6 {
		workers = 6
	}
	jobs := make(chan [3]string)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				size, count := CalculateDirSize(job[0])
				if size > 0 || count > 0 {
					mu.Lock()
					results = append(results, &FoundArtifact{Path: job[0], Name: job[1], Description: job[2], Size: size, ItemCount: count})
					mu.Unlock()
				}
			}
		}()
	}

	rootDepth := strings.Count(absRoot, string(os.PathSeparator))
	_ = filepath.WalkDir(absRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !entry.IsDir() {
			return nil
		}
		if isProtectedSystemPath(path) {
			return filepath.SkipDir
		}
		name := entry.Name()
		if name == ".git" || name == "$RECYCLE.BIN" || name == "System Volume Information" {
			return filepath.SkipDir
		}
		if maxDepth > 0 && strings.Count(path, string(os.PathSeparator))-rootDepth > maxDepth {
			return filepath.SkipDir
		}
		lowerName := strings.ToLower(name)
		if desc, ok := DefaultArtifactDirs[lowerName]; ok {
			if (lowerName == "bin" || lowerName == "obj") && !isDevProjectDirectory(filepath.Dir(path)) {
				return nil
			}
			jobs <- [3]string{path, name, desc}
			return filepath.SkipDir
		}
		return nil
	})
	close(jobs)
	wg.Wait()
	sort.Slice(results, func(i, j int) bool { return results[i].Path < results[j].Path })
	return results, nil
}

// DeleteArtifact removes the artifact directory.
func DeleteArtifact(artifact *FoundArtifact, dryRun bool) (int64, error) {
	if dryRun {
		return artifact.Size, nil
	}

	err := os.RemoveAll(artifact.Path)
	if err != nil {
		return 0, err
	}

	return artifact.Size, nil
}

// PurgeAll removes all passed artifacts.
func PurgeAll(artifacts []*FoundArtifact, dryRun bool) (reclaimedBytes int64, errCount int64) {
	var totalReclaimed int64
	var totalErrors int64

	for _, a := range artifacts {
		reclaimed, err := DeleteArtifact(a, dryRun)
		if err != nil {
			atomic.AddInt64(&totalErrors, 1)
		} else {
			atomic.AddInt64(&totalReclaimed, reclaimed)
		}
	}

	return totalReclaimed, totalErrors
}
