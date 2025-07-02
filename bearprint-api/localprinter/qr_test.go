package localprinter

import (
	"testing"

	"github.com/ian-antking/bear-print/shared/printer"
	"github.com/stretchr/testify/assert"
)

type qrCodeTest struct {
	t *testing.T
}

func newQRCodeTest(t *testing.T) *qrCodeTest {
	return &qrCodeTest{t: t}
}

func (qt *qrCodeTest) testBuildQRCodeCmd() {
	data := "Hello, QR!"
	align := printer.AlignCenter

	cmds := buildQRCodeCmd(data, align, 64)

	assert.NotEmpty(qt.t, cmds, "Expected commands bytes, got empty slice")

	assert.Equal(qt.t, byte(0x1B), cmds[0], "First byte should be 0x1B")
	assert.Equal(qt.t, byte(0x61), cmds[1], "Second byte should be 0x61")
	assert.Equal(qt.t, byte(0x01), cmds[2], "Third byte should be 0x01 for center alignment")

	foundData := false
	dataBytes := []byte(data)
	for i := 0; i <= len(cmds)-len(dataBytes); i++ {
		if string(cmds[i:i+len(dataBytes)]) == data {
			foundData = true
			break
		}
	}
	assert.True(qt.t, foundData, "Data string not found in QR code command bytes")
}

func TestBuildQRCodeCmd(t *testing.T) {
	qt := newQRCodeTest(t)
	t.Run("buildQRCodeCmd produces correct command bytes including alignment and data", func(t *testing.T) {
		qt.testBuildQRCodeCmd()
	})
}
