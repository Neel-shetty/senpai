package core

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestHashObjectNoWrite(t *testing.T) {
	content := []byte("hello world")
	expectedHash := "95d09f2b10159347eece71399a7e2e907ea3df4f"

	hash, err := HashObject(content, "blob", false)

	if err != nil {
		t.Fatalf("HashObject() with write=false returned an unexpected error: %v", err)
	}

	if hash != expectedHash {
		t.Errorf("expected hash %s, got %s", expectedHash, hash)
	}
}

func TestHashObjectWrite(t *testing.T) {
	tempDir := t.TempDir()

	originalRepoDir := repoDirName
	repoDirName = tempDir
	defer func() { repoDirName = originalRepoDir }()

	content := []byte("what is up, doc?")
	objectType := "blob"
	expectedHeader := fmt.Sprintf("%s %d\u0000", objectType, len(content))
	expectedFullContent := append([]byte(expectedHeader), content...)
	expectedHash := "bd9dbf5aae1a3862dd1526723246b20206e5fc37"
	expectedPath := filepath.Join(tempDir, "objects", expectedHash[:2], expectedHash[2:])

	t.Run("SuccessfulWrite", func(t *testing.T) {
		hash, err := HashObject(content, objectType, true)

		if err != nil {
			t.Fatalf("HashObject() with write=true returned an unexpected error: %v", err)
		}
		if hash != expectedHash {
			t.Errorf("expected hash %s, got %s", expectedHash, hash)
		}

		file, err := os.Open(expectedPath)
		if err != nil {
			t.Fatalf("could not open expected object file at %s: %v", expectedPath, err)
		}
		defer file.Close()

		// Verify the file's contents are correct by decompressing them.
		r, err := zlib.NewReader(file)
		if err != nil {
			t.Fatalf("could not create zlib reader for object file: %v", err)
		}
		defer r.Close()

		decompressed, err := io.ReadAll(r)
		if err != nil {
			t.Fatalf("failed to decompress object file content: %v", err)
		}

		if !bytes.Equal(decompressed, expectedFullContent) {
			t.Errorf("decompressed content does not match expected content")
		}
	})

	t.Run("ObjectAlreadyExists", func(t *testing.T) {
		hash, err := HashObject(content, objectType, true)
		if err != nil {
			t.Fatalf("HashObject() returned an error when object already exists: %v", err)
		}
		if hash != expectedHash {
			t.Errorf("expected hash %s, got %s", expectedHash, hash)
		}
	})
}
