package logger

import (
	"github.com/impl0x/go-utils/logger/color"
	"github.com/impl0x/go-utils/logger/utils"
	"bytes"
)

const colorStartStr = "\033[38;2;"
const colorEndStr = "\033[0m"

func writeColorStart(b *bytes.Buffer, c color.RGB) {
	b.WriteString(colorStartStr)
	utils.WriteInteger(b, c.Red)
	b.WriteByte(';')
	utils.WriteInteger(b, c.Green)
	b.WriteByte(';')
	utils.WriteInteger(b, c.Blue)
	b.WriteByte('m')
}

func writeColorEnd(b *bytes.Buffer) {
	b.WriteString(colorEndStr)
}
