package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Checkout(repoPath, target string) error {
	repoDir := filepath.Join(repoPath, RepoDirName)

	branchExists, err := BranchExists(repoPath, target)
	if err != nil {
		return err
	}

	var commitHash string

	headPath := filepath.Join(repoDir, "HEAD")
	if branchExists {
		commitHash, err = ResolveBranchCommit(repoPath, target)
		if err != nil {
			return fmt.Errorf("failed to resolve branch: %w", err)
		}

		refPath := fmt.Sprintf("ref: refs/heads/%s", target)
		if err := os.WriteFile(headPath, []byte(refPath+"\n"), 0644); err != nil {
			return fmt.Errorf("failed to update HEAD: %w", err)
		}
	} else {
		objectPath := filepath.Join(repoDir, "objects", target[:2], target[2:])
		if _, err := os.Stat(objectPath); os.IsNotExist(err) {
			return fmt.Errorf("reference '%s' not found (not a branch or commit)", target)
		}

		commitHash = target
		if err := os.WriteFile(headPath, []byte(commitHash+"\n"), 0644); err != nil {
			return fmt.Errorf("failed to update HEAD: %w", err)
		}
	}

	if err := updateWorkingTreeToCommit(repoPath, commitHash); err != nil {
		return fmt.Errorf("failed to update working tree: %w", err)
	}
	return nil
}

func CheckoutNewBranch(repoPath, branchName string) error {
	if err := CreateBranch(repoPath, branchName); err != nil {
		return err
	}

	return Checkout(repoPath, branchName)
}

func updateWorkingTreeToCommit(repoPath, commitHash string) error {
	treeFiles, err := getCommitTree(repoPath, commitHash)
	if err != nil {
		return fmt.Errorf("failed to get commit tree: %w", err)
	}

	if err := clearWorkingDirectory(repoPath); err != nil {
		return fmt.Errorf("failed to cleear working directory: %w", err)
	}

	for filepath, fileHash := range treeFiles {
		if err := restoreFile(repoPath, filepath, fileHash); err != nil {
			return fmt.Errorf("failed to resotre file %s: %w", filepath, err)
		}
	}

	if err := updateIndexToTree(repoPath, treeFiles); err != nil {
		return fmt.Errorf("failed to update index: %w", err)
	}
	return nil
}

func getCommitTree(repoPath, commitHash string) (map[string]string, error) {
	commitContent, err := readObject(repoPath, commitHash)
	if err != nil {
		return nil, fmt.Errorf("failed to read commit object: %w", err)
	}

	lines := strings.Split(string(commitContent), "\n")
	var treeHash string
	for _, line := range lines {
		if strings.HasPrefix(line, "tree ") {
			treeHash = strings.TrimSpace(strings.TrimPrefix(line, "tree "))
			break
		}
	}

	if treeHash == "" {
		return nil, fmt.Errorf("invalid commit: no tree found")
	}

	return readTreeRecursive(repoPath, treeHash, "")
}

func clearWorkingDirectory(repoPath string) error {
	entries, err := os.ReadDir(repoPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.Name() == RepoDirName {
			continue
		}

		fullPath := filepath.Join(repoPath, entry.Name())
		if err := os.RemoveAll(fullPath); err != nil {
			return fmt.Errorf("failed to remove %s: %w", fullPath, err)
		}
	}

	return nil
}

func restoreFile(repoPath, filePath, blobHash string) error {
	blobContent, err := readObject(repoPath, blobHash)
	if err != nil {
		return fmt.Errorf("failed to read blob: %w", err)
	}

	fullPath := filepath.Join(repoPath, filePath)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(fullPath, blobContent, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func updateIndexToTree(repoPath string, treeFiles map[string]string) error {
	indexPath := filepath.Join(repoPath, RepoDirName, "index")

	var content strings.Builder
	for path, hash := range treeFiles {
		content.WriteString(fmt.Sprintf("100644 %s %s\n", path, hash))
	}

	return os.WriteFile(indexPath, []byte(content.String()), 0644)
}
