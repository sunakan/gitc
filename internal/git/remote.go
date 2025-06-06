package git

import (
	"fmt"
	"strings"
	"time"
)

// Pull はgit pull操作を実行します
func Pull() error {
	result, err := ExecuteCommand("pull")
	if err != nil {
		// マージコンフリクトか確認
		if strings.Contains(result.Error, "conflict") || strings.Contains(result.Output, "conflict") {
			return NewGitError("pull", ErrMergeConflict)
		}
		return NewGitError("pull", err)
	}
	return nil
}

// PullWithRebase はrebaseを伴うgit pullを実行します
func PullWithRebase() error {
	result, err := ExecuteCommand("pull", "--rebase")
	if err != nil {
		// マージコンフリクトか確認
		if strings.Contains(result.Error, "conflict") || strings.Contains(result.Output, "conflict") {
			return NewGitError("pull", ErrMergeConflict).WithMessage("conflict during rebase")
		}
		return NewGitError("pull", err).WithMessage("rebase failed")
	}
	return nil
}

// Fetch はリモート参照を更新します
func Fetch() error {
	_, err := ExecuteCommand("fetch", "--all", "--prune")
	if err != nil {
		return NewGitError("fetch", err)
	}
	return nil
}

// CheckRemoteAccess はリモートリポジトリにアクセスできるか確認します
func CheckRemoteAccess() error {
	// リモートチェックのタイムアウトを設定
	result, err := ExecuteCommandWithTimeout(10*time.Second, "ls-remote", "--heads", "origin")
	if err != nil {
		return NewGitError("check-remote", ErrRemoteAccessFailed).WithMessage(fmt.Sprintf("failed to access remote: %v", err))
	}
	
	// 出力がない場合、リモートは空だがアクセス可能かもしれない
	if result.Output == "" {
		// リモートが存在するか確認するためURLを取得
		urlResult, urlErr := ExecuteCommand("remote", "get-url", "origin")
		if urlErr != nil {
			return NewGitError("check-remote", ErrRemoteAccessFailed).WithMessage("no remote 'origin' configured")
		}
		if urlResult.Output == "" {
			return NewGitError("check-remote", ErrRemoteAccessFailed).WithMessage("remote 'origin' has no URL")
		}
	}
	
	return nil
}

// ExecuteCommandWithTimeout はタイムアウト付きでGitコマンドを実行します
func ExecuteCommandWithTimeout(timeout time.Duration, args ...string) (*CommandResult, error) {
	type resultWrapper struct {
		result *CommandResult
		err    error
	}
	
	done := make(chan resultWrapper, 1)
	
	go func() {
		result, err := ExecuteCommand(args...)
		done <- resultWrapper{result: result, err: err}
	}()
	
	select {
	case <-time.After(timeout):
		return nil, fmt.Errorf("command timed out after %v", timeout)
	case wrapper := <-done:
		return wrapper.result, wrapper.err
	}
}

// HasRemote は指定された名前のリモートが存在するか確認します
func HasRemote(name string) (bool, error) {
	result, err := ExecuteCommand("remote")
	if err != nil {
		return false, NewGitError("check-remote", err)
	}
	
	remotes := strings.Split(result.Output, "\n")
	for _, remote := range remotes {
		if strings.TrimSpace(remote) == name {
			return true, nil
		}
	}
	
	return false, nil
}

// GetRemoteURL は指定されたリモートのURLを返します
func GetRemoteURL(name string) (string, error) {
	result, err := ExecuteCommand("remote", "get-url", name)
	if err != nil {
		return "", NewGitError("get-remote-url", err).WithPath(name)
	}
	return result.Output, nil
}