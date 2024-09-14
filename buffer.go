package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"unicode/utf8"
)

type BufferWriter struct {
	buffer *bytes.Buffer
}

func newBufferWriter() *BufferWriter {
	buff := new(BufferWriter)
	buff.buffer = new(bytes.Buffer)

	return buff
}

func (b BufferWriter) writeInt(val int32) {
	err := binary.Write(b.buffer, binary.BigEndian, val)

	if err != nil {
		fmt.Println(err)
	}
}

func (b BufferWriter) writeString(val string) {
	byteLength := utf8.RuneCountInString(val)
	b.writeInt(int32(byteLength))

	_, err := b.buffer.WriteString(val)

	if err != nil {
		fmt.Println(err)
	}
}

func (b BufferWriter) writeFixBuffer(buffer []byte) {
	bufferLen := len(buffer)

	b.writeInt(int32(bufferLen))
	_, err := b.buffer.Write(buffer)

	if err != nil {
		fmt.Println(err)
	}
}

func (b BufferWriter) writeBuffer(buffer []byte) {
	_, err := b.buffer.Write(buffer)

	if err != nil {
		fmt.Println(err)
	}
}

func (b BufferWriter) crop(offset int, length int) []byte {
	return b.buffer.Bytes()[offset:length]
}

func (b BufferWriter) getBytes() []byte {
	return b.buffer.Bytes()
}

type BufferReader struct {
	reader *bytes.Reader
}

func newBufferReader(byteArray []byte) *BufferReader {
	buf := new(BufferReader)
	buf.reader = bytes.NewReader(byteArray)

	return buf
}

func (b BufferReader) readInt() int32 {
	var val int32

	err := binary.Read(b.reader, binary.BigEndian, &val)
	if nil != err {
		fmt.Print(err)
	}

	return val
}

func (b BufferReader) readString() string {
	strLen := b.readInt()

	buffer := make([]byte, strLen)

	err := binary.Read(b.reader, binary.BigEndian, buffer)
	if nil != err {
		fmt.Print(err)
	}

	return string(buffer)
}

func (b BufferReader) readFixBuffer() []byte {
	length := b.readInt()
	out := make([]byte, length)

	_, err := b.reader.Read(out)
	if nil != err {
		fmt.Println(err)
	}

	return out
}

func (b BufferReader) readRemain() []byte {
	res, err := io.ReadAll(b.reader)

	if nil != err {
		fmt.Println(err)
	}

	return res
}
