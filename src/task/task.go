package task

import (
    pb "galaxy_walker/src/proto"
)

type Task interface {
    /*
    用来描述Task， taskid贯穿始终
    */
    TaskId() string
    GetJobDescription() *JobDescription
    /*
    发现下一级 新的连接
    */
    Process(rtype pb.RequestType, doc *pb.CrawlDoc) []*pb.CrawlDoc
}

