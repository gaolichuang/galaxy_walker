package leveldb

import (
	"galaxy_walker/internal/github.com/syndtr/goleveldb/leveldb"
	"galaxy_walker/internal/github.com/syndtr/goleveldb/leveldb/util"
	pb "galaxy_walker/src/proto"
	LOG "galaxy_walker/internal/gcodebase/log"
	"galaxy_walker/internal/gcodebase/time_util"
	"galaxy_walker/internal/github.com/golang/protobuf/proto"
	"galaxy_walker/internal/gcodebase/persistent"
)

const (
	kMaxDocReplicationNum = 5
	// DD > CT
	kContentDbEndKey        = "DD0000000000000000000015325011333709056136"
	kScanContentDBBatchSize = 500
	kMergeDBIteratorTag     = "db"
	kmergeDBIteratorKey     = "mergedDBIterator"
)

var ps persistent.Persistent

func init() {
	persistent.InitPersistentTag(kMergeDBIteratorTag)
	ps = persistent.NewPersistent()
}

func MergeContentDbProcess(restart bool) {
	/*
	   scan from ContentDB. use persistentkey
	   Merge into MergedDb. only save kMaxDocReplicationNum for each task-docid
	   only process key.
	*/
	t1 := time_util.GetTimeInMs()
	err, skey := ps.GetWithCategory(kMergeDBIteratorTag, kmergeDBIteratorKey)
	if err != nil {
		skey = ""
	}
	ekey, docsMap := scanContentDBWithIterator(skey, kScanContentDBBatchSize)
	insnum, delnum := 0, 0
	for task, docs := range docsMap {
		docIDMap := make(map[uint32][]*pb.CrawlDoc)
		for _, doc := range docs {
			if _, ok := docIDMap[doc.Docid]; !ok {
				docIDMap[doc.Docid] = make([]*pb.CrawlDoc, 0)
			}
			docIDMap[doc.Docid] = append(docIDMap[doc.Docid], doc)
		}
		ins := make([]*pb.CrawlDoc, 0)
		del := make([][]byte, 0)
		for docid, dcs := range docIDMap {
			indb := getKeyDocById(task, docid)
			if len(dcs) >= kMaxDocReplicationNum {
				del = append(del, indb...)
			} else if len(dcs)+len(indb) > kMaxDocReplicationNum {
				del = append(del, indb[:(len(dcs)+len(indb)-kMaxDocReplicationNum)]...)
			}
			if len(dcs) >= kMaxDocReplicationNum {
				ins = append(ins, dcs[:kMaxDocReplicationNum]...)
			} else {
				ins = append(ins, dcs...)
			}
		}
		insnum += len(ins)
		saveIntoMergeDb(task, ins)
		delnum += len(del)
		deleteBatch(del)
	}
	// persistent contentdb iterator.
	err = ps.SetWithCategory(kMergeDBIteratorTag, kmergeDBIteratorKey, ekey)
	if err != nil {
		LOG.Errorf("MergeContentDbProcess set end key %s err %v", ekey, err)
	}
	LOG.VLog(3).DebugTag("MergeContentDbProcess", " ins %d,del %d record use %d ms, iterator %s=>%s",
		insnum, delnum, time_util.GetTimeInMs()-t1, skey, ekey)
}

func scanContentDBWithIterator(iterator string, maxnum int) (string, map[string][]*pb.CrawlDoc) {
	// support multi task
	var errr error
	Db, errr = leveldb.OpenFile(DbConf, nil)
	if errr != nil {
		LOG.Fatalf("Open Leveldb Err:%v, use %s", errr, DbConf)
	}
	defer Db.Close()
	t1 := time_util.GetTimeInMs()
	skey := iterator
	if skey == "" {
		skey = contentDbKey("", 0, 0)
	}
	iter := Db.NewIterator(&util.Range{
		Start: []byte(skey),
		Limit: []byte(kContentDbEndKey),
	}, nil)

	ret := make(map[string][]*pb.CrawlDoc)

	num := 0
	if maxnum <= 0 {
		maxnum = kScanContentDBBatchSize
	}
	var key []byte = []byte(skey)
	for iter.Next() && num < maxnum {
		key = iter.Key()
		value := iter.Value()
		doc := new(pb.CrawlDoc)
		err := proto.Unmarshal(value, doc)
		if err != nil {
			LOG.VLog(2).DebugTag("LevelDb", "ScanWithIterator %s proto marshal err %v", string(key), err)
			continue
		}
		err, task, _, _ := parseContentDbKey(string(key))
		if err != nil {
			LOG.VLog(2).DebugTag("LevelDb", "ScanWithIterator %s parse key error %v", string(key), err)
			continue
		}
		if _, ok := ret[task]; !ok {
			ret[task] = make([]*pb.CrawlDoc, 0)
		}
		ret[task] = append(ret[task], doc)
		num += 1
	}
	LOG.VLog(3).DebugTag("Leveldb", "scanContentDBWithIterator %d, nextkey:%s record use %d ms", len(ret), key, time_util.GetTimeInMs()-t1)
	key[len(key)-1] += 1
	return string(key), ret
}

func saveIntoMergeDb(task string, docs []*pb.CrawlDoc) (error, int) {
	var errr error
	Db, errr = leveldb.OpenFile(DbConf, nil)
	if errr != nil {
		LOG.Fatalf("Open Leveldb Err:%v, use %s", errr, DbConf)
	}
	defer Db.Close()
	t1 := time_util.GetTimeInMs()
	batch := new(leveldb.Batch)
	for _, doc := range docs {
		var tm int64
		if doc.CrawlRecord != nil && doc.CrawlRecord.FetchTime > 0 {
			tm = doc.CrawlRecord.FetchTime
		} else {
			tm = time_util.GetCurrentTimeStamp()
		}
		key := mergedDbKey(task, doc.Docid, tm)
        // record db key
        if doc.CrawlRecord == nil {
            doc.CrawlRecord = &pb.CrawlRecord{}
        }
        doc.CrawlRecord.DbKey = key
		value, e := proto.Marshal(doc)
		if e != nil {
			LOG.VLog(2).DebugTag("LevelDb", "saveIntoMergeDb %s proto marshal err %v", string(key), e)
			continue
		}
		batch.Put([]byte(key), value)
		LOG.VLog(5).DebugTag("LevelDb", "insert key(%s)", string(key))
	}
	err := Db.Write(batch, nil)
	if err != nil {
		return err, -1
	}
	LOG.VLog(4).DebugTag("LevelDb", "saveIntoMergeDb %d record use %d ms", batch.Len(), time_util.GetTimeInMs()-t1)
	return nil, batch.Len()
}

// delete batch by docid. delete from mergeddb
func deleteBatch(keys [][]byte) int {
	var errr error
	Db, errr = leveldb.OpenFile(DbConf, nil)
	if errr != nil {
		LOG.Fatalf("Open Leveldb Err:%v, use %s", errr, DbConf)
	}
	defer Db.Close()

	t1 := time_util.GetTimeInMs()
	batch := new(leveldb.Batch)
	for _, k := range keys {
		batch.Delete(k)
	}
	err := Db.Write(batch, nil)
	if err != nil {
		return -1
	}
	LOG.VLog(3).DebugTag("Leveldb", "deleteBatch %d record use %d ms", batch.Len(), time_util.GetTimeInMs()-t1)
	return batch.Len()
}
