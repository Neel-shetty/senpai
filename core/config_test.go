package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseConfig(t *testing.T) {
	tmpDir := t.TempDir()

	if err := InitRepo(tmpDir, "master"); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	cfg, err := ParseConfig(tmpDir)
	if err != nil {
		t.Fatalf("ParseConfig failed: %v", err)
	}

	if cfg.Sections == nil {
		t.Fatal("cfg.Sections is nil")
	}

	if _, ok := cfg.Sections["core"]; !ok {
		t.Error("core section not found")
	}
}

func TestGetConfig(t *testing.T) {
	tmpDir := t.TempDir()

	if err := InitRepo(tmpDir, "master"); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	value, err := GetConfig(tmpDir, "core", "repositoryformatversion")
	if err != nil {
		t.Fatalf("GetConfig failed: %v", err)
	}

	if value != "0" {
		t.Errorf("expected repositoryformatversion=0, got %s", value)
	}

	_, err = GetConfig(tmpDir, "nonexistent", "key")
	if err == nil {
		t.Error("expected error for nonexistent section")
	}

	_, err = GetConfig(tmpDir, "core", "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent key")
	}
}

func TestSetConfig(t *testing.T) {
	tmpDir := t.TempDir()

	if err := InitRepo(tmpDir, "master"); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	if err := SetConfig(tmpDir, "user", "name", "Test User"); err != nil {
		t.Fatalf("SetConfig failed: %v", err)
	}

	value, err := GetConfig(tmpDir, "user", "name")
	if err != nil {
		t.Fatalf("GetConfig after SetConfig failed: %v", err)
	}

	if value != "Test User" {
		t.Errorf("expected 'Test User', got %s", value)
	}

	if err := SetConfig(tmpDir, "user", "email", "test@example.com"); err != nil {
		t.Fatalf("SetConfig email failed: %v", err)
	}

	email, err := GetConfig(tmpDir, "user", "email")
	if err != nil {
		t.Fatalf("GetConfig email failed: %v", err)
	}

	if email != "test@example.com" {
		t.Errorf("expected 'test@example.com', got %s", email)
	}
}

func TestSetConfigUpdateExisting(t *testing.T) {
	tmpDir := t.TempDir()

	if err := InitRepo(tmpDir, "master"); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	if err := SetConfig(tmpDir, "core", "bare", "true"); err != nil {
		t.Fatalf("SetConfig failed: %v", err)
	}

	value, err := GetConfig(tmpDir, "core", "bare")
	if err != nil {
		t.Fatalf("GetConfig failed: %v", err)
	}

	if value != "true" {
		t.Errorf("expected 'true', got %s", value)
	}

	if err := SetConfig(tmpDir, "core", "bare", "false"); err != nil {
		t.Fatalf("SetConfig update failed: %v", err)
	}

	value, err = GetConfig(tmpDir, "core", "bare")
	if err != nil {
		t.Fatalf("GetConfig after update failed: %v", err)
	}

	if value != "false" {
		t.Errorf("expected 'false', got %s", value)
	}
}

func TestListConfig(t *testing.T) {
	tmpDir := t.TempDir()

	if err := InitRepo(tmpDir, "master"); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	if err := SetConfig(tmpDir, "user", "name", "Test User"); err != nil {
		t.Fatalf("SetConfig failed: %v", err)
	}

	if err := ListConfig(tmpDir); err != nil {
		t.Fatalf("ListConfig failed: %v", err)
	}
}

func TestCreateDefaultConfig(t *testing.T) {
	tmpDir := t.TempDir()

	repoPath := filepath.Join(tmpDir, RepoDirName)
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("failed to create repo dir: %v", err)
	}

	if err := CreateDefaultConfig(tmpDir, false); err != nil {
		t.Fatalf("CreateDefaultConfig failed: %v", err)
	}

	configPath := filepath.Join(repoPath, "config")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("config file was not created")
	}

	cfg, err := ParseConfig(tmpDir)
	if err != nil {
		t.Fatalf("ParseConfig failed: %v", err)
	}

	coreSection, ok := cfg.Sections["core"]
	if !ok {
		t.Fatal("core section not found")
	}

	if coreSection["repositoryformatversion"] != "0" {
		t.Errorf("expected repositoryformatversion=0, got %s", coreSection["repositoryformatversion"])
	}

	if coreSection["bare"] != "false" {
		t.Errorf("expected bare=false, got %s", coreSection["bare"])
	}
}

func TestCreateDefaultConfigBare(t *testing.T) {
	tmpDir := t.TempDir()

	repoPath := filepath.Join(tmpDir, RepoDirName)
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("failed to create repo dir: %v", err)
	}

	if err := CreateDefaultConfig(tmpDir, true); err != nil {
		t.Fatalf("CreateDefaultConfig failed: %v", err)
	}

	cfg, err := ParseConfig(tmpDir)
	if err != nil {
		t.Fatalf("ParseConfig failed: %v", err)
	}

	if cfg.Sections["core"]["bare"] != "true" {
		t.Errorf("expected bare=true, got %s", cfg.Sections["core"]["bare"])
	}
}

func TestParseConfigWithComments(t *testing.T) {
	tmpDir := t.TempDir()

	repoPath := filepath.Join(tmpDir, RepoDirName)
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("failed to create repo dir: %v", err)
	}

	configContent := `# This is a comment
[core]
	repositoryformatversion = 0
	# Another comment
	filemode = true
; Semicolon comment
	bare = false

[user]
	name = Test User
	email = test@example.com
`
	configPath := filepath.Join(repoPath, "config")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := ParseConfig(tmpDir)
	if err != nil {
		t.Fatalf("ParseConfig failed: %v", err)
	}

	if len(cfg.Sections) != 2 {
		t.Errorf("expected 2 sections, got %d", len(cfg.Sections))
	}

	if cfg.Sections["user"]["name"] != "Test User" {
		t.Errorf("expected 'Test User', got %s", cfg.Sections["user"]["name"])
	}

	if cfg.Sections["user"]["email"] != "test@example.com" {
		t.Errorf("expected 'test@example.com', got %s", cfg.Sections["user"]["email"])
	}
}

func TestSetConfigNewSection(t *testing.T) {
	tmpDir := t.TempDir()

	if err := InitRepo(tmpDir, "master"); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	if err := SetConfig(tmpDir, "remote \"origin\"", "url", "https://github.com/test/repo.git"); err != nil {
		t.Fatalf("SetConfig failed: %v", err)
	}

	value, err := GetConfig(tmpDir, "remote \"origin\"", "url")
	if err != nil {
		t.Fatalf("GetConfig failed: %v", err)
	}

	if value != "https://github.com/test/repo.git" {
		t.Errorf("expected 'https://github.com/test/repo.git', got %s", value)
	}
}
