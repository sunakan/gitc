package git

import (
	"fmt"
	"strings"
)

// DetectDefaultBranch detects the default branch of the repository
func DetectDefaultBranch() (string, error) {
	// Try to get the default branch from remote HEAD
	result, err := ExecuteCommand("symbolic-ref", "refs/remotes/origin/HEAD")
	if err == nil && result.Output != "" {
		// Extract branch name from refs/remotes/origin/main format
		parts := strings.Split(result.Output, "/")
		if len(parts) > 0 {
			return parts[len(parts)-1], nil
		}
	}
	
	// Fallback: check common default branch names
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
	
	// If still not found, check remote branches
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

// GetCurrentBranch returns the name of the current branch
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

// ListLocalBranches returns a list of all local branches
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

// ListRemoteBranches returns a list of all remote branches
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

// CheckoutBranch switches to the specified branch
func CheckoutBranch(branch string) error {
	_, err := ExecuteCommand("checkout", branch)
	if err != nil {
		return NewGitError("checkout", err).WithPath(branch)
	}
	return nil
}

// DeleteBranch deletes the specified local branch
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

// filterEmptyStrings removes empty strings from a slice
func filterEmptyStrings(strings []string) []string {
	var filtered []string
	for _, s := range strings {
		if s = strings.TrimSpace(s); s != "" {
			filtered = append(filtered, s)
		}
	}
	return filtered
}