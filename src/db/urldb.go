package db

import (
    pb "galaxy_walker/src/proto"
    "galaxy_walker/src/db/sqlite"
)
//////////////////////////////////
type UrlDbItf interface {
    /*
        新发现的url，如果task表不存在，则创建表
        需要去重复。。。使用task+docid
    */
    SetFreshUrls(taskid string, parentType int32, parentDocid int32, docs []*pb.CrawlDoc) (error,int)

    /*
        抓取完成，标记成功失败等; 失败次数更新
    */
    MarkCrawlFinishUrls(taskid string,docs []*pb.CrawlDoc)
    MarkCrawlFailUrls(taskid string, docs []*pb.CrawlDoc)
    ScanFreshUrls(taskid string,num int) (error,[]*pb.CrawlDoc)

    // TODO for task statistic
    //Statistic() *TaskStatistic
}
func NewUrlDbItf() UrlDbItf {
    return sqlite.NewUrlDbBySQLite()
}
//////////////////////////////////
