package crawl


import (
    LOG "galaxy_walker/internal/gcodebase/log"
    "galaxy_walker/src/task"
    "time"
    "strings"
)
/*

每个TaskProcessor对应一个Task

支持多个task, 从TaskDb读取Task信息，然后为每一个Task创建TaskProcessor

调用DoFresh，并将结果发出去调度，不适用sender，将dispatcher地址填充到receiver地址即可。
依赖response_handler发送出去

需要对速度做一个控制，可以看后边的队列是否有内容，如果没有了在做下一次调度。

频繁的重启会对fail retry num有影响，这个如何修补？？？？ 判断output chain 大小
*/

type TaskSchedulerHandler struct {
    CrawlHandler
    taskProcessors []*task.TaskProcessor
}
func (h *TaskSchedulerHandler) register(taskname string) {
    tp := task.NewTaskProcessor()
    err := tp.Init(taskname)
    if err != nil {
        LOG.Fatalf("Get Task Error %s,%v",taskname,err)
    }
    h.taskProcessors = append(h.taskProcessors,tp)
    LOG.VLog(2).DebugTag("TaskSchedulerHandler","Add TaskProcessor %s",taskname)
}
func (h *TaskSchedulerHandler) Init() bool {
    LOG.VLog(3).Debug("TaskSchedulerHandler Init Finish")
    h.taskProcessors = make([]*task.TaskProcessor,0)
    for _,t := range strings.Split(*CONF.Crawler.SupportTasks,":") {
        h.register(t)
    }
    return true
}
func (h *TaskSchedulerHandler) Run(p CrawlProcessor) {
    for {
        time.Sleep(time.Second)
        if len(h.output_chan) > 0 {
            continue
        }
        for _,t := range h.taskProcessors {
            for _,doc := range t.DoFresh() {
                h.Output(doc)
            }
        }
    }
}

// use for create instance from a string
func init() {
    registerCrawlTaskType(&TaskSchedulerHandler{})
}
