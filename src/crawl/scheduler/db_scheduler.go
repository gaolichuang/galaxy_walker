/*
基于db，未抓取的发送出去开抓
抓完的标记状态
bloom filter
*/
package scheduler

import (
	pb "galaxy_walker/src/proto"
	"galaxy_walker/internal/github.com/willf/bloom"
)

type DBScheduler struct {
	bf *bloom.BloomFilter
}

func (s *DBScheduler) Run() {

}
func (s *DBScheduler) Fresh(taskid string, num int) []*pb.CrawlDoc {
	/*
	   undo的
	   失败重试
	   超时
	*/
	return nil
}
