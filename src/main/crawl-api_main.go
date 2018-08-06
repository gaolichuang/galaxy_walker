package main

import (
    "time"
    "net"
    "strconv"
    LOG "galaxy_walker/internal/gcodebase/log"
    "galaxy_walker/internal/gcodebase/conf"
    "runtime"
    "runtime/pprof"
    "os"
    "fmt"
    "galaxy_walker/internal/gcodebase/babysitter"
    "galaxy_walker/src/api"
)

var CONF = conf.Conf

func main() {
    currentProc := runtime.GOMAXPROCS(-1)
    LOG.Infof("==============Use MaxProcNum : %d, CPUNum:%d===========", currentProc, runtime.NumCPU())
    if *CONF.Crawler.EnableCPUPProf {
        f, _ := os.Create("./api.profile_file")
        pprof.StartCPUProfile(f) // 开始cpu profile，结果写到文件f中
    }
    //LOG.Infof("Set GCPercent from %d => 20",debug.SetGCPercent(20))

    LOG.DumpFlags()

    apiService := &api.APIService{}
    apiService.Init()
    h, p, e := net.SplitHostPort(*CONF.Crawler.APIListenAddr)
    if e != nil {
        LOG.Fatal("%s Err %v", *CONF.Crawler.APIListenAddr, e)
    }
    port, _ := strconv.Atoi(p)

    sh, sp, e := net.SplitHostPort(*CONF.Crawler.APIHTTPSListenAddr)
    if e != nil {
        LOG.Errorf("%s Err %v", *CONF.Crawler.APIHTTPSListenAddr, e)
    }
    sport, _ := strconv.Atoi(sp)
    go apiService.Serve(h, port, sh, sport)

    if *CONF.Crawler.HttpPort > 0 {
        var http_server babysitter.MonitorServer
        http_server.Init("dndns.api")

        http_server.AddMonitor(apiService)
        http_server.AddHandleFunc("/config", LOG.HandleFlagsFunc)
        go http_server.Serve(
            fmt.Sprintf(":%d", *CONF.Crawler.HttpPort),
            fmt.Sprintf(":%d", *CONF.Crawler.HttpsPort))
    }

    for {
        LOG.VLog(5).Debugf("APISERVICE:%s", apiService.Status())
        time.Sleep(time.Second * time.Duration(60))
        apiService.Reload()
    }
}
