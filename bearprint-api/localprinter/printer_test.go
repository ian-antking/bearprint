package localprinter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ian-antking/bearprint/bearprint-api/localprinter"
	"github.com/ian-antking/bearprint/shared/printer"
	"github.com/stretchr/testify/assert"
)

const (
	lineWidth = 48
)

func newPrinterTest() (*bytes.Buffer, *localprinter.Printer) {
	buf := &bytes.Buffer{}
	p := localprinter.NewPrinter(buf)
	return buf, p
}

func TestTextMethod(t *testing.T) {
	buf, p := newPrinterTest()

	item := printer.PrintItem{Content: "hello world", Align: "left"}
	err := p.Text(item)
	assert.NoError(t, err)

	expected := "hello world" + strings.Repeat(" ", lineWidth-len("hello world")) + "\n"
	assert.Equal(t, expected, buf.String())
}

func TestPrintJob(t *testing.T) {
	buf, p := newPrinterTest()

	job := []printer.PrintItem{
		{Type: "text", Content: "line1", Align: "left"},
		{Type: "blank", Count: 1},
		{Type: "line"},
		{Type: "cut"},
	}

	err := p.PrintJob(job)
	assert.NoError(t, err)

	expected := &bytes.Buffer{}
	// 1. Text item: "line1"
	expected.WriteString("line1" + strings.Repeat(" ", lineWidth-len("line1")) + "\n")
	// 2. Blank item
	expected.WriteString("\n")
	// 3. Line item
	expected.WriteString(strings.Repeat("-", lineWidth) + "\n")
	// 4. Cut item
	expected.WriteString(strings.Repeat("\n", 6))
	expected.Write([]byte{0x1D, 0x56, 0x00})

	assert.Equal(t, expected.String(), buf.String())
}

func TestTextNormalization(t *testing.T) {
	buf, p := newPrinterTest()

	item := printer.PrintItem{
		Content: "Temp: 23°C\nFraction: ½\nDash: —\nQuote: “hello”\nApostrophe: ’",
		Align:   "left",
	}
	err := p.Text(item)
	assert.NoError(t, err)

	expectedStr := "Temp: 23degC" + strings.Repeat(" ", lineWidth-len("Temp: 23degC")) + "\n" +
		"Fraction: 1/2" + strings.Repeat(" ", lineWidth-len("Fraction: 1/2")) + "\n" +
		"Dash: --" + strings.Repeat(" ", lineWidth-len("Dash: --")) + "\n" +
		"Quote: \"hello\"" + strings.Repeat(" ", lineWidth-len("Quote: \"hello\"")) + "\n" +
		"Apostrophe: '" + strings.Repeat(" ", lineWidth-len("Apostrophe: '")) + "\n"


	assert.Equal(t, expectedStr, buf.String())
}
