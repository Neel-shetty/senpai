package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type StatusType int

const (
	Unmodified StatusType = iota
	Modified
	Staged
	Untracked
)

type FileStatus struct {
	Path   string
	Status StatusType
}

func Status(repoPath string) ([]FileStatus, error) {
	indexPath := filepath.Join(repoPath, RepoDirName, "index")

	index := map[string]string{}
	if data, err := os.ReadFile(indexPath); err == nil {
		lines := strings.Split(strings.TrimSpace(string(data)), "\n")

		for _, line := range lines {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				index[parts[1]] = parts[2]
			}
		}
	}

	ig, loadErr := LoadIgnore(repoPath)
	if loadErr != nil {
		return nil, loadErr
	}

	var statuses []FileStatus
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, _ := filepath.Rel(repoPath, path)

		if info.IsDir() {
			if info.Name() == RepoDirName {
				return filepath.SkipDir
			}
			if ig != nil && ig.Ignored(relPath, true) {
				return filepath.SkipDir
			}
			return nil
		}

		if strings.HasPrefix(relPath, RepoDirName+string(os.PathSeparator)) || relPath == RepoDirName {
			return nil
		}

		if ig != nil && ig.Ignored(relPath, false) {
			return nil
		}

		// relPath, _ := filepath.Rel(repoPath, path)
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		hash, err := HashObject(content, "blob", false)
		if err != nil {
			return err
		}

		if idxHash, ok := index[relPath]; ok {
			if idxHash == hash {
				// File is in index and matches working directory - it's staged
				statuses = append(statuses, FileStatus{Path: relPath, Status: Staged})
			} else {
				// File is in index but working directory has different content
				statuses = append(statuses, FileStatus{Path: relPath, Status: Modified})
			}
		} else {
			statuses = append(statuses, FileStatus{Path: relPath, Status: Untracked})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return statuses, nil
}

func PrettyPrint(statuses []FileStatus) {
	for _, s := range statuses {
		switch s.Status {
		case Modified:
			fmt.Printf("M\t%s\n", s.Path)
		case Staged:
			fmt.Printf("A\t%s\n", s.Path)
		case Untracked:
			fmt.Printf("??\t%s\n", s.Path)
		}
	}
}
