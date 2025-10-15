package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommitTree(t *testing.T) {
	tmp := t.TempDir()
	oldRepo := RepoDirName
	RepoDirName = tmp
	defer func() { RepoDirName = oldRepo }()

	if err := InitRepo(tmp, "master"); err != nil {
		t.Fatalf("failed to init repo: %v", err)
	}

	blobHash, err := HashObject([]byte("Hello commit!"), "blob", true)
	if err != nil {
		t.Fatalf("failed to hash blob: %v", err)
	}

	hex, err := hexToBytes(blobHash)
	if err != nil {
		t.Fatalf("failed to hex bytes: %v", err)
	}
	treeData := []byte("100644 file.txt\x00" + string(hex))
	treeHash, err := HashObject(treeData, "tree", true)
	if err != nil {
		t.Fatalf("failed to hash tree: %v", err)
	}

	commitHash, err := CommitTree(treeHash, nil, "Initial commit", "Neel", "neel@neel.com")
	if err != nil {
		t.Fatalf("CommitTree failed: %v", err)
	}

	if len(commitHash) != 40 {
		t.Errorf("invalid commit hash: %s", commitHash)
	}

	objPath := filepath.Join(tmp, "objects", commitHash[:2], commitHash[2:])
	if _, err := os.Stat(objPath); err != nil {
		t.Errorf("commit object not written: %v", err)
	}

	output := captureOutput(func() {
		_ = CatFile(commitHash, true, false, true, false)
	})

	if !strings.Contains(output, "tree "+treeHash) {
		t.Errorf("commit missing tree reference")
	}
	if !strings.Contains(output, "Initial commit") {
		t.Errorf("commit missing message")
	}
}
