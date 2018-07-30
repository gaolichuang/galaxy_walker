package sqlite
import (
    _ "galaxy_walker/internal/github.com/mattn/go-sqlite3"
    pb "galaxy_walker/src/proto"
    "database/sql"
    LOG "galaxy_walker/internal/gcodebase/log"
    "galaxy_walker/internal/gcodebase/file"
)
/*
name
status
expire
description

*/
const (
    kCreateTaskTable = `
CREATE TABLE IF NOT EXISTS task (
    name VARCHAR(255) PRIMARY KEY,
    status VARCHAR(255),
    createTimeStamp int,
    expireTimeStamp int,
    desc TEXT
);`
    kInsertTask = `insert into task (name,status,createTimeStamp,expireTimeStamp,desc) values `
)

func createTaskTable(dbname string) {
    execSQL(dbname,kCreateTaskTable)
}

type TaskDbBySQLite struct {
    dbname string
}

func NewTaskDbBySQLite() *TaskDbBySQLite {
    dbname := file.GetConfFile(*CONF.Crawler.UrlDbSQLiteFile)
    createTaskTable(dbname)
    return &TaskDbBySQLite{dbname:dbname}
}

func (t *TaskDbBySQLite)Put(task *pb.TaskDescription) error {
    db, err := sql.Open(kDriverName, t.dbname)
    if err != nil {
        LOG.Errorf("SQL Err %v", err)
        return err
    }
    defer db.Close()
}
func (t *TaskDbBySQLite) Update(task string,status string, des *pb.JobDescription) error {

}
func (t *TaskDbBySQLite) Get(task string) *pb.TaskDescription {
    return nil
}