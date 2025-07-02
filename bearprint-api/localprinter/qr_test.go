package localprinter

import (
	"testing"

	"github.com/ian-antking/bear-print/shared/printer"
	"github.com/stretchr/testify/assert"
)

func TestBuildQRCodeCmd(t *testing.T) {
	data := "Hello, QR!"
	align := printer.AlignCenter

	cmds := buildQRCodeCmd(data, align, 64)

	assert.NotEmpty(t, cmds, "Expected commands bytes, got empty slice")

	assert.Equal(t, byte(0x1B), cmds[0], "First byte should be 0x1B")
	assert.Equal(t, byte(0x61), cmds[1], "Second byte should be 0x61")
	assert.Equal(t, byte(0x01), cmds[2], "Third byte should be 0x01 for center alignment")

	foundData := false
	dataBytes := []byte(data)
	for i := 0; i <= len(cmds)-len(dataBytes); i++ {
		if string(cmds[i:i+len(dataBytes)]) == data {
			foundData = true
			break
		}
	}
	assert.True(t, foundData, "Data string not found in QR code command bytes")
}
