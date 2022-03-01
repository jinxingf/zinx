package znet

import (
	"fmt"
	"net"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
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
	for {
		// 3. block wait for client connection, process client buss
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Accept error :", err)
			continue
		}

		// client finished connection create, do something. write back some data
		go func() {
			for {
				buf := make([]byte, 512)
				cnt, err := conn.Read(buf)
				if err != nil {
					fmt.Println("receive buf error :", err)
					continue
				}
				fmt.Printf("receive data from client buf : %s, cnt = %d\n", buf, cnt)

				if _, err := conn.Write(buf[:cnt]); err != nil {
					fmt.Println("write buf back error :", err)
					continue
				}
			}
		}()
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
