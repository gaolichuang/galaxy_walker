package crawl

import (
    LOG "galaxy_walker/internal/gcodebase/log"
    "galaxy_walker/src/proto"
    "time"
)

type DummyRequestProcessor struct {
    CrawlHandler
}

func (request *DummyRequestProcessor) Run(p CrawlProcessor) {
    for {
        doc := new(proto.CrawlDoc)
        doc.Url = "http://roll.sohu.com/sports/"
        doc.CrawlParam = new(proto.CrawlParam)
        doc.CrawlParam.FetchHint = new(proto.FetchHint)
        doc.CrawlParam.FetchHint.Host = "roll.sohu.com"
        doc.CrawlParam.Hostload = 5
        doc.CrawlRecord = new(proto.CrawlRecord)
        request.Output(doc)
        LOG.Info("Send one request")
        time.Sleep(time.Minute)
    }
}

// use for create instance from a string
func init() {
    registerCrawlTaskType(&DummyRequestProcessor{})
}
