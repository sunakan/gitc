package git

import (
	"fmt"
	"strings"
)

// DetectDefaultBranch はリポジトリのデフォルトブランチを検出します
func DetectDefaultBranch() (string, error) {
	// リモートHEADからデフォルトブランチを取得してみる
	result, err := ExecuteCommand("symbolic-ref", "refs/remotes/origin/HEAD")
	if err == nil && result.Output != "" {
		// refs/remotes/origin/main形式からブランチ名を抽出
		parts := strings.Split(result.Output, "/")
		if len(parts) > 0 {
			return parts[len(parts)-1], nil
		}
	}
	
	// フォールバック: 一般的なデフォルトブランチ名を確認
	commonDefaults := []string{"main", "master", "develop", "dev"}
	branches, err := ListLocalBranches()
	if err != nil {
		return "", NewGitError("detect-default-branch", err).WithMessage("failed to list branches")
	}
	
	for _, defaultName := range commonDefaults {
		for _, branch := range branches {
			if branch == defaultName {
				return branch, nil
			}
		}
	}
	
	// まだ見つからない場合は、リモートブランチを確認
	remoteBranches, err := ListRemoteBranches()
	if err == nil {
		for _, defaultName := range commonDefaults {
			for _, branch := range remoteBranches {
				if strings.HasSuffix(branch, "/"+defaultName) {
					return defaultName, nil
				}
			}
		}
	}
	
	return "", NewGitError("detect-default-branch", ErrNoDefaultBranch)
}

// GetCurrentBranch は現在のブランチ名を返します
func GetCurrentBranch() (string, error) {
	result, err := ExecuteCommand("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", NewGitError("get-current-branch", err)
	}
	
	if result.Output == "" {
		return "", NewGitError("get-current-branch", fmt.Errorf("no output from git command"))
	}
	
	return result.Output, nil
}

// ListLocalBranches はすべてのローカルブランチの一覧を返します
func ListLocalBranches() ([]string, error) {
	result, err := ExecuteCommand("branch", "--format=%(refname:short)")
	if err != nil {
		return nil, NewGitError("list-local-branches", err)
	}
	
	if result.Output == "" {
		return []string{}, nil
	}
	
	branches := strings.Split(result.Output, "\n")
	return filterEmptyStrings(branches), nil
}

// ListRemoteBranches はすべてのリモートブランチの一覧を返します
func ListRemoteBranches() ([]string, error) {
	result, err := ExecuteCommand("branch", "-r", "--format=%(refname:short)")
	if err != nil {
		return nil, NewGitError("list-remote-branches", err)
	}
	
	if result.Output == "" {
		return []string{}, nil
	}
	
	branches := strings.Split(result.Output, "\n")
	return filterEmptyStrings(branches), nil
}

// CheckoutBranch は指定されたブランチに切り替えます
func CheckoutBranch(branch string) error {
	_, err := ExecuteCommand("checkout", branch)
	if err != nil {
		return NewGitError("checkout", err).WithPath(branch)
	}
	return nil
}

// DeleteBranch は指定されたローカルブランチを削除します
func DeleteBranch(branch string, force bool) error {
	args := []string{"branch", "-d", branch}
	if force {
		args[1] = "-D"
	}
	
	_, err := ExecuteCommand(args...)
	if err != nil {
		return NewGitError("delete-branch", err).WithPath(branch)
	}
	return nil
}

// BranchExists は指定されたブランチが存在するかチェックします
func BranchExists(branch string) (bool, error) {
	// ローカルブランチをチェック
	localBranches, err := ListLocalBranches()
	if err != nil {
		return false, NewGitError("check-branch-exists", err).WithMessage("failed to list local branches")
	}

	for _, localBranch := range localBranches {
		if localBranch == branch {
			return true, nil
		}
	}

	// リモートブランチもチェック
	remoteBranches, err := ListRemoteBranches()
	if err != nil {
		// リモートブランチの取得に失敗した場合は警告として扱い、ローカルのみの結果を返す
		return false, nil
	}

	for _, remoteBranch := range remoteBranches {
		// origin/branch 形式のリモートブランチ名から branch 部分を抽出
		if strings.HasSuffix(remoteBranch, "/"+branch) {
			return true, nil
		}
	}

	return false, nil
}

// filterEmptyStrings はスライスから空文字列を除去します
func filterEmptyStrings(strs []string) []string {
	var filtered []string
	for _, s := range strs {
		if s = strings.TrimSpace(s); s != "" {
			filtered = append(filtered, s)
		}
	}
	return filtered
}