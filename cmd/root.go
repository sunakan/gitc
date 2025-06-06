package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sunakan/gitc/internal/git"
)

var (
	// フラグ変数
	flagDryRun bool
	flagYes    bool
	flagVerbose bool
)

// newRootCmd creates a new root command
func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gitc",
		Short: "Git repository cleanup tool",
		Long: `gitc is a CLI tool that automates Git repository cleanup.
It switches to the default branch, pulls the latest changes,
and removes unnecessary local branches.`,
		RunE: runCleanup,
	}

	// フラグの定義
	cmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "Perform a dry run without making actual changes")
	cmd.Flags().BoolVarP(&flagYes, "yes", "y", false, "Skip confirmation prompts")
	cmd.Flags().BoolVarP(&flagVerbose, "verbose", "v", false, "Show detailed logs")

	return cmd
}

var rootCmd = newRootCmd()

func runCleanup(cmd *cobra.Command, args []string) error {
	// ドライランモードの表示
	if flagDryRun {
		cmd.Println("🔍 Dry-run mode: No actual changes will be made")
		cmd.Println()
	}

	// クリーンアップオプションの設定
	options := git.CleanupOptions{
		DryRun:  flagDryRun,
		Verbose: flagVerbose,
		Yes:     flagYes,
		NoPull:  true, // 最小実装ではプルをスキップ
	}

	// クリーンアップ実行
	result, err := git.ExecuteCleanup(options)
	if err != nil {
		return fmt.Errorf("cleanup failed: %w", err)
	}

	// 結果の表示（最小限）
	cmd.Printf("Default branch: %s\n", result.DefaultBranch)
	
	if len(result.DeletedBranches) > 0 {
		cmd.Println("\nDeleted branches:")
		for _, branch := range result.DeletedBranches {
			cmd.Printf("  - %s\n", branch)
		}
	}

	if flagDryRun {
		cmd.Println("\n✨ Dry-run completed. Run without --dry-run to perform actual cleanup.")
	} else {
		cmd.Printf("\n✨ Cleanup completed! Deleted %d branches.\n", len(result.DeletedBranches))
	}

	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}