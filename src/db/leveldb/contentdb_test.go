package leveldb

import (
    "testing"
)

func TestContentDbKey(t *testing.T) {
    key1 := contentDbKey("",123,1532070173)
    err,task,docid,timestamp := parseContentDbKey(key1)
    if err != nil || task != "" || docid != 123 || timestamp != 1532070173 {
        t.Error("parse content db error")
    }
    t.Logf("key1:%s",key1)
    key2 := contentDbKey("xxx",3709056136,1532070173)
    err,task,docid,timestamp = parseContentDbKey(key2)
    if err != nil || task != "xxx" || docid != 3709056136 || timestamp != 1532070173 {
        t.Error("parse content db error")
    }
    t.Logf("key2:%s",key2)
}
func TestMergedDbKey(t *testing.T) {
    key1 := mergedDbKey("",123,1532070173)
    err,task,docid,timestamp := parseMergedDbKey(key1)
    if err != nil || task != "" || docid != 123 || timestamp != 1532070173 {
        t.Errorf("parse mergedb db error,%s,%s,%d,%d",key1,task,docid,timestamp)
    }
    t.Logf("key1:%s",key1)
    key2 := mergedDbKey("xxx",3709056136,1532070173)
    err,task,docid,timestamp = parseMergedDbKey(key2)
    if err != nil || task != "xxx" || docid != 3709056136 || timestamp != 1532070173 {
        t.Error("parse mergedb db error")
    }
    t.Logf("key2:%s",key2)
}