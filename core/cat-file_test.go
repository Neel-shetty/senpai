package core

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(f func()) string {
	var buf bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	io.Copy(&buf, r)
	os.Stdout = stdout
	return buf.String()
}

func TestCatFile(t *testing.T) {
	oldRepo := repoDirName
    repoDirName = t.TempDir()
    t.Cleanup(func() { repoDirName = oldRepo })

	content := []byte("Hello, test object!")
	objectType := "blob"

	hash, err := HashObject(content, objectType, true)
	if err != nil {
		t.Fatalf("failed to write test object: %v", err)
	}

	tests := []struct {
		name       string
		showType   bool
		showSize   bool
		pretty     bool
		exists     bool
		wantOutput string
	}{
		{"ShowType", true, false, false, false, objectType + "\n"},
		{"ShowSize", false, true, false, false, fmt.Sprintf("%d\n", len(content))},
		{"Pretty", false, false, true, false, string(content)},
		{"Exists", false, false, false, true, ""},
		{"Default", false, false, false, false, string(content)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := captureOutput(func() {
				if err := CatFile(hash, tt.showType, tt.showSize, tt.pretty, tt.exists); err != nil {
					t.Errorf("CatFile() error = %v", err)
				}
			})

			out = strings.TrimSpace(out)
			want := strings.TrimSpace(tt.wantOutput)

			if out != want {
				t.Errorf("output = %q, want %q", out, want)
			}
		})
	}
}

func TestCatFileInvalidHash(t *testing.T) {
	oldRepo := repoDirName
    repoDirName = t.TempDir()
    t.Cleanup(func() { repoDirName = oldRepo })

	err := CatFile("deadbeef", false, false, false, false)
	if err == nil {
		t.Errorf("expected error for invalid hash, got nil")
	}
}
