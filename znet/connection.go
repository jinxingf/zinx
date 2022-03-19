package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	// current connection TCP socket
	Conn *net.TCPConn
	// connection id
	ConnID uint32
	// current connection status
	isClosed bool
	// notify current connection exist channel
	// reader tell writer exist by channel
	ExitChan chan bool

	// no buf channel for reader and writer transform data
	msgChan chan []byte

	MsgHandler ziface.IMsgHandler
}

func NewConnection(conn *net.TCPConn, connID uint32, handler ziface.IMsgHandler) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		MsgHandler: handler,
	}

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("[Reader goroutine is running.]")
	defer fmt.Println("connID = ", c.ConnID, " Reader is exit, remoteAddr is ", c.RemoteAddr())
	defer c.Stop()

	for {
		// read msg header, 8 byte
		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg header err ", err)
			break
		}
		// unpack header, get msg id and length
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack msg err ", err)
			break
		}
		// read msg data by msg length
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data err ", err)
				break
			}
		}
		msg.SetMsgData(data)

		// get current request
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalConf.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTask(&req)
		} else {
			// get register binding function from msg handler
			c.MsgHandler.SendMsgToTask(&req)
		}

	}
}

// StartWriter send data to client
func (c *Connection) StartWriter() {
	fmt.Println("[Writer goroutine is running.]")
	defer fmt.Println(c.RemoteAddr().String(), " [Writer goroutine exist.]")

	// block for reader data from channel and send to client
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.GetTCPConnection().Write(data); err != nil {
				fmt.Println("Send data to client failed, err ", err)
				return
			}
		case <-c.ExitChan:
			// reader has existed, writer goroutine also need to exist
			return
		}
	}

}

// Start connection, current connection start to work
func (c *Connection) Start() {
	fmt.Println("Connection start(), ConnID = ", c.ConnID)
	go c.StartReader()
	// start to read data from current connection, and do some business
	go c.StartWriter()

}

// Stop connection, close current connection, end current work
func (c *Connection) Stop() {
	fmt.Println("Connection stop(), ConnID = ", c.ConnID)
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	err := c.Conn.Close()
	if err != nil {

	}
	// close writer goroutine
	c.ExitChan <- true
	// recycle channel resource
	close(c.ExitChan)
	close(c.msgChan)
}

// GetTCPConnection get current connection binding socket
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID get connection id
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr get remote client connection info, such as TCP status
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SendMsg Send data, send data to remote client
func (c *Connection) SendMsg(id uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection is closed when send msg")
	}
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(id, data))
	if err != nil {
		fmt.Println("pack msg err ", err)
		return errors.New("pack msg err")
	}

	c.msgChan <- binaryMsg

	return nil
}
