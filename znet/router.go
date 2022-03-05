package znet

import "zinx/ziface"

/*
when implement router, first insert the BaseRouter, and overwrite the BaseRouter
*/

type BaseRouter struct {
}

// the inherit function is None, because some other router
// perhaps not need one of function.

// PreHandle PreHandle before request process
func (b *BaseRouter) PreHandle(request ziface.IRequest) {}

// Handle request process
func (b *BaseRouter) Handle(request ziface.IRequest) {}

// PostHandle after request process
func (b *BaseRouter) PostHandle(request ziface.IRequest) {}
