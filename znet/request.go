package znet

import "zinx/ziface"

type Request struct {
	conn ziface.IConnection
	data []byte
}

// GetConnection get current connection
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// GetData get request data
func (r *Request) GetData() []byte {
	return r.data
}
