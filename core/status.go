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
	indexPath := filepath.Join(repoPath, repoDirName, "index")

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

	var statuses []FileStatus
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		relPath, _ := filepath.Rel(repoPath, path)
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
				statuses = append(statuses, FileStatus{Path: relPath, Status: Unmodified})
			} else {
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
