/*
处理task的process
*/
package crawl

import (
    "galaxy_walker/src/crawl/scheduler"
    pb "galaxy_walker/src/proto"
    "galaxy_walker/src/task"
)

type TaskReceiverHandler struct {
    CrawlHandler
    dbScheduler *scheduler.TaskDbScheduler
    taskItf task.TaskItf
}
func (h *TaskReceiverHandler) RegisterTask(t task.TaskItf) {
    h.taskItf = t
}

func (handler *TaskReceiverHandler) Accept(crawlDoc *pb.CrawlDoc) bool {
    return true
}

func (handler *TaskReceiverHandler) Process(crawlDoc *pb.CrawlDoc) {
    fdocs := handler.taskItf.Process(crawlDoc.CrawlParam.Rtype,crawlDoc)
    handler.dbScheduler.SetFresh(crawlDoc.CrawlParam.Taskid,fdocs)
    handler.dbScheduler.MarkFinishAndFail(crawlDoc.CrawlParam.Taskid,[]*pb.CrawlDoc{crawlDoc})
}

// use for create instance from a string
func init() {
    registerCrawlTaskType(&TaskReceiverHandler{})
}
