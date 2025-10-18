package core

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type TreeNode struct {
	name     string
	isTree   bool
	hash     string
	children map[string]*TreeNode
}

func Commit(repoPath, message, author, email string) (string, error) {
	repoDir := filepath.Join(repoPath, RepoDirName)
	if _, err := os.Stat(repoDir); errors.Is(err, fs.ErrNotExist) {
		return "", fmt.Errorf("repository not initialized yet")
	}

	indexPath := filepath.Join(repoDir, "index")
	indexData, err := os.ReadFile(indexPath)
	if err != nil {
		return "", fmt.Errorf("nothing to commit (staging area is empty)")
	}

	indexEntries, err := parseIndex(string(indexData))
	if err != nil {
		return "", fmt.Errorf("could not parse index: %w", err)
	}
	if len(indexEntries) == 0 {
		return "", fmt.Errorf("nothing to commit (staging area is empty)")
	}

	parentHashes, err := getParentCommit(repoDir)
	if err != nil {
		return "", err
	}

	// Merge index with parent commit tree (to include unchanged files)
	mergedEntries, err := mergeIndexWithParent(repoPath, indexEntries, parentHashes)
	if err != nil {
		return "", fmt.Errorf("failed to merge with parent: %w", err)
	}

	treeHash, err := writeTreeFromIndex(mergedEntries)
	if err != nil {
		return "", fmt.Errorf("failed to create tree: %w", err)
	}

	commitHash, err := CommitTree(treeHash, parentHashes, message, author, email)
	if err != nil {
		return "", fmt.Errorf("failed to create commit: %w", err)
	}

	if err := updateHEAD(repoDir, commitHash); err != nil {
		return "", fmt.Errorf("failed to update HEAD: %w", err)
	}

	// Update index to reflect the full commit tree (not just what was staged)
	if err := updateIndexAfterCommit(repoPath, mergedEntries); err != nil {
		return "", fmt.Errorf("failed to update index after commit: %w", err)
	}

	return commitHash, nil
}

type IndexEntry struct {
	Mode string
	Path string
	Hash string
}

func parseIndex(data string) ([]IndexEntry, error) {
	lines := strings.Split(strings.TrimSpace(data), "\n")
	var entries []IndexEntry

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 3 {
			return nil, fmt.Errorf("invalid index entry: %s", line)
		}
		entries = append(entries, IndexEntry{
			Mode: parts[0],
			Path: parts[1],
			Hash: parts[2],
		})
	}
	return entries, nil
}

func writeTreeFromIndex(entries []IndexEntry) (string, error) {
	root := &TreeNode{
		name:     "",
		isTree:   true,
		children: make(map[string]*TreeNode),
	}

	for _, entry := range entries {
		parts := strings.Split(entry.Path, string(os.PathSeparator))
		current := root

		for i, part := range parts {
			if i == len(parts)-1 {
				current.children[part] = &TreeNode{
					name:   part,
					isTree: false,
					hash:   entry.Hash,
				}
			} else {
				if _, exists := current.children[part]; !exists {
					current.children[part] = &TreeNode{
						name:     part,
						isTree:   true,
						children: make(map[string]*TreeNode),
					}
				}
				current = current.children[part]
			}
		}

	}
	return createTreeObject(root)
}

func createTreeObject(node *TreeNode) (string, error) {
	var buf bytes.Buffer

	names := make([]string, 0, len(node.children))
	for name := range node.children {
		names = append(names, name)
	}

	sort.Strings(names)

	for _, name := range names {
		child := node.children[name]
		var hash string
		var mode string

		if child.isTree {
			var err error
			hash, err = createTreeObject(child)
			if err != nil {
				return "", err
			}
			mode = "40000"
		} else {
			hash = child.hash
			mode = "100644"
		}

		modeName := fmt.Sprintf("%s %s", mode, child.name)
		hashBytes, err := hexToBytes(hash)
		if err != nil {
			return "", err
		}

		buf.Write([]byte(modeName))
		buf.WriteByte(0)
		buf.Write(hashBytes)
	}

	treeHash, err := HashObject(buf.Bytes(), "tree", true)
	if err != nil {
		return "", fmt.Errorf("failed to create tree object: %w", err)
	}

	return treeHash, nil
}

func getParentCommit(repoDir string) ([]string, error) {
	headPath := filepath.Join(repoDir, "HEAD")
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read HEAD: %w", err)
	}
	headStr := strings.TrimSpace(string(headContent))

	if strings.HasPrefix(headStr, "ref: ") {
		refPath := strings.TrimPrefix(headStr, "ref: ")
		refFullPath := filepath.Join(repoDir, refPath)

		refContent, err := os.ReadFile(refFullPath)
		if err != nil {
			return []string{}, nil
		}
		parentHash := strings.TrimSpace(string(refContent))
		return []string{parentHash}, nil
	}

	return []string{headStr}, nil
}

func updateHEAD(repoDir string, commitHash string) error {
	headPath := filepath.Join(repoDir, "HEAD")
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return fmt.Errorf("failed to read HEAD: %w", err)
	}

	headStr := strings.TrimSpace(string(headContent))

	if strings.HasPrefix(headStr, "ref: ") {
		refPath := strings.TrimPrefix(headStr, "ref: ")
		refFullPath := filepath.Join(repoDir, refPath)

		if err := os.MkdirAll(filepath.Dir(refFullPath), 0755); err != nil {
			return fmt.Errorf("failed to create ref directory: %w", err)
		}

		if err := os.WriteFile(refFullPath, []byte(commitHash+"\n"), 0644); err != nil {
			return fmt.Errorf("failed to update ref: %w", err)
		}
	} else {
		if err := os.WriteFile(headPath, []byte(commitHash+"\n"), 0644); err != nil {
			return fmt.Errorf("failed to update HEAD: %w", err)
		}
	}

	return nil
}

func updateIndexAfterCommit(repoPath string, entries []IndexEntry) error {
	indexPath := filepath.Join(repoPath, RepoDirName, "index")

	var buf bytes.Buffer
	for _, entry := range entries {
		buf.WriteString(fmt.Sprintf("%s %s %s\n", entry.Mode, entry.Path, entry.Hash))
	}

	return os.WriteFile(indexPath, buf.Bytes(), 0644)
}

func mergeIndexWithParent(repoPath string, indexEntries []IndexEntry, parentHashes []string) ([]IndexEntry, error) {
	// If no parent, just return index entries
	if len(parentHashes) == 0 || parentHashes[0] == "" {
		return indexEntries, nil
	}

	// Get the tree from parent commit - reuse the function from status.go
	parentTree, err := getLastCommitTree(repoPath)
	if err != nil {
		// If we can't read parent tree, just use index entries
		return indexEntries, nil
	}

	// Create a map of paths in the index for quick lookup
	indexMap := make(map[string]IndexEntry)
	for _, entry := range indexEntries {
		indexMap[entry.Path] = entry
	}

	// Merge: start with all index entries
	merged := make([]IndexEntry, 0, len(indexEntries)+len(parentTree))
	merged = append(merged, indexEntries...)

	// Add files from parent that aren't in the index
	for path, hash := range parentTree {
		if _, inIndex := indexMap[path]; !inIndex {
			merged = append(merged, IndexEntry{
				Mode: "100644",
				Path: path,
				Hash: hash,
			})
		}
	}

	return merged, nil
}
