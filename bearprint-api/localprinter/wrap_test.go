package localprinter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type wrapTextTest struct {
	t *testing.T
}

func newWrapTextTest(t *testing.T) *wrapTextTest {
	return &wrapTextTest{t: t}
}

func (w *wrapTextTest) runTests() {
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
		// Use subtests for clarity/reporting
		w.t.Run(tt.name, func(t *testing.T) {
			got := wrapText(tt.text, tt.width) // assuming wrapText is exported as WrapText
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWrapText(t *testing.T) {
	wt := newWrapTextTest(t)
	wt.runTests()
}
