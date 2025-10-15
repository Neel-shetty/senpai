package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStatus_UntrackedAndIgnored(t *testing.T) {
	tmpDir := t.TempDir()

	oldRepo := RepoDirName
	RepoDirName = ".senpai"
	t.Cleanup(func() { RepoDirName = oldRepo })

	if err := InitRepo(tmpDir, "master"); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	if err := os.WriteFile(filepath.Join(tmpDir, "a.txt"), []byte("A"), 0644); err != nil {
		t.Fatalf("write a.txt: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "b.log"), []byte("B"), 0644); err != nil {
		t.Fatalf("write b.log: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, GitIgnoreFile), []byte("*.log\n"), 0644); err != nil {
		t.Fatalf("write .gitignore: %v", err)
	}

	st, err := Status(tmpDir)
	if err != nil {
		t.Fatalf("Status failed: %v", err)
	}

	var sawATxt bool
	for _, s := range st {
		if s.Path == "a.txt" {
			if s.Status != Untracked {
				t.Fatalf("a.txt expected Untracked, got %#v", s)
			}
			sawATxt = true
		}
		if s.Path == "b.log" {
			t.Fatalf("b.log should be ignored by .gitignore, but was present: %#v", s)
		}
	}
	if !sawATxt {
		t.Fatalf("did not find a.txt in status: %#v", st)
	}
}
