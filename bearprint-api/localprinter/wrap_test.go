package localprinter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapText(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		width int
		want  []string
	}{
		{
			name:  "basic wrap",
			text:  "The quick brown fox jumps over the lazy dog",
			width: 10,
			want: []string{
				"The quick",
				"brown fox",
				"jumps over",
				"the lazy",
				"dog",
			},
		},
		{
			name:  "single long word",
			text:  "supercalifragilisticexpialidocious",
			width: 10,
			want:  []string{"supercalifragilisticexpialidocious"},
		},
		{
			name:  "empty string",
			text:  "",
			width: 10,
			want:  []string{},
		},
		{
			name:  "words exactly width",
			text:  "abc def ghi",
			width: 3,
			want: []string{
				"abc",
				"def",
				"ghi",
			},
		},
		{
			name:  "width larger than text",
			text:  "short text",
			width: 50,
			want:  []string{"short text"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wrapText(tt.text, tt.width)
			assert.Equal(t, tt.want, got)
		})
	}
}
