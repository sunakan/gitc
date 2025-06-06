package git

import (
	"fmt"
	"os"
	"path/filepath"
)

// IsGitRepository は指定されたディレクトリがGitリポジトリかどうかを確認します
func IsGitRepository(path string) error {
	gitDir := filepath.Join(path, ".git")
	info, err := os.Stat(gitDir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("not a git repository: %s", path)
		}
		return fmt.Errorf("failed to check git directory: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf(".git exists but is not a directory: %s", gitDir)
	}

	return nil
}

// GetCurrentDirectory は現在の作業ディレクトリを返します
func GetCurrentDirectory() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}
	return cwd, nil
}