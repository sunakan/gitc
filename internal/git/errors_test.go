package git

import (
	"errors"
	"testing"
)

func TestGitError(t *testing.T) {
	t.Run("基本的なGitError", func(t *testing.T) {
		baseErr := errors.New("base error")
		gitErr := NewGitError("test-op", baseErr)

		expected := "git test-op: base error"
		if gitErr.Error() != expected {
			t.Errorf("GitError.Error() = %v, want %v", gitErr.Error(), expected)
		}

		if gitErr.Unwrap() != baseErr {
			t.Errorf("GitError.Unwrap() = %v, want %v", gitErr.Unwrap(), baseErr)
		}
	})

	t.Run("パス付きGitError", func(t *testing.T) {
		baseErr := errors.New("base error")
		gitErr := NewGitError("test-op", baseErr).WithPath("/test/path")

		expected := "git test-op /test/path: base error"
		if gitErr.Error() != expected {
			t.Errorf("GitError.Error() = %v, want %v", gitErr.Error(), expected)
		}
	})

	t.Run("メッセージ付きGitError", func(t *testing.T) {
		baseErr := errors.New("base error")
		gitErr := NewGitError("test-op", baseErr).WithMessage("additional context")

		expected := "git test-op: additional context: base error"
		if gitErr.Error() != expected {
			t.Errorf("GitError.Error() = %v, want %v", gitErr.Error(), expected)
		}
	})
}

func TestErrorCheckers(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		checkFn  func(error) bool
		expected bool
	}{
		{
			name:     "IsNotGitRepository - 正しいエラー",
			err:      NewGitError("test", ErrNotGitRepository),
			checkFn:  IsNotGitRepository,
			expected: true,
		},
		{
			name:     "IsNotGitRepository - 違うエラー",
			err:      NewGitError("test", ErrNoDefaultBranch),
			checkFn:  IsNotGitRepository,
			expected: false,
		},
		{
			name:     "IsNoDefaultBranch - 正しいエラー",
			err:      NewGitError("test", ErrNoDefaultBranch),
			checkFn:  IsNoDefaultBranch,
			expected: true,
		},
		{
			name:     "IsRemoteAccessFailed - 正しいエラー",
			err:      NewGitError("test", ErrRemoteAccessFailed),
			checkFn:  IsRemoteAccessFailed,
			expected: true,
		},
		{
			name:     "IsMergeConflict - 正しいエラー",
			err:      NewGitError("test", ErrMergeConflict),
			checkFn:  IsMergeConflict,
			expected: true,
		},
		{
			name:     "nil エラー",
			err:      nil,
			checkFn:  IsNotGitRepository,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.checkFn(tt.err); got != tt.expected {
				t.Errorf("%s() = %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}

func TestErrorConstants(t *testing.T) {
	// エラー定数が期待通りのメッセージを持つことを確認
	errorMessages := map[error]string{
		ErrNotGitRepository:   "not a git repository",
		ErrNoDefaultBranch:    "could not detect default branch",
		ErrRemoteAccessFailed: "failed to access remote repository",
		ErrMergeConflict:      "merge conflict detected",
		ErrBranchNotFound:     "branch not found",
		ErrCannotDeleteCurrent: "cannot delete current branch",
	}

	for err, expectedMsg := range errorMessages {
		if err.Error() != expectedMsg {
			t.Errorf("Error constant %v has message %q, want %q", err, err.Error(), expectedMsg)
		}
	}
}