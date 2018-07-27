/*
基于db，未抓取的发送出去开抓
抓完的标记状态
bloom filter
*/
package scheduler

import (
    "galaxy_walker/internal/gcodebase/conf"
    "galaxy_walker/src/utils"
    pb "galaxy_walker/src/proto"
    LOG "galaxy_walker/internal/gcodebase/log"
)

var CONF = conf.Conf

// one taskscheduler process one task.
type TaskDbScheduler struct {
    // support fresh duplicate, load from sqlite...

    sender *CrawlDocSender
}

func (s *TaskDbScheduler) processFinish() {
    failDocs := make([]*pb.CrawlDoc, 0)
    successDocs := make([]*pb.CrawlDoc, 0)
    for _, doc := range docs {
        if !utils.IsCrawlSuccess(doc) {
            failDocs = append(failDocs, doc)
            continue
        }
        successDocs = append(successDocs, doc)

        freshDocs := s.taskItf.Process(doc.CrawlParam.Rtype, doc)
        // derepeat.

        // fill parenttype,parentdocid; set into urldb
        s.urlDbItf.SetFreshUrls(s.taskName, doc.CrawlParam.Rtype, doc.Docid, freshDocs)
    }
    // mark fail/success docs
    s.urlDbItf.MarkCrawlFailUrls(s.taskName, failDocs)
    s.urlDbItf.MarkCrawlFinishUrls(s.taskName, successDocs)
}
func (s *TaskDbScheduler) processFresh() {
    /*
        scan db 根据requesttype，调用TaskItf.Process处理
        将处理结果写回到db
    */
    err, docs := s.urlDbItf.ScanFreshUrls(s.taskName, *CONF.Crawler.TaskDBScanNum)
    if err != nil {
        LOG.Errorf("ScanFreshUrls Err %s %v",s.taskName,err)
        return
    }
    // fill parent and sendout.
    for _,doc := range docs {
        s.sender.Flush(doc)
    }
}
func (s *TaskDbScheduler) Run() {

}
