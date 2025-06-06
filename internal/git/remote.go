package git

import (
	"fmt"
	"strings"
	"time"
)

// Pull performs a git pull operation
func Pull() error {
	result, err := ExecuteCommand("pull")
	if err != nil {
		// Check if it's a merge conflict
		if strings.Contains(result.Error, "conflict") || strings.Contains(result.Output, "conflict") {
			return NewGitError("pull", ErrMergeConflict)
		}
		return NewGitError("pull", err)
	}
	return nil
}

// PullWithRebase performs a git pull with rebase
func PullWithRebase() error {
	result, err := ExecuteCommand("pull", "--rebase")
	if err != nil {
		// Check if it's a merge conflict
		if strings.Contains(result.Error, "conflict") || strings.Contains(result.Output, "conflict") {
			return NewGitError("pull", ErrMergeConflict).WithMessage("conflict during rebase")
		}
		return NewGitError("pull", err).WithMessage("rebase failed")
	}
	return nil
}

// Fetch updates remote references
func Fetch() error {
	_, err := ExecuteCommand("fetch", "--all", "--prune")
	if err != nil {
		return NewGitError("fetch", err)
	}
	return nil
}

// CheckRemoteAccess verifies that we can access the remote repository
func CheckRemoteAccess() error {
	// Set a timeout for the remote check
	result, err := ExecuteCommandWithTimeout(10*time.Second, "ls-remote", "--heads", "origin")
	if err != nil {
		return NewGitError("check-remote", ErrRemoteAccessFailed).WithMessage(fmt.Sprintf("failed to access remote: %v", err))
	}
	
	// If we get here but no output, the remote might be empty but accessible
	if result.Output == "" {
		// Try to get remote URL to ensure remote exists
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

// ExecuteCommandWithTimeout executes a git command with a timeout
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

// HasRemote checks if a remote with the given name exists
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

// GetRemoteURL returns the URL of the specified remote
func GetRemoteURL(name string) (string, error) {
	result, err := ExecuteCommand("remote", "get-url", name)
	if err != nil {
		return "", NewGitError("get-remote-url", err).WithPath(name)
	}
	return result.Output, nil
}