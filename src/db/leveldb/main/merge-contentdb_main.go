package main

import "galaxy_walker/src/db/leveldb"

func main() {
    leveldb.InitDb()
    leveldb.MergeContentDbProcess(false)
}
