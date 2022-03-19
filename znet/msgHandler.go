package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

type MsgHandle struct {
	// msg id map to handler
	Apis map[uint32]ziface.IRouter
	// message queue
	TaskQueue []chan ziface.IRequest
	// worker pool size
	WorkerPoolSize uint32
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalConf.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalConf.WorkerPoolSize),
	}
}

// DoMsgHandler scheduler or exec router msg handler
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgId(), " is not found, please register first")
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// AddRouter add msg handler
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeat msg id = " + strconv.Itoa(int(msgID)))
	}
	mh.Apis[msgID] = router
	fmt.Println("Add msg id ", msgID, " success")
}

// StartWorkerPool start worker pool
func (mh *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(utils.GlobalConf.WorkerPoolSize); i++ {
		// Start a worker

		// 1. prepare task queue for each worker
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalConf.MaxWorkerTaskSize)

		// 2. start worker, block for read message from channel
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// StartOneWorker start worker to process request
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Start one worker, worker id : ", workerID)

	// block to read from channel

	for {
		select {
		case req := <-taskQueue:
			mh.DoMsgHandler(req)
		}
	}
}

func (mh *MsgHandle) SendMsgToTask(request ziface.IRequest) {
	// get worker id by connection
	workerId := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add conn id ", request.GetConnection().GetConnID(),
		" request id ", request.GetMsgId(), " to worker id ", workerId)
	// Send request to taskQueue
	mh.TaskQueue[workerId] <- request
}
