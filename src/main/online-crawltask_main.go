/*
单独一个进程，全部包含

*/
package main


import (
    "galaxy_walker/internal/gcodebase/conf"
    "galaxy_walker/src/crawl"
    "galaxy_walker/internal/gcodebase/babysitter"
    "runtime"
    "time"
    "fmt"
    "galaxy_walker/src/task"
    LOG "galaxy_walker/internal/gcodebase/log"
)

var CONF = conf.Conf

func main() {
    runtime.GOMAXPROCS(2)
    t := task.GetTaskItfByName(*CONF.Crawler.CrawlTaskName)
    if t == nil {
        LOG.Fatalf("Can not get Crawl Handler %s", *CONF.Crawler.CrawlTaskName)
    }
    //////////////
    h := crawl.GetCrawlHandlerByName("TaskSchedulerHandler")
    if h == nil {
        LOG.Fatalf("Can not get Crawl Handler TaskSchedulerHandler")
    }
    schedulerHandler,ok := h.(*crawl.TaskSchedulerHandler)
    if !ok {
        LOG.Fatalf("Can translate to TaskSchedulerHandler")
    }
    schedulerHandler.RegisterTask(t)
    ///////////////
    //////////////
    h = crawl.GetCrawlHandlerByName("TaskReceiverHandler")
    if h == nil {
        LOG.Fatalf("Can not get Crawl Handler TaskReceiverHandler")
    }
    receiverHandler,ok := h.(*crawl.TaskReceiverHandler)
    if !ok {
        LOG.Fatalf("Can translate to TaskReceiverHandler")
    }
    receiverHandler.RegisterTask(t)
    //////////////

    c := crawl.CrawlHandlerController{}
    c.InitCrawlService()

    var http_server babysitter.MonitorServer
    http_server.Init("fetcher")

    http_server.AddMonitor(&c)
    go http_server.Serve(fmt.Sprintf(":%d",*CONF.Crawler.HttpPort),"")

    for {
        c.PrintStatus()
        time.Sleep(time.Second * 10)
    }
}

