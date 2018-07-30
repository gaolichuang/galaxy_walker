package task
/*
处理task，调用task interface
*/
import (
    "galaxy_walker/internal/github.com/willf/bloom"
    pb "galaxy_walker/src/proto"
    LOG "galaxy_walker/internal/gcodebase/log"
    "galaxy_walker/src/db"
    "fmt"
    "galaxy_walker/internal/gcodebase/conf"
    "galaxy_walker/src/utils"
)
var CONF = conf.Conf
const (
    kBloomSize = 10000
    kBloomFalsePositive = 0.001
)

type TaskProcessor struct {
    bf *bloom.BloomFilter

    taskItf  TaskItf
    taskName string
    taskDes *pb.TaskDescription

    urlDbItf db.UrlDbItf
    taskDbItf db.TaskDbItf
}
/////
func genBloomKey(task,url string) string {
    return task + "^_^" + url
}
func NewTaskProcessor() *TaskProcessor {
    return &TaskProcessor{}
}
func (t *TaskProcessor)Init(task string) error {
    t.taskItf = GetTaskItfByName(task)
    if t.taskItf == nil {
        return fmt.Errorf("%s not valid TaskItf",task)
    }
    t.urlDbItf = db.NewUrlDbItf()
    t.taskDbItf = db.NewTaskItf()

    t.taskDes = t.taskDbItf.Get(task)
    if t.taskDes == nil {
        return fmt.Errorf("%s not init to db",task)
    }
    t.taskName = t.taskDes.Name
    //  init bf
    bf := bloom.NewWithEstimates(kBloomSize,kBloomFalsePositive)
    num := 0
    for _,url := range t.urlDbItf.ListUrls(t.taskName,-1,0) {
        bf.AddString(genBloomKey(t.taskName,url))
        num += 1
    }
    t.bf = bf
    LOG.VLog(2).DebugTag("TaskProcessor","%s load %d url in bf",t.taskName,num)
    return nil
}
func (t *TaskProcessor)DoFresh() []*pb.CrawlDoc{
    /*
    如果初始每调用，则调用初始状态的taskItf.Process
    否则就是scan db，把doc拿出来。
    */
    if pb.IsFreshTask(t.taskDes) {
        docs := t.taskItf.Process(pb.RequestType_WEB_StartUp,nil)
        err,num := t.urlDbItf.SetFreshUrls(t.taskName,pb.RequestType_WEB_StartUp,0,docs)
        if err != nil {
            LOG.Errorf("Set StartUp FreshUrls Error,%v,%v",err,docs)
            return nil
        }
        LOG.VLog(2).DebugTag("TaskProcess","StartUp Urls Set %d num",num)
        t.taskDbItf.Update(t.taskName,pb.KTaskStatusStarting,nil)
    }
    err,docs := t.urlDbItf.ScanFreshUrls(t.taskName,*CONF.Crawler.ScanFreshEachNumber)
    if err != nil {
        LOG.Errorf("%s ScanFreshUrls Error,%v",t.taskName,err)
        return nil
    }
    return docs
}
func (t *TaskProcessor)DoFinish(doc *pb.CrawlDoc) {
    /*
    调用process处理
    将返回的url，使用bf去重
    存入urldb返回
    */
    if !utils.IsCrawlSuccess(doc) {
        t.urlDbItf.MarkCrawlFailUrls(t.taskName,[]*pb.CrawlDoc{doc})
        return
    }
    docs := t.taskItf.Process(doc.CrawlParam.Rtype,doc)
    freshDocs := make([]*pb.CrawlDoc,0)
    //derepeat.
    for _,doc := range docs {
        if !t.bf.TestAndAddString(genBloomKey(t.taskName,doc.Url)) {
            freshDocs = append(freshDocs,doc)
        }
    }
    err,num := t.urlDbItf.SetFreshUrls(t.taskName,doc.CrawlParam.Rtype,doc.Docid,freshDocs)
    if err != nil {
        LOG.Errorf("TaskProcessor %s SetFreshUrls err %v",t.taskName,err)
        return
    }
    LOG.VLog(3).DebugTag("TaskProcessor","%s SetFreshUrls %d nums",num)
    t.urlDbItf.MarkCrawlFinishUrls(t.taskName,[]*pb.CrawlDoc{doc})
}

////////////////////////
func CreateTask(name string,job *pb.JobDescription) error {
    return nil
}