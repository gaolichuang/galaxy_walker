package task

import (
    pb "galaxy_walker/src/proto"
)

type DummyTask struct {

}

func (d *DummyTask)TaskId() string {
    return "dummy"
}
func (d *DummyTask)GetJobDescription() *JobDescription {
    return nil
}
func (d *DummyTask)Process(rtype pb.RequestType, doc *pb.CrawlDoc) []*pb.CrawlDoc {
    switch rtype {
    case pb.RequestType_WEB_StartUp:
        // start url, no need doc. just return docs.
    case pb.RequestType_WEB_MAIN:
        // startup docs. parse and return web hub. mark RequestType.
    case pb.RequestType_WEB_HUB:
    case pb.RequestType_WEB_CONTENT:
    case pb.RequestType_WEB_DETAIL:
    default:
    }
    return nil
}

// use for create instance from a string
func init() {
    registerTaskItfType(&DummyTask{})
}