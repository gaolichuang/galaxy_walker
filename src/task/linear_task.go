package task

import (
    pb "galaxy_walker/src/proto"
    LOG "galaxy_walker/internal/gcodebase/log"
    "runtime"
    "reflect"
)
/*

*/

// input: current RequestType, crawlSuccessDoc
// output: extract links,urls.
type  TaskProcessFn func(pb.RequestType, *pb.CrawlDoc) []string

type taskChain struct {
    current pb.RequestType
    next pb.RequestType
    fn TaskProcessFn
    fnName string
}

type LinearTask struct {
    rtypeCallBack map[pb.RequestType][]taskChain
}
func (l *LinearTask)RegisterRequestTypeCallBack(current,next pb.RequestType,fn TaskProcessFn) {
    if l.rtypeCallBack == nil {
        l.rtypeCallBack = make(map[pb.RequestType][]taskChain)
    }
    if _,ok:=l.rtypeCallBack[current];!ok {
        l.rtypeCallBack[current] = make([]taskChain,0)
    }
    l.rtypeCallBack[current]=append(l.rtypeCallBack[current],taskChain{current,next,fn,runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()})
    // TODO. check circle. check all node can connect from zero.
}

func (l *LinearTask)Process(rtype pb.RequestType, doc *pb.CrawlDoc) []*pb.CrawlDoc {
    if l.rtypeCallBack == nil {return nil}
    if _,ok := l.rtypeCallBack[rtype];!ok {
        return nil
    }
    docs := make([]*pb.CrawlDoc,0)
    for _,taskc := range l.rtypeCallBack[rtype] {
        LOG.VLog(4).DebugTag("LinearTask","Process Current:%s,Next:%s,fn:%s",
            pb.RequestTypeToString(rtype),pb.RequestTypeToString(taskc.next),taskc.fnName)
        if doc != nil {
            LOG.VLog(4).DebugTag("LinearTask","Process doc:%s",doc.Url)
        }
        urls := taskc.fn(rtype,doc)
        for _,url := range urls {
            docs = append(docs, &pb.CrawlDoc{
                RequestUrl:url,
                CrawlParam:&pb.CrawlParam{
                    Rtype:taskc.next, // next level
                },
            })
            LOG.VLog(5).DebugTag("LinearTask","Fresh Current:%s,Next:%s,fn:%s,doc:%s",
                pb.RequestTypeToString(rtype),pb.RequestTypeToString(taskc.next),taskc.fnName,url)
        }
    }
    return docs
}