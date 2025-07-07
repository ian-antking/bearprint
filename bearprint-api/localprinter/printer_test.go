package localprinter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ian-antking/bear-print/bearprint-api/localprinter"
	"github.com/ian-antking/bear-print/shared/printer"
)

type printerTest struct {
	t   *testing.T
	buf *bytes.Buffer
	p   *localprinter.Printer
}

func newPrinterTest(t *testing.T) *printerTest {
	buf := &bytes.Buffer{}
	p := localprinter.NewPrinter(buf)
	return &printerTest{t: t, buf: buf, p: p}
}

func (pt *printerTest) testText() {
	err := pt.p.Text("hello world", "left")
	assert.NoError(pt.t, err)

	output := pt.buf.String()
	expected := "hello world" + strings.Repeat(" ", 64-len("hello world")) + "\n"
	assert.Equal(pt.t, expected, output)
}

func (pt *printerTest) testPrintJob() {
	job := []printer.PrintItem{
		{Type: "text", Content: "line1", Align: "left"},
		{Type: "blank", Count: 1},
		{Type: "line"},
		{Type: "cut"},
	}

	err := pt.p.PrintJob(job)
	assert.NoError(pt.t, err)

	strOut := pt.buf.String()

	assert.Contains(pt.t, strOut, "line1")
	assert.Contains(pt.t, strOut, "\n\n")
	assert.Contains(pt.t, strOut, strings.Repeat("-", 64))
	assert.Contains(pt.t, strOut, string([]byte{0x1D, 0x56, 0x00}))
}

func TestPrinter(t *testing.T) {
	pt := newPrinterTest(t)

	t.Run("Text prints single line left aligned with padding", func(t *testing.T) {
		pt.testText()
	})

	t.Run("PrintJob executes job with text, blank line, line, and cut commands", func(t *testing.T) {
		pt.testPrintJob()
	})
}
