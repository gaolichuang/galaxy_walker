/*
处理所有跟抓取相关内容。依赖late pageanalysis;proto;db

fetcher: 抓取，创建连接部分

prepare： 转码

doc: 抽链 a

storage: 存储到db

response： 发送给crawlDoc.CrawlParam.Receivers
*/
package crawl
