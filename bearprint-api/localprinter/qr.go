package localprinter

import "github.com/ian-antking/bear-print/shared/printer"

func buildQRCodeCmd(data string, align printer.Alignment, lineWidth int) []byte {
	alignMap := map[printer.Alignment]byte{
		printer.AlignLeft:   0,
		printer.AlignCenter: 1,
		printer.AlignRight:  2,
	}

	alignVal, ok := alignMap[align]
	if !ok {
		alignVal = 0
	}

	alignCmd := []byte{0x1B, 0x61, alignVal}

	qrData := []byte(data)
	length := len(qrData) + 3
	lenLow := byte(length & 0xFF)
	lenHigh := byte((length >> 8) & 0xFF)

	cmds := []byte{}
	cmds = append(cmds, alignCmd...)
	cmds = append(cmds, []byte{
		0x1D, 0x28, 0x6B, 0x04, 0x00,
		0x31, 0x41, 0x32, 0x00,
	}...)
	cmds = append(cmds, []byte{
		0x1D, 0x28, 0x6B, 0x03, 0x00,
		0x31, 0x43, 0x06,
	}...)
	cmds = append(cmds, []byte{
		0x1D, 0x28, 0x6B, 0x03, 0x00,
		0x31, 0x45, 0x30,
	}...)
	cmds = append(cmds, []byte{
		0x1D, 0x28, 0x6B, lenLow, lenHigh,
		0x31, 0x50, 0x30,
	}...)
	cmds = append(cmds, qrData...)
	cmds = append(cmds, []byte{
		0x1D, 0x28, 0x6B, 0x03, 0x00,
		0x31, 0x51, 0x30, '\n',
	}...)

	return cmds
}
