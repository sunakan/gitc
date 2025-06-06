// +build integration

package git

import (
	"os/exec"
	"testing"
)

func TestIntegrationGitOperations(t *testing.T) {
	// git コマンドが利用可能かチェック
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git command not found, skipping integration tests")
	}
	
	t.Run("実際のGitリポジトリでの操作", func(t *testing.T) {
		dir, cleanup := createTestGitRepo(t)
		defer cleanup()
		defer changeDir(t, dir)()
		
		// リポジトリ検証
		if err := IsGitRepository(dir); err != nil {
			t.Errorf("IsGitRepository() failed: %v", err)
		}
		
		// 現在のブランチ取得
		branch, err := GetCurrentBranch()
		if err != nil {
			t.Fatalf("GetCurrentBranch() failed: %v", err)
		}
		
		// デフォルトではmainまたはmasterブランチのはず
		if branch != "main" && branch != "master" {
			t.Errorf("GetCurrentBranch() = %v, want main or master", branch)
		}
		
		// ローカルブランチ一覧
		branches, err := ListLocalBranches()
		if err != nil {
			t.Fatalf("ListLocalBranches() failed: %v", err)
		}
		
		if len(branches) == 0 {
			t.Error("ListLocalBranches() returned empty list")
		}
		
		// 新しいブランチを作成してチェックアウト
		result, err := ExecuteCommand("checkout", "-b", "test-branch")
		if err != nil {
			t.Fatalf("Failed to create test branch: %v", err)
		}
		
		// 新しいブランチに切り替わったことを確認
		currentBranch, err := GetCurrentBranch()
		if err != nil {
			t.Fatalf("GetCurrentBranch() after checkout failed: %v", err)
		}
		
		if currentBranch != "test-branch" {
			t.Errorf("GetCurrentBranch() = %v, want test-branch", currentBranch)
		}
		
		// 元のブランチに戻る
		if err := CheckoutBranch(branch); err != nil {
			t.Fatalf("CheckoutBranch() failed: %v", err)
		}
		
		// テストブランチを削除
		if err := DeleteBranch("test-branch", false); err != nil {
			t.Errorf("DeleteBranch() failed: %v", err)
		}
	})
}