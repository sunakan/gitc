// +build integration

package git

import (
	"os/exec"
	"testing"
)

func TestCleanupIntegrationFlow(t *testing.T) {
	// git コマンドが利用可能かチェック
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git command not found, skipping integration tests")
	}

	t.Run("完全なクリーンアップフロー", func(t *testing.T) {
		dir, cleanup := createTestGitRepo(t)
		defer cleanup()
		defer changeDir(t, dir)()

		// テスト用ブランチを複数作成
		branches := []string{"feature/test1", "feature/test2", "bugfix/issue123"}
		for _, branch := range branches {
			_, err := ExecuteCommand("checkout", "-b", branch)
			if err != nil {
				t.Fatalf("Failed to create branch %s: %v", branch, err)
			}
		}

		// mainブランチに戻る
		if err := CheckoutBranch("main"); err != nil {
			// masterブランチを試す
			if err := CheckoutBranch("master"); err != nil {
				t.Fatalf("Failed to checkout main/master branch: %v", err)
			}
		}

		// クリーンアップ実行
		options := CleanupOptions{
			DryRun:  false,
			Verbose: true,
			Yes:     true,
			NoPull:  true, // テスト環境ではプルをスキップ
		}

		result, err := ExecuteCleanup(options)
		if err != nil {
			t.Fatalf("ExecuteCleanup() failed: %v", err)
		}

		// 結果の検証
		if len(result.DeletedBranches) != len(branches) {
			t.Errorf("Expected %d deleted branches, got %d", len(branches), len(result.DeletedBranches))
		}

		// デフォルトブランチが設定されていることを確認
		if result.DefaultBranch == "" {
			t.Error("Default branch should be detected")
		}

		// 現在のブランチがデフォルトブランチになっていることを確認
		currentBranch, err := GetCurrentBranch()
		if err != nil {
			t.Fatalf("Failed to get current branch: %v", err)
		}

		if currentBranch != result.DefaultBranch {
			t.Errorf("Current branch %s does not match default branch %s", currentBranch, result.DefaultBranch)
		}
	})

	t.Run("処理順序の確認", func(t *testing.T) {
		dir, cleanup := createTestGitRepo(t)
		defer cleanup()
		defer changeDir(t, dir)()

		// テスト用ブランチを作成し、そのブランチにいる状態でクリーンアップを実行
		_, err := ExecuteCommand("checkout", "-b", "test-branch")
		if err != nil {
			t.Fatalf("Failed to create test branch: %v", err)
		}

		// 現在のブランチを確認
		currentBranch, err := GetCurrentBranch()
		if err != nil {
			t.Fatalf("Failed to get current branch: %v", err)
		}

		if currentBranch != "test-branch" {
			t.Fatalf("Expected to be on test-branch, got %s", currentBranch)
		}

		// クリーンアップ実行
		options := CleanupOptions{
			DryRun: false,
			Yes:    true,
			NoPull: true,
		}

		result, err := ExecuteCleanup(options)
		if err != nil {
			t.Fatalf("ExecuteCleanup() failed: %v", err)
		}

		// クリーンアップ後、デフォルトブランチに切り替わっていることを確認
		finalBranch, err := GetCurrentBranch()
		if err != nil {
			t.Fatalf("Failed to get final branch: %v", err)
		}

		if finalBranch != result.DefaultBranch {
			t.Errorf("Expected to be on default branch %s, got %s", result.DefaultBranch, finalBranch)
		}

		// test-branchが削除されていることを確認
		found := false
		for _, deleted := range result.DeletedBranches {
			if deleted == "test-branch" {
				found = true
				break
			}
		}

		if !found {
			t.Error("test-branch should have been deleted")
		}
	})

	t.Run("エラーハンドリング", func(t *testing.T) {
		dir, cleanup := createTestGitRepo(t)
		defer cleanup()
		defer changeDir(t, dir)()

		// 存在しないデフォルトブランチを指定
		options := CleanupOptions{
			DefaultBranch: "nonexistent-branch",
			DryRun:        false,
			Yes:           true,
			NoPull:        true,
		}

		_, err := ExecuteCleanup(options)
		if err == nil {
			t.Error("Expected error when specifying nonexistent default branch")
		}
	})
}