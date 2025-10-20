package core

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type CommitInfo struct {
	Hash      string
	Tree      string
	Parents   []string
	Author    string
	Email     string
	Timestamp int64
	Timezone  string
	Committer string
	Message   string
}

func Log(repoPath string) ([]CommitInfo, error) {
	repoDir := filepath.Join(repoPath, RepoDirName)
	if _, err := os.Stat(repoDir); errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("repository not initialized")
	}

	currentHash, err := getCurrentCommitHash(repoDir)
	if err != nil {
		return nil, err
	}

	if currentHash == "" {
		return nil, fmt.Errorf("no commits yet")
	}

	var commits []CommitInfo
	visited := make(map[string]bool)

	err = walkCommits(currentHash, &commits, visited)
	if err != nil {
		return nil, err
	}
	return commits, nil
}

func getCurrentCommitHash(repoDir string) (string, error) {
	headPath := filepath.Join(repoDir, "HEAD")
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return "", fmt.Errorf("failed to read HEAD: %w", err)
	}
	headStr := strings.TrimSpace(string(headContent))

	if strings.HasPrefix(headStr, "ref: ") {
		refPath := strings.TrimPrefix(headStr, "ref: ")
		refFullPath := filepath.Join(repoDir, refPath)

		refContent, err := os.ReadFile(refFullPath)

		if err != nil {
			return "", fmt.Errorf("failed to read ref: %w", err)
		}

		return strings.TrimSpace(string(refContent)), nil
	}

	return headStr, nil
}

func walkCommits(hash string, commits *[]CommitInfo, visited map[string]bool) error {
	if hash == "" || visited[hash] {
		return nil
	}

	visited[hash] = true

	commitInfo, err := readCommitObject(hash)
	if err != nil {
		return err
	}

	commitInfo.Hash = hash
	*commits = append(*commits, commitInfo)

	for _, parent := range commitInfo.Parents {
		err := walkCommits(parent, commits, visited)
		if err != nil {
			return err
		}
	}
	return err
}

func readCommitObject(hash string) (CommitInfo, error) {
	objectDir := filepath.Join(RepoDirName, "objects", hash[:2])
	objectPath := filepath.Join(objectDir, hash[2:])

	file, err := os.Open(objectPath)
	if err != nil {
		return CommitInfo{}, fmt.Errorf("error reading commit object: %w", err)
	}
	defer file.Close()

	r, err := zlib.NewReader(file)
	if err != nil {
		return CommitInfo{}, fmt.Errorf("error decompressing object: %w", err)
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		return CommitInfo{}, fmt.Errorf("error reading decompressed object: %w", err)
	}

	sepIndex := bytes.IndexByte(data, 0)
	if sepIndex < 0 {
		return CommitInfo{}, fmt.Errorf("invalid object format")
	}

	content := string(data[sepIndex+1:])

	return parseCommitContent(content)
}

func parseCommitContent(content string) (CommitInfo, error) {
	var commit CommitInfo
	lines := strings.Split(content, "\n")

	messageStart := -1
	for i, line := range lines {
		if line == "" {
			messageStart = i + 1
			break
		}

		if strings.HasPrefix(line, "tree ") {
			commit.Tree = strings.TrimPrefix(line, "tree ")
		} else if strings.HasPrefix(line, "parent ") {
			commit.Parents = append(commit.Parents, strings.TrimPrefix(line, "parent "))
		} else if strings.HasPrefix(line, "author ") {
			authorLine := strings.TrimPrefix(line, "author ")
			name, email, timestamp, timezone := parseAuthorLine(authorLine)
			commit.Author = name
			commit.Email = email
			commit.Timestamp = timestamp
			commit.Timezone = timezone
		} else if strings.HasPrefix(line, "committer ") {
			commit.Committer = strings.TrimPrefix(line, "committer ")
		}
	}

	if messageStart >= 0 && messageStart < len(lines) {
		commit.Message = strings.TrimSpace(strings.Join(lines[messageStart:], "\n"))
	}

	return commit, nil
}

func parseAuthorLine(line string) (name, email string, timestamp int64, timezone string) {
	parts := strings.Split(line, "<")
	if len(parts) < 2 {
		return line, "", 0, ""
	}

	name = strings.TrimSpace(parts[0])

	emailAndRest := parts[1]
	emailEnd := strings.Index(emailAndRest, ">")
	if emailEnd == -1 {
		return name, "", 0, ""
	}

	email = emailAndRest[:emailEnd]
	rest := strings.TrimSpace(emailAndRest[emailEnd+1:])

	restParts := strings.Fields(rest)
	if len(restParts) >= 2 {
		fmt.Sscanf(restParts[0], "%d", &timestamp)
		timezone = restParts[1]
	}

	return name, email, timestamp, timezone
}
