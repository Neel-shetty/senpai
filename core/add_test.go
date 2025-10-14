package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAddFiles(t *testing.T) {
	tmpDir := t.TempDir()
	oldRepo := repoDirName
	repoDirName = ".senpai"
	t.Cleanup(func() { repoDirName = oldRepo })

	repoPath := tmpDir
	if err := InitRepo(repoPath, "master"); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	// Create test files
	file1 := filepath.Join(repoPath, "hello.txt")
	file2 := filepath.Join(repoPath, "world.txt")
	if err := os.WriteFile(file1, []byte("Hello world!"), 0644); err != nil {
		t.Fatalf("failed to create file1: %v", err)
	}
	if err := os.WriteFile(file2, []byte("Another file"), 0644); err != nil {
		t.Fatalf("failed to create file2: %v", err)
	}

	// Add single file
	if err := Add(repoPath, file1); err != nil {
		t.Fatalf("Add single file failed: %v", err)
	}

	// Add multiple files
	if err := Add(repoPath, file1, file2); err != nil {
		t.Fatalf("Add multiple files failed: %v", err)
	}

	indexPath := filepath.Join(repoPath, ".senpai", "index")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("failed to read index: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 entries in index, got %d lines", len(lines))
	}

	content := string(data)
	for _, f := range []string{"hello.txt", "world.txt"} {
		if !strings.Contains(content, f) {
			t.Errorf("index missing expected entry %q", f)
		}
	}

	// Add same files again, ensure no duplicates
	if err := Add(repoPath, file1, file2); err != nil {
		t.Fatalf("Add duplicate files failed: %v", err)
	}
	data2, _ := os.ReadFile(indexPath)
	lines2 := strings.Split(strings.TrimSpace(string(data2)), "\n")
	if len(lines2) != 2 {
		t.Errorf("expected 2 entries after duplicate add, got %d lines", len(lines2))
	}
}
