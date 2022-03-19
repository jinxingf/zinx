package ziface

type IMsgHandler interface {
	// DoMsgHandler scheduler or exec router msg handler
	DoMsgHandler(IRequest)
	// AddRouter add msg handler
	AddRouter(uint32, IRouter)
	StartWorkerPool()
	SendMsgToTask(IRequest)
}
