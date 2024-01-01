package cmd

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrNotDirectory error = errors.New("Not a directory")
	ErrFileNotFound error = errors.New("File not found")
)

func Ls(directory string) error {
	dirStat, err := os.Stat(directory)
	if err != nil {
		return fmt.Errorf("Failed to stat directory: %w", err)
	}
	if !dirStat.IsDir() {
		return ErrNotDirectory
	}

	dirEntries, err := os.ReadDir(directory)
	if err != nil {
		return err
	}

	for _, entry := range dirEntries {
		fmt.Println(entry.Name())
	}

	return nil
}
