package utils

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"os"
	"path/filepath"
	"strings"
)

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !errors.Is(err, os.ErrNotExist)
}

func DirectoryHasYAML(dirPath string) bool {
	files, err := os.ReadDir(dirPath)
	if errors.Is(err, os.ErrNotExist) || err != nil {
		return false
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(file.Name()))
		if ext == ".yaml" || ext == ".yml" {
			return true
		}
	}

	return false
}

func PrintError(err error) {
	// Only colorize if output is a terminal
	fmt.Fprintf(os.Stderr, "%s %v\n", color.RedString("✖"), err)
}
