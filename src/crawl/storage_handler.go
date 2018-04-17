package crawl

import (
        LOG "galaxy_walker/internal/gcodebase/log"
        pb "galaxy_walker/src/proto"
        "galaxy_walker/internal/gcodebase/string_util"
        "galaxy_walker/internal/gcodebase/time_util"
	"sync"
        "github.com/golang/protobuf/proto"
)

const (
	kSaveBatchSize     = 2
	kSaveBatchInterval = 300 // second
)

type StorageHandler struct {
	CrawlHandler
	docs []*pb.CrawlDoc
	sync.RWMutex
	last_db_time int64
}

func (handler *StorageHandler) saveDocs() {
	handler.Lock()
	defer handler.Unlock()
	if len(handler.docs) == 0 {
		return
	}
	if len(handler.docs) > kSaveBatchSize ||
		time_util.GetCurrentTimeStamp()-handler.last_db_time > kSaveBatchInterval {
		t1 := time_util.GetTimeInMs()
		num, err := STORAGE_ENGINE_IMPL.WithDb(*CONF.Crawler.ContentDbName).
			WithTable(*CONF.Crawler.ContentDbTable).
			SaveBatch(handler.docs)
		if err != nil {
			LOG.VLog(2).Debugf("Save Content to db error %s", err.Error())
		}
		handler.docs = make([]*pb.CrawlDoc, 0)
		handler.last_db_time = time_util.GetCurrentTimeStamp()
		LOG.VLog(3).Debugf("Flush %d using time %d ms.", num, time_util.GetTimeInMs()-t1)
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
	STORAGE_ENGINE_IMPL.Init(*CONF.Crawler.ContentDBServers)
	handler.docs = make([]*pb.CrawlDoc, 0)
	go handler.DBThread()
	return true
}
func (handler *StorageHandler) Accept(crawlDoc *pb.CrawlDoc) bool {
	return true
}

// save doc to content db
func (handler *StorageHandler) Process(crawlDoc *pb.CrawlDoc) {
	handler.Lock()
	defer handler.Unlock()
	// deepcopy crawldoc
        nd := proto.Clone(crawlDoc)
        newDoc := nd.(*pb.CrawlDoc)

	// compress and save
	compressContent, err := string_util.Compress(newDoc.Content)
	if err == nil {
		newDoc.Content = compressContent
		newDoc.ContentCompressed = true
	} else {
		LOG.VLog(2).Debugf("Compress Error url:%s,docid:%s,error:%s",
			newDoc.Url, newDoc.Docid, err.Error())
	}
	handler.docs = append(handler.docs, &newDoc)
}

// use for create instance from a string
func init() {
	registerCrawlTaskType(&StorageHandler{})
}
