package fetcher

import (
	LOG "galaxy_walker/internal/gcodebase/log"
	"galaxy_walker/internal/gcodebase/time_util"
	"galaxy_walker/src/proto"
	"galaxy_walker/src/utils"
    "sync"
)

const (
	kConnectionPoolRecoverInterval = 3600
	kConnectionPoolTimeOut = 3600
	kConnectionPoolNewDozenSize = 10
)

type ConnectionPool struct {
	record                 map[string]int64
	hold                   map[string]bool // only one host could get in connection pool
	free                   []*Connection
	busy                   map[*Connection]bool // make it could delete
	output_chan            chan<- *proto.CrawlDoc
	last_recover_timestamp int64
	sync.Mutex
}

func (c *ConnectionPool) SetOutChan(output_chan chan<- *proto.CrawlDoc) {
	c.output_chan = output_chan
}

func (c *ConnectionPool) GetCrawlHostMap() map[string]int64 {
	return c.record
}
func (c *ConnectionPool) RecordNum() int {
	return len(c.record)
}
func (c *ConnectionPool) FreeConnectionNum() int {
	return len(c.free)
}
func (c *ConnectionPool) BusyConnectionNum() int {
	return len(c.busy)
}
func (c *ConnectionPool) releaseRecordAndHold() {
	if time_util.GetCurrentTimeStamp() - c.last_recover_timestamp < kConnectionPoolRecoverInterval {
		return
	}
	release := make([]string, 0)
	for k, v := range c.record {
		if time_util.GetCurrentTimeStamp() - v > kConnectionPoolTimeOut {
			release = append(release, k)
		}
	}
	for _, k := range release {
		delete(c.record, k)
		delete(c.hold, k)
		LOG.VLog(3).Debugf("Release Connection Pool Size: %d", len(release))
	}
}

// return false: connection all busy, can not fetch
// single thread
func (c *ConnectionPool) Fetch(doc *proto.CrawlDoc) bool {
    c.Lock()
    defer c.Unlock()
	// check hold or not
	host := utils.GetHostName(doc)
	if c.hold[host] == true {
		return false
	}
	if len(c.free) == 0 {
		if len(c.free) + len(c.busy) < *CONF.Crawler.FetchConnectionNum {
			// new dozen conns
			for i := 0; i < kConnectionPoolNewDozenSize; i++ {
				conn := NewConnection()
				c.free = append(c.free, conn)
			}
		} else {
			LOG.VLog(2).Debugf("Connection Pool full %s/%s", len(c.free), len(c.busy))
			return false
		}
	}
	conn := c.free[0]
	c.free = c.free[1:]
	c.busy[conn] = true
	c.hold[host] = true

	// use goroutine to fetch.
	t1 := time_util.GetTimeInMs()
	doc.CrawlRecord.FetchTime = time_util.GetCurrentTimeStamp()
	go conn.FetchOne(doc,func(doc *proto.CrawlDoc, conn *Connection) {
	    c.Lock()
	    defer c.Unlock()
		c.free = append(c.free, conn)
		delete(c.busy, conn)
		c.record[utils.GetHostName(doc)] = time_util.GetCurrentTimeStamp()
		doc.CrawlRecord.FetchUseInms = time_util.GetTimeInMs() - t1
		c.hold[utils.GetHostName(doc)] = false
		c.output_chan <- doc
	})
	c.releaseRecordAndHold()
	return true
}

// constructor
func NewConnectionPool(output_chan chan<- *proto.CrawlDoc) *ConnectionPool {
	return &ConnectionPool{
		record:      make(map[string]int64),
		hold:        make(map[string]bool),
		busy:        make(map[*Connection]bool),
		output_chan: output_chan,
	}
}
