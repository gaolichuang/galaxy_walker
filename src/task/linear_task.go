package task

import (
    pb "galaxy_walker/src/proto"
    LOG "galaxy_walker/internal/gcodebase/log"
    "runtime"
    "reflect"
    "fmt"
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
}
func (l *LinearTask)CheckIsLinearTopology() error {
    /*
    TODO
    1.必须有startup
    2.是全联通图
    3.不能有环
    */
    if _,ok := l.rtypeCallBack[pb.RequestType_WEB_StartUp];!ok {
        return fmt.Errorf("Not exist RequestType_WEB_StartUp Node")
    }
    visit := make(map[pb.RequestType]bool)
    queue := make([]pb.RequestType,0)
    visit[pb.RequestType_WEB_StartUp]=true
    queue = append(queue,pb.RequestType_WEB_StartUp)
    for len(queue) > 0 {
        top := queue[0]
        for _,n := range l.rtypeCallBack[top] {
            queue = append(queue,n.next)
            if _,ok := visit[n.next];ok {
                return fmt.Errorf("Exist Circle %s",pb.RequestTypeToString(n.next))
            }
            visit[n.next]=true
        }
        queue = queue[:len(queue)-1]
    }
    if len(visit)!=len(l.rtypeCallBack) {
        return fmt.Errorf("not one island")
    }
    return nil
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