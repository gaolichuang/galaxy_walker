/*
处理task的process
*/
package crawl

import (
    "galaxy_walker/src/proto"
    "galaxy_walker/src/crawl/scheduler"
    pb "galaxy_walker/src/proto"
    "galaxy_walker/src/task"
)

type TaskReceiverHandler struct {
    CrawlHandler
    dbScheduler *scheduler.DBScheduler
    docs         []*pb.CrawlDoc
    taskItf task.TaskItf
}
func (h *TaskReceiverHandler) RegisterTask(t task.TaskItf) {
    h.taskItf = t
}

func (handler *TaskReceiverHandler) Accept(crawlDoc *proto.CrawlDoc) bool {
    return true
}

func (handler *TaskReceiverHandler) Process(crawlDoc *proto.CrawlDoc) {


}

// use for create instance from a string
func init() {
    registerCrawlTaskType(&TaskReceiverHandler{})
}
