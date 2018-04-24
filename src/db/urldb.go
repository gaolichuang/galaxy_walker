package db

import (
    pb "galaxy_walker/src/proto"
)
const (
    KUrlStatusUndo = 1
    KUrlStatusDoing = 2
    KUrlStatusDone = 3
)
//////////////////////////////////
func SetFreshUrls(taskid string, docs []*pb.CrawlDoc) {
    /*
    抓取完成，标记成功失败等; 失败次数更新
    支持多个taskid
    */
}

func MarkCrawlFinishUrls(taskid string, docs []*pb.CrawlDoc) {
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