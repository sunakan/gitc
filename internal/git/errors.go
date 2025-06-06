package git

import (
	"errors"
	"fmt"
)

// Common errors
var (
	ErrNotGitRepository     = errors.New("not a git repository")
	ErrNoDefaultBranch      = errors.New("could not detect default branch")
	ErrRemoteAccessFailed   = errors.New("failed to access remote repository")
	ErrMergeConflict        = errors.New("merge conflict detected")
	ErrBranchNotFound       = errors.New("branch not found")
	ErrCannotDeleteCurrent  = errors.New("cannot delete current branch")
)

// GitError represents a git-specific error with context
type GitError struct {
	Op      string // Operation that failed
	Path    string // Path where the error occurred
	Err     error  // Underlying error
	Message string // Additional context message
}

// Error implements the error interface
func (e *GitError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("git %s: %s: %v", e.Op, e.Message, e.Err)
	}
	if e.Path != "" {
		return fmt.Sprintf("git %s %s: %v", e.Op, e.Path, e.Err)
	}
	return fmt.Sprintf("git %s: %v", e.Op, e.Err)
}

// Unwrap returns the underlying error
func (e *GitError) Unwrap() error {
	return e.Err
}

// NewGitError creates a new GitError
func NewGitError(op string, err error) *GitError {
	return &GitError{
		Op:  op,
		Err: err,
	}
}

// WithPath adds a path to the GitError
func (e *GitError) WithPath(path string) *GitError {
	e.Path = path
	return e
}

// WithMessage adds a message to the GitError
func (e *GitError) WithMessage(msg string) *GitError {
	e.Message = msg
	return e
}

// IsNotGitRepository checks if the error indicates not a git repository
func IsNotGitRepository(err error) bool {
	return errors.Is(err, ErrNotGitRepository)
}

// IsNoDefaultBranch checks if the error indicates no default branch found
func IsNoDefaultBranch(err error) bool {
	return errors.Is(err, ErrNoDefaultBranch)
}

// IsRemoteAccessFailed checks if the error indicates remote access failure
func IsRemoteAccessFailed(err error) bool {
	return errors.Is(err, ErrRemoteAccessFailed)
}

// IsMergeConflict checks if the error indicates a merge conflict
func IsMergeConflict(err error) bool {
	return errors.Is(err, ErrMergeConflict)
}