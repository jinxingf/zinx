package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

func main() {
	fmt.Println("client 1 starting ...")
	time.Sleep(1 * time.Second)
	// 1. get connection with remote server
	conn, err := net.Dial("tcp", "0.0.0.0:8999")
	if err != nil {
		fmt.Println("Failed to connect to server")
		return
	}

	for {
		// message struct
		dp := znet.DataPack{}
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(1, []byte("Zinx0.7 client 1 test message")))
		if err != nil {
			fmt.Println("pack message error ", err)
			return
		}
		if _, err = conn.Write(binaryMsg); err != nil {
			fmt.Println("write message error ", err)
			return
		}

		// server will write back some message
		// read header data, get msg id and data length
		// read data by data length
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read data err ", err)
			break
		}
		// unpack bianary header
		msgHeader, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("unpack header data err ", err)
			break
		}
		if msgHeader.GetMsgLen() > 0 {
			// read msg data second
			msg := msgHeader.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read message data err ", err)
				return
			}

			fmt.Println("Recv server msg, Msg ID = ", msg.GetMsgId(), ", length is ", msg.GetMsgLen(),
				" data is ", string(msg.GetData()))
		}
		// sleep for break
		time.Sleep(1 * time.Second)
	}
}
