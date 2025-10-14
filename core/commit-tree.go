package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func CommitTree(treeHash string, parentHashes []string, message string, author string, email string) (string, error) {
	repoPath := filepath.Join(repoDirName, "objects")
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return "", fmt.Errorf("repository not initialized")
	}

	timestamp := time.Now().Unix()
	timezone, _ := time.Now().Zone()

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("tree %s\n", treeHash))
	for _, parent := range parentHashes {
		sb.WriteString(fmt.Sprintf("parent %s\n", parent))
	}
	sb.WriteString(fmt.Sprintf("author %s <%s> %d %s\n", author, email, timestamp, timezone))
	sb.WriteString(fmt.Sprintf("committer %s <%s> %d %s\n", author, email, timestamp, timezone))
	sb.WriteString("\n" + strings.TrimSpace(message) + "\n")

	content := []byte(sb.String())

	hash, err := HashObject(content, "commit", true)
	if err != nil {
		return "", fmt.Errorf("failed to write commit object: %w", err)
	}
	return hash, nil
}
