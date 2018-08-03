package task

import (
    pb "galaxy_walker/src/proto"
    "galaxy_walker/src/crawl/scheduler"
    LOG "galaxy_walker/internal/gcodebase/log"
)

type DummyTask struct {

}

func (d *DummyTask) Init() error {
    return nil
}
func (d *DummyTask) GetJobDescription() *pb.JobDescription {
    return &scheduler.NormalJobD
}
func (d *DummyTask)Process(rtype pb.RequestType, doc *pb.CrawlDoc) []*pb.CrawlDoc {
    switch rtype {
    case pb.RequestType_WEB_StartUp:
        // start url, no need doc. just return docs.
        urls := []string{
            "http://roll.sohu.com/sports/",
        }
        docs := make([]*pb.CrawlDoc,0)
        for _,url := range urls {
            docs = append(docs, &pb.CrawlDoc{
                RequestUrl:url,
                CrawlParam:&pb.CrawlParam{
                    Rtype:pb.RequestType_WEB_MAIN, // next level
                },
            })
        }
        return docs
    case pb.RequestType_WEB_MAIN:
        // startup docs. parse and return web hub. mark RequestType.
        LOG.VLog(2).DebugTag("XXXXXX","RequestType_WEB_MAIN %s",pb.FromProtoToString(doc))
        docs := make([]*pb.CrawlDoc,0)
        for _,link := range doc.IndomainOutlinks {
            docs = append(docs,&pb.CrawlDoc{
                RequestUrl:link.Url,
                CrawlParam:&pb.CrawlParam{
                    Rtype:pb.RequestType_WEB_CONTENT, // next level
                },
            })
            if len(docs)>3 {
                break
            }
            LOG.VLog(2).DebugTag("XXXXXX","RequestType_WEB_MAIN set FreshDoc %s",link.Url)
        }
        return docs
    case pb.RequestType_WEB_HUB:
    case pb.RequestType_WEB_CONTENT:
        LOG.VLog(2).DebugTag("XXXXXX","RequestType_WEB_CONTENT %s",pb.FromProtoToString(doc))
    case pb.RequestType_WEB_DETAIL:
    default:
    }
    return nil
}

// use for create instance from a string
func init() {
    registerTaskItfType(&DummyTask{})
}