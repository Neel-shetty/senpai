package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Add(repoPath string, filePath string) error {
	repoIndex := filepath.Join(repoPath, repoDirName, "index")

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read the file: %w", err)
	}

	hash, err := HashObject(content, "blob", true)

	if err != nil {
		return fmt.Errorf("failed to hash file: %w", err)
	}

	entry := fmt.Sprintf("100644 %s %s\n", filePath, hash)

	if err := os.MkdirAll(filepath.Dir(repoIndex), 0755); err != nil {
		return fmt.Errorf("failed to create index dir: %w", err)
	}

	var lines []string
	if data, err := os.ReadFile(repoIndex); err != nil {
		lines = strings.Split(strings.TrimSpace(string(data)), "\n")
	}

	found := false
	for i, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 2 && parts[1] == filePath {
			lines[i] = strings.TrimSpace(entry)
			found = true
			break
		}
	}
	if !found {
		lines = append(lines, strings.TrimSpace(entry))
	}

	return os.WriteFile(repoIndex, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}
