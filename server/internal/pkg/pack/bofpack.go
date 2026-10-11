package pack

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"unicode/utf16"
)

func BofPack(format string, args []string) ([]byte, error) {
	var buf bytes.Buffer
	index := 0

	for _, f := range format {
		if index >= len(args) {
			return nil, fmt.Errorf("not enough arguments provided for format string '%s'", format)
		}

		arg := args[index]
		switch f {
		case 'i':
			val, err := strconv.ParseInt(arg, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("arg %d ('i') failed to parse as decimal integer: %v", index, err)
			}
			binary.Write(&buf, binary.LittleEndian, uint32(val))

		case 's':
			val, err := strconv.ParseInt(arg, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("arg %d ('s') failed to parse as decimal short: %v", index, err)
			}
			binary.Write(&buf, binary.LittleEndian, uint16(val))

		case 'b':
			data, err := hex.DecodeString(arg)
			if err != nil {
				data = []byte(arg)
			}

			binary.Write(&buf, binary.LittleEndian, uint32(len(data)))
			buf.Write(data)
		case 'z':
			bytearray := []byte(arg)
			length := uint32(len(bytearray) + 1)

			binary.Write(&buf, binary.LittleEndian, length)
			buf.Write(bytearray)
			buf.WriteByte(0)
		case 'Z':
			utf16arr := utf16.Encode([]rune(arg))
			length := uint32(len(utf16arr)*2 + 2)

			binary.Write(&buf, binary.LittleEndian, length)
			for _, r := range utf16arr {
				binary.Write(&buf, binary.LittleEndian, r)
			}
			binary.Write(&buf, binary.LittleEndian, uint16(0))

		default:
			return nil, fmt.Errorf("unknown format character: '%c'", f)
		}

		index++
	}

	packedbuf := buf.Bytes()
	total := uint32(len(packedbuf))
	data := make([]byte, 4+total)
	binary.LittleEndian.PutUint32(data[0:4], total)
	copy(data[4:], packedbuf)
	return data, nil
}
