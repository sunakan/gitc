package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsGitRepository(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) string
		wantErr bool
		errMsg  string
	}{
		{
			name: "有効なGitリポジトリ",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				gitDir := filepath.Join(dir, ".git")
				if err := os.Mkdir(gitDir, 0755); err != nil {
					t.Fatalf("Failed to create .git directory: %v", err)
				}
				return dir
			},
			wantErr: false,
		},
		{
			name: "Gitリポジトリではないディレクトリ",
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			wantErr: true,
			errMsg:  "not a git repository",
		},
		{
			name: ".gitがファイルの場合",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				gitFile := filepath.Join(dir, ".git")
				if err := os.WriteFile(gitFile, []byte("gitdir: /path/to/git"), 0644); err != nil {
					t.Fatalf("Failed to create .git file: %v", err)
				}
				return dir
			},
			wantErr: true,
			errMsg:  ".git exists but is not a directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setup(t)
			err := IsGitRepository(dir)

			if (err != nil) != tt.wantErr {
				t.Errorf("IsGitRepository() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errMsg != "" && err.Error() != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("IsGitRepository() error message = %v, want to contain %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestGetCurrentDirectory(t *testing.T) {
	// 現在のディレクトリを取得
	got, err := GetCurrentDirectory()
	if err != nil {
		t.Fatalf("GetCurrentDirectory() error = %v", err)
	}

	// 結果が空でないことを確認
	if got == "" {
		t.Error("GetCurrentDirectory() returned empty string")
	}

	// os.Getwd()の結果と一致することを確認
	expected, _ := os.Getwd()
	if got != expected {
		t.Errorf("GetCurrentDirectory() = %v, want %v", got, expected)
	}
}

// contains はstrがsubstrを含むかチェックするヘルパー関数
func contains(str, substr string) bool {
	return len(str) >= len(substr) && str[:len(substr)] == substr || 
		len(str) >= len(substr) && contains(str[1:], substr)
}