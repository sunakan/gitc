package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootCmd(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantErr   bool
		wantOut   string
	}{
		{
			name:    "ヘルプ表示",
			args:    []string{"--help"},
			wantErr: false,
			wantOut: "gitc is a CLI tool",
		},
		{
			name:    "ドライランモード（Gitリポジトリ外）",
			args:    []string{"--dry-run"},
			wantErr: true,
			wantOut: "not a git repository",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newRootCmd()
			buf := bytes.Buffer{}
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			output := buf.String()
			if !strings.Contains(output, tt.wantOut) {
				t.Errorf("Expected output to contain %q, got %q", tt.wantOut, output)
			}
		})
	}
}