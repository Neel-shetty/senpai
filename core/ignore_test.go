package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIgnoreMatcher_SimplePatterns(t *testing.T) {
	m := &IgnoreMatcher{rules: []ignoreRule{
		{pattern: "*.log"},
	}}

	if !m.Ignored("app.log", false) {
		t.Fatalf("expected app.log to be ignored")
	}
	if m.Ignored("readme.md", false) {
		t.Fatalf("expected readme.md not to be ignored")
	}
	if m.Ignored(".", true) {
		t.Fatalf("root should never be ignored")
	}
}

func TestIgnoreMatcher_Negation(t *testing.T) {
	m := &IgnoreMatcher{rules: []ignoreRule{
		{pattern: "*.log"},
		{pattern: "keep.log", negate: true},
	}}

	if !m.Ignored("app.log", false) {
		t.Fatalf("expected app.log to be ignored by *.log rule")
	}
	if m.Ignored("keep.log", false) {
		t.Fatalf("expected keep.log to be un-ignored by negation rule")
	}
}

func TestIgnoreMatcher_RootAnchored(t *testing.T) {
	m := &IgnoreMatcher{rules: []ignoreRule{
		{pattern: "root.txt", rootAnchored: true},
	}}

	if !m.Ignored("root.txt", false) {
		t.Fatalf("expected root.txt at repo root to be ignored")
	}
	if m.Ignored(filepath.ToSlash(filepath.Join("sub", "root.txt")), false) {
		t.Fatalf("expected sub/root.txt not to match root-anchored pattern")
	}
}

func TestIgnoreMatcher_DirOnly(t *testing.T) {
	m := &IgnoreMatcher{rules: []ignoreRule{
		{pattern: "docs", dirOnly: true},
	}}

	if !m.Ignored("docs", true) {
		t.Fatalf("expected docs directory to be ignored")
	}
	if m.Ignored("docs", false) {
		t.Fatalf("dir-only rule should not ignore a file named docs")
	}
}

func TestLoadIgnore_ParsesFileAndMatches(t *testing.T) {
	repo := t.TempDir()
	// Create a .gitignore file in the repo root
	gi := []byte("" +
		"# comment should be ignored\n" +
		"/root.txt\n" +
		"logs/*.log\n" +
		"!logs/keep.log\n" +
		"docs/\n" +
		"build\n",
	)
	if err := os.WriteFile(filepath.Join(repo, GitIgnoreFile), gi, 0o644); err != nil {
		t.Fatalf("failed writing .gitignore: %v", err)
	}

	m, err := LoadIgnore(repo)
	if err != nil {
		t.Fatalf("LoadIgnore failed: %v", err)
	}

	// Root anchored
	if !m.Ignored("root.txt", false) {
		t.Errorf("expected root.txt to be ignored by /root.txt")
	}
	if m.Ignored("sub/root.txt", false) {
		t.Errorf("expected sub/root.txt not ignored by /root.txt")
	}

	// Glob and negation
	if !m.Ignored("logs/app.log", false) {
		t.Errorf("expected logs/app.log to be ignored by logs/*.log")
	}
	if m.Ignored("logs/keep.log", false) {
		t.Errorf("expected logs/keep.log to be un-ignored by negation rule")
	}

	// Dir-only rule only applies to the directory itself
	if !m.Ignored("docs", true) {
		t.Errorf("expected docs directory to be ignored by docs/")
	}
	if m.Ignored("docs/file.txt", false) {
		t.Errorf("dir-only rule should not directly ignore files under docs when matching a file path")
	}

	// Non-dir-only directory name should still ignore the directory itself
	if !m.Ignored("build", true) {
		t.Errorf("expected build directory to be ignored by build rule")
	}
}
