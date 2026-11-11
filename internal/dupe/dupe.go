// Package dupe finds files with identical content.
package dupe

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// Group is a set of files with identical content.
type Group struct {
	Hash  string
	Size  int64
	Paths []string
}

// Find walks roots and returns groups of duplicate files (2+ copies).
// Files are first bucketed by size, and only same-size files are hashed.
func Find(roots []string, minSize int64) ([]Group, []error) {
	var errs []error
	bySize := map[int64][]string{}
	for _, root := range roots {
		err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				errs = append(errs, err)
				return nil
			}
			if !d.Type().IsRegular() {
				return nil
			}
			info, err := d.Info()
			if err != nil {
				errs = append(errs, err)
				return nil
			}
			if info.Size() >= minSize {
				bySize[info.Size()] = append(bySize[info.Size()], p)
			}
			return nil
		})
		if err != nil {
			errs = append(errs, err)
		}
	}
	byHash := map[string][]string{}
	sizes := map[string]int64{}
	for size, paths := range bySize {
		if len(paths) < 2 {
			continue
		}
		for _, p := range paths {
			h, err := HashFile(p)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			byHash[h] = append(byHash[h], p)
			sizes[h] = size
		}
	}
	var groups []Group
	for h, paths := range byHash {
		if len(paths) >= 2 {
			groups = append(groups, Group{Hash: h, Size: sizes[h], Paths: paths})
		}
	}
	return groups, errs
}

// HashFile returns the hex SHA-256 of a file's contents.
func HashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("hash %s: %w", path, err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
