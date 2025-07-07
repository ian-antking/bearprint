package localprinter

import (
	"io"
	"strings"
	"unicode/utf8"

	"github.com/ian-antking/bear-print/shared/printer"
)

type Printer struct {
	encoding  string
	lineWidth int
	writer    io.Writer
}

func NewPrinter(writer io.Writer) *Printer {
	return &Printer{
		writer:    writer,
		encoding:  "cp437",
		lineWidth: 48,
	}
}

func (p *Printer) write(data []byte) error {
	_, err := p.writer.Write(data)
	return err
}

func (p *Printer) writeLine(line string) error {
	return p.write(append([]byte(line), '\n'))
}

func (p *Printer) formatLine(line string, align printer.Alignment) string {
	pad := p.lineWidth - utf8.RuneCountInString(line)
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
	default:
		return line + strings.Repeat(" ", pad)
	}
}

func (p *Printer) Text(item printer.PrintItem) error {
  lines := strings.Split(item.Content, "\n")
  for _, rawLine := range lines {
    wrapped := wrapText(rawLine, p.lineWidth)
    if len(wrapped) == 0 {
      wrapped = []string{""}
    }
    for _, line := range wrapped {
      err := p.writeLine(p.formatLine(line, item.Align))
      if err != nil {
        return err
      }
    }
  }
  return nil
}

func (p *Printer) BlankLine(count int) error {
	return p.write([]byte(strings.Repeat("\n", count)))
}

func (p *Printer) Cut() error {
	err := p.BlankLine(6)
	if err != nil {
		return err
	}
	return p.write([]byte{0x1D, 0x56, 0x00})
}

func (p *Printer) printQRCode(data string, align printer.Alignment) error {
	cmds := buildQRCodeCmd(data, align, p.lineWidth)
	return p.write(cmds)
}

func (p *Printer) PrintJob(items []printer.PrintItem) error {
  for _, item := range items {
    switch item.Type {
    case "text":
      if err := p.Text(item); err != nil {
        return err
      }
    case "blank":
      if err := p.BlankLine(item.Count); err != nil {
        return err
      }
    case "line":
      lineItem := printer.PrintItem{Content: strings.Repeat("-", p.lineWidth), Align: "left"}
      if err := p.Text(lineItem); err != nil {
        return err
      }
    case "cut":
      if err := p.Cut(); err != nil {
        return err
      }
    case "qrcode":
      if err := p.printQRCode(item.Content, item.Align); err != nil {
        return err
      }
    }
  }
  return nil
}
