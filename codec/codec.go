package codec

import (
	"bytes"
	"encoding/binary"
	_ "github.com/go-redis/redis/v8"
	"log"
)

type Error struct {}

func (e Error) Error() string {
	panic("implement me")
	return "Magic error."
}


func Parse(bufferContext []byte) (n uint32, arr []Proto, err error) {
	var parseBuffer []byte
	parseBuffer = bufferContext[:]
	bufferContextLen := len(bufferContext)
	for {
		if n > uint32(bufferContextLen) {
			break
		}
		msg := Proto{}
		msg.Header = parseHeader(parseBuffer[n:])
		// 长度不足, 等待下一次解析
		if n + msg.Header.ProtoLen + headerLen > uint32(bufferContextLen) {
			return n, arr, err
		}
		msg.Context = parseBuffer[n + headerLen:n + msg.Header.ProtoLen + headerLen]
		// 检查报文magic
		if !checkMagic(&msg.Header) {
			log.Printf("Magic %d != %d", msg.Header.Magic, magicNum)
			return n, arr, Error{}
		}
		arr = append(arr, msg)
		n += headerLen
		n += msg.Header.ProtoLen
	}
	return n, arr, err
}

func checkMagic(msg *ProtoHeader) bool {
	return msg.Magic == magicNum
}

func parseHeader(bufferContext []byte) (header ProtoHeader) {
	buffer := bytes.NewBuffer(bufferContext)
	readBinary(buffer, &header.Magic)
	readBinary(buffer, &header.From)
	readBinary(buffer, &header.ProtoLen)
	readBinary(buffer, &header.Version)
	readBinary(buffer, &header.Dst)
	return header
}

func readBinary(buffer *bytes.Buffer, data interface{}) {
	_ = binary.Read(buffer, binary.BigEndian, data)
}
