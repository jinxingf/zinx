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

func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call router Handle...")
	fmt.Println("recv from client, MsgId: ", request.GetMsgId(), ", msgData ", string(request.GetData()))
	// request args has connection and data, and write back something
	if err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping")); err != nil {
		fmt.Println(err)
	}

}

type HelloRouter struct {
	// inherit
	znet.BaseRouter
}

func (h *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call hello router Handle...")
	fmt.Println("recv from client, MsgId: ", request.GetMsgId(), ", msgData ", string(request.GetData()))
	// request args has connection and data, and write back something
	if err := request.GetConnection().SendMsg(201, []byte("hello...hello...hello")); err != nil {
		fmt.Println(err)
	}

}

func main() {
	s := znet.NewServer("[Zinx V0.8]")

	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	s.Server()
}
