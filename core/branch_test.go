package core

import (
	"os"
	"path/filepath"
	"testing"
)

func createTestCommit(t *testing.T, repoPath string) string {
	t.Helper()

	testFile := filepath.Join(repoPath, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	if err := Add(repoPath, testFile); err != nil {
		t.Fatalf("failed to add file: %v", err)
	}

	commitHash, err := Commit(repoPath, "test commit", "Test Author", "test@example.com")
	if err != nil {
		t.Fatalf("failed to commit: %v", err)
	}

	return commitHash
}

func TestListBranches(t *testing.T) {
	tmpDir := t.TempDir()

	if err := InitRepo(tmpDir, "main"); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	branches, err := ListBranches(tmpDir)
	if err != nil {
		t.Fatalf("ListBranches failed: %v", err)
	}
	if len(branches) != 0 {
		t.Errorf("expected 0 branches, got %d", len(branches))
	}

	createTestCommit(t, tmpDir)

	branches, err = ListBranches(tmpDir)
	if err != nil {
		t.Fatalf("ListBranches failed: %v", err)
	}
	if len(branches) != 1 {
		t.Errorf("expected 1 branch, got %d", len(branches))
	}
	if len(branches) > 0 && branches[0] != "main" {
		t.Errorf("expected branch 'main', got '%s'", branches[0])
	}

	if err := CreateBranch(tmpDir, "feature"); err != nil {
		t.Fatalf("CreateBranch failed: %v", err)
	}

	branches, err = ListBranches(tmpDir)
	if err != nil {
		t.Fatalf("ListBranches failed: %v", err)
	}
	if len(branches) != 2 {
		t.Errorf("expected 2 branches, got %d", len(branches))
	}
}

func TestCreateBranch(t *testing.T) {
	tmpDir := t.TempDir()

	if err := InitRepo(tmpDir, "main"); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	err := CreateBranch(tmpDir, "feature")
	if err == nil {
		t.Error("expected error when creating branch without commits")
	}

	commitHash := createTestCommit(t, tmpDir)

	if err := CreateBranch(tmpDir, "feature"); err != nil {
		t.Fatalf("CreateBranch failed: %v", err)
	}

	exists, err := BranchExists(tmpDir, "feature")
	if err != nil {
		t.Fatalf("BranchExists failed: %v", err)
	}
	if !exists {
		t.Error("expected branch 'feature' to exist")
	}

	branchHash, err := ResolveBranchCommit(tmpDir, "feature")
	if err != nil {
		t.Fatalf("ResolveBranchCommit failed: %v", err)
	}
	if branchHash != commitHash {
		t.Errorf("expected branch to point to %s, got %s", commitHash, branchHash)
	}

	err = CreateBranch(tmpDir, "feature")
	if err == nil {
		t.Error("expected error when creating duplicate branch")
	}
}

func TestDeleteBranch(t *testing.T) {
	tmpDir := t.TempDir()

	if err := InitRepo(tmpDir, "main"); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	createTestCommit(t, tmpDir)

	if err := CreateBranch(tmpDir, "feature"); err != nil {
		t.Fatalf("CreateBranch failed: %v", err)
	}

	if err := DeleteBranch(tmpDir, "feature"); err != nil {
		t.Fatalf("DeleteBranch failed: %v", err)
	}

	exists, err := BranchExists(tmpDir, "feature")
	if err != nil {
		t.Fatalf("BranchExists failed: %v", err)
	}
	if exists {
		t.Error("expected branch 'feature' to not exist after deletion")
	}

	err = DeleteBranch(tmpDir, "nonexistent")
	if err == nil {
		t.Error("expected error when deleting non-existent branch")
	}

	err = DeleteBranch(tmpDir, "main")
	if err == nil {
		t.Error("expected error when deleting current branch")
	}
}

func TestBranchExists(t *testing.T) {
	tmpDir := t.TempDir()

	if err := InitRepo(tmpDir, "main"); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	exists, err := BranchExists(tmpDir, "feature")
	if err != nil {
		t.Fatalf("BranchExists failed: %v", err)
	}
	if exists {
		t.Error("expected branch 'feature' to not exist")
	}

	createTestCommit(t, tmpDir)
	if err := CreateBranch(tmpDir, "feature"); err != nil {
		t.Fatalf("CreateBranch failed: %v", err)
	}

	exists, err = BranchExists(tmpDir, "feature")
	if err != nil {
		t.Fatalf("BranchExists failed: %v", err)
	}
	if !exists {
		t.Error("expected branch 'feature' to exist")
	}
}

func TestGetCurrentBranch(t *testing.T) {
	tmpDir := t.TempDir()

	if err := InitRepo(tmpDir, "main"); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	current, err := GetCurrentBranch(tmpDir)
	if err != nil {
		t.Fatalf("GetCurrentBranch failed: %v", err)
	}
	if current != "main" {
		t.Errorf("expected current branch 'main', got '%s'", current)
	}

	createTestCommit(t, tmpDir)

	current, err = GetCurrentBranch(tmpDir)
	if err != nil {
		t.Fatalf("GetCurrentBranch failed: %v", err)
	}
	if current != "main" {
		t.Errorf("expected current branch 'main', got '%s'", current)
	}
}

func TestResolveBranchCommit(t *testing.T) {
	tmpDir := t.TempDir()

	if err := InitRepo(tmpDir, "main"); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	_, err := ResolveBranchCommit(tmpDir, "main")
	if err == nil {
		t.Error("expected error when resolving branch without commits")
	}

	commitHash := createTestCommit(t, tmpDir)

	branchHash, err := ResolveBranchCommit(tmpDir, "main")
	if err != nil {
		t.Fatalf("ResolveBranchCommit failed: %v", err)
	}
	if branchHash != commitHash {
		t.Errorf("expected branch to point to %s, got %s", commitHash, branchHash)
	}

	if err := CreateBranch(tmpDir, "feature"); err != nil {
		t.Fatalf("CreateBranch failed: %v", err)
	}

	featureHash, err := ResolveBranchCommit(tmpDir, "feature")
	if err != nil {
		t.Fatalf("ResolveBranchCommit failed: %v", err)
	}
	if featureHash != commitHash {
		t.Errorf("expected feature branch to point to %s, got %s", commitHash, featureHash)
	}

	_, err = ResolveBranchCommit(tmpDir, "nonexistent")
	if err == nil {
		t.Error("expected error when resolving non-existent branch")
	}
}

func TestMultipleBranches(t *testing.T) {
	tmpDir := t.TempDir()

	if err := InitRepo(tmpDir, "main"); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	commitHash := createTestCommit(t, tmpDir)

	branches := []string{"feature1", "feature2", "bugfix", "dev"}
	for _, branch := range branches {
		if err := CreateBranch(tmpDir, branch); err != nil {
			t.Fatalf("CreateBranch(%s) failed: %v", branch, err)
		}
	}

	allBranches, err := ListBranches(tmpDir)
	if err != nil {
		t.Fatalf("ListBranches failed: %v", err)
	}
	if len(allBranches) != 5 {
		t.Errorf("expected 5 branches, got %d", len(allBranches))
	}

	for _, branch := range append(branches, "main") {
		hash, err := ResolveBranchCommit(tmpDir, branch)
		if err != nil {
			t.Fatalf("ResolveBranchCommit(%s) failed: %v", branch, err)
		}
		if hash != commitHash {
			t.Errorf("expected branch %s to point to %s, got %s", branch, commitHash, hash)
		}
	}

	if err := DeleteBranch(tmpDir, "feature1"); err != nil {
		t.Fatalf("DeleteBranch failed: %v", err)
	}

	allBranches, err = ListBranches(tmpDir)
	if err != nil {
		t.Fatalf("ListBranches failed: %v", err)
	}
	if len(allBranches) != 4 {
		t.Errorf("expected 4 branches after deletion, got %d", len(allBranches))
	}
}
