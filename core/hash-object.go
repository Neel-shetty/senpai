package core

import (
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

func HashObject(fileContent []byte, objectType string, write bool) (string, error) {
	hash, fullObjectData := calculateObjectHash(fileContent, objectType)

	if !write {
		return hash, nil
	}

	if err := writeObject(hash, fullObjectData); err != nil {
		return "", err
	}

	return hash, nil
}

func calculateObjectHash(content []byte, objectType string) (string, []byte) {
	header := fmt.Sprintf("%s %d\u0000", objectType, len(content))
	store := append([]byte(header), content...)

	h := sha1.New()
	h.Write(store)
	hash := hex.EncodeToString(h.Sum(nil))

	return hash, store
}

func writeObject(hash string, data []byte) error {
	objectDir := filepath.Join(repoDirName, "objects", hash[:2])
	objectPath := filepath.Join(objectDir, hash[2:])

	// skip writing to disk if it already exists
	if _, err := os.Stat(objectPath); err == nil {
		return nil
	}

	if err := os.MkdirAll(objectDir, 0755); err != nil {
		return fmt.Errorf("failed to create object: %w", err)
	}

	f, err := os.Create(objectPath)
	if err != nil {
		return fmt.Errorf("failed to create object file: %w", err)
	}
	defer f.Close()

	w := zlib.NewWriter(f)
	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("failed to write compressed object: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to close zlib writer: %w", err)
	}
	return nil

}
