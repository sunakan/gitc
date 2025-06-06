package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sunakan/gitc/internal/git"
)

var (
	// フラグ変数
	flagYes           bool
	flagVerbose       bool
	flagDryRun        bool
	flagForce         bool
	flagDefaultBranch string
	flagExclude       string
	flagNoPull        bool
	flagVersion       bool
)

// NewRootCmd は新しいrootコマンドを作成します（テスト用）
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gitc",
		Short: "Git repository cleanup tool",
		Long: `gitc is a CLI tool that automates Git repository cleanup.
It switches to the default branch, pulls the latest changes,
and removes unnecessary local branches.`,
		RunE: runGitCleanup,
	}

	// フラグの定義
	cmd.Flags().BoolVarP(&flagYes, "yes", "y", false, "確認プロンプトをスキップ")
	cmd.Flags().BoolVarP(&flagVerbose, "verbose", "v", false, "詳細な実行ログを表示")
	cmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "実際の処理は行わず、実行予定の内容のみ表示")
	cmd.Flags().BoolVar(&flagForce, "force", false, "強制実行（マージされていないブランチも削除）")
	cmd.Flags().StringVar(&flagDefaultBranch, "default-branch", "", "デフォルトブランチを手動指定")
	cmd.Flags().StringVar(&flagExclude, "exclude", "", "削除から除外するブランチのパターン指定")
	cmd.Flags().BoolVar(&flagNoPull, "no-pull", false, "プル処理をスキップ")
	cmd.Flags().BoolVar(&flagVersion, "version", false, "バージョン情報を表示")

	return cmd
}

var rootCmd = NewRootCmd()

// runGitCleanup はメインの処理を実行します
func runGitCleanup(cmd *cobra.Command, args []string) error {
	// バージョン表示
	if flagVersion {
		cmd.Println("gitc version 0.1.0")
		return nil
	}

	// フラグの検証
	flagMap := map[string]interface{}{
		"dry-run": flagDryRun,
		"force":   flagForce,
	}

	if err := validateFlags(flagMap); err != nil {
		return err
	}

	// クリーンアップオプションの設定
	options := git.CleanupOptions{
		DryRun:         flagDryRun,
		Verbose:        flagVerbose,
		Yes:            flagYes,
		Force:          flagForce,
		DefaultBranch:  flagDefaultBranch,
		ExcludePattern: flagExclude,
		NoPull:         flagNoPull,
	}

	if flagDryRun {
		cmd.Printf("🔍 ドライランモード: 実際の変更は行われません\n\n")
	}

	if flagVerbose {
		cmd.Printf("📋 実行オプション:\n")
		cmd.Printf("  - ドライラン: %t\n", options.DryRun)
		cmd.Printf("  - 詳細モード: %t\n", options.Verbose)
		cmd.Printf("  - 確認スキップ: %t\n", options.Yes)
		cmd.Printf("  - 強制実行: %t\n", options.Force)
		if options.DefaultBranch != "" {
			cmd.Printf("  - デフォルトブランチ: %s\n", options.DefaultBranch)
		}
		if options.ExcludePattern != "" {
			cmd.Printf("  - 除外パターン: %s\n", options.ExcludePattern)
		}
		cmd.Printf("  - プルスキップ: %t\n", options.NoPull)
		cmd.Printf("\n")
	}

	// クリーンアップ実行
	result, err := git.ExecuteCleanup(options)
	if err != nil {
		return fmt.Errorf("クリーンアップに失敗しました: %w", err)
	}

	// 結果の表示
	displayResult(cmd, result, options)
	return nil
}

// validateFlags はフラグの組み合わせをバリデーションします
func validateFlags(flags map[string]interface{}) error {
	dryRun, _ := flags["dry-run"].(bool)
	force, _ := flags["force"].(bool)

	if dryRun && force {
		return fmt.Errorf("--dry-run と --force は同時に使用できません")
	}

	return nil
}

// displayResult は実行結果を表示します
func displayResult(cmd *cobra.Command, result *git.CleanupResult, options git.CleanupOptions) {
	cmd.Printf("🔍 デフォルトブランチ: %s\n", result.DefaultBranch)

	if len(result.DeletedBranches) > 0 {
		cmd.Printf("\n🗑️  削除されたブランチ:\n")
		for _, branch := range result.DeletedBranches {
			cmd.Printf("  ✅ %s\n", branch)
		}
	}

	if len(result.SkippedBranches) > 0 {
		cmd.Printf("\n⏭️  スキップされたブランチ:\n")
		for _, branch := range result.SkippedBranches {
			cmd.Printf("  ⚠️  %s\n", branch)
		}
	}

	if len(result.Errors) > 0 {
		cmd.Printf("\n❌ エラー:\n")
		for _, err := range result.Errors {
			cmd.Printf("  ⚠️  %v\n", err)
		}
	}

	if !options.DryRun {
		cmd.Printf("\n🎉 クリーンアップ完了! %d個のブランチを削除しました。\n", len(result.DeletedBranches))
	} else {
		cmd.Printf("\n💡 ドライランモードで実行しました。実際の変更を行うには --dry-run フラグを外してください。\n")
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}