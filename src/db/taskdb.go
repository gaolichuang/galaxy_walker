package db

import (
    pb "galaxy_walker/src/proto"
    "galaxy_walker/src/db/sqlite"
)

/*
每一次都是一个抓取任务，contentdb,urldb都是跟task相关的。
contentdb使用key schema标记task
urldb sqlite实现方式是使用不同的table，即task是表名
*/

type TaskDbItf interface {
    Put(task *pb.TaskDescription) error
    Update(task string,status string, des *pb.JobDescription) error
    Get(task string) *pb.TaskDescription
}
func NewTaskItf() TaskDbItf {
    return sqlite.NewTaskDbBySQLite()
}