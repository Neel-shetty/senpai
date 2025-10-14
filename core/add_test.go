package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAddFile(t *testing.T) {
	tmpDir := t.TempDir()
	oldRepo := repoDirName
	repoDirName = ".senpai"
	t.Cleanup(func() { repoDirName = oldRepo })

	repoPath := tmpDir
	if err := InitRepo(repoPath, "master"); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	filePath := filepath.Join(repoPath, "hello.txt")
	if err := os.WriteFile(filePath, []byte("Hello world!"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	if err := Add(repoPath, filePath); err != nil {
		t.Fatalf("AddFile failed: %v", err)
	}

	indexPath := filepath.Join(repoPath, ".senpai", "index")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("failed to read index: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "100644") || !strings.Contains(content, "hello.txt") {
		t.Errorf("index missing expected entry, got: %q", content)
	}

	if err := Add(repoPath, filePath); err != nil {
		t.Fatalf("AddFile failed second time: %v", err)
	}

	data2, _ := os.ReadFile(indexPath)
	lines := strings.Split(strings.TrimSpace(string(data2)), "\n")
	if len(lines) != 1 {
		t.Errorf("expected single entry in index, got %d lines", len(lines))
	}
}
