package listen

import (
	"log"
	"net"
	"proxy/configure"
	"proxy/handle"
	"strconv"
)

type Server struct {
	strIp string
	iPort int
	handle func(net.Conn)
}

func (s *Server) Init() {
	s.strIp = configure.GetConfigInstance().Server.SvrIp
	s.iPort = int(configure.GetConfigInstance().Server.SvrPort)
}

func (s *Server) Listen() {
	addr := s.strIp + ":" + strconv.Itoa(s.iPort)
	tcpAdd,err:= net.ResolveTCPAddr("tcp",addr)
	if err!=nil{
		log.Fatalln("net.ResolveTCPAddr error:",err)
		return
	}
	listener, err := net.ListenTCP("tcp", tcpAdd)
	defer listener.Close()
	if err != nil {
		log.Fatalf("Listen err: %v", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Conn err: %v", err)
		}
		go handle.Handle(conn)
	}
}