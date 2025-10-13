package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitRepoCreateDirs(t *testing.T) {
	tmpDir := t.TempDir()

	err := InitRepo(tmpDir, "master")
	repoPath := filepath.Join(tmpDir, ".senpai")
	if err != nil {
		t.Fatalf("InitRepo Failed: %v", err)
	}

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		t.Fatalf(".senpai directory not created")
	}

	expectedDirs := []string{
		"objects",
		"refs/heads",
	}

	for _, d := range expectedDirs {
		path := filepath.Join(repoPath, d)
		if fileInfo, err := os.Stat(path); err != nil || !fileInfo.IsDir() {
			t.Errorf("expected directory %s missing", d)
		}
	}

	headPath := filepath.Join(repoPath, "HEAD")
	data, _ := os.ReadFile(headPath)
	if string(data) != "ref: refs/heads/master\n" {
		t.Errorf("HEAD file incorrect: %s", data)
	}

}

func TestInitRepoTwice(t *testing.T) {
	tmpDir := t.TempDir()

	if err := InitRepo(tmpDir, "master"); err != nil {
		t.Fatalf("first InitRepo failed: %v", err)
	}

	repoPath := filepath.Join(tmpDir, ".senpai")
	headPath := filepath.Join(repoPath, "HEAD")

	data, _ := os.ReadFile(headPath)
	if string(data) != "ref: refs/heads/master\n" {
		t.Errorf("HEAD file incorrect after first init: %s", data)
	}

	if err := InitRepo(tmpDir, "master"); err != nil {
		t.Fatalf("second InitRepo failed: %v", err)
	}

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		t.Fatalf(".senpai directory missing after second init")
	}

	data2, _ := os.ReadFile(headPath)
	if string(data2) != "ref: refs/heads/master\n" {
		t.Errorf("HEAD file incorrect after second init: %s", data2)
	}

	expectedDirs := []string{"objects", "refs/heads"}
	for _, d := range expectedDirs {
		path := filepath.Join(repoPath, d)
		if fi, err := os.Stat(path); err != nil || !fi.IsDir() {
			t.Errorf("expected directory %s missing after second init", d)
		}
	}
}
