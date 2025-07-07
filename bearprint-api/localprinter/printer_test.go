package localprinter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ian-antking/bear-print/bearprint-api/localprinter"
	"github.com/ian-antking/bear-print/shared/printer"
	"github.com/stretchr/testify/assert"
)

// Define ESC/POS constants here to create expected byte sequences for tests.
var (
	escAt         = []byte{0x1B, 0x40}       // Initialize/Reset printer
	escE_on       = []byte{0x1B, 0x45, 0x01} // Bold ON
	escE_off      = []byte{0x1B, 0x45, 0x00} // Bold OFF
	escMinus_on   = []byte{0x1B, 0x2D, 0x01} // Underline ON
	escMinus_off  = []byte{0x1B, 0x2D, 0x00} // Underline OFF
	gsExcl_normal = []byte{0x1D, 0x21, 0x00} // Normal size (also used to turn off size changes)
	gsExcl_dw_dh  = []byte{0x1D, 0x21, 0x11} // Double Width & Double Height ON
	escCut        = []byte{0x1D, 0x56, 0x00} // Full cut
)

// newTestPrinter is a helper to set up the printer and buffer for each test.
func newTestPrinter() (*bytes.Buffer, *localprinter.Printer) {
	buf := &bytes.Buffer{}
	p := localprinter.NewPrinter(buf)
	return buf, p
}

func TestText_Simple(t *testing.T) {
	buf, p := newTestPrinter()
	buf.Reset() // Reset buffer to ignore the initial `escAt` command from NewPrinter.

	item := printer.PrintItem{
		Type:    printer.Text,
		Content: "hello",
		Align:   printer.AlignLeft,
	}

	err := p.Text(item)
	assert.NoError(t, err)

	// Build the expected sequence of bytes
	expected := &bytes.Buffer{}
	// CORRECTED: Expect the font size command at the start, as Text() always calls setFontSize().
	expected.Write(gsExcl_normal)
	expected.WriteString("hello" + strings.Repeat(" ", 64-len("hello")) + "\n")
	expected.Write(escE_off)      // Style reset
	expected.Write(escMinus_off)  // Style reset
	expected.Write(gsExcl_normal) // Style reset

	assert.Equal(t, expected.String(), buf.String(), "Should print simple text and reset styles")
}

func TestText_Styled(t *testing.T) {
	buf, p := newTestPrinter()
	buf.Reset()

	item := printer.PrintItem{
		Type:      printer.Text,
		Content:   "bold text",
		Align:     printer.AlignCenter,
		Bold:      true,
		Underline: true,
	}

	err := p.Text(item)
	assert.NoError(t, err)

	expected := &bytes.Buffer{}
	// 1. Set styles
	expected.Write(escE_on)
	expected.Write(escMinus_on)
	expected.Write(gsExcl_normal) // Font size is set even if default

	// 2. Write formatted line
	line := "bold text"
	padding := 64 - len(line)
	leftPad := strings.Repeat(" ", padding/2)
	rightPad := strings.Repeat(" ", padding-(padding/2))
	expected.WriteString(leftPad + line + rightPad + "\n")

	// 3. Reset styles
	expected.Write(escE_off)
	expected.Write(escMinus_off)
	expected.Write(gsExcl_normal)

	assert.Equal(t, expected.Bytes(), buf.Bytes(), "Should wrap text with style on/off commands")
}

func TestText_DoubleWidth(t *testing.T) {
	buf, p := newTestPrinter()
	buf.Reset()

	item := printer.PrintItem{
		Type:         printer.Text,
		Content:      "wide",
		Align:        printer.AlignRight,
		DoubleWidth:  true,
		DoubleHeight: true,
	}

	err := p.Text(item)
	assert.NoError(t, err)

	expected := &bytes.Buffer{}
	// 1. Set style for double width and height
	expected.Write(gsExcl_dw_dh)

	// 2. Write formatted line (note: line width is halved to 32)
	line := "wide"
	lineWidth := 32 // 64 / 2
	padding := lineWidth - len(line)
	expected.WriteString(strings.Repeat(" ", padding) + line + "\n")

	// 3. Reset styles
	expected.Write(escE_off)
	expected.Write(escMinus_off)
	expected.Write(gsExcl_normal)

	assert.Equal(t, expected.Bytes(), buf.Bytes(), "Should handle halved line width for double-width text")
}

func TestPrintJob(t *testing.T) {
	buf, p := newTestPrinter()
	// In a full job, we test the entire buffer, including the initial `escAt`

	job := []printer.PrintItem{
		{Type: printer.Text, Content: "Title", Align: printer.AlignCenter, Bold: true, DoubleHeight: true},
		{Type: printer.Blank, Count: 1},
		{Type: printer.Text, Content: "Some normal text.", Align: printer.AlignLeft},
		{Type: printer.Line},
		{Type: printer.Cut},
	}

	err := p.PrintJob(job)
	assert.NoError(t, err)

	output := buf.Bytes()

	// Build the full expected output piece by piece
	expected := &bytes.Buffer{}
	// Initial reset from NewPrinter and start of job
	expected.Write(escAt)
	expected.Write(escAt)

	// Item 1: Styled Title
	expected.Write(escE_on)                               // Bold on
	expected.Write([]byte{0x1D, 0x21, 0x01})              // Double Height on
	titleLine := "Title"
	padding := 64 - len(titleLine)
	expected.WriteString(strings.Repeat(" ", padding/2) + titleLine + strings.Repeat(" ", padding-(padding/2)) + "\n")
	expected.Write(escE_off)
	expected.Write(escMinus_off)
	expected.Write(gsExcl_normal)

	// Item 2: Blank line
	expected.WriteString("\n")

	// Item 3: Normal text
	// CORRECTED: Added the expected font size command for the normal text block.
	expected.Write(gsExcl_normal)
	normalLine := "Some normal text."
	expected.WriteString(normalLine + strings.Repeat(" ", 64-len(normalLine)) + "\n")
	expected.Write(escE_off)
	expected.Write(escMinus_off)
	expected.Write(gsExcl_normal)

	// Item 4: Horizontal line (which is just a text item)
	// CORRECTED: Added the expected font size command for the line block.
	expected.Write(gsExcl_normal)
	line := strings.Repeat("-", 64)
	expected.WriteString(line + "\n")
	expected.Write(escE_off)
	expected.Write(escMinus_off)
	expected.Write(gsExcl_normal)

	// Item 5: Cut
	expected.WriteString(strings.Repeat("\n", 6))
	expected.Write(escCut)

	// Final reset from the 'defer' in PrintJob
	expected.Write(escE_off)
	expected.Write(escMinus_off)
	expected.Write(gsExcl_normal)

	// Using string representation for easier diffs in test failures
	assert.Equal(t, expected.String(), string(output))
}