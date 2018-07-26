package main

import (
    "database/sql"
    _ "galaxy_walker/internal/github.com/mattn/go-sqlite3"
    LOG "galaxy_walker/internal/gcodebase/log"
    "galaxy_walker/src/db/sqlite"
    "fmt"
)
const (
    kDriverName = "sqlite3"
)


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
func main() {
    dbfile := "/Users/gaolichuang/workspace/go/src/galaxy_walker/logs/1.q"
    ts := sqlite.ListTable(dbfile,"")
    fmt.Println(ts)
}
