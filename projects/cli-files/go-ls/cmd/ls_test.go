package cmd_test

import (
	"errors"
	"go-ls/cmd"
	"os"
	"testing"
)

func TestLsEmptyDirectory(t *testing.T) {
	// GIVEN
	tempDir, err := os.MkdirTemp("", "test-dir-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempDir)
	// WHEN
	err = cmd.Ls(tempDir)
	// THEN
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestLsInvalidDirectory(t *testing.T) {
	// GIVEN
	tempDir := "invalid-directory"
	defer os.Remove(tempDir)
	// WHEN
	err := cmd.Ls(tempDir)
	// THEN
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Expected error: %v, but got: %v", os.ErrNotExist, err)
	}
}

func TestLsDirectoryWithFiles(t *testing.T) {
	// GIVEN
	tempDir, err := os.MkdirTemp("", "test-dir-*")
	if err != nil {
		t.Fatal(err)
	}
	files := []string{"file1", "file2", "file3"}
	for _, file := range files {
		filePath := tempDir + "/" + file
		err = os.WriteFile(filePath, []byte(file), 0o644)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(filePath)
	}

	defer os.Remove(tempDir)
	// WHEN
	err = cmd.Ls(tempDir)
	// THEN
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
