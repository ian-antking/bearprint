package localprinter

import (
	"io"
	"strings"
	"unicode/utf8"

	"github.com/ian-antking/bear-print/shared/printer"
)

// ESC/POS Command Constants
var (
	escAt        = []byte{0x1B, 0x40}       // Initialize/Reset printer
	escE_on      = []byte{0x1B, 0x45, 0x01} // Bold ON
	escE_off     = []byte{0x1B, 0x45, 0x00} // Bold OFF
	escMinus_on  = []byte{0x1B, 0x2D, 0x01} // Underline ON (1-dot thick)
	escMinus_off = []byte{0x1B, 0x2D, 0x00} // Underline OFF
	escCut       = []byte{0x1D, 0x56, 0x00} // Full cut
)

type Printer struct {
	encoding  string
	lineWidth int
	writer    io.Writer
}

func NewPrinter(writer io.Writer) *Printer {
	p := &Printer{
		writer:    writer,
		encoding:  "cp437",
		lineWidth: 64,
	}

	_ = p.write(escAt)
	return p
}

func (p *Printer) write(data []byte) error {
	_, err := p.writer.Write(data)
	return err
}

func (p *Printer) writeLine(line string) error {
	return p.write(append([]byte(line), '\n'))
}

func (p *Printer) setFontSize(doubleWidth, doubleHeight bool) error {
	var size byte = 0x00
	if doubleWidth {
		size |= 0x10 // Set bit 4 for double width
	}
	if doubleHeight {
		size |= 0x01 // Set bit 0 for double height
	}
	return p.write([]byte{0x1D, 0x21, size})
}

func (p *Printer) resetStyles() error {
	if err := p.write(escE_off); err != nil {
		return err
	}
	if err := p.write(escMinus_off); err != nil {
		return err
	}
	return p.setFontSize(false, false)
}

func (p *Printer) formatLine(line string, align printer.Alignment, lineWidth int) string {
	pad := lineWidth - utf8.RuneCountInString(line)
	if pad <= 0 {
		return line
	}
	switch align {
	case "center":
		left := pad / 2
		right := pad - left
		return strings.Repeat(" ", left) + line + strings.Repeat(" ", right)
	case "right":
		return strings.Repeat(" ", pad) + line
	default: // "left"
		return line + strings.Repeat(" ", pad)
	}
}

func (p *Printer) Text(item printer.PrintItem) error {
	if item.Bold {
		if err := p.write(escE_on); err != nil {
			return err
		}
	}
	if item.Underline {
		if err := p.write(escMinus_on); err != nil {
			return err
		}
	}
	if err := p.setFontSize(item.DoubleWidth, item.DoubleHeight); err != nil {
		return err
	}

	lineWidth := p.lineWidth
	if item.DoubleWidth {
		lineWidth /= 2
	}

	lines := strings.Split(item.Content, "\n")
	for _, rawLine := range lines {
		wrapped := wrapText(rawLine, lineWidth)
		if len(wrapped) == 0 {
			wrapped = []string{""}
		}
		for _, line := range wrapped {
			if err := p.writeLine(p.formatLine(line, item.Align, lineWidth)); err != nil {
				return err
			}
		}
	}

	return p.resetStyles()
}

func (p *Printer) BlankLine(count int) error {
	if count <= 0 {
		return nil
	}
	return p.write([]byte(strings.Repeat("\n", count)))
}

func (p *Printer) Cut() error {
	if err := p.BlankLine(6); err != nil {
		return err
	}
	return p.write(escCut)
}

func (p *Printer) printQRCode(data string, align printer.Alignment) error {
	return p.Text(printer.PrintItem{Content: "[QR Code: " + data + "]", Align: align})
}

func (p *Printer) PrintJob(items []printer.PrintItem) (err error) {
	if err = p.write(escAt); err != nil {
		return err
	}

	defer func() {
		resetErr := p.resetStyles()
		if err == nil {
			err = resetErr
		}
	}()

	for _, item := range items {
		switch item.Type {
		case "text":
			if err = p.Text(item); err != nil {
				return err
			}
		case "blank":
			if err = p.BlankLine(item.Count); err != nil {
				return err
			}
		case "line":
			lineItem := printer.PrintItem{Content: strings.Repeat("-", p.lineWidth)}
			if err = p.Text(lineItem); err != nil {
				return err
			}
		case "cut":
			if err = p.Cut(); err != nil {
				return err
			}
		case "qrcode":
			if err = p.printQRCode(item.Content, item.Align); err != nil {
				return err
			}
		}
	}
	return nil
}
