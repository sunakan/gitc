package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootCmdExecute(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		expectOut   string
	}{
		{
			name:      "ヘルプ表示",
			args:      []string{"--help"},
			wantErr:   false,
			expectOut: "gitc is a CLI tool that automates Git repository cleanup",
		},
		{
			name:      "バージョン表示",
			args:      []string{"--version"},
			wantErr:   false,
			expectOut: "gitc version 0.1.0",
		},
		{
			name:      "ドライランモード（Gitリポジトリ外）",
			args:      []string{"--dry-run"},
			wantErr:   true,
			expectOut: "not a git repository",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewRootCmd()
			
			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			output := buf.String()
			if !strings.Contains(output, tt.expectOut) {
				t.Errorf("Expected output to contain %q, got %q", tt.expectOut, output)
			}
		})
	}
}

func TestRootCmdFlags(t *testing.T) {
	cmd := NewRootCmd()

	flags := []string{"yes", "verbose", "dry-run", "version", "force", "default-branch", "exclude", "no-pull"}
	
	for _, flagName := range flags {
		t.Run(flagName+" フラグ", func(t *testing.T) {
			flag := cmd.Flags().Lookup(flagName)
			if flag == nil {
				t.Errorf("Flag %s should be defined", flagName)
			}
		})
	}
}

func TestValidateFlags(t *testing.T) {
	tests := []struct {
		name    string
		flags   map[string]interface{}
		wantErr bool
	}{
		{
			name: "有効なフラグの組み合わせ",
			flags: map[string]interface{}{
				"dry-run": false,
				"force":   false,
			},
			wantErr: false,
		},
		{
			name: "無効なフラグの組み合わせ",
			flags: map[string]interface{}{
				"dry-run": true,
				"force":   true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFlags(tt.flags)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFlags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}