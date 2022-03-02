package main

import (
	"zinx/znet"
)

func main() {
	s := znet.NewServer("[Zinx V0.0.1]")
	s.Server()
}
