package leveldb

import (
    pb "galaxy_walker/src/proto"
    "github.com/syndtr/goleveldb/leveldb"
    "path/filepath"
    "galaxy_walker/internal/gcodebase/file"
    LOG "galaxy_walker/internal/gcodebase/log"
    "galaxy_walker/internal/gcodebase/time_util"
    "galaxy_walker/internal/gcodebase/conf"

    "fmt"
    "strings"
    "strconv"
    "galaxy_walker/internal/github.com/golang/protobuf/proto"
    "github.com/syndtr/goleveldb/leveldb/util"
    "math"
)

/*
    leveldb支持content需要两张表
    contentdb: schema : [task + timestamp + docid],顺序写
        支持的操作： SaveBatch  ScanByTimeRange

    mergeddb: schema : [task + docid + timestamp]
        通过离线任务，定期将contentdb合成到mergeddb,记录处理顺序处理位置；顺序读contentdb，随机读/写mergeddb
        通过scan contentdb；将处理结果按照docid merge在一起，按照时间顺序保留N副本
        支持操作：GetDocById     ScanWithIterator    DeleteBatch Purge

    schema
        timestamp second 长度 10
        docid   长度 2^32-1=4294967295  10 需要前补位0
        task    长度 暂定20 需要前补位0

*/
/////////////////////////////////
const (
    kContentDbTaskFieldLength = 20
    kContentDBPrefix          = "CT"
    kMergedDBPrefix           = "MD"
    kSubFixPath       = "leveldb"
    kMaxTimeStampKey = 3000000000
    kMaxDocId = 1<<32-1
    kConstMergedbBatchSize = 2000
)
var CONF = conf.Conf

var Db *leveldb.DB
var DbConf string

func InitDb() {
    dbfile := filepath.Join(*CONF.ConfPathPrefix,*CONF.Crawler.ContentDbLevelDbFile,kSubFixPath)
    if !file.Exist(dbfile) {
        err := file.MkDirAll(dbfile)
        if err != nil {
            LOG.Fatalf("Mkdir Leveldb Path Err:%v, use %s", err, dbfile)
        }
    }
    DbConf = dbfile

    var err error
    Db, err = leveldb.OpenFile(DbConf, nil)
    if err != nil {
        LOG.Fatalf("Open Leveldb Err:%v, use %s", err, DbConf)
    }
    defer Db.Close()
    LOG.Infof("Init Leveldb use %s", DbConf)
}
/////////////////////////////////
func contentDbKey(task string, docid uint32, timestamp int64) string {
    // CT + task(20) + timestamp(10) + docid(10) = 42
    if len(task) > kContentDbTaskFieldLength {
        task = task[:kContentDbTaskFieldLength]
    } else {
        pad := strings.Repeat("0", kContentDbTaskFieldLength-len(task))
        task = pad + task
    }
    return fmt.Sprintf("%s%s%.10d%.10d", kContentDBPrefix, task, timestamp, docid)
}
func parseContentDbKey(key string) (err error, task string, docid uint32, timestamp int64) {
    if len(key) != 42 {
        err = fmt.Errorf("invalid length %d",len(key))
        return
    }
    if !strings.HasPrefix(key,kContentDBPrefix) {
        err = fmt.Errorf("not contentdb prefix %s",key)
        return
    }
    i := 2
    for ;i<22&&key[i] == '0';i++ {}
    task = key[i:22]
    t,e := strconv.ParseInt(key[22:32],10,0)
    if e != nil {
        err = e
        return
    }
    timestamp = int64(t)
    d,e := strconv.ParseUint(key[32:],10,0)
    if e != nil {
        err = e
        return
    }
    docid = uint32(d)
    return
}
func mergedDbKey(task string, docid uint32, timestamp int64) string {
    // MD + task(20) + docid(10) + timestamp(10) = 42
    if len(task) > kContentDbTaskFieldLength {
        task = task[:kContentDbTaskFieldLength]
    } else {
        pad := strings.Repeat("0", kContentDbTaskFieldLength-len(task))
        task = pad + task
    }
    return fmt.Sprintf("%s%s%.10d%.10d", kMergedDBPrefix, task, docid, timestamp)
}
func parseMergedDbKey(key string) (err error, task string, docid uint32, timestamp int64) {
    if len(key) != 42 {
        err = fmt.Errorf("invalid length %d",len(key))
        return
    }
    if !strings.HasPrefix(key,kMergedDBPrefix) {
        err = fmt.Errorf("not mergeddb prefix %s",key)
        return
    }
    i := 2
    for ;i<22&&key[i] == '0';i++ {}
    task = key[i:22]
    d,e := strconv.ParseUint(key[22:32],10,0)
    if e != nil {
        err = e
        return
    }
    docid = uint32(d)
    t,e := strconv.ParseInt(key[32:],10,0)
    if e != nil {
        err = e
        return
    }
    timestamp = int64(t)
    return
}

/////////////////////////////////
func NewContentDbByLevelDB() *ContentDBByLevelDB {
    InitDb()
    return &ContentDBByLevelDB{}
}
type ContentDBByLevelDB struct {
}

// write into ContentDb
func (db *ContentDBByLevelDB) SaveBatch(task string, docs []*pb.CrawlDoc) (error, int) {
    var errr error
    Db, errr = leveldb.OpenFile(DbConf, nil)
    if errr != nil {
        LOG.Fatalf("Open Leveldb Err:%v, use %s", errr, DbConf)
    }
    defer Db.Close()
    t1 := time_util.GetTimeInMs()
    batch := new(leveldb.Batch)
    for _,doc := range docs {
        var tm int64
        if doc.CrawlRecord != nil && doc.CrawlRecord.FetchTime > 0 {
            tm = doc.CrawlRecord.FetchTime
        } else {
            tm = time_util.GetCurrentTimeStamp()
        }
        key := contentDbKey(task,doc.Docid,tm)
        value,e := proto.Marshal(doc)
        if e != nil {
            LOG.VLog(2).DebugTag("LevelDb", "SaveBatch %s proto marshal err %v", string(key), e)
            continue
        }
        batch.Put([]byte(key),value)
        LOG.VLog(5).DebugTag("LevelDb", "insert key(%s)", string(key))
    }
    err := Db.Write(batch, nil)
    if err != nil {
        return err, -1
    }
    LOG.VLog(4).DebugTag("LevelDb", "SaveBatch %d record use %d ms", batch.Len(), time_util.GetTimeInMs()-t1)
    return nil, batch.Len()
}

/*
    scan by timestamp, if cross maxnum, only return maxnumber
    result will be sorted by crawl timestamp.
    read from contentdb
*/
func (db *ContentDBByLevelDB) ScanByTimeRange(task string, start,end int64, maxnum int) (error, int64, []*pb.CrawlDoc) {
    var errr error
    Db, errr = leveldb.OpenFile(DbConf, nil)
    if errr != nil {
        LOG.Fatalf("Open Leveldb Err:%v, use %s", errr, DbConf)
    }
    defer Db.Close()
    t1 := time_util.GetTimeInMs()
    skey := contentDbKey(task,0,start)
    if end <= 0 {end = kMaxTimeStampKey}
    ekey := contentDbKey(task,kMaxDocId,end)
    iter := Db.NewIterator(&util.Range{
        Start:[]byte(skey),
        Limit:[]byte(ekey),
    },nil)
    ret := make([]*pb.CrawlDoc,0)
    num := 0
    if maxnum <=0 {maxnum = math.MaxInt32}
    var nextkey []byte
    for iter.Next()&&num<maxnum {
        nextkey = iter.Key()
        value := iter.Value()
        doc := new(pb.CrawlDoc)
        err := proto.Unmarshal(value,doc)
        if err != nil {
            LOG.VLog(2).DebugTag("LevelDb", "ScanByTimeRange %s proto marshal err %v", string(nextkey), err)
            continue
        }
        ret = append(ret,doc)
        num += 1
    }
    err,_,_,tm := parseContentDbKey(string(nextkey))
    if err != nil {tm = end}
    LOG.VLog(3).DebugTag("Leveldb", "ScanByTimeRange %s %d,%d, %d record use %d ms",task,start,end,len(ret), time_util.GetTimeInMs()-t1)
    return nil,tm,ret
}
// content db
func (db *ContentDBByLevelDB) ScanKeyByTimeRange(task string,start,end int64) []string {
    var errr error
    Db, errr = leveldb.OpenFile(DbConf, nil)
    if errr != nil {
        LOG.Fatalf("Open Leveldb Err:%v, use %s", errr, DbConf)
    }
    defer Db.Close()
    t1 := time_util.GetTimeInMs()
    skey := contentDbKey(task,0,start)
    if end <= 0 {end = kMaxTimeStampKey}
    ekey := contentDbKey(task,kMaxDocId,end)
    iter := Db.NewIterator(&util.Range{
        Start:[]byte(skey),
        Limit:[]byte(ekey),
    },nil)
    ret := make([]string,0)
    for iter.Next() {
        ret = append(ret,string(iter.Key()))
    }
    LOG.VLog(3).DebugTag("Leveldb", "ScanKeyByTimeRange %s %d,%d, %d record use %d ms",task,start,end,len(ret), time_util.GetTimeInMs()-t1)
    return ret
}

/////////////////////////////////////////////////////////////
// read one doc by docid from mergeddb.
func (db *ContentDBByLevelDB) GetDocById(task string, docid uint32) []*pb.CrawlDoc {
    var errr error
    Db, errr = leveldb.OpenFile(DbConf, nil)
    if errr != nil {
        LOG.Fatalf("Open Leveldb Err:%v, use %s", errr, DbConf)
    }
    defer Db.Close()
    t1 := time_util.GetTimeInMs()
    skey := mergedDbKey(task,docid,0)
    ekey := mergedDbKey(task,docid,kMaxTimeStampKey)
    iter := Db.NewIterator(&util.Range{
        Start:[]byte(skey),
        Limit:[]byte(ekey),
    },nil)
    ret := make([]*pb.CrawlDoc,0)
    for iter.Next() {
        value := iter.Value()
        doc := new(pb.CrawlDoc)
        err := proto.Unmarshal(value,doc)
        if err != nil {
            LOG.VLog(2).DebugTag("LevelDb", "ScanByTimeRange %s proto marshal err %v", string(iter.Key()), err)
            continue
        }
        ret = append(ret,doc)
    }
    LOG.VLog(3).DebugTag("Leveldb", "GetDocById %s %d record use %d ms",task,len(ret), time_util.GetTimeInMs()-t1)
    return ret
}
func getKeyDocById(task string, docid uint32) [][]byte {
    var errr error
    Db, errr = leveldb.OpenFile(DbConf, nil)
    if errr != nil {
        LOG.Fatalf("Open Leveldb Err:%v, use %s", errr, DbConf)
    }
    defer Db.Close()
    t1 := time_util.GetTimeInMs()
    skey := mergedDbKey(task,docid,0)
    ekey := mergedDbKey(task,docid,kMaxTimeStampKey)
    iter := Db.NewIterator(&util.Range{
        Start:[]byte(skey),
        Limit:[]byte(ekey),
    },nil)
    ret := make([][]byte,0)
    for iter.Next() {
        ret = append(ret,iter.Key())
    }
    LOG.VLog(3).DebugTag("Leveldb", "getKeyDocById %s %d record use %d ms",task,len(ret), time_util.GetTimeInMs()-t1)
    return ret
}

/*
from merged db
*/
func (db *ContentDBByLevelDB) ScanWithIterator(task, iterator string, maxnum int) (string, []*pb.CrawlDoc) {
    var errr error
    Db, errr = leveldb.OpenFile(DbConf, nil)
    if errr != nil {
        LOG.Fatalf("Open Leveldb Err:%v, use %s", errr, DbConf)
    }
    defer Db.Close()
    t1 := time_util.GetTimeInMs()
    skey := iterator
    if skey == "" {
        skey = mergedDbKey(task,0,0)
    }
    ekey := mergedDbKey(task,kMaxDocId,kMaxTimeStampKey)
    iter := Db.NewIterator(&util.Range{
        Start:[]byte(skey),
        Limit:[]byte(ekey),
    },nil)
    ret := make([]*pb.CrawlDoc,0)
    num := 0
    if maxnum <=0 {maxnum = kConstMergedbBatchSize}
    var nextkey []byte
    for iter.Next()&&num<maxnum {
        nextkey = iter.Key()
        value := iter.Value()
        doc := new(pb.CrawlDoc)
        err := proto.Unmarshal(value,doc)
        if err != nil {
            LOG.VLog(2).DebugTag("LevelDb", "ScanWithIterator %s proto marshal err %v", string(nextkey), err)
            continue
        }
        ret = append(ret,doc)
        num += 1
    }
    LOG.VLog(3).DebugTag("Leveldb", "ScanWithIterator %s %d, nextkey:%s record use %d ms",task,len(ret),nextkey,time_util.GetTimeInMs()-t1)
    nextkey[len(nextkey)-1] += 1
    return string(nextkey),ret
}
// from mergeddb
func (db *ContentDBByLevelDB)ScanKey(task string) []string {
    var errr error
    Db, errr = leveldb.OpenFile(DbConf, nil)
    if errr != nil {
        LOG.Fatalf("Open Leveldb Err:%v, use %s", errr, DbConf)
    }
    defer Db.Close()
    t1 := time_util.GetTimeInMs()
    skey := mergedDbKey(task,0,0)
    ekey := mergedDbKey(task,kMaxDocId,kMaxTimeStampKey)
    iter := Db.NewIterator(&util.Range{
        Start:[]byte(skey),
        Limit:[]byte(ekey),
    },nil)
    ret := make([]string,0)
    for iter.Next() {
        ret = append(ret,string(iter.Key()))
    }
    LOG.VLog(3).DebugTag("Leveldb", "ScanKey %s %d record use %d ms",task,len(ret), time_util.GetTimeInMs()-t1)
    return ret
}

// delete batch by docid. delete from mergeddb
func (db *ContentDBByLevelDB) DeleteBatch(task string, docids []uint32) (error, int) {
    var errr error
    Db, errr = leveldb.OpenFile(DbConf, nil)
    if errr != nil {
        LOG.Fatalf("Open Leveldb Err:%v, use %s", errr, DbConf)
    }
    defer Db.Close()

    t1 := time_util.GetTimeInMs()
    batch := new(leveldb.Batch)
    for _,docid := range docids {
        keys := getKeyDocById(task,docid)
        for _,k := range keys {
            batch.Delete(k)
        }
    }
    err := Db.Write(batch, nil)
    if err != nil {
        return err, -1
    }
    LOG.VLog(3).DebugTag("Leveldb", "deleta %s batch %d record use %d ms",task, batch.Len(), time_util.GetTimeInMs()-t1)
    return nil, batch.Len()
}

// delete all both contentdb,and mergeddb
func (db *ContentDBByLevelDB) Purge(task string) (error,int) {
    var errr error
    Db, errr = leveldb.OpenFile(DbConf, nil)
    if errr != nil {
        LOG.Fatalf("Open Leveldb Err:%v, use %s", errr, DbConf)
    }
    defer Db.Close()

    t1 := time_util.GetTimeInMs()
    batch := new(leveldb.Batch)
    for _,k := range db.ScanKey(task) {
            batch.Delete([]byte(k))
    }
    err := Db.Write(batch, nil)
    if err != nil {
        return err, -1
    }
    LOG.VLog(3).DebugTag("Leveldb", "Purge %s %d record use %d ms",task, batch.Len(), time_util.GetTimeInMs()-t1)
    return nil, batch.Len()
}
