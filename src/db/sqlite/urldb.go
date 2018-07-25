package sqlite

/*
each task use one table.


统计需求
1.已发现未抓取    status=0
2.抓取失败 / 成功  status=1,2
3.抓取失败不再重试的  retry>X
4.N次失败的统计   status + retry
*/
const (
    kCreateUrlTableVersion = `
CREATE TABLE IF NOT EXISTS url_%s (
    url VARCHAR(255),
    parentType int,
    parentDocid int,
    status int,
    createTimeStamp int,
    updateTimeStamp int,
    retryNum int,
);`
)
