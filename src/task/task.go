package task

import (
    pb "galaxy_walker/src/proto"
    "reflect"
)

// should support multi goroutine
type TaskItf interface {
    /*
    用来描述Task， taskid贯穿始终
    */
    TaskId() string
    GetJobDescription() *JobDescription
    /*
    发现下一级 新的连接
    需要标记下级的 requesttype
    第一级type默认是 pb.RequestType_WEB_StartUp

    通过requesttype和process中针对不同requesttype的处理，形成状态机
    */
    Process(rtype pb.RequestType, doc *pb.CrawlDoc) []*pb.CrawlDoc
}


//////////////////////////////////////////////////////
var TaskItfRegistry = make(map[string]TaskItf)

func registerTaskItfType(task TaskItf) {
    t := reflect.TypeOf(task).Elem()
    TaskItfRegistry[t.Name()] = task
}

func GetTaskItfByName(name string) TaskItf {
    return TaskItfRegistry[name]
}
