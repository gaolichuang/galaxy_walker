package crawl


import (
    "galaxy_walker/src/crawl/scheduler"
    pb "galaxy_walker/src/proto"
    LOG "galaxy_walker/internal/gcodebase/log"
    "galaxy_walker/src/task"
)

type TaskSchedulerHandler struct {
    CrawlHandler
    dbScheduler *scheduler.DBScheduler
    docs         []*pb.CrawlDoc
    taskItf task.TaskItf
}
func (h *TaskSchedulerHandler) RegisterTask(t task.TaskItf) {
    h.taskItf = t
}
func (h *TaskSchedulerHandler) Init() bool {
    LOG.VLog(3).Debug("TaskSchedulerHandler Init Finish")
    return true
}
func (h *TaskSchedulerHandler) Run(p CrawlProcessor) {

}

// use for create instance from a string
func init() {
    registerCrawlTaskType(&TaskSchedulerHandler{})
}
