package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// CommandResult represents the result of a git command execution
type CommandResult struct {
	Output string
	Error  string
}

// ExecuteCommand executes a git command and returns the result
func ExecuteCommand(args ...string) (*CommandResult, error) {
	cmd := exec.Command("git", args...)
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	
	result := &CommandResult{
		Output: strings.TrimSpace(stdout.String()),
		Error:  strings.TrimSpace(stderr.String()),
	}
	
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return result, fmt.Errorf("git command failed with exit code %d: %s", exitErr.ExitCode(), result.Error)
		}
		return result, fmt.Errorf("failed to execute git command: %w", err)
	}
	
	return result, nil
}

// ExecuteCommandWithInput executes a git command with input and returns the result
func ExecuteCommandWithInput(input string, args ...string) (*CommandResult, error) {
	cmd := exec.Command("git", args...)
	cmd.Stdin = strings.NewReader(input)
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	
	result := &CommandResult{
		Output: strings.TrimSpace(stdout.String()),
		Error:  strings.TrimSpace(stderr.String()),
	}
	
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return result, fmt.Errorf("git command failed with exit code %d: %s", exitErr.ExitCode(), result.Error)
		}
		return result, fmt.Errorf("failed to execute git command: %w", err)
	}
	
	return result, nil
}