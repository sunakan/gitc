package git

import (
	"os/exec"
	"strings"
	"testing"
)

func TestBranchExists(t *testing.T) {
	if !isIntegrationTest() {
		t.Skip("統合テスト環境でのみ実行")
	}

	// テスト用のGitリポジトリをセットアップ
	repoPath, cleanup := createTestGitRepo(t)
	defer cleanup()
	restoreDir := changeDir(t, repoPath)
	defer restoreDir()

	// テスト用ブランチを作成
	cmd := exec.Command("git", "checkout", "-b", "test-branch")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create test branch: %v", err)
	}

	// mainブランチに戻る
	cmd = exec.Command("git", "checkout", "main")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to checkout main: %v", err)
	}

	tests := []struct {
		name       string
		branch     string
		wantExists bool
		wantError  bool
	}{
		{
			name:       "存在するブランチ（main）",
			branch:     "main",
			wantExists: true,
			wantError:  false,
		},
		{
			name:       "存在するブランチ（test-branch）",
			branch:     "test-branch",
			wantExists: true,
			wantError:  false,
		},
		{
			name:       "存在しないブランチ",
			branch:     "non-existent-branch",
			wantExists: false,
			wantError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := BranchExists(tt.branch)

			if tt.wantError && err == nil {
				t.Errorf("BranchExists() expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("BranchExists() unexpected error: %v", err)
			}
			if exists != tt.wantExists {
				t.Errorf("BranchExists() = %v, want %v", exists, tt.wantExists)
			}
		})
	}
}

func TestDefaultBranchOption(t *testing.T) {
	if !isIntegrationTest() {
		t.Skip("統合テスト環境でのみ実行")
	}

	// テスト用のGitリポジトリをセットアップ
	repoPath, cleanup := createTestGitRepo(t)
	defer cleanup()
	restoreDir := changeDir(t, repoPath)
	defer restoreDir()

	// テスト用ブランチを作成
	cmd := exec.Command("git", "checkout", "-b", "custom-default")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create custom-default branch: %v", err)
	}

	// mainブランチに戻る
	cmd = exec.Command("git", "checkout", "main")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to checkout main: %v", err)
	}

	tests := []struct {
		name            string
		defaultBranch   string
		wantError       bool
		expectedBranch  string
	}{
		{
			name:            "存在するカスタムブランチを指定",
			defaultBranch:   "custom-default",
			wantError:       false,
			expectedBranch:  "custom-default",
		},
		{
			name:            "存在しないブランチを指定",
			defaultBranch:   "non-existent",
			wantError:       true,
			expectedBranch:  "",
		},
		{
			name:            "空文字列（自動検出）",
			defaultBranch:   "",
			wantError:       false,
			expectedBranch:  "main", // 自動検出でmainが選ばれる
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := CleanupOptions{
				DryRun:        true, // ドライランで安全にテスト
				DefaultBranch: tt.defaultBranch,
				Yes:           true,
			}

			result, err := ExecuteCleanup(options)

			if tt.wantError {
				if err == nil {
					t.Errorf("ExecuteCleanup() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ExecuteCleanup() unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Fatal("ExecuteCleanup() returned nil result")
			}

			if result.DefaultBranch != tt.expectedBranch {
				t.Errorf("ExecuteCleanup() DefaultBranch = %v, want %v", result.DefaultBranch, tt.expectedBranch)
			}
		})
	}
}

func TestDefaultBranchValidation(t *testing.T) {
	if !isIntegrationTest() {
		t.Skip("統合テスト環境でのみ実行")
	}

	// テスト用のGitリポジトリをセットアップ
	repoPath, cleanup := createTestGitRepo(t)
	defer cleanup()
	restoreDir := changeDir(t, repoPath)
	defer restoreDir()

	// 存在しないブランチでのテスト
	options := CleanupOptions{
		DryRun:        true,
		DefaultBranch: "definitely-does-not-exist",
		Yes:           true,
	}

	_, err := ExecuteCleanup(options)
	if err == nil {
		t.Error("ExecuteCleanup() should fail with non-existent branch")
	}

	// エラーメッセージに指定したブランチ名が含まれているか確認
	errorStr := err.Error()
	if !strings.Contains(errorStr, "definitely-does-not-exist") {
		t.Errorf("Error message should contain branch name: %s", errorStr)
	}
}