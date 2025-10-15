package core

import (
	"fmt"
	"os"
	"path/filepath"
)

// var RepoDirName = ".senpai"

func InitRepo(path string, initialBranch string) error {
	repoPath := filepath.Join(path, RepoDirName)

	if _, err := os.Stat(repoPath); err == nil {
		fmt.Println("Reinitialized existing repository at", repoPath)
		return nil
	}

	if err := os.MkdirAll(repoPath, 0755); err != nil {
		return fmt.Errorf("failed to create repo folder: %w", err)
	}

	dirs := []string{
		"objects",
		"refs/heads",
	}

	for _, d := range dirs {
		fullPath := filepath.Join(repoPath, d)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("failed to create subdirectory %s: %w", d, err)
		}
	}

	headPath := filepath.Join(repoPath, "HEAD")
	headContents := fmt.Sprintf("ref: refs/heads/%s\n", initialBranch)
	if err := os.WriteFile(headPath, []byte(headContents), 0644); err != nil {
		return fmt.Errorf("failed to write HEAD file: %w", err)
	}

	return nil
}
