package ziface

import "net"

// IConnection define connection interface
type IConnection interface {
	// Start connection, current connection start to work
	Start()
	// Stop connection, close current connection, end current work
	Stop()
	// GetTCPConnection get current connection binding socket
	GetTCPConnection() *net.TCPConn
	// GetConnID get connection id
	GetConnID() uint32
	// RemoteAddr get remote client connection info, such as TCP status
	RemoteAddr() net.Addr
	// Send data, send data to remote client
	Send(date []byte) error
}

// HandleFunc define process business function,
// use net.TCPConn to process len of []byte data,
type HandleFunc func(conn *net.TCPConn, data []byte, len int) error
