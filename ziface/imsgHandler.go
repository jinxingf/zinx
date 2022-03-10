package ziface

type IMsgHandler interface {
	// DoMsgHandler scheduler or exec router msg handler
	DoMsgHandler(request IRequest)
	// AddRouter add msg handler
	AddRouter(uint32, IRouter)
}
