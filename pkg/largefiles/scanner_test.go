package largefiles

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestFindLargeFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "tole_large_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	subDir := filepath.Join(tempDir, "sub")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create sub dir: %v", err)
	}

	big := filepath.Join(tempDir, "big.bin")
	if err := os.WriteFile(big, make([]byte, 2048), 0644); err != nil {
		t.Fatalf("Failed to write big file: %v", err)
	}

	small := filepath.Join(tempDir, "small.txt")
	if err := os.WriteFile(small, []byte("tiny"), 0644); err != nil {
		t.Fatalf("Failed to write small file: %v", err)
	}

	deep := filepath.Join(subDir, "deep.bin")
	if dErr := os.WriteFile(deep, make([]byte, 4096), 0644); dErr != nil {
		t.Fatalf("Failed to write deep file: %v", dErr)
	}

	files, totalSize, err := FindLargeFiles(tempDir, 1024, 10)
	if err != nil {
		t.Fatalf("FindLargeFiles failed: %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("Expected 2 large files, got %d", len(files))
	}
	if files[0].Name != "deep.bin" {
		t.Errorf("Expected largest first (deep.bin), got %s", files[0].Name)
	}
	if files[0].Size != 4096 || files[1].Size != 2048 {
		t.Errorf("Unexpected sizes: %d, %d", files[0].Size, files[1].Size)
	}
	if totalSize != 2048+4096+4 {
		t.Errorf("Expected total scanned 6148, got %d", totalSize)
	}
}

func TestFindLargeFilesTopN(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "tole_large_top_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	for i := 0; i < 5; i++ {
		name := filepath.Join(tempDir, fmt.Sprintf("f%d.bin", i))
		if wErr := os.WriteFile(name, make([]byte, 1024+i*100), 0644); wErr != nil {
			t.Fatalf("Failed to write file: %v", wErr)
		}
	}

	files, _, err := FindLargeFiles(tempDir, 0, 3)
	if err != nil {
		t.Fatalf("FindLargeFiles failed: %v", err)
	}
	if len(files) != 3 {
		t.Errorf("Expected top 3 files, got %d", len(files))
	}
}
