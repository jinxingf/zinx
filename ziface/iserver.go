package ziface

type IServer interface {
	// Start server
	Start()
	// Stop server
	Stop()
	// Server Run server
	Server()
	// AddRouter function, register route function for current server, for client connection to use.
	AddRouter(route IRouter)
}
