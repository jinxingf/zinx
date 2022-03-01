package znet

import (
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
	// current connection binding process business function
	HandleAPI ziface.HandleFunc
	// notify current connection exist channel
	ExitChan chan bool
}

func NewConnection(conn *net.TCPConn, connID uint32, callback_api ziface.HandleFunc) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		isClosed:  false,
		HandleAPI: callback_api,
		ExitChan:  make(chan bool, 1),
	}

	return c
}
