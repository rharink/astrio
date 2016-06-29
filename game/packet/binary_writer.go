package packet

import (
	"bytes"
	"encoding/binary"
)

//BinaryWriter ...
type BinaryWriter struct {
	buf *bytes.Buffer
}

//NewBinaryWriter returns a binary writer
func NewBinaryWriter() *BinaryWriter {
	return &BinaryWriter{
		buf: bytes.NewBuffer(make([]byte, 0)),
	}
}

func (b *BinaryWriter) Bytes() []byte {
	return b.buf.Bytes()
}

func (b *BinaryWriter) WriteUint8(v uint8) {
	b.write(v)
}

func (b *BinaryWriter) WriteUint16(v uint16) {
	b.write(v)
}

func (b *BinaryWriter) WriteUint32(v uint32) {
	b.write(v)
}

func (b *BinaryWriter) WriteFloat(v float32) {
	b.write(v)
}

func (b *BinaryWriter) WriteBytes(v []byte) {
	b.write(v)
}

func (b *BinaryWriter) write(v interface{}) {
	binary.Write(b.buf, binary.LittleEndian, v)
}
