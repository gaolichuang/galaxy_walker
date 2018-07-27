package sqlite

import (
    _ "galaxy_walker/internal/github.com/mattn/go-sqlite3"
    pb "galaxy_walker/src/proto"
    "database/sql"
    LOG "galaxy_walker/internal/gcodebase/log"
    "galaxy_walker/internal/gcodebase/time_util"
    "galaxy_walker/internal/gcodebase/conf"
    "strings"
    "fmt"
    "gcodebase/file"
)

var CONF = conf.Conf

/*
each task use one table.
统计需求
1.已发现未抓取    status=0
2.抓取失败 / 成功  status=1,2
3.抓取失败不再重试的  retry>X
4.N次失败的统计   status + retry
*/
const (
    kDriverName      = "sqlite3"
    kWriteBatchSize  = 100
    kUrlTablePrefix  = "urldb_"
    kRetryMaxNum = 5

    /*
                         +--------+
                         |  Undo  |
                         +--------+
                             |       Scan
                         Scan   +------------+
                             |  |            |
                             |  |            |
       +------+   Success+---v--v+  Fail  +-----+
       | Done <----------+ Doing +--------> Fail|
       +------+          +-------+        +-----+
    */
    // copy from urldb.
    KUrlStatusUndo  = 0 //  not crawl
    KUrlStatusDoing = 1 // crawling.
    KUrlStatusDone  = 2 // success
    KUrlStatusFail  = 3 // fail

    kCreateUrlTableVersion = `
CREATE TABLE IF NOT EXISTS urldb_%s (
    url VARCHAR(255) PRIMARY KEY,
    rtype int,
    parentType int,
    parentDocid int,
    status int default 0,
    createTimeStamp int,
    updateTimeStamp int,
    retryNum int default 0
);`
    kListTable                      = `select name from sqlite_master where type = 'table';`
    kInsertIntoTaskTable            = `insert into urldb_%s (url,rtype,parentType,parentDocid,createTimeStamp,updateTimeStamp) values `
    kUpdateTaskStatus               = `update urldb_%s set status="%d",updateTimeStamp="%d" where url="%s"`
    kUpdateTaskStatusAndRetry       = `update urldb_%s set status="%d",updateTimeStamp="%d",retryNum="%d where url="%s"`
    kSelectFromTaskByStatus         = `select url,rtype,parentType,parentDocid,status,createTimeStamp,updateTimeStamp,retryNum from urldb_%s where status=%d limit %d`
    kSelectFromTaskByStatusAndRetry = `select url,rtype,parentType,parentDocid,status,createTimeStamp,updateTimeStamp,retryNum from urldb_%s where status=%d and retryNum<%d limit %d`
    kSelectFromTaskByUrl = `select url from urldb_%s where %s`
    kSelectAllFromTaskByUrl = `select url from urldb_%s `
)

func ListTable(dbname string, prefix string) []string {
    tables := make([]string, 0)

    db, err := sql.Open(kDriverName, dbname)
    if err != nil {
        LOG.Errorf("SQL Err %v from %s", err, dbname)
        return nil
    }
    defer db.Close()
    t1 := time_util.GetTimeInMs()
    rows, err := db.Query(kListTable)
    if err != nil {
        LOG.Errorf("SQL %s Err %v from %s", kListTable, err, dbname)
        return nil
    }
    defer rows.Close()
    recordNum := 0
    for rows.Next() {
        var t string
        err = rows.Scan(&t)
        if err != nil {
            LOG.Errorf("SQL Err %v", err)
            continue
        }
        if prefix == "" || strings.HasPrefix(t, kUrlTablePrefix) {
            tables = append(tables, t)
        }
    }
    LOG.VLog(4).DebugTag("SQL", "%d record read use %d ms by %s from %s", recordNum, time_util.GetTimeInMs()-t1, kListTable, dbname)
    return tables
}
func execSQL(dbfile, smt string) (error, int) {
    db, err := sql.Open(kDriverName, dbfile)
    if err != nil {
        LOG.Errorf("SQL Open %s Err %v", dbfile, err)
        return err, -1
    }
    defer db.Close()
    r, err := db.Exec(smt)
    if err != nil {
        LOG.Errorf("SQL %s Err %v :\n%s", dbfile, err, smt)
        return err, -1
    }
    LOG.VLog(2).DebugTag("SQLRECORD", smt)
    num, _ := r.RowsAffected()
    if num == 0 {
        LOG.VLog(4).DebugTag("XXXXXX", "%s exec affect 0 record.", smt)
    }
    return nil, int(num)
}

//////////////////////////////////////////////
type UrlDBBySQLite struct {
    tables map[string]bool
    dbname string
}

func NewUrlDbBySQLite() *UrlDBBySQLite {
    dbname := file.GetConfFile(*CONF.Crawler.UrlDbSQLiteFile)
    tables := make(map[string]bool)
    for _, t := range ListTable(dbname, kUrlTablePrefix) {
        tables[t] = true
    }
    u := &UrlDBBySQLite{
        dbname: dbname,
        tables: tables,
    }
    return u
}
func (u *UrlDBBySQLite) createTaskTableIfNotExist(task string) {
    if _, ok := u.tables[task]; ok {
        return
    }
    execSQL(u.dbname, fmt.Sprintf(kCreateUrlTableVersion, task))
    tables := make(map[string]bool)
    for _, t := range ListTable(u.dbname, kUrlTablePrefix) {
        tables[t] = true
    }
    u.tables = tables
}
func (u *UrlDBBySQLite) ListUrls(task string,status int,lastTimeStamp int64) []string {
    // status == -1, lasttimestamp <=0 will list all urls for task
    listUrls := func(sqlsmt string) []string {
        dbname := u.dbname
        db, err := sql.Open(kDriverName, dbname)
        if err != nil {
            LOG.Errorf("SQL Err %v from %s", err, dbname)
            return nil
        }
        defer db.Close()
        t1 := time_util.GetTimeInMs()
        rows, err := db.Query(sqlsmt)
        if err != nil {
            LOG.Errorf("SQL %s Err %v from %s", sqlsmt, err, dbname)
            return nil
        }
        defer rows.Close()
        ret := make([]string,0)
        for rows.Next() {
            var url string
            err = rows.Scan(&url)
            if err != nil {
                LOG.Errorf("SQL Err %v", err)
                continue
            }
            ret=append(ret,url)
        }
        LOG.VLog(4).DebugTag("SQL", "%d record read use %d ms by %s from %s", len(ret), time_util.GetTimeInMs() - t1, sqlsmt, dbname)
        return ret
    }
    sqlsmt := fmt.Sprintf(kSelectAllFromTaskByUrl,task)
    filter := make([]string,0)
    if status >= 0 {
        filter = append(filter,fmt.Sprintf(" status==%d ",status))
    }
    if lastTimeStamp > 0 {
        filter = append(filter,fmt.Sprintf(" updateTimeStamp>%d ",lastTimeStamp))
    }
    if len(filter)>0 {
        sqlsmt += " where " + strings.Join(filter,"and")
    }
    return listUrls(sqlsmt)
}
func (u *UrlDBBySQLite) listUrlsByWhiteList(task string, urls[]string) []string {
    listUrls := func(task string,urls []string) []string {
        dbname := u.dbname
        db, err := sql.Open(kDriverName, dbname)
        if err != nil {
            LOG.Errorf("SQL Err %v from %s", err, dbname)
            return nil
        }
        defer db.Close()
        t1 := time_util.GetTimeInMs()
        record := make([]string,0)
        for _,url := range urls {
            record = append(record,fmt.Sprintf(" url=\"%s\" ",url))
        }
        sqlsmt := fmt.Sprintf(kSelectFromTaskByUrl,task,strings.Join(record, "or"))
        rows, err := db.Query(sqlsmt)
        if err != nil {
            LOG.Errorf("SQL %s Err %v from %s", sqlsmt, err, dbname)
            return nil
        }
        defer rows.Close()
        ret := make([]string,0)
        for rows.Next() {
            var url string
            err = rows.Scan(&url)
            if err != nil {
                LOG.Errorf("SQL Err %v", err)
                continue
            }
            ret=append(ret,url)
        }
        LOG.VLog(4).DebugTag("SQL", "%d record read use %d ms by %s from %s", len(ret), time_util.GetTimeInMs() - t1, sqlsmt, dbname)
        return ret

    }

    ret := make([]string,0)
    times := len(urls)/kWriteBatchSize
    for i:=0;i<times;i++ {
        r := listUrls(task,urls[i*kWriteBatchSize:(i+1)*kWriteBatchSize])
        ret = append(ret,r...)
    }
    r:=listUrls(task,urls[times*kWriteBatchSize:])
    ret = append(ret,r...)
    return ret
}
func (u *UrlDBBySQLite) SetFreshUrls(task string, parentType pb.RequestType, parentDocid uint32, docs []*pb.CrawlDoc) (error, int) {
    /*
       新发现的url，如果task表不存在，则创建表
       需要去重复。。。使用task+docid
    // TODO. insert 之前先检查是否存在。
    */
    u.createTaskTableIfNotExist(task)
    urls := make([]string,0)
    for _,doc := range docs {
        urls = append(urls,doc.Url)
    }
    urlsBlack := make(map[string]bool)
    for _,url := range u.listUrlsByWhiteList(task,urls) {
        urlsBlack[url]=true
    }
    //    kInsertIntoTaskTable = `insert into urldb_%s (url,parentType,parentDocid,createTimeStamp,updateTimeStamp) values `
    db, err := sql.Open(kDriverName, u.dbname)
    if err != nil {
        LOG.Errorf("SQL Err %v", err)
        return err, -1
    }
    defer db.Close()
    insertSmt := fmt.Sprintf(kInsertIntoTaskTable, task)
    now := time_util.GetCurrentTimeStamp()
    t1 := time_util.GetTimeInMs()
    record := make([]string, 0)
    var affectNum int64 = 0
    for _, doc := range docs {
        if urlsBlack[doc.Url] {
            continue
        }
        rtype := 0
        if doc.CrawlParam != nil {rtype = int(doc.CrawlParam.Rtype)}
        record = append(record,
            fmt.Sprintf(`("%s","%d","%d","%d","%d",%d)`,
                doc.Url,rtype, parentType, parentDocid, now, now))
        if len(record) > kWriteBatchSize {
            sqlsmt := fmt.Sprintf("%s %s", insertSmt, strings.Join(record, ","))
            r, err := db.Exec(sqlsmt)
            if err != nil {
                LOG.Errorf("Inset SQL Err %v \n%s", err, sqlsmt)
                continue
            }
            record = nil
            num, _ := r.RowsAffected()
            affectNum += num
            LOG.VLog(3).DebugTag("SQLRECORD", sqlsmt)
            LOG.VLog(5).DebugTag("SQL", "Insert %s records into %s ", num, task)
        }
    }
    if len(record) > 0 {
        sqlsmt := fmt.Sprintf("%s %s", insertSmt, strings.Join(record, ","))
        r, err := db.Exec(sqlsmt)
        if err != nil {
            LOG.Errorf("Inset SQL Err %v", err)
        } else {
            record = nil
            num, _ := r.RowsAffected()
            affectNum += num
            LOG.VLog(4).DebugTag("SQLRECORD", sqlsmt)
            LOG.VLog(5).DebugTag("SQL", "Insert into %s %d", task, num)
        }
    }
    LOG.VLog(4).DebugTag("SQL", "Insert SetFreshUrls %d record affect %d use %d ms", len(docs), affectNum, time_util.GetTimeInMs()-t1)
    return nil, int(affectNum)

}

func (u *UrlDBBySQLite) MarkCrawlFinishUrls(taskid string, docs []*pb.CrawlDoc) {
    /*
       抓取完成，标记成功失败等; 失败次数更新
       支持多个taskid
    */
    now := time_util.GetCurrentTimeStamp()
    for _, doc := range docs {
        param := doc.CrawlParam
        if param == nil {
            continue
        }
        execSQL(u.dbname, fmt.Sprintf(kUpdateTaskStatusAndRetry, KUrlStatusDone, now, param.RetryNumber, doc.RequestUrl))
    }
}
func (u *UrlDBBySQLite) MarkCrawlFailUrls(taskid string, docs []*pb.CrawlDoc) {
    /*
       抓取完成，标记成功失败等; 失败次数更新
       支持多个taskid
    */
    now := time_util.GetCurrentTimeStamp()
    for _, doc := range docs {
        param := doc.CrawlParam
        if param == nil {
            continue
        }
        execSQL(u.dbname, fmt.Sprintf(kUpdateTaskStatusAndRetry, KUrlStatusFail, now, param.RetryNumber+1, doc.RequestUrl))
    }
}

func (u *UrlDBBySQLite) ScanFreshUrls(task string, num int) (error,[]*pb.CrawlDoc) {
    /*
    1.已发现未抓取的url
    2.抓取失败，仍在重试范围内的url
    */
    t1 := time_util.GetTimeInMs()
    // fresh
    err,docs,urls := u.scanFreshAndUpdateStatus(task,num)
    if err != nil{
        return err,nil
    }
    now := time_util.GetCurrentTimeStamp()
    for _,url := range urls {
        execSQL(u.dbname, fmt.Sprintf(kUpdateTaskStatus, KUrlStatusDoing, now, url))
    }
    // fail
    err,fdocs,furls := u.scanFailAndUpdateStatus(task,num,kRetryMaxNum)
    if err == nil {
        docs = append(docs,fdocs...)
        for _,url := range furls {
            execSQL(u.dbname, fmt.Sprintf(kUpdateTaskStatus, KUrlStatusDoing, now, url))
        }
    }
    LOG.VLog(3).DebugTag("UrlDBBySQLite","scan %d,%d urls from %s use %d ms",len(docs),len(fdocs),task,time_util.GetTimeInMs()-t1)
    return nil,docs
}
func (u *UrlDBBySQLite) scanFreshAndUpdateStatus(task string, num int) (error,[]*pb.CrawlDoc,[]string) {
    dbname := u.dbname
    db, err := sql.Open(kDriverName, dbname)
    if err != nil {
        LOG.Errorf("SQL Err %v from %s", err, dbname)
        return err, nil, nil
    }
    defer db.Close()
    t1 := time_util.GetTimeInMs()
    //kSelectFromTaskByStatus= `select url,parentType,parentDocid,status,createTimeStamp,updateTimeStamp,retryNum from urldb_%s where status="%d"`
    sqlsmt := fmt.Sprintf(kSelectFromTaskByStatus,task,KUrlStatusFail,num)
    rows, err := db.Query(sqlsmt)
    if err != nil {
        LOG.Errorf("SQL %s Err %v from %s", sqlsmt, err, dbname)
        return err, nil, nil
    }
    defer rows.Close()
    ret := make([]*pb.CrawlDoc,0)
    urls := make([]string,0)
    for rows.Next() {
        var url string
        var rtype,parentType,parentDocid,status,createTimeStamp,updateTimeStamp,retrynum int
        err = rows.Scan(&url,&rtype,&parentType,&parentDocid,&status,&createTimeStamp,&updateTimeStamp,&retrynum)
        if err != nil {
            LOG.Errorf("SQL Err %v", err)
            continue
        }
        urls=append(urls,url)
        doc := &pb.CrawlDoc{
            RequestUrl:url,
        }
        doc.CrawlParam = &pb.CrawlParam{
            Taskid:task,
            Rtype:pb.RequestType(rtype),
            ParentDocid:int32(parentDocid),
            ParentRtype:pb.RequestType(parentType),
            RetryNumber:0,
            DiscoverTime:int64(createTimeStamp),
        }
        ret = append(ret,doc)
    }
    LOG.VLog(4).DebugTag("SQL", "%d record read use %d ms by %s from %s", len(ret), time_util.GetTimeInMs() - t1, sqlsmt, dbname)
    return nil, ret,urls
}
func (u *UrlDBBySQLite) scanFailAndUpdateStatus(task string, num int,retry int) (error,[]*pb.CrawlDoc,[]string) {
    dbname := u.dbname
    db, err := sql.Open(kDriverName, dbname)
    if err != nil {
        LOG.Errorf("SQL Err %v from %s", err, dbname)
        return err, nil, nil
    }
    defer db.Close()
    t1 := time_util.GetTimeInMs()
    //kSelectFromTaskByStatus= `select url,parentType,parentDocid,status,createTimeStamp,updateTimeStamp,retryNum from urldb_%s where status="%d"`
    sqlsmt := fmt.Sprintf(kSelectFromTaskByStatusAndRetry,task,KUrlStatusFail,retry,num)
    rows, err := db.Query(sqlsmt)
    if err != nil {
        LOG.Errorf("SQL %s Err %v from %s", sqlsmt, err, dbname)
        return err, nil, nil
    }
    defer rows.Close()
    ret := make([]*pb.CrawlDoc,0)
    urls := make([]string,0)
    for rows.Next() {
        var url string
        var rtype,parentType,parentDocid,status,createTimeStamp,updateTimeStamp,retrynum int
        err = rows.Scan(&url,&rtype,&parentType,&parentDocid,&status,&createTimeStamp,&updateTimeStamp,&retrynum)
        if err != nil {
            LOG.Errorf("SQL Err %v", err)
            continue
        }
        urls=append(urls,url)
        doc := &pb.CrawlDoc{
            RequestUrl:url,
        }
        doc.CrawlParam = &pb.CrawlParam{
            Taskid:task,
            Rtype:pb.RequestType(rtype),
            ParentDocid:int32(parentDocid),
            ParentRtype:pb.RequestType(parentType),
            RetryNumber:int32(retrynum),
            DiscoverTime:int64(createTimeStamp),
        }
        ret = append(ret,doc)
    }
    LOG.VLog(4).DebugTag("SQL", "%d record read use %d ms by %s from %s", len(ret), time_util.GetTimeInMs() - t1, sqlsmt, dbname)
    return nil, ret,urls
}
