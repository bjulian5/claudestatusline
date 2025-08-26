package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetGitBranch(t *testing.T) {
	t.Run("branch ref in HEAD", func(t *testing.T) {
		tmpDir := t.TempDir()
		gitDir := filepath.Join(tmpDir, ".git")
		require.NoError(t, os.Mkdir(gitDir, 0755))

		headFile := filepath.Join(gitDir, "HEAD")
		err := os.WriteFile(headFile, []byte("ref: refs/heads/main\n"), 0644)
		require.NoError(t, err)

		branch, err := GetGitBranch(tmpDir)
		assert.NoError(t, err)
		assert.Equal(t, "main", branch)
	})

	t.Run("detached HEAD with long commit hash", func(t *testing.T) {
		tmpDir := t.TempDir()
		gitDir := filepath.Join(tmpDir, ".git")
		require.NoError(t, os.Mkdir(gitDir, 0755))

		headFile := filepath.Join(gitDir, "HEAD")
		err := os.WriteFile(headFile, []byte("abc123456789def\n"), 0644)
		require.NoError(t, err)

		branch, err := GetGitBranch(tmpDir)
		assert.NoError(t, err)
		assert.Equal(t, "abc1234...", branch)
	})

	t.Run("detached HEAD with short commit hash", func(t *testing.T) {
		tmpDir := t.TempDir()
		gitDir := filepath.Join(tmpDir, ".git")
		require.NoError(t, os.Mkdir(gitDir, 0755))

		headFile := filepath.Join(gitDir, "HEAD")
		err := os.WriteFile(headFile, []byte("abc123\n"), 0644)
		require.NoError(t, err)

		branch, err := GetGitBranch(tmpDir)
		assert.NoError(t, err)
		assert.Equal(t, "abc123", branch)
	})

	t.Run("finds git directory in parent", func(t *testing.T) {
		tmpDir := t.TempDir()
		gitDir := filepath.Join(tmpDir, ".git")
		require.NoError(t, os.Mkdir(gitDir, 0755))

		headFile := filepath.Join(gitDir, "HEAD")
		err := os.WriteFile(headFile, []byte("ref: refs/heads/develop\n"), 0644)
		require.NoError(t, err)

		subDir := filepath.Join(tmpDir, "subdir", "deep")
		require.NoError(t, os.MkdirAll(subDir, 0755))

		branch, err := GetGitBranch(subDir)
		assert.NoError(t, err)
		assert.Equal(t, "develop", branch)
	})

	t.Run("not a git repository", func(t *testing.T) {
		tmpDir := t.TempDir()
		
		branch, err := GetGitBranch(tmpDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not a git repository")
		assert.Empty(t, branch)
	})

	t.Run("HEAD file unreadable", func(t *testing.T) {
		tmpDir := t.TempDir()
		gitDir := filepath.Join(tmpDir, ".git")
		require.NoError(t, os.Mkdir(gitDir, 0755))

		headFile := filepath.Join(gitDir, "HEAD")
		require.NoError(t, os.WriteFile(headFile, []byte("ref: refs/heads/main\n"), 0000))

		branch, err := GetGitBranch(tmpDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read HEAD file")
		assert.Empty(t, branch)
	})
}