package git

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestVerboseMode(t *testing.T) {
	tests := []struct {
		name    string
		options CleanupOptions
		wantLog bool
	}{
		{
			name: "verboseモード有効",
			options: CleanupOptions{
				DryRun:  true,
				Verbose: true,
				Yes:     true,
			},
			wantLog: true,
		},
		{
			name: "verboseモード無効",
			options: CleanupOptions{
				DryRun:  true,
				Verbose: false,
				Yes:     true,
			},
			wantLog: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !isIntegrationTest() {
				t.Skip("統合テスト環境でのみ実行")
			}

			// ログ出力をキャプチャ
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(nil) // テスト後にリセット

			// テスト用のGitリポジトリをセットアップ
			repoPath, cleanup := createTestGitRepo(t)
			defer cleanup()
			restoreDir := changeDir(t, repoPath)
			defer restoreDir()

			// テスト実行
			result, err := ExecuteCleanup(tt.options)
			
			// エラーは許可（リモート接続エラーなど）
			if err != nil {
				t.Logf("Expected error in test environment: %v", err)
			}
			
			if result == nil {
				t.Fatal("ExecuteCleanup() returned nil result")
			}

			// ログ出力の確認
			logOutput := buf.String()
			hasVerboseLog := strings.Contains(logOutput, "[VERBOSE]")

			if tt.wantLog && !hasVerboseLog {
				t.Errorf("verboseモードが有効なのに詳細ログが出力されていません")
			}
			
			if !tt.wantLog && hasVerboseLog {
				t.Errorf("verboseモードが無効なのに詳細ログが出力されています")
			}

			// verboseモード有効時の詳細確認
			if tt.wantLog {
				expectedLogs := []string{
					"クリーンアップ処理を開始",
					"Gitリポジトリの確認を開始",
					"デフォルトブランチの検出を開始",
					"フェッチ処理を開始",
				}
				
				for _, expectedLog := range expectedLogs {
					if !strings.Contains(logOutput, expectedLog) {
						t.Errorf("期待される詳細ログが見つかりません: %s", expectedLog)
					}
				}
			}
		})
	}
}

func TestVerboseLogContent(t *testing.T) {
	if !isIntegrationTest() {
		t.Skip("統合テスト環境でのみ実行")
	}

	// ログ出力をキャプチャ
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	// テスト用のGitリポジトリをセットアップ
	repoPath, cleanup := createTestGitRepo(t)
	defer cleanup()
	restoreDir := changeDir(t, repoPath)
	defer restoreDir()

	options := CleanupOptions{
		DryRun:  true,
		Verbose: true,
		Yes:     true,
	}

	result, err := ExecuteCleanup(options)
	
	// エラーは許可
	if err != nil {
		t.Logf("Expected error in test environment: %v", err)
	}
	
	if result == nil {
		t.Fatal("ExecuteCleanup() returned nil result")
	}

	logOutput := buf.String()
	
	// オプション情報がログに含まれているか確認
	if !strings.Contains(logOutput, "DryRun=true") {
		t.Error("オプション情報がログに含まれていません")
	}
	
	if !strings.Contains(logOutput, "Verbose=true") {
		t.Error("verboseオプション情報がログに含まれていません")
	}

	// 処理ステップがログに含まれているか確認
	expectedSteps := []string{
		"現在のディレクトリ:",
		"検出されたデフォルトブランチ:",
		"現在のブランチ:",
		"git fetch --all --prune",
		"ドライランモードのため",
	}
	
	for _, step := range expectedSteps {
		if !strings.Contains(logOutput, step) {
			t.Errorf("期待される処理ステップがログに含まれていません: %s", step)
		}
	}
}

func TestVerboseWithNormalMode(t *testing.T) {
	// verboseフラグが通常の動作に影響しないことを確認
	if !isIntegrationTest() {
		t.Skip("統合テスト環境でのみ実行")
	}

	repoPath, cleanup := createTestGitRepo(t)
	defer cleanup()
	restoreDir := changeDir(t, repoPath)
	defer restoreDir()

	// verbose無効での実行
	options1 := CleanupOptions{
		DryRun:  true,
		Verbose: false,
		Yes:     true,
	}

	result1, err1 := ExecuteCleanup(options1)

	// verbose有効での実行（ログ出力をキャプチャ）
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	options2 := CleanupOptions{
		DryRun:  true,
		Verbose: true,
		Yes:     true,
	}

	result2, err2 := ExecuteCleanup(options2)

	// 結果が同じであることを確認（ログ出力以外）
	if (err1 == nil) != (err2 == nil) {
		t.Error("verboseモードの有無でエラー状態が異なります")
	}

	if result1 != nil && result2 != nil {
		if result1.DefaultBranch != result2.DefaultBranch {
			t.Error("verboseモードの有無でDefaultBranchが異なります")
		}
		
		if result1.WasDryRun != result2.WasDryRun {
			t.Error("verboseモードの有無でWasDryRunが異なります")
		}
	}
}