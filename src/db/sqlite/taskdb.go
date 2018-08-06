package sqlite
import (
    _ "galaxy_walker/internal/github.com/mattn/go-sqlite3"
    pb "galaxy_walker/src/proto"
    "database/sql"
    LOG "galaxy_walker/internal/gcodebase/log"
    "galaxy_walker/internal/gcodebase/file"
    "galaxy_walker/internal/gcodebase/time_util"
    "fmt"
    "strings"
    "encoding/json"
    "path/filepath"
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
    kInsertTaskSQL = `insert into task (name,status,createTimeStamp,expireTimeStamp,desc) values `
    kUpdateTaskSQL = `update task set %s where name="%s"`
    kSelectTaskByNameSQL = `select name,status,createTimeStamp,expireTimeStamp,desc from task where name="%s"`
    kSelectTaskSQL = `select name,status,createTimeStamp,expireTimeStamp,desc from task`
    kDeleteTaskSQL = `delete from task where name="%s"`
)

func createTaskTable(dbname string) {
    execSQL(dbname,kCreateTaskTable)
}

type TaskDbBySQLite struct {
    dbname string
}

func NewTaskDbBySQLite() *TaskDbBySQLite {
    dbfile := filepath.Join(*CONF.ConfPathPrefix, *CONF.Crawler.TaskDbSQLiteFile)
    if !file.Exist(dbfile) {
        err := file.WriteStringToFile([]byte(""),dbfile)
        if err != nil {
            LOG.Fatalf("Touch UrlDB Err:%v, use %s", err, dbfile)
        }
    }
    createTaskTable(dbfile)
    return &TaskDbBySQLite{dbname:dbfile}
}

func (t *TaskDbBySQLite)Delete(name string) error  {
    err,_ := execSQL(t.dbname,fmt.Sprintf(kDeleteTaskSQL,name))
    return err
}
func (t *TaskDbBySQLite)Put(task *pb.TaskDescription) error {
    db, err := sql.Open(kDriverName, t.dbname)
    if err != nil {
        LOG.Errorf("SQL Err %v", err)
        return err
    }
    defer db.Close()
    sqlsmt := fmt.Sprintf(`%s ("%s","%s","%d","%d",'%s')`,kInsertTaskSQL,task.Name,task.Status,task.CreateAt,task.ExpireAt,task.Desc.ToString())
    _, err = db.Exec(sqlsmt)
    if err != nil {
        LOG.Errorf("Inset SQL Err %v \n%s", err, sqlsmt)
        return err
    }
    LOG.VLog(4).DebugTag("SQLRECORD", sqlsmt)
    return nil
}
func (t *TaskDbBySQLite) Update(task string,expire int64,status string, des *pb.JobDescription) error {
    sets := make([]string,0)
    if status != "" {
        sets = append(sets,fmt.Sprintf(`status="%s"`,status))
    }
    if expire > 0 {
        sets = append(sets,fmt.Sprintf(`expireTimeStamp=%d`,expire))
    }
    if des != nil {
        sets = append(sets,fmt.Sprintf(`des="%s"`,des.ToString()))
    }
    if len(sets) <= 0 {return nil}
    sqlsmt := fmt.Sprintf(kUpdateTaskSQL,strings.Join(sets,","),task)
    err,_ := execSQL(t.dbname,sqlsmt)
    return err
}
func (t *TaskDbBySQLite) Get(task string) *pb.TaskDescription {
    dbname := t.dbname
    db, err := sql.Open(kDriverName, dbname)
    if err != nil {
        LOG.Errorf("SQL Err %v from %s", err, dbname)
        return nil
    }
    defer db.Close()
    //    kSelectTaskByNameSQL = `select name,status,createTimeStamp,expireTimeStamp,desc where name=%s`
    sqlsmt := fmt.Sprintf(kSelectTaskByNameSQL,task)
    rows, err := db.Query(sqlsmt)
    if err != nil {
        LOG.Errorf("SQL %s Err %v from %s", sqlsmt, err, dbname)
        return nil
    }
    defer rows.Close()
    var obj *pb.TaskDescription
    for rows.Next() {
        var name,status,desc string
        var createat,expireat int
        err = rows.Scan(&name,&status,&createat,&expireat,&desc)
        if err != nil {
            LOG.Errorf("SQL Err %v", err)
            continue
        }
        obj = &pb.TaskDescription{
            Name:name,
            Status:status,
            CreateAt:createat,
            ExpireAt:expireat,
        }
        desObj := &pb.JobDescription{}
        err = json.Unmarshal([]byte(desc),desObj)
        if err != nil {
            LOG.Errorf("UnMarshal %s desc err %v",task,err)
            continue
        }
        obj.Desc=desObj
    }
    return obj
}
func (t *TaskDbBySQLite) List() ([]string,[]*pb.TaskDescription) {
    dbname := t.dbname
    db, err := sql.Open(kDriverName, dbname)
    if err != nil {
        LOG.Errorf("SQL Err %v from %s", err, dbname)
        return nil,nil
    }
    defer db.Close()
    //    kSelectTaskByNameSQL = `select name,status,createTimeStamp,expireTimeStamp,desc where name=%s`
    sqlsmt := kSelectTaskSQL
    rows, err := db.Query(sqlsmt)
    if err != nil {
        LOG.Errorf("SQL %s Err %v from %s", sqlsmt, err, dbname)
        return nil,nil
    }
    defer rows.Close()
    t1 := time_util.GetTimeInMs()
    tasks := make([]string,0)
    descs := make([]*pb.TaskDescription,0)
    for rows.Next() {
        var name,status,desc string
        var createat,expireat int
        err = rows.Scan(&name,&status,&createat,&expireat,&desc)
        if err != nil {
            LOG.Errorf("SQL Err %v", err)
            continue
        }
        obj := &pb.TaskDescription{
            Name:name,
            Status:status,
            CreateAt:createat,
            ExpireAt:expireat,
        }
        desObj := &pb.JobDescription{}
        err = json.Unmarshal([]byte(desc),desObj)
        if err != nil {
            LOG.Errorf("UnMarshal %s desc err %v",name,err)
        } else {
            obj.Desc=desObj
        }
        descs = append(descs,obj)
        tasks = append(tasks,name)
    }
    LOG.VLog(4).DebugTag("SQL", "%d record read use %d ms by %s from %s", len(tasks), time_util.GetTimeInMs() - t1, sqlsmt, dbname)
    return tasks,descs
}
