package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
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

func ExpandHome(p string) string {
	if p == "" || p[0] != '~' {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return p
	}
	if p == "~" {
		return home
	}
	if strings.HasPrefix(p, "~/") {
		return filepath.Join(home, p[2:])
	}
	return p
}

// moveSmart tries os.Rename first, and if it fails with EXDEV (cross-device),
// it copies the file and then removes the source.
func MoveSmart(src, dst string) error {
	// Fast path
	if err := os.Rename(src, dst); err == nil {
		return nil
	} else if !IsCrossDevice(err) {
		return err
	}
	// Cross-device: copy then remove
	if err := CopyFile(src, dst); err != nil {
		return err
	}
	return os.Remove(src)
}

func IsCrossDevice(err error) bool {
	// OS-specific checks; for MVP, a coarse check:
	// On Unix, os.Rename across devices returns EXDEV
	return strings.Contains(strings.ToLower(err.Error()), "cross-device")
}

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// Ensure parent exists
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
