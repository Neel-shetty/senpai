package core

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type TreeEntry struct {
	Mode string
	Name string
	Hash string
}

func WriteTree(dir string) (string, error) {
	entries := []TreeEntry{}

	files, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	for _, f := range files {
		fullPath := filepath.Join(dir, f.Name())

		if f.IsDir() {
			subTreeHash, err := WriteTree(fullPath)
			if err != nil {
				return "", err
			}
			entries = append(entries, TreeEntry{"40000", f.Name(), subTreeHash})
		} else {
			content, err := os.ReadFile(fullPath)
			if err != nil {
				return "", fmt.Errorf("failed to read file %s: %w", fullPath, err)
			}
			blobHash, err := HashObject(content, "blob", true)
			if err != nil {
				return "", err
			}
			entries = append(entries, TreeEntry{"100644", f.Name(), blobHash})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})

	var buf bytes.Buffer
	for _, e := range entries {
		modeName := fmt.Sprintf("%s %s", e.Mode, e.Name)
		hashBytes, _ := hexToBytes(e.Hash)
		buf.Write([]byte(modeName))
		buf.WriteByte(0)
		buf.Write(hashBytes)
	}

	treeHash, err := HashObject(buf.Bytes(), "tree", true)
	if err != nil {
		return "", fmt.Errorf("failed to hash tree object %w", err)
	}
	return treeHash, nil
}

func hexToBytes(s string) ([]byte, error) {
	b := make([]byte, 20)
	_, err := fmt.Sscanf(s, "%40x", &b)
	return b, err
}
