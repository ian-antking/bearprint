package localprinter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ian-antking/bear-print/bearprint-api/localprinter"
	"github.com/ian-antking/bear-print/shared/printer"
)

func TestPrinter_Text(t *testing.T) {
	buf := &bytes.Buffer{}
	p := localprinter.NewPrinter(buf)

	err := p.Text("hello world", "left")
	assert.NoError(t, err)

	output := buf.String()

	expected := "hello world" + strings.Repeat(" ", 64-len("hello world")) + "\n"
	assert.Equal(t, expected, output)
}

func TestPrinter_PrintJob(t *testing.T) {
	buf := &bytes.Buffer{}
	p := localprinter.NewPrinter(buf)

	job := []printer.PrintItem{
		{Type: "text", Content: "line1", Align: "left"},
		{Type: "blank", Count: 1},
		{Type: "line"},
		{Type: "cut"},
	}

	err := p.PrintJob(job)
	assert.NoError(t, err)

	output := buf.Bytes()

	strOut := string(output)

	assert.Contains(t, strOut, "line1")
	assert.Contains(t, strOut, "\n\n")
	assert.Contains(t, strOut, strings.Repeat("-", 64))

	assert.Contains(t, strOut, string([]byte{0x1D, 0x56, 0x00}))
}
