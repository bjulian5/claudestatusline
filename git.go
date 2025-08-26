package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetGitBranch(dir string) (string, error) {
	for {
		gitDir := filepath.Join(dir, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			headFile := filepath.Join(gitDir, "HEAD")
			content, err := os.ReadFile(headFile)
			if err != nil {
				return "", fmt.Errorf("failed to read HEAD file: %w", err)
			}

			ref := strings.TrimSpace(string(content))
			if strings.HasPrefix(ref, "ref: refs/heads/") {
				return strings.TrimPrefix(ref, "ref: refs/heads/"), nil
			}

			if len(ref) >= 7 {
				return ref[:7] + "...", nil
			}
			return ref, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir || parent == "/" {
			return "", fmt.Errorf("not a git repository")
		}
		dir = parent
	}
}
