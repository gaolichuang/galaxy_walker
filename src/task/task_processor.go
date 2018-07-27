package task
/*
处理task，调用task interface
*/
import (
    "galaxy_walker/internal/github.com/willf/bloom"
    pb "galaxy_walker/src/proto"
    "galaxy_walker/src/db"
)


type TaskProcessor struct {
    bf *bloom.BloomFilter

    urlDbItf db.UrlDbItf
    taskItf  TaskItf
    taskName string
}
func (t *TaskProcessor)DoFresh() []*pb.CrawlDoc{
    /*
    如果初始每调用，则调用初始状态的taskItf.Process
    否则就是scan db，把doc拿出来。
    */
}
func (t *TaskProcessor)DoFinish(doc *pb.CrawlDoc) {
    /*
    调用process处理
    将返回的url，使用bf去重
    存入urldb返回
    */

}
