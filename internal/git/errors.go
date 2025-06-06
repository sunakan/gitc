package git

import (
	"errors"
	"fmt"
)

// 一般的なエラー
var (
	ErrNotGitRepository     = errors.New("not a git repository")
	ErrNoDefaultBranch      = errors.New("could not detect default branch")
	ErrRemoteAccessFailed   = errors.New("failed to access remote repository")
	ErrMergeConflict        = errors.New("merge conflict detected")
	ErrBranchNotFound       = errors.New("branch not found")
	ErrCannotDeleteCurrent  = errors.New("cannot delete current branch")
)

// GitError はGit固有のエラーとコンテキストを表します
type GitError struct {
	Op      string // 失敗した操作
	Path    string // エラーが発生したパス
	Err     error  // 内部エラー
	Message string // 追加のコンテキストメッセージ
}

// Error はerrorインターフェースを実装します
func (e *GitError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("git %s: %s: %v", e.Op, e.Message, e.Err)
	}
	if e.Path != "" {
		return fmt.Sprintf("git %s %s: %v", e.Op, e.Path, e.Err)
	}
	return fmt.Sprintf("git %s: %v", e.Op, e.Err)
}

// Unwrap は内部エラーを返します
func (e *GitError) Unwrap() error {
	return e.Err
}

// NewGitError は新しいGitErrorを作成します
func NewGitError(op string, err error) *GitError {
	return &GitError{
		Op:  op,
		Err: err,
	}
}

// WithPath はGitErrorにパスを追加します
func (e *GitError) WithPath(path string) *GitError {
	e.Path = path
	return e
}

// WithMessage はGitErrorにメッセージを追加します
func (e *GitError) WithMessage(msg string) *GitError {
	e.Message = msg
	return e
}

// IsNotGitRepository はエラーがGitリポジトリではないことを示しているか確認します
func IsNotGitRepository(err error) bool {
	return errors.Is(err, ErrNotGitRepository)
}

// IsNoDefaultBranch はエラーがデフォルトブランチが見つからないことを示しているか確認します
func IsNoDefaultBranch(err error) bool {
	return errors.Is(err, ErrNoDefaultBranch)
}

// IsRemoteAccessFailed はエラーがリモートアクセス失敗を示しているか確認します
func IsRemoteAccessFailed(err error) bool {
	return errors.Is(err, ErrRemoteAccessFailed)
}

// IsMergeConflict はエラーがマージコンフリクトを示しているか確認します
func IsMergeConflict(err error) bool {
	return errors.Is(err, ErrMergeConflict)
}