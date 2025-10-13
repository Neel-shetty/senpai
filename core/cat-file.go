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

func CatFile(hash string, showType, showSize, pretty, exists bool) error {
	objectDir := filepath.Join(repoDirName, "objects", hash[:2])
	objectPath := filepath.Join(objectDir, hash[2:])

	file, err := os.Open(objectPath)
	if err != nil {
		return fmt.Errorf("error reading object: %s", err)
	}
	defer file.Close()

	r, err := zlib.NewReader(file)
	if err != nil {
		return fmt.Errorf("error decompressing object: %w", err)
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("error reading decompressed object: %w", err)
	}

	sepIndex := bytes.IndexByte(data, 0)
	if sepIndex < 0 {
		return fmt.Errorf("invalid object format")
	}

	header := string(data[:sepIndex])
	content := data[sepIndex+1:]

	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		return fmt.Errorf("invalid object format")
	}

	objectType := parts[0]
	objectSize := parts[1]

	if exists {
		return nil
	}

	if showType {
		fmt.Println(objectType)
	}
	if showSize {
		fmt.Println(objectSize)
	}
	if pretty || (!showType && !showSize && !exists) {
		fmt.Print(string(content))
	}

	return nil
}
