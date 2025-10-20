package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ListBranches(repoPath string) ([]string, error) {
	refsHeadsPath := filepath.Join(repoPath, RepoDirName, "refs", "heads")

	if _, err := os.Stat(refsHeadsPath); os.IsNotExist(err) {
		return []string{}, nil
	}

	var branches []string

	err := filepath.Walk(refsHeadsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(refsHeadsPath, path)
		if err != nil {
			return err
		}

		branches = append(branches, relPath)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}

	return branches, nil
}

func CreateBranch(repoPath, branchName string) error {
	repoDir := filepath.Join(repoPath, RepoDirName)

	exists, err := BranchExists(repoPath, branchName)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("branch '%s' already exists", branchName)
	}

	headPath := filepath.Join(repoDir, "HEAD")
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return fmt.Errorf("failed to read HEAD: %w", err)
	}

	headStr := strings.TrimSpace(string(headContent))
	var commitHash string

	if strings.HasPrefix(headStr, "ref: ") {
		refPath := strings.TrimPrefix(headStr, "ref: ")
		refFullPath := filepath.Join(repoDir, refPath)

		refContent, err := os.ReadFile(refFullPath)
		if err != nil {
			return fmt.Errorf("cannot create branch: no commits yet")
		}
		commitHash = strings.TrimSpace(string(refContent))
	} else {
		commitHash = headStr
	}

	if commitHash == "" {
		return fmt.Errorf("cannot create branch: no commits yet")
	}

	branchPath := filepath.Join(repoDir, "refs", "heads", branchName)
	if err := os.MkdirAll(filepath.Dir(branchPath), 0755); err != nil {
		return fmt.Errorf("failed to create branch directory: %w", err)
	}

	if err := os.WriteFile(branchPath, []byte(commitHash+"\n"), 0644); err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	return nil
}

func DeleteBranch(repoPath, branchName string) error {
	repoDir := filepath.Join(repoPath, RepoDirName)

	exists, err := BranchExists(repoPath, branchName)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("branch '%s' not found", branchName)
	}

	currentBranch, err := GetCurrentBranch(repoPath)
	if err == nil && currentBranch == branchName {
		return fmt.Errorf("cannot delete branch '%s': currently checked out", branchName)
	}

	branchPath := filepath.Join(repoDir, "refs", "heads", branchName)
	if err := os.Remove(branchPath); err != nil {
		return fmt.Errorf("failed to delete branch: %w", err)
	}

	return nil
}

func BranchExists(repoPath, branchName string) (bool, error) {
	branchPath := filepath.Join(repoPath, RepoDirName, "refs", "heads", branchName)
	_, err := os.Stat(branchPath)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, fmt.Errorf("failed to check branch existence: %w", err)
}

func GetCurrentBranch(repoPath string) (string, error) {
	headPath := filepath.Join(repoPath, RepoDirName, "HEAD")
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return "", fmt.Errorf("failed to read HEAD: %w", err)
	}

	headStr := strings.TrimSpace(string(headContent))

	if strings.HasPrefix(headStr, "ref: ") {
		refPath := strings.TrimPrefix(headStr, "ref: ")
		if strings.HasPrefix(refPath, "refs/heads/") {
			branchName := strings.TrimPrefix(refPath, "refs/heads/")
			return branchName, nil
		}
		return "", fmt.Errorf("HEAD is not pointing to a branch")
	}

	return "", fmt.Errorf("HEAD is detached")
}

func ResolveBranchCommit(repoPath, branchName string) (string, error) {
	branchPath := filepath.Join(repoPath, RepoDirName, "refs", "heads", branchName)

	content, err := os.ReadFile(branchPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("branch '%s' not found", branchName)
		}
		return "", fmt.Errorf("failed to read branch: %w", err)
	}

	commitHash := strings.TrimSpace(string(content))
	if commitHash == "" {
		return "", fmt.Errorf("branch '%s' has no commit", branchName)
	}

	return commitHash, nil
}
