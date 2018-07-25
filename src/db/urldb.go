package db

import (
    pb "galaxy_walker/src/proto"
)
const (
    KUrlStatusUndo = 1  //  not crawl
    KUrlStatusDoing = 2  // scan out
    KUrlStatusDone = 3  // success
)
//////////////////////////////////
func SetFreshUrls(taskid string, parentType int32, parentDocid int32, docs []*pb.CrawlDoc) {
    /*
    新发现的url，如果task表不存在，则创建表
    需要去重复。。。使用task+docid
    */
}

func MarkCrawlFinishUrls(taskid string,docs []*pb.CrawlDoc) {
    /*
    抓取完成，标记成功失败等; 失败次数更新
    支持多个taskid
    */
}
func MarkCrawlFailUrls(taskid string, docs []*pb.CrawlDoc) {
    /*
    抓取完成，标记成功失败等; 失败次数更新
    支持多个taskid
    */
}

func ScanFreshUrls(taskid string,num int) []*pb.CrawlDoc {
    /*
    抓取完成，标记成功失败等
    支持多个taskid
    是否按照level区分
    */
    return nil
}
//////////////////////////////////
func UrlNumbersByType(taskid string,status int) int {
    return 0
}
func ScanUrlsByTask(taskid string,status,num int) []*pb.CrawlDoc {
    return nil
}