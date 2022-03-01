package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("client starting ...")
	time.Sleep(1 * time.Second)
	// 1. get connection with remote server
	conn, err := net.Dial("tcp", "0.0.0.0:8999")
	if err != nil {
		fmt.Println("Failed to connect to server")
		return
	}

	for {
		_, err := conn.Write([]byte("hello, Zinx v0.1 ..."))
		if err != nil {
			fmt.Println("write to server error, ", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read from server error, ", err)
			return
		}
		fmt.Printf("server call back: %s, cnt : %d\n", buf, cnt)

		// sleep for break
		time.Sleep(1 * time.Second)
	}
}
