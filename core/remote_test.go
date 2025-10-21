package core

import (
	"testing"
)

func TestListRemotesEmpty(t *testing.T) {
	tmpDir := t.TempDir()

	if err := InitRepo(tmpDir, "master"); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	remotes, err := ListRemotes(tmpDir)
	if err != nil {
		t.Fatalf("ListRemotes failed: %v", err)
	}

	if len(remotes) != 0 {
		t.Errorf("expected 0 remotes, got %d", len(remotes))
	}
}

func TestAddRemote(t *testing.T) {
	tmpDir := t.TempDir()

	if err := InitRepo(tmpDir, "master"); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	name := "origin"
	url := "https://github.com/user/repo.git"

	if err := AddRemote(tmpDir, name, url); err != nil {
		t.Fatalf("AddRemote failed: %v", err)
	}

	remotes, err := ListRemotes(tmpDir)
	if err != nil {
		t.Fatalf("ListRemotes failed: %v", err)
	}

	if len(remotes) != 1 {
		t.Fatalf("expected 1 remote, got %d", len(remotes))
	}

	if remotes[0].Name != name {
		t.Errorf("expected name '%s', got '%s'", name, remotes[0].Name)
	}

	if remotes[0].URL != url {
		t.Errorf("expected URL '%s', got '%s'", url, remotes[0].URL)
	}
}

func TestAddRemoteDuplicate(t *testing.T) {
	tmpDir := t.TempDir()

	if err := InitRepo(tmpDir, "master"); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	name := "origin"
	url := "https://github.com/user/repo.git"

	if err := AddRemote(tmpDir, name, url); err != nil {
		t.Fatalf("AddRemote failed: %v", err)
	}

	err := AddRemote(tmpDir, name, url)
	if err == nil {
		t.Error("expected error when adding duplicate remote")
	}
}

func TestAddMultipleRemotes(t *testing.T) {
	tmpDir := t.TempDir()

	if err := InitRepo(tmpDir, "master"); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	if err := AddRemote(tmpDir, "origin", "https://github.com/user/repo.git"); err != nil {
		t.Fatalf("AddRemote origin failed: %v", err)
	}

	if err := AddRemote(tmpDir, "upstream", "https://github.com/upstream/repo.git"); err != nil {
		t.Fatalf("AddRemote upstream failed: %v", err)
	}

	remotes, err := ListRemotes(tmpDir)
	if err != nil {
		t.Fatalf("ListRemotes failed: %v", err)
	}

	if len(remotes) != 2 {
		t.Fatalf("expected 2 remotes, got %d", len(remotes))
	}
}

func TestGetRemoteURL(t *testing.T) {
	tmpDir := t.TempDir()

	if err := InitRepo(tmpDir, "master"); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	name := "origin"
	expectedURL := "https://github.com/user/repo.git"

	if err := AddRemote(tmpDir, name, expectedURL); err != nil {
		t.Fatalf("AddRemote failed: %v", err)
	}

	url, err := GetRemoteURL(tmpDir, name)
	if err != nil {
		t.Fatalf("GetRemoteURL failed: %v", err)
	}

	if url != expectedURL {
		t.Errorf("expected URL '%s', got '%s'", expectedURL, url)
	}
}

func TestGetRemoteURLNonexistent(t *testing.T) {
	tmpDir := t.TempDir()

	if err := InitRepo(tmpDir, "master"); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	_, err := GetRemoteURL(tmpDir, "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent remote")
	}
}

func TestSetRemoteURL(t *testing.T) {
	tmpDir := t.TempDir()

	if err := InitRepo(tmpDir, "master"); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	name := "origin"
	oldURL := "https://github.com/user/repo.git"
	newURL := "git@github.com:user/repo.git"

	if err := AddRemote(tmpDir, name, oldURL); err != nil {
		t.Fatalf("AddRemote failed: %v", err)
	}

	if err := SetRemoteURL(tmpDir, name, newURL); err != nil {
		t.Fatalf("SetRemoteURL failed: %v", err)
	}

	url, err := GetRemoteURL(tmpDir, name)
	if err != nil {
		t.Fatalf("GetRemoteURL failed: %v", err)
	}

	if url != newURL {
		t.Errorf("expected URL '%s', got '%s'", newURL, url)
	}
}

func TestSetRemoteURLNonexistent(t *testing.T) {
	tmpDir := t.TempDir()

	if err := InitRepo(tmpDir, "master"); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	err := SetRemoteURL(tmpDir, "nonexistent", "https://github.com/user/repo.git")
	if err == nil {
		t.Error("expected error when setting URL for nonexistent remote")
	}
}

func TestRemoveRemote(t *testing.T) {
	tmpDir := t.TempDir()

	if err := InitRepo(tmpDir, "master"); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	name := "origin"
	url := "https://github.com/user/repo.git"

	if err := AddRemote(tmpDir, name, url); err != nil {
		t.Fatalf("AddRemote failed: %v", err)
	}

	remotes, err := ListRemotes(tmpDir)
	if err != nil {
		t.Fatalf("ListRemotes failed: %v", err)
	}

	if len(remotes) != 1 {
		t.Fatalf("expected 1 remote before removal, got %d", len(remotes))
	}

	if err := RemoveRemote(tmpDir, name); err != nil {
		t.Fatalf("RemoveRemote failed: %v", err)
	}

	remotes, err = ListRemotes(tmpDir)
	if err != nil {
		t.Fatalf("ListRemotes after removal failed: %v", err)
	}

	if len(remotes) != 0 {
		t.Errorf("expected 0 remotes after removal, got %d", len(remotes))
	}
}

func TestRemoveRemoteNonexistent(t *testing.T) {
	tmpDir := t.TempDir()

	if err := InitRepo(tmpDir, "master"); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	err := RemoveRemote(tmpDir, "nonexistent")
	if err == nil {
		t.Error("expected error when removing nonexistent remote")
	}
}

func TestRemoveOneOfMultipleRemotes(t *testing.T) {
	tmpDir := t.TempDir()

	if err := InitRepo(tmpDir, "master"); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	if err := AddRemote(tmpDir, "origin", "https://github.com/user/repo.git"); err != nil {
		t.Fatalf("AddRemote origin failed: %v", err)
	}

	if err := AddRemote(tmpDir, "upstream", "https://github.com/upstream/repo.git"); err != nil {
		t.Fatalf("AddRemote upstream failed: %v", err)
	}

	if err := RemoveRemote(tmpDir, "origin"); err != nil {
		t.Fatalf("RemoveRemote failed: %v", err)
	}

	remotes, err := ListRemotes(tmpDir)
	if err != nil {
		t.Fatalf("ListRemotes failed: %v", err)
	}

	if len(remotes) != 1 {
		t.Fatalf("expected 1 remote, got %d", len(remotes))
	}

	if remotes[0].Name != "upstream" {
		t.Errorf("expected remaining remote to be 'upstream', got '%s'", remotes[0].Name)
	}
}

func TestRemoteWithSSHURL(t *testing.T) {
	tmpDir := t.TempDir()

	if err := InitRepo(tmpDir, "master"); err != nil {
		t.Fatalf("InitRepo failed: %v", err)
	}

	name := "origin"
	sshURL := "git@github.com:user/repo.git"

	if err := AddRemote(tmpDir, name, sshURL); err != nil {
		t.Fatalf("AddRemote failed: %v", err)
	}

	url, err := GetRemoteURL(tmpDir, name)
	if err != nil {
		t.Fatalf("GetRemoteURL failed: %v", err)
	}

	if url != sshURL {
		t.Errorf("expected URL '%s', got '%s'", sshURL, url)
	}
}
