package git

import (
	"reflect"
	"testing"
)

func TestFilterEmptyStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "空文字列を含むスライス",
			input:    []string{"main", "", "develop", "  ", "feature"},
			expected: []string{"main", "develop", "feature"},
		},
		{
			name:     "空文字列のみのスライス",
			input:    []string{"", "  ", "\t", "\n"},
			expected: []string{},
		},
		{
			name:     "空文字列を含まないスライス",
			input:    []string{"main", "develop", "feature"},
			expected: []string{"main", "develop", "feature"},
		},
		{
			name:     "空のスライス",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "nilスライス",
			input:    nil,
			expected: []string{},
		},
		{
			name:     "前後に空白を含む文字列",
			input:    []string{"  main  ", "\tdevelop\t", " feature\n"},
			expected: []string{"main", "develop", "feature"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterEmptyStrings(tt.input)
			
			// nilと空スライスの比較を適切に行う
			if tt.expected == nil {
				if got != nil {
					t.Errorf("filterEmptyStrings() = %v, want nil", got)
				}
				return
			}
			
			if got == nil {
				got = []string{}
			}
			
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("filterEmptyStrings() = %v, want %v", got, tt.expected)
			}
		})
	}
}