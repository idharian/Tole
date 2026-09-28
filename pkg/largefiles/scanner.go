package largefiles

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type LargeFile struct {
	Path string
	Name string
	Size int64
}

// FindLargeFiles walks root recursively and returns the largest regular files
// above minSize bytes, sorted descending and capped at topN entries.
// Symlinked directories are skipped to avoid cycles and double counting.
func FindLargeFiles(root string, minSize int64, topN int) ([]*LargeFile, int64, error) {
	var all []*LargeFile
	var scanned int64

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // unreadable entries are skipped, not fatal
		}
		if d.Type()&os.ModeSymlink != 0 {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			return nil
		}
		info, ierr := d.Info()
		if ierr != nil {
			return nil
		}
		size := info.Size()
		scanned += size
		if size < minSize {
			return nil
		}
		all = append(all, &LargeFile{Path: path, Name: d.Name(), Size: size})
		return nil
	})
	if err != nil {
		return nil, scanned, err
	}

	sort.Slice(all, func(i, j int) bool {
		if all[i].Size != all[j].Size {
			return all[i].Size > all[j].Size
		}
		return strings.ToLower(all[i].Name) < strings.ToLower(all[j].Name)
	})

	if topN > 0 && len(all) > topN {
		all = all[:topN]
	}
	return all, scanned, nil
}
