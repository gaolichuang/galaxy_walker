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
    return nil
}

// use for create instance from a string
func init() {
    registerTaskItfType(&DummyTask{})
}