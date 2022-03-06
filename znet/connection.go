package znet

import (
	"fmt"
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
	ExitChan chan bool

	// process router
	Router ziface.IRouter
}

func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		Router:   router,
		ExitChan: make(chan bool, 1),
	}

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("reader goroutine running.")
	defer fmt.Println("connID = ", c.ConnID, " Reader is exit, remoteAddr is ", c.RemoteAddr())
	defer c.Stop()

	for {
		// read data into buffer, 512 byte
		buf := make([]byte, utils.GlobalConf.MaxPackageSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err ", err)
			continue
		}

		// get current request
		req := Request{
			conn: c,
			data: buf,
		}
		// get register binding function from router
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
		//// call handleAPI to process data read from buf
		//if err := c.HandleAPI(c.Conn, buf, cnt); err != nil {
		//	fmt.Println("ConnID = ", c.Conn, " handle err ", err)
		//	break
		//}
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

// Send data, send data to remote client
func (c *Connection) Send(date []byte) error {
	return nil
}
