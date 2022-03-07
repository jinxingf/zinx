package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack_Pack(t *testing.T) {
	/*
		Mock Server
	*/
	// 1. create socketTCP
	listener, err := net.Listen("tcp4", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen 127.0.0.1:7777 err ", err)
		return
	}
	// create goroutine to accept data from client
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept err ", err)
			}

			go func(conn net.Conn) {
				// process client request
				dp := NewDataPack()
				// use for loop to continue read data from connection
				for {
					// 1. first unpack package to get message header: length and type
					headData := make([]byte, dp.GetHeadLen())
					// read all message
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read header err ", err)
						break
					}
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack err in first time", err)
						return
					}

					if msgHead.GetMsgLen() > 0 {
						// message data is not nil
						// 2. second read data by message header length
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						// base on the msg length, read data from io again
						// the length for msg is defined by split in msg.Data
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack err in second time", err)
							return
						}
						fmt.Println("---> Recv MsgId: ", msg.Id, " data len: ", msg.DataLen, " data : ", string(msg.Data))
					}
				}
			}(conn)
		}
	}()
	// 2. read data from client, and Pack

	/*
		Mock Client
	*/

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("conn to server err ", err)
		return
	}
	dp := NewDataPack()

	// two msg send together
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte("zinx"),
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err", err)
		return
	}

	// two msg send together
	msg2 := &Message{
		Id:      2,
		DataLen: 12,
		Data:    []byte("hello, zinx!"),
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg1 err", err)
		return
	}

	sendData1 = append(sendData1, sendData2...)

	_, err = conn.Write(sendData1)
	if err != nil {
		fmt.Println("client write sendData1 err", err)
		return
	}

	select {}
}
