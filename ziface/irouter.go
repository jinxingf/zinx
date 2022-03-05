package ziface

type IRouter interface {
	// PreHandle PreHandle before request process
	PreHandle(request IRequest)
	// Handle request process
	Handle(request IRequest)
	// PostHandle after request process
	PostHandle(request IRequest)
}
