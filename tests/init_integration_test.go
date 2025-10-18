package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func getBinaryPath(t *testing.T) string {
	t.Helper()

	projectRoot, err := filepath.Abs(filepath.Join("..", "."))
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	binaryPath := filepath.Join(projectRoot, "senpai")

	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	cmd.Dir = projectRoot
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\nOutput: %s", err, output)
	}

	return binaryPath
}

func runSenpaiCommand(t *testing.T, binaryPath string, dir string, args ...string) (string, string, error) {
	t.Helper()

	cmd := exec.Command(binaryPath, args...)
	cmd.Dir = dir

	stdout, err := cmd.Output()
	stderr := ""
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr = string(exitErr.Stderr)
		}
	}

	return string(stdout), stderr, err
}

func TestInitCommand(t *testing.T) {
	binaryPath := getBinaryPath(t)

	tmpDir, err := os.MkdirTemp("", "senpai-init-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	stdout, stderr, err := runSenpaiCommand(t, binaryPath, tmpDir, "init")
	if err != nil {
		t.Fatalf("senpai init failed: %v\nStdout: %s\nStderr: %s", err, stdout, stderr)
	}

	expectedOutput := "Initialized empty repository on branch master"
	if !strings.Contains(stdout, expectedOutput) {
		t.Errorf("Expected output to contain %q, got: %s", expectedOutput, stdout)
	}

	senpaiDir := filepath.Join(tmpDir, ".senpai")
	if _, err := os.Stat(senpaiDir); os.IsNotExist(err) {
		t.Fatalf(".senpai directory was not created")
	}

	requiredDirs := []string{
		"objects",
		"refs/heads",
	}

	for _, dir := range requiredDirs {
		dirPath := filepath.Join(senpaiDir, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Errorf("Required directory %q was not created", dir)
		}
	}

	headPath := filepath.Join(senpaiDir, "HEAD")
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		t.Fatalf("Failed to read HEAD file: %v", err)
	}

	expectedHead := "ref: refs/heads/master\n"
	if string(headContent) != expectedHead {
		t.Errorf("HEAD content incorrect. Expected %q, got %q", expectedHead, string(headContent))
	}
}

func TestInitCommandWithCustomBranch(t *testing.T) {
	binaryPath := getBinaryPath(t)

	tmpDir, err := os.MkdirTemp("", "senpai-init-branch-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	customBranch := "main"

	stdout, stderr, err := runSenpaiCommand(t, binaryPath, tmpDir, "init", "--initial-branch", customBranch)
	if err != nil {
		t.Fatalf("senpai init failed: %v\nStdout: %s\nStderr: %s", err, stdout, stderr)
	}

	expectedOutput := "Initialized empty repository on branch " + customBranch
	if !strings.Contains(stdout, expectedOutput) {
		t.Errorf("Expected output to contain %q, got: %s", expectedOutput, stdout)
	}

	headPath := filepath.Join(tmpDir, ".senpai", "HEAD")
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		t.Fatalf("Failed to read HEAD file: %v", err)
	}

	expectedHead := "ref: refs/heads/" + customBranch + "\n"
	if string(headContent) != expectedHead {
		t.Errorf("HEAD content incorrect. Expected %q, got %q", expectedHead, string(headContent))
	}
}

func TestInitCommandReinitialize(t *testing.T) {
	binaryPath := getBinaryPath(t)

	tmpDir, err := os.MkdirTemp("", "senpai-reinit-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	stdout1, stderr1, err := runSenpaiCommand(t, binaryPath, tmpDir, "init")
	if err != nil {
		t.Fatalf("First senpai init failed: %v\nStdout: %s\nStderr: %s", err, stdout1, stderr1)
	}

	testFilePath := filepath.Join(tmpDir, ".senpai", "test-marker")
	if err := os.WriteFile(testFilePath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test marker file: %v", err)
	}

	stdout2, stderr2, err := runSenpaiCommand(t, binaryPath, tmpDir, "init")
	if err != nil {
		t.Fatalf("Second senpai init failed: %v\nStdout: %s\nStderr: %s", err, stdout2, stderr2)
	}

	expectedOutput := "Reinitialized existing repository"
	if !strings.Contains(stdout2, expectedOutput) {
		t.Errorf("Expected output to contain %q, got: %s", expectedOutput, stdout2)
	}

	if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
		t.Error("Reinitialization destroyed existing repository content")
	}
}

func TestInitAddStatusWorkflow(t *testing.T) {
	binaryPath := getBinaryPath(t)

	tmpDir, err := os.MkdirTemp("", "senpai-workflow-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize repository
	stdout, stderr, err := runSenpaiCommand(t, binaryPath, tmpDir, "init")
	if err != nil {
		t.Fatalf("senpai init failed: %v\nStdout: %s\nStderr: %s", err, stdout, stderr)
	}

	expectedOutput := "Initialized empty repository on branch master"
	if !strings.Contains(stdout, expectedOutput) {
		t.Errorf("Expected init output to contain %q, got: %s", expectedOutput, stdout)
	}

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := "Hello, Senpai!"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Check status before adding (should show untracked file)
	stdout, stderr, err = runSenpaiCommand(t, binaryPath, tmpDir, "status")
	if err != nil {
		t.Fatalf("senpai status failed: %v\nStdout: %s\nStderr: %s", err, stdout, stderr)
	}

	if !strings.Contains(stdout, "test.txt") {
		t.Errorf("Expected status to show untracked file 'test.txt', got: %s", stdout)
	}
	if !strings.Contains(stdout, "??") {
		t.Errorf("Expected status to show '??' for untracked file, got: %s", stdout)
	}

	// Add the file
	stdout, stderr, err = runSenpaiCommand(t, binaryPath, tmpDir, "add", "test.txt")
	if err != nil {
		t.Fatalf("senpai add failed: %v\nStdout: %s\nStderr: %s", err, stdout, stderr)
	}

	// Check status after adding (should show staged file)
	stdout, stderr, err = runSenpaiCommand(t, binaryPath, tmpDir, "status")
	if err != nil {
		t.Fatalf("senpai status after add failed: %v\nStdout: %s\nStderr: %s", err, stdout, stderr)
	}

	if !strings.Contains(stdout, "test.txt") {
		t.Errorf("Expected status to show staged file 'test.txt', got: %s", stdout)
	}
	// Status should now show the file as staged (A for added) or unmodified, not untracked (??)
	if strings.Contains(stdout, "??") {
		t.Errorf("File should not be untracked after adding, got: %s", stdout)
	}

	// Verify index file was created
	indexPath := filepath.Join(tmpDir, ".senpai", "index")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Error("Index file was not created after adding file")
	}
}
