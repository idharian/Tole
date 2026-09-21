package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanDirectory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "tole_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	subDir := filepath.Join(tempDir, "sub_folder")
	os.Mkdir(subDir, 0755)

	file1 := filepath.Join(tempDir, "file1.txt")
	os.WriteFile(file1, []byte("Hello World 12345"), 0644)

	file2 := filepath.Join(subDir, "file2.txt")
	os.WriteFile(file2, []byte("Hello inside subfolder! More bytes here!"), 0644)

	items, totalSize, err := ScanDirectory(tempDir)
	if err != nil {
		t.Fatalf("ScanDirectory failed: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}

	if totalSize <= 0 {
		t.Errorf("Expected totalSize > 0, got %d", totalSize)
	}
}
