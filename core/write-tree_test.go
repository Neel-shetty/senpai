package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteTree(t *testing.T) {
	tmpDir := t.TempDir()
	oldRepo := RepoDirName
	RepoDirName = tmpDir
	t.Cleanup(func() { RepoDirName = oldRepo })

	os.MkdirAll(filepath.Join(tmpDir, "dirA"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, "dirB"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("file1"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "dirA", "file2.txt"), []byte("file2"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "dirB", "file3.txt"), []byte("file3"), 0644)

	treeHash, err := WriteTree(tmpDir)
	if err != nil {
		t.Fatalf("WriteTree() failed: %v", err)
	}

	if len(treeHash) != 40 {
		t.Errorf("expected tree hash length 40, got %d", len(treeHash))
	}
}
