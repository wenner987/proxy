package handle

import (
	"bytes"
	"context"
	"log"
	"net"
	"proxy/codec"
)

func Handle(conn net.Conn) {
	buffer := make([]byte, 1024)
	var cacheBuf []byte
	var cache = false
	var bufferCache []byte
	readCh := make(chan codec.Proto)
	writeCh := make(chan []byte)
	defer conn.Close()
	go readCon(context.Background(), conn, readCh, writeCh)
	go writeCon(context.Background(), conn, writeCh)
	for {
		n, err := conn.Read(buffer)
		if err != nil || n == 0{
			log.Printf("Read err: %v", err.Error())
			return
		}
		bWriter := bytes.NewBuffer([]byte{})
		if cache {
			bWriter.Write(cacheBuf)
			bWriter.Write(buffer[:n])
		} else {
			bWriter.Write(buffer[:n])
		}
		bufferCache = bWriter.Bytes()
		decodeLen, proto, err := codec.Parse(bufferCache)
		if err != nil {
			return
		}
		if decodeLen < uint32(len(bufferCache)) {
			cacheBuf = bufferCache[decodeLen:]
			cache = true
		} else {
			cache = false
		}
		for i := 0;i < len(proto); i++ {
			readCh <- proto[i]
		}
	}
}
