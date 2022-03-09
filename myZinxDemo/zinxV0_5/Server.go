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
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println("Call back ping...ping...ping")
	}
	fmt.Println("recv from client, MsgId: ", request.GetMsgId(), ", msgData ", string(request.GetData()))
	// request args has connection and data, and write back something
	if err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping")); err != nil {
		fmt.Println(err)
	}

}

func main() {
	s := znet.NewServer("[Zinx V0.5]")

	s.AddRouter(&PingRouter{})

	s.Server()
}
