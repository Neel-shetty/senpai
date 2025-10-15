package core

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

type ignoreRule struct {
	pattern      string
	negate       bool
	dirOnly      bool
	rootAnchored bool
}

type IgnoreMatcher struct {
	rules []ignoreRule
}

func LoadIgnore(repoPath string) (*IgnoreMatcher, error) {
	fp := filepath.Join(repoPath, GitIgnoreFile)
	data, err := os.ReadFile(fp)
	if err != nil {
		return &IgnoreMatcher{rules: nil}, nil
	}

	var rules []ignoreRule
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		r := ignoreRule{}
		if strings.HasPrefix(line, "!") {
			r.negate = true
			line = line[1:]
		}
		if strings.HasSuffix(line, "/") {
			r.dirOnly = true
			line = strings.TrimSuffix(line, "/")
		}
		if strings.HasPrefix(line, "/") {
			r.rootAnchored = true
			line = strings.TrimPrefix(line, "/")
		}
		r.pattern = filepath.ToSlash(line)
		if r.pattern != "" {
			rules = append(rules, r)
		}
	}
	return &IgnoreMatcher{rules: rules}, nil
}

func (m *IgnoreMatcher) Ignored(relPath string, isDir bool) bool {
	rel := filepath.ToSlash(relPath)
	if rel == "." {
		return false
	}
	var matched bool
	for _, r := range m.rules {
		if r.dirOnly && !isDir {
			continue
		}
		if ruleMatches(r, rel) {
			matched = !r.negate
		}
	}
	return matched
}

func ruleMatches(r ignoreRule, rel string) bool {
	if r.rootAnchored {
		ok, _ := path.Match(r.pattern, rel)
		return ok
	}

	parts := strings.Split(rel, "/")
	for i := range parts {
		candidate := strings.Join(parts[i:], "/")
		ok, _ := path.Match(r.pattern, candidate)
		if ok {
			return true
		}
	}
	return false
}
