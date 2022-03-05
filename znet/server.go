package znet

import (
	"errors"
	"fmt"
	"net"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
}

// CallBackToClient currently call back function is hardcode, call back function
// should be as args of connection
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	// write back
	fmt.Println("[Conn Handle] CallBackToClient ...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err ", err)
		return errors.New("CallBackToClient error")
	}

	return nil
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server listenner at IP : %s, Port : %d\n", s.IP, s.Port)

	// 1. get tcp addr
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("resolve addr error : ", err)
		return
	}
	// 2. listen server addr
	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("listen ", s.IPVersion, " tcp error : ", err)
		return
	}
	fmt.Println("start Zinx server succ, ", s.Name, " success, Listening...")

	var cid uint32
	cid = 0

	for {
		// 3. block wait for client connection, process client buss
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Accept error :", err)
			continue
		}

		// client finished connection create, do something. write back some data
		NewConnection(conn, cid, CallBackToClient)
		cid++
	}

}

func (s *Server) Stop() {
	// release resource, stop server

}

func (s *Server) Server() {
	// start server
	go s.Start()

	// do something

	// block forever to wait connection
	select {}
}

func NewServer(name string) *Server {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
