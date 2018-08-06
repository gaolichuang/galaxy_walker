package db

import (
    pb "galaxy_walker/src/proto"
    "galaxy_walker/src/db/leveldb"
)
/*
根据抓取时间将content存储下来，支持查询。顺序写
如果后续有Merge的工作，离线在处理
*/
type ContentDBItf interface {
    // write
    SaveBatch(task string, docs []*pb.CrawlDoc) (error,int)
    /*
    scan by timestamp, if cross maxnum, only return maxnumber
    result will be sorted by crawl timestamp.
    */
    ScanByTimeRange(task string, start,end int64, maxnum int) (error,int64,[]*pb.CrawlDoc)
    ScanKeyByTimeRange(task string, start,end int64) []string

    // read one doc by docid
    GetDocById(task string, docid uint32) []*pb.CrawlDoc
    /*
        scan with iterator. start with "", num: read number,
        if maxnum < 0, will read all doc
        return value: error, iterator, batch doc.
    */
    ScanWithIterator(task, iterator string, maxnum int) (string,[]*pb.CrawlDoc)
    // for statistic.
    ScanKey(task string) []string

    // delete batch by docid.
    DeleteBatch(task string, docids []uint32) (error,int)
    // delete all
    Purge(task string) (error,int)
}

func NewContentDBItf() ContentDBItf {
    return leveldb.NewContentDbByLevelDB()
}