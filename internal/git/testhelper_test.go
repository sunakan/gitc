package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// テスト用のGitリポジトリを作成するヘルパー関数
func createTestGitRepo(t *testing.T) (string, func()) {
	t.Helper()
	
	dir := t.TempDir()
	
	// Git初期化
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}
	
	// ユーザー設定（テスト用）
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = dir
	cmd.Run()
	
	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = dir
	cmd.Run()
	
	// 初期コミット
	testFile := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	cmd = exec.Command("git", "add", ".")
	cmd.Dir = dir
	cmd.Run()
	
	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = dir
	cmd.Run()
	
	cleanup := func() {
		// クリーンアップは t.TempDir() が自動的に行う
	}
	
	return dir, cleanup
}

// 現在のディレクトリを変更して、テスト後に元に戻すヘルパー関数
func changeDir(t *testing.T, dir string) func() {
	t.Helper()
	
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	
	return func() {
		if err := os.Chdir(oldDir); err != nil {
			t.Errorf("Failed to restore directory: %v", err)
		}
	}
}