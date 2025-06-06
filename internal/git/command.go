package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// CommandResult はGitコマンドの実行結果を表します
type CommandResult struct {
	Output string
	Error  string
}

// ExecuteCommand はGitコマンドを実行し、結果を返します
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

// ExecuteCommandWithInput は入力を伴うGitコマンドを実行し、結果を返します
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