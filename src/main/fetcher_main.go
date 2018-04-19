package main

import (
	"galaxy_walker/internal/gcodebase/conf"
	"galaxy_walker/src/crawl"
	"galaxy_walker/internal/gcodebase/babysitter"
	"runtime"
	"time"
    "fmt"
)

var CONF = conf.Conf

func main() {
	runtime.GOMAXPROCS(2)
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
