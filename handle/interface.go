package handle

import (
	"context"
	"fmt"
	"log"
	"net"
	"proxy/codec"
	"proxy/data"
	"strconv"
)

const port = 28000
var conCache = map[string]net.Conn{}

func readCon(ctx context.Context, conn net.Conn, readCh chan codec.Proto, writeCh chan []byte) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("readCon exit.")
			return
		case proto := <-readCh:
			writeProxy(proto, writeCh)
		}
	}
}

func writeProxy(source codec.Proto, writeCh chan []byte) {
	var conn net.Conn
	var err error
	ip := data.GetIpFromZSet(int(source.Header.Dst))
	conn = conCache[ip + ":" + strconv.Itoa(port)]
	if conn == nil {
		conn, err = connect(ip, port)
		if err != nil {
			writeError(writeCh, err)
			return
		}
	}
	_, err = conn.Write(source.Context)
	if err != nil {
		writeError(writeCh, err)
		return
	}
}

func writeError(writeCh chan []byte, err error) {
	log.Println(err)
	writeCh <- []byte(fmt.Sprintf("conn err: %v", err))
	return
}

func writeCon(ctx context.Context, conn net.Conn, writeCh chan []byte) {
	for {
		select {
			case <-ctx.Done():
				return
			case writeContext := <-writeCh:
				_, _ = conn.Write(writeContext)
		}
	}
}

func connect(ip string, port int) (net.Conn, error) {
	addr := ip + ":" + strconv.Itoa(port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		conCache[addr] = conn
	}
	return conn, err
}
