package actions

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/karthikkashyap98/sweeperd/internal/utils"
)

type Matcher interface {
	Match(path string, d fs.DirEntry) bool
}

type MoveInstruction struct {
	Source      string
	Destination string
	Matcher     Matcher
}

func (mv *MoveInstruction) Plan(ctx context.Context) ([]string, error) {
	if mv.Source == "" || mv.Destination == "" {
		return nil, errors.New("Source and destination must be provided")
	}

	srcAbs, err := filepath.Abs(mv.Source)
	if err != nil {
		return nil, err
	}
	dstAbs, err := filepath.Abs(mv.Destination)
	if err != nil {
		return nil, err
	}

	var out []string
	err = fetchFiles(ctx, srcAbs, dstAbs, mv.Matcher, func(path string, d fs.DirEntry) error {
		out = append(out, path)
		return nil
	})
	return out, err
}

func (mv *MoveInstruction) Execute(ctx context.Context, files []string) error {
	dstAbs, err := filepath.Abs(utils.ExpandHome(mv.Destination))
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dstAbs, 0o755); err != nil {
		return err
	}

	for _, src := range files {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		dst := filepath.Join(dstAbs, filepath.Base(src))
		if err := utils.MoveSmart(src, dst); err != nil {
			return err
		}
	}
	return nil
}

func fetchFiles(
	ctx context.Context,
	srcAbs string,
	dstAbs string,
	matcher Matcher,
	emit func(path string, d fs.DirEntry) error,
) error {
	// Normalize trailing separators for safe prefix checks
	dstPrefix := dstAbs + string(os.PathSeparator)

	return filepath.WalkDir(srcAbs, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		// Respect cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if d.IsDir() {
			// Same dir or a child of destination?
			if path == dstAbs || strings.HasPrefix(path+string(os.PathSeparator), dstPrefix) {
				return fs.SkipDir
			}
			return nil
		}

		// Skip if the file itself is in (or is) the destination path.
		// (Normally caught by the directory check above, but keep for safety.)
		if path == dstAbs || strings.HasPrefix(path+string(os.PathSeparator), dstPrefix) {
			return nil
		}

		if matcher != nil && !matcher.Match(path, d) {
			return nil
		}
		return emit(path, d)
	})
}
