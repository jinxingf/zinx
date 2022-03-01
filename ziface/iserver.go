package ziface

type IServer interface {
	// Start server
	Start()
	// Stop server
	Stop()
	// Server Run server
	Server()
}
