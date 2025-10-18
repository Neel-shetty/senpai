package core

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
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

	lastCommitTree, err := getLastCommitTree(repoPath)
	if err != nil {
		lastCommitTree = map[string]string{}
	}

	ig, loadErr := LoadIgnore(repoPath)
	if loadErr != nil {
		return nil, loadErr
	}

	var statuses []FileStatus
	err = filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
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

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		hash, err := HashObject(content, "blob", false)
		if err != nil {
			return err
		}

		idxHash, inIndex := index[relPath]
		commitHash, inCommit := lastCommitTree[relPath]

		if inIndex {
			if inCommit {
				indexMatchesCommit := idxHash == commitHash
				workingTreeMatchesIndex := idxHash == hash

				if indexMatchesCommit && workingTreeMatchesIndex {
					// skip unmodified files
				} else if workingTreeMatchesIndex {
					statuses = append(statuses, FileStatus{Path: relPath, Status: Staged})
				} else {
					statuses = append(statuses, FileStatus{Path: relPath, Status: Modified})
				}
			} else {
				if idxHash == hash {
					statuses = append(statuses, FileStatus{Path: relPath, Status: Staged})
				} else {
					statuses = append(statuses, FileStatus{Path: relPath, Status: Modified})
				}
			}
		} else if inCommit {
			if commitHash != hash {
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

func getLastCommitTree(repoPath string) (map[string]string, error) {
	repoDir := filepath.Join(repoPath, RepoDirName)

	headPath := filepath.Join(repoDir, "HEAD")
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return nil, err
	}

	headStr := strings.TrimSpace(string(headContent))
	var commitHash string

	if strings.HasPrefix(headStr, "ref: ") {
		refPath := strings.TrimPrefix(headStr, "ref: ")
		refFullPath := filepath.Join(repoDir, refPath)

		refContent, err := os.ReadFile(refFullPath)
		if err != nil {
			return nil, fmt.Errorf("no commits yet")
		}
		commitHash = strings.TrimSpace(string(refContent))
	} else {
		commitHash = headStr
	}

	commitContent, err := readObject(repoPath, commitHash)
	if err != nil {
		return nil, err
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

func readObject(repoPath string, hash string) ([]byte, error) {
	objectPath := filepath.Join(repoPath, RepoDirName, "objects", hash[:2], hash[2:])

	file, err := os.Open(objectPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r, err := zlib.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	sepIndex := bytes.IndexByte(data, 0)
	if sepIndex < 0 {
		return nil, fmt.Errorf("invalid object format")
	}

	return data[sepIndex+1:], nil
}

func readTreeRecursive(repoPath string, treeHash string, prefix string) (map[string]string, error) {
	treeContent, err := readObject(repoPath, treeHash)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	i := 0

	for i < len(treeContent) {
		spaceIdx := bytes.IndexByte(treeContent[i:], ' ')
		if spaceIdx < 0 {
			break
		}
		mode := string(treeContent[i : i+spaceIdx])
		i += spaceIdx + 1

		nullIdx := bytes.IndexByte(treeContent[i:], 0)
		if nullIdx < 0 {
			break
		}
		name := string(treeContent[i : i+nullIdx])
		i += nullIdx + 1

		// Read hash (size depends on hash algorithm, currently SHA-1: 20 bytes)
		if i+HashSize > len(treeContent) {
			break
		}
		hashBytes := treeContent[i : i+HashSize]
		hash := fmt.Sprintf("%x", hashBytes)
		i += HashSize

		var fullPath string
		if prefix == "" {
			fullPath = name
		} else {
			fullPath = filepath.Join(prefix, name)
		}

		if mode == "40000" {
			subTree, err := readTreeRecursive(repoPath, hash, fullPath)
			if err != nil {
				return nil, err
			}
			for path, h := range subTree {
				result[path] = h
			}
		} else {
			result[fullPath] = hash
		}
	}

	return result, nil
}
