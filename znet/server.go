package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
	Route     ziface.IRouter
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] zinx global config : %+v", utils.GlobalConf)
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
		dealConn := NewConnection(conn, cid, s.Route)
		cid++

		// process business
		go dealConn.Start()
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

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Route = router
	fmt.Println("Add router success")
}

func NewServer(name string) *Server {
	s := &Server{
		Name:      utils.GlobalConf.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalConf.Host,
		Port:      utils.GlobalConf.TCPPort,
		Route:     nil,
	}
	return s
}
