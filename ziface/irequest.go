package ziface

/*
Request connection and data binding
*/

type IRequest interface {
	// GetConnection get current connection
	GetConnection() IConnection

	// GetData get request data
	GetData() []byte
}
