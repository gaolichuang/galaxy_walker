package task

import (
    pb "galaxy_walker/src/proto"
    "galaxy_walker/src/crawl/scheduler"
    LOG "galaxy_walker/internal/gcodebase/log"
)

type LinearDummyTask struct {
    LinearTask
}
func (d *LinearDummyTask)Init() error {
    d.RegisterRequestTypeCallBack(pb.RequestType_WEB_StartUp,pb.RequestType_WEB_MAIN,d.startUpCallBack)
    d.RegisterRequestTypeCallBack(pb.RequestType_WEB_MAIN,pb.RequestType_WEB_CONTENT,d.webmainCallBack)
    d.RegisterRequestTypeCallBack(pb.RequestType_WEB_CONTENT,pb.RequestType_WEB_End,d.webcontentCallBack)
    return nil
}
func (d *LinearDummyTask) GetJobDescription() *pb.JobDescription {
    return &scheduler.NormalJobD
}

func (d *LinearDummyTask)startUpCallBack(rtype pb.RequestType, doc *pb.CrawlDoc) []string {
    return  []string{
        "http://roll.sohu.com/sports/",
    }
}
func (d *LinearDummyTask)webmainCallBack(rtype pb.RequestType, doc *pb.CrawlDoc) []string {
    LOG.VLog(2).DebugTag("XXXXXX","RequestType_WEB_MAIN %s",pb.FromProtoToString(doc))
    urls := make([]string,0)
    for _,link := range doc.IndomainOutlinks {
        urls = append(urls,link.Url)
        if len(urls) > 3 {
            break
        }
        LOG.VLog(2).DebugTag("XXXXXX","RequestType_WEB_MAIN set FreshDoc %s",link.Url)
    }
    return urls
}
func (d *LinearDummyTask)webcontentCallBack(rtype pb.RequestType, doc *pb.CrawlDoc) []string {
    LOG.VLog(2).DebugTag("XXXXXX","RequestType_WEB_CONTENT %s",pb.FromProtoToString(doc))
    return nil
}

// use for create instance from a string
func init() {
    registerTaskItfType(&LinearDummyTask{})
}