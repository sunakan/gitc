package git

import (
	"testing"
)

func TestExecuteCleanup(t *testing.T) {
	tests := []struct {
		name    string
		options CleanupOptions
		setup   func(t *testing.T) string
		wantErr bool
		errType error
	}{
		{
			name: "正常なクリーンアップ処理",
			options: CleanupOptions{
				DryRun:  false,
				Verbose: false,
				Yes:     true,
			},
			setup: func(t *testing.T) string {
				// テスト用Gitリポジトリを作成
				dir, _ := createTestGitRepo(t)
				return dir
			},
			wantErr: false,
		},
		{
			name: "Gitリポジトリではないディレクトリでの実行",
			options: CleanupOptions{
				DryRun:  false,
				Verbose: false,
				Yes:     true,
			},
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			wantErr: true,
			errType: ErrNotGitRepository,
		},
		{
			name: "ドライランモードでの実行",
			options: CleanupOptions{
				DryRun:  true,
				Verbose: false,
				Yes:     true,
			},
			setup: func(t *testing.T) string {
				dir, _ := createTestGitRepo(t)
				return dir
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setup(t)
			defer changeDir(t, dir)()

			result, err := ExecuteCleanup(tt.options)

			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteCleanup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errType != nil {
				gitErr, ok := err.(*GitError)
				if !ok {
					t.Errorf("Expected GitError, got %T", err)
					return
				}
				if !IsNotGitRepository(gitErr) && tt.errType == ErrNotGitRepository {
					t.Errorf("Expected ErrNotGitRepository, got %v", gitErr.Err)
				}
			}

			if !tt.wantErr && result == nil {
				t.Error("ExecuteCleanup() returned nil result on success")
			}
		})
	}
}

func TestCleanupOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		options CleanupOptions
		wantErr bool
	}{
		{
			name: "有効なオプション",
			options: CleanupOptions{
				DryRun:  false,
				Verbose: false,
				Yes:     true,
			},
			wantErr: false,
		},
		{
			name: "ドライランとフォースの組み合わせ（無効）",
			options: CleanupOptions{
				DryRun: true,
				Force:  true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.options.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("CleanupOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}