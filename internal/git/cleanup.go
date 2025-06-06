package git

import (
	"fmt"
)

// CleanupOptions はクリーンアップ処理のオプションを表します
type CleanupOptions struct {
	DryRun        bool   // 実行のシミュレーションのみ
	Verbose       bool   // 詳細ログの表示
	Yes           bool   // 確認プロンプトのスキップ
	Force         bool   // 強制実行（未マージブランチも削除）
	DefaultBranch string // 手動指定のデフォルトブランチ
	ExcludePattern string // 除外パターン
	NoPull        bool   // プル処理のスキップ
}

// CleanupResult はクリーンアップ処理の結果を表します
type CleanupResult struct {
	DefaultBranch    string   // 検出されたデフォルトブランチ
	DeletedBranches  []string // 削除されたブランチのリスト
	SkippedBranches  []string // スキップされたブランチのリスト
	Errors          []error  // 発生したエラーのリスト
	WasDryRun       bool     // ドライランモードだったかどうか
}

// Validate はオプションの妥当性をチェックします
func (opts *CleanupOptions) Validate() error {
	if opts.DryRun && opts.Force {
		return fmt.Errorf("--dry-run and --force cannot be used together")
	}
	return nil
}

// ExecuteCleanup はメインのクリーンアップ処理を実行します
func ExecuteCleanup(options CleanupOptions) (*CleanupResult, error) {
	// オプションのバリデーション
	if err := options.Validate(); err != nil {
		return nil, err
	}

	result := &CleanupResult{
		WasDryRun: options.DryRun,
	}

	// 1. Gitリポジトリかどうかの確認
	cwd, err := GetCurrentDirectory()
	if err != nil {
		return nil, NewGitError("cleanup", err)
	}

	if err := IsGitRepository(cwd); err != nil {
		return nil, NewGitError("cleanup", ErrNotGitRepository).WithPath(cwd)
	}

	// 2. デフォルトブランチの検出
	var defaultBranch string
	if options.DefaultBranch != "" {
		// 手動指定されたブランチの存在確認
		exists, err := BranchExists(options.DefaultBranch)
		if err != nil {
			return nil, NewGitError("cleanup", err).WithMessage("failed to check branch existence")
		}
		if !exists {
			return nil, NewGitError("cleanup", fmt.Errorf("specified branch '%s' does not exist", options.DefaultBranch))
		}
		defaultBranch = options.DefaultBranch
	} else {
		defaultBranch, err = DetectDefaultBranch()
		if err != nil {
			return nil, NewGitError("cleanup", err)
		}
	}
	result.DefaultBranch = defaultBranch

	// 3. デフォルトブランチへの切り替え
	currentBranch, err := GetCurrentBranch()
	if err != nil {
		return nil, NewGitError("cleanup", err)
	}

	if currentBranch != defaultBranch {
		if options.DryRun {
			// ドライランモードではブランチ切り替えはシミュレーションのみ
		} else {
			if err := CheckoutBranch(defaultBranch); err != nil {
				return nil, NewGitError("cleanup", err).WithMessage("failed to switch to default branch")
			}
		}
	}

	// 4. フェッチ処理（必須・ドライランでも実行）
	if err := Fetch(); err != nil {
		// フェッチ失敗は警告として扱い、処理を継続
		result.Errors = append(result.Errors, NewGitError("cleanup", err).WithMessage("fetch failed"))
	}

	if options.DryRun {
		// ドライランモードの場合はfetch以外の実際の処理は行わない
		return result, nil
	}

	// 5. プル処理（--no-pullが指定されていない場合）
	if !options.NoPull {
		if err := Pull(); err != nil {
			// プル失敗は警告として扱い、処理を継続
			result.Errors = append(result.Errors, NewGitError("cleanup", err).WithMessage("pull failed"))
		}
	}

	// 6. ローカルブランチの一覧取得
	branches, err := ListLocalBranches()
	if err != nil {
		return nil, NewGitError("cleanup", err)
	}

	// 7. ブランチの削除
	for _, branch := range branches {
		if branch == defaultBranch {
			// デフォルトブランチはスキップ
			result.SkippedBranches = append(result.SkippedBranches, branch)
			continue
		}

		// 除外パターンのチェック（簡単な実装）
		if options.ExcludePattern != "" && branch == options.ExcludePattern {
			result.SkippedBranches = append(result.SkippedBranches, branch)
			continue
		}

		// ブランチ削除
		if err := DeleteBranch(branch, options.Force); err != nil {
			result.Errors = append(result.Errors, NewGitError("cleanup", err).WithPath(branch))
			result.SkippedBranches = append(result.SkippedBranches, branch)
		} else {
			result.DeletedBranches = append(result.DeletedBranches, branch)
		}
	}

	return result, nil
}