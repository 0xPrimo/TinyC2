// Package packer
package packer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unicode/utf16"
)

type WString string

func Pack(args ...any) ([]byte, error) {
	var buf bytes.Buffer

	for index, arg := range args {
		switch v := arg.(type) {

		case int:
			binary.Write(&buf, binary.LittleEndian, uint32(v))
		case int32:
			binary.Write(&buf, binary.LittleEndian, uint32(v))
		case uint32:
			binary.Write(&buf, binary.LittleEndian, v)
		case int16:
			binary.Write(&buf, binary.LittleEndian, uint16(v))
		case uint16:
			binary.Write(&buf, binary.LittleEndian, v)

		case []byte:
			binary.Write(&buf, binary.LittleEndian, uint32(len(v)))
			buf.Write(v)

		case string:
			bytearray := []byte(v)
			length := uint32(len(bytearray) + 1)

			binary.Write(&buf, binary.LittleEndian, length)
			buf.Write(bytearray)
			buf.WriteByte(0)

		case WString:
			utf16arr := utf16.Encode([]rune(v))
			length := uint32(len(utf16arr)*2 + 2)

			binary.Write(&buf, binary.LittleEndian, length)
			for _, r := range utf16arr {
				binary.Write(&buf, binary.LittleEndian, r)
			}
			binary.Write(&buf, binary.LittleEndian, uint16(0))
		default:
			return nil, fmt.Errorf("element at index %d has unsupported type: %T", index, v)
		}
	}

	packedbuf := buf.Bytes()
	total := uint32(len(packedbuf))
	data := make([]byte, 4+total)

	binary.LittleEndian.PutUint32(data[0:4], total)
	copy(data[4:], packedbuf)

	return data, nil
}
