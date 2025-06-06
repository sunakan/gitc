package git

import (
	"fmt"
	"log"
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
	// verboseログ出力用のヘルパー関数
	logVerbose := func(format string, args ...interface{}) {
		if options.Verbose {
			log.Printf("[VERBOSE] "+format, args...)
		}
	}

	// オプションのバリデーション
	if err := options.Validate(); err != nil {
		return nil, err
	}

	logVerbose("クリーンアップ処理を開始します")
	logVerbose("オプション: DryRun=%t, Verbose=%t, Yes=%t, Force=%t", options.DryRun, options.Verbose, options.Yes, options.Force)

	result := &CleanupResult{
		WasDryRun: options.DryRun,
	}

	// 1. Gitリポジトリかどうかの確認
	logVerbose("Gitリポジトリの確認を開始")
	cwd, err := GetCurrentDirectory()
	if err != nil {
		return nil, NewGitError("cleanup", err)
	}
	logVerbose("現在のディレクトリ: %s", cwd)

	if err := IsGitRepository(cwd); err != nil {
		return nil, NewGitError("cleanup", ErrNotGitRepository).WithPath(cwd)
	}
	logVerbose("Gitリポジトリであることを確認")

	// 2. デフォルトブランチの検出
	logVerbose("デフォルトブランチの検出を開始")
	var defaultBranch string
	if options.DefaultBranch != "" {
		defaultBranch = options.DefaultBranch
		logVerbose("手動指定されたデフォルトブランチ: %s", defaultBranch)
	} else {
		defaultBranch, err = DetectDefaultBranch()
		if err != nil {
			return nil, NewGitError("cleanup", err)
		}
		logVerbose("検出されたデフォルトブランチ: %s", defaultBranch)
	}
	result.DefaultBranch = defaultBranch

	// 3. デフォルトブランチへの切り替え
	logVerbose("現在のブランチ確認と切り替えを開始")
	currentBranch, err := GetCurrentBranch()
	if err != nil {
		return nil, NewGitError("cleanup", err)
	}
	logVerbose("現在のブランチ: %s", currentBranch)

	if currentBranch != defaultBranch {
		if options.DryRun {
			logVerbose("ドライランモード: ブランチ切り替えをシミュレーション (%s -> %s)", currentBranch, defaultBranch)
		} else {
			logVerbose("デフォルトブランチに切り替え: %s -> %s", currentBranch, defaultBranch)
			if err := CheckoutBranch(defaultBranch); err != nil {
				return nil, NewGitError("cleanup", err).WithMessage("failed to switch to default branch")
			}
			logVerbose("ブランチ切り替え完了")
		}
	} else {
		logVerbose("すでにデフォルトブランチにいます")
	}

	// 4. フェッチ処理（必須・ドライランでも実行）
	logVerbose("フェッチ処理を開始 (git fetch --all --prune)")
	if err := Fetch(); err != nil {
		// フェッチ失敗は警告として扱い、処理を継続
		logVerbose("フェッチエラー: %v", err)
		result.Errors = append(result.Errors, NewGitError("cleanup", err).WithMessage("fetch failed"))
	} else {
		logVerbose("フェッチ完了")
	}

	if options.DryRun {
		// ドライランモードの場合はfetch以外の実際の処理は行わない
		logVerbose("ドライランモードのため、フェッチ以外の処理をスキップ")
		return result, nil
	}

	// 5. プル処理（--no-pullが指定されていない場合）
	if !options.NoPull {
		logVerbose("プル処理を開始 (git pull)")
		if err := Pull(); err != nil {
			// プル失敗は警告として扱い、処理を継続
			logVerbose("プルエラー: %v", err)
			result.Errors = append(result.Errors, NewGitError("cleanup", err).WithMessage("pull failed"))
		} else {
			logVerbose("プル完了")
		}
	} else {
		logVerbose("プル処理をスキップ (--no-pull指定)")
	}

	// 6. ローカルブランチの一覧取得
	logVerbose("ローカルブランチ一覧を取得")
	branches, err := ListLocalBranches()
	if err != nil {
		return nil, NewGitError("cleanup", err)
	}
	logVerbose("検出されたブランチ: %v", branches)

	// 7. ブランチの削除
	logVerbose("ブランチ削除処理を開始")
	for _, branch := range branches {
		if branch == defaultBranch {
			// デフォルトブランチはスキップ
			logVerbose("デフォルトブランチをスキップ: %s", branch)
			result.SkippedBranches = append(result.SkippedBranches, branch)
			continue
		}

		// 除外パターンのチェック（簡単な実装）
		if options.ExcludePattern != "" && branch == options.ExcludePattern {
			logVerbose("除外パターンにマッチするためスキップ: %s", branch)
			result.SkippedBranches = append(result.SkippedBranches, branch)
			continue
		}

		// ブランチ削除
		logVerbose("ブランチ削除を試行: %s", branch)
		if err := DeleteBranch(branch, options.Force); err != nil {
			logVerbose("ブランチ削除エラー: %s - %v", branch, err)
			result.Errors = append(result.Errors, NewGitError("cleanup", err).WithPath(branch))
			result.SkippedBranches = append(result.SkippedBranches, branch)
		} else {
			logVerbose("ブランチ削除成功: %s", branch)
			result.DeletedBranches = append(result.DeletedBranches, branch)
		}
	}

	logVerbose("クリーンアップ処理完了 - 削除: %d, スキップ: %d, エラー: %d", len(result.DeletedBranches), len(result.SkippedBranches), len(result.Errors))

	return result, nil
}