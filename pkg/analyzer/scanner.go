package analyzer

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
)

type ItemType int

const (
	ItemDir ItemType = iota
	ItemFile
)

type FileItem struct {
	Name       string
	Path       string
	Type       ItemType
	Size       int64
	ItemCount  int64
	Percentage float64
}

// CalculateSize uses WalkDir and skips symlinked directories.
func CalculateSize(path string) (int64, int64) {
	var totalSize, count int64
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
				count++
			}
		}
		return nil
	})
	return totalSize, count
}

// ScanDirectory reads immediate children, with bounded directory workers.
func ScanDirectory(dir string) ([]*FileItem, int64, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, 0, err
	}
	items := make([]*FileItem, 0, len(entries))
	for _, entry := range entries {
		info, ierr := entry.Info()
		if ierr != nil {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		item := &FileItem{Name: entry.Name(), Path: path, Type: ItemFile, Size: info.Size(), ItemCount: 1}
		if entry.IsDir() {
			item.Type = ItemDir
			item.Size = 0
			item.ItemCount = 0
		}
		items = append(items, item)
	}
	workers := runtime.GOMAXPROCS(0)
	if workers > 6 {
		workers = 6
	}
	if workers < 1 {
		workers = 1
	}
	jobs := make(chan *FileItem)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range jobs {
				item.Size, item.ItemCount = CalculateSize(item.Path)
			}
		}()
	}
	for _, item := range items {
		if item.Type == ItemDir {
			jobs <- item
		}
	}
	close(jobs)
	wg.Wait()
	var totalSize int64
	for _, item := range items {
		totalSize += item.Size
	}
	for _, item := range items {
		if totalSize > 0 {
			item.Percentage = float64(item.Size) / float64(totalSize) * 100
		}
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Size != items[j].Size {
			return items[i].Size > items[j].Size
		}
		return items[i].Name < items[j].Name
	})
	return items, totalSize, nil
}
