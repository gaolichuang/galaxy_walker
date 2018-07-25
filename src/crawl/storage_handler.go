package crawl

import (
    LOG "galaxy_walker/internal/gcodebase/log"
    pb "galaxy_walker/src/proto"
    "galaxy_walker/internal/gcodebase/string_util"
    "galaxy_walker/internal/gcodebase/time_util"
    "sync"
    "galaxy_walker/src/db"
)

const (
    kSaveBatchSize     = 2
    kSaveBatchInterval = 33 // second
)

type StorageHandler struct {
    CrawlHandler
    // key is task.
    docs         map[string][]*pb.CrawlDoc
    docsNum     int
    contentDb db.ContentDBItf

    sync.RWMutex
    last_db_time int64
}

func (handler *StorageHandler) saveDocs() {
    handler.Lock()
    defer handler.Unlock()
    if len(handler.docs) == 0 {
        return
    }
    if handler.docsNum > kSaveBatchSize ||
        time_util.GetCurrentTimeStamp()-handler.last_db_time > kSaveBatchInterval {
        t1 := time_util.GetTimeInMs()
        num := 0
        for task,docs := range handler.docs {
            err,n := handler.contentDb.SaveBatch(task,docs)
            if err != nil {
                LOG.VLog(2).Debugf("Save Content to db error %v", err)
            } else {
                num += n
            }
        }
        handler.docs = make(map[string][]*pb.CrawlDoc)
        handler.last_db_time = time_util.GetCurrentTimeStamp()
        LOG.VLog(3).Debugf("Flush %d(%d) using time %d ms.", num, handler.docsNum, time_util.GetTimeInMs() - t1)
        handler.docsNum = 0
    }
}
func (handler *StorageHandler) DBThread() {
    for true {
        handler.saveDocs()
        time_util.Sleep(10)
    }
}
func (handler *StorageHandler) Init() bool {
    // TODO: change to init URI
    handler.docs = make(map[string][]*pb.CrawlDoc)
    handler.contentDb = db.NewContentDBItf()
    go handler.DBThread()
    return handler.contentDb != nil
}
func (handler *StorageHandler) Accept(crawlDoc *pb.CrawlDoc) bool {
    return true
}

// save doc to content db
func (handler *StorageHandler) Process(crawlDoc *pb.CrawlDoc) {
    handler.Lock()
    defer handler.Unlock()
    // deepcopy crawldoc
    newDoc := pb.CloneCrawlDoc(crawlDoc)

    // compress and save
    compressContent, err := string_util.Compress(newDoc.Content)
    if err == nil {
        newDoc.Content = compressContent
        newDoc.ContentCompressed = true
    } else {
        LOG.VLog(2).Debugf("Compress Error url:%s,docid:%s,error:%s",
            newDoc.Url, newDoc.Docid, err.Error())
    }
    task := crawlDoc.CrawlParam.Taskid
    handler.docsNum += 1
    if _,ok := handler.docs[task];!ok {
        handler.docs[task] = make([]*pb.CrawlDoc,0)
    }
    handler.docs[task] = append(handler.docs[task],newDoc)
}

// use for create instance from a string
func init() {
    registerCrawlTaskType(&StorageHandler{})
}
