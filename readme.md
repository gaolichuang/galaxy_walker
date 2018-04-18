# spider
## init
> govender init

## 抓取部分-fetch
 - hostload
 - connectionpool【总并发连接数量】
 - proxy
 - follow302
 - 


框架
1.任务管理
sqlite : save task and docid, level... 已发现,未抓取。。。 初始化,新建table
result : leveldb,用来去重复。。。
2.web
storage提供api,抽取level,已发现等case。
想想翻页如何关联的?
层级是个更好的选择

抓完之后的解析通过travel leveldb完成。

用代码描述每一层级之间的关联,并且是可以断点继续执行的

1.开始页 = 抓取 = 解析翻页获得新的url【个数】
2.翻页 = 抓取 = 解析内容页
3.内容页 = 抓取
==每个抓取任务单独的程序,能够做到重新run的时候标记是否重新抓取就行了;还得支持去重
==支持debug模式,只抓一个;并且能检测报警。记录产生url个数等

proto
1.贯穿始终,处理单位

抓取部分[支持lib,支持实例,multi,goroutine,channel]
1.支持代理
2.follow 302
3.内容死链检测
4.封禁检测,验证码。。
5.去重复
6.抽链; 内部链接;外部链接


页面分析部分
1.html dom 抽取有效区域
2.regex 提取需要字段  http://www.cnblogs.com/golove/p/3269099.html
3.url链接抽取
4.翻页解析
5.encode utf8 dos2unix

持久化
1.遍历
2.读写
3.文件形式; 列表页;内容页;防止重复抓取
