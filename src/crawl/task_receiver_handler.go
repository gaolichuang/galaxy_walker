package crawl
/*
可以处理多个task

调用taskprocessor的DoFinish方法处理
*/

import (
    pb "galaxy_walker/src/proto"
    "galaxy_walker/src/task"
    LOG "galaxy_walker/internal/gcodebase/log"
    "strings"
)

type TaskReceiverHandler struct {
    CrawlHandler
    taskProcessors map[string]*task.TaskProcessor
}

func (h *TaskReceiverHandler) Init() bool {
    LOG.VLog(3).Debug("TaskReceiverHandler Init Finish")
    h.taskProcessors = make(map[string]*task.TaskProcessor,0)
    for _,t := range strings.Split(*CONF.Crawler.SupportTasks,":") {
        tp := task.NewTaskProcessor()
        err := tp.Init(t)
        if err != nil {
            LOG.Fatalf("Get Task Error %s,%v",t,err)
        }
        h.taskProcessors[t]=tp
        LOG.VLog(2).DebugTag("TaskReceiverHandler","Add TaskProcessor %s",t)
    }
    return true
}
func (handler *TaskReceiverHandler) Accept(crawlDoc *pb.CrawlDoc) bool {
    return true
}

func (handler *TaskReceiverHandler) Process(crawlDoc *pb.CrawlDoc) {
    if _,ok := handler.taskProcessors[crawlDoc.CrawlParam.Taskid];!ok {
        LOG.Errorf("%s not belong any task",crawlDoc.Url)
        return
    }
    processor := handler.taskProcessors[crawlDoc.CrawlParam.Taskid]
    processor.DoFinish(crawlDoc)
}

// use for create instance from a string
func init() {
    registerCrawlTaskType(&TaskReceiverHandler{})
}
