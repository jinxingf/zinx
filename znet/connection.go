package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
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
	ExitChan chan bool

	MsgHandler ziface.IMsgHandler
}

func NewConnection(conn *net.TCPConn, connID uint32, handler ziface.IMsgHandler) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: handler,
		ExitChan:   make(chan bool, 1),
	}

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("reader goroutine running.")
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
		// get register binding function from msg handler
		go c.MsgHandler.DoMsgHandler(&req)
	}
}

// Start connection, current connection start to work
func (c *Connection) Start() {
	fmt.Println("Connection start(), ConnID = ", c.ConnID)
	go c.StartReader()
	// start to read data from current connection, and do some business
	// TODO: add write data from current connection adn do some business

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
	// recycle channel resource
	close(c.ExitChan)
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
	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("write msg id  ", id, " msg data ", string(data), " err ", err)
		return errors.New("conn Write err")
	}
	return nil
}
