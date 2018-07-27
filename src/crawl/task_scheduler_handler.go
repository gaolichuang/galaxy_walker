package crawl


import (
    "galaxy_walker/src/crawl/scheduler"
    pb "galaxy_walker/src/proto"
    LOG "galaxy_walker/internal/gcodebase/log"
    "galaxy_walker/src/task"
    "time"
)

type TaskSchedulerHandler struct {
    CrawlHandler
    dbScheduler *scheduler.TaskDbScheduler
    docs         []*pb.CrawlDoc
    taskItf task.TaskItf
}
func (h *TaskSchedulerHandler) RegisterTask(t task.TaskItf) {
    h.taskItf = t
}
func (h *TaskSchedulerHandler) Init() bool {
    LOG.VLog(3).Debug("TaskSchedulerHandler Init Finish")
    h.dbScheduler.SetFresh(h.taskItf.TaskId(),
        // pb.RequestType_WEB_StartUp is starting.
        h.taskItf.Process(pb.RequestType_WEB_StartUp,nil))
    return true
}
func (h *TaskSchedulerHandler) Run(p CrawlProcessor) {
    for {
        docs := h.dbScheduler.ScanFresh(h.taskItf.TaskId(),100)
        for _,doc := range docs {

            h.output_chan <- doc
        }
        time.Sleep(time.Second)
    }

}

// use for create instance from a string
func init() {
    registerCrawlTaskType(&TaskSchedulerHandler{})
}
