package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

// ping test router

type PingRouter struct {
	// inherit
	znet.BaseRouter
}

func (p *PingRouter) PreHandle(request ziface.IRequest) {

	fmt.Println("Call router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping\n"))
	if err != nil {
		fmt.Println("Call back before ping error")
	}
}

func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call router Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println("Call back ping...ping...ping")
	}
}

func (p *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("post ping\n"))
	if err != nil {
		fmt.Println("Call back post ping error")
	}
}

func main() {
	s := znet.NewServer("[Zinx V0.3]")

	s.AddRouter(&PingRouter{})

	s.Server()
}
