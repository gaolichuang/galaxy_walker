/*
Page Analysis and extract link.
*/
package crawl

import (
    "fmt"
    LOG "galaxy_walker/internal/gcodebase/log"
    "galaxy_walker/src/proto"
    "galaxy_walker/internal/github.com/PuerkitoBio/goquery"
    "galaxy_walker/internal/gcodebase/string_util"
    "net/url"
    "reflect"
    "strings"
    "galaxy_walker/src/utils/url_parser"
    "galaxy_walker/src/utils/page_analysis"
    "galaxy_walker/src/utils"
)

type DocHandler struct {
    CrawlHandler
    htmlParser *page_analysis.HtmlParser
    doc        *proto.CrawlDoc
    domain     string
    // Distinct: key is OutLink.url + OutLink.text
    indomainDict map[string]*proto.OutLink
    outdomainDict map[string]*proto.OutLink
}

func (handler *DocHandler) Init() bool {
    handler.htmlParser = page_analysis.NewHtmlParser()
    handler.indomainDict = make(map[string]*proto.OutLink)
    handler.outdomainDict = make(map[string]*proto.OutLink)
    handler.htmlParser.RegisterSelector("a", handler.extractLinkCallBack)
    return true
}
func (handler *DocHandler) extractLinkCallBack(i int, s *goquery.Selection) {
    href, hrefexit := s.Attr("href")
    if !hrefexit {
        return
    }
    if !(strings.HasPrefix(href, "/") || strings.HasPrefix(href, "http") || strings.HasPrefix(href, ".")) {
        LOG.VLog(5).Debugf("Not Avaliable link %s", href)
        return
    }
    nofollow, nofollowexit := s.Attr("rel")
    if (!handler.doc.CrawlParam.Nofollow) && nofollowexit && nofollow == "nofollow" {
        LOG.VLog(4).Debugf("NoFollow link doc.crawlparam.nofollow:%t, link:%s,text:%s",
            handler.doc.CrawlParam.Nofollow,
            href,
            s.Text())
        return
    }
    text := string_util.Purify(s.Text(), "\n", "\t", " ")
    href = string_util.Purify(href, "\n", "\t", "")
    if strings.HasPrefix(href, "/") {
        requrl, _ := url.Parse(handler.doc.RequestUrl)
        fixHref := fmt.Sprintf("%s://%s%s", requrl.Scheme, requrl.Host, href)
        LOG.VLog(6).Debugf("InDomainLinkFill %s,text:%s,fixhref %s", href, text,fixHref)
        key := fixHref + "^_^" + text
        if _,ok := handler.indomainDict[key];ok {
            handler.indomainDict[key].Num += 1
        } else {
            handler.indomainDict[key] = &proto.OutLink{
                Url:  fixHref,
                Text: text,
                Num:1,
            }
        }
    } else {
        newdomain := utils.GetDomainFromHost(url_parser.GetHost(href))
        key := href + "^_^" + text
        if utils.IsSameDomain(newdomain, handler.domain) {
            if _,ok := handler.indomainDict[key];ok {
                LOG.VLog(5).Debugf("InDomainLink Add %s,text:%s", href, text)
                handler.indomainDict[key].Num += 1
            } else {
                handler.indomainDict[key] = &proto.OutLink{
                    Url:  href,
                    Text: text,
                }
            }
        } else {
            if _,ok := handler.outdomainDict[key];ok {
                LOG.VLog(5).Debugf("OutDomainLink Add %s,text:%s", href, text)
                handler.outdomainDict[key].Num += 1
            } else {
                handler.outdomainDict[key] = &proto.OutLink{
                    Url:  href,
                    Text: text,
                }
            }
        }
    }
}
func (handler *DocHandler) Accept(crawlDoc *proto.CrawlDoc) bool {
    return utils.IsCrawlSuccess(crawlDoc.Code) && !crawlDoc.CrawlParam.NoExtractLink
}
func (handler *DocHandler) Process(crawlDoc *proto.CrawlDoc) {
    LOG.VLog(3).Debugf("[%s]Process One Doc %s ",
        reflect.Indirect(reflect.ValueOf(handler)).Type().Name(),
        crawlDoc.Url)
    handler.doc = crawlDoc
    handler.domain = utils.GetDomainFromHost(url_parser.GetHost(crawlDoc.Url))
    handler.htmlParser.Parse(handler.doc.Url, handler.doc.Content)
    // merge indomain,outdomain
    for _,v := range handler.indomainDict {
        handler.doc.IndomainOutlinks = append(handler.doc.IndomainOutlinks,v)
    }
    for _,v := range handler.outdomainDict {
        handler.doc.OutdomainOutlinks = append(handler.doc.OutdomainOutlinks,v)
    }
    handler.indomainDict = nil
    handler.outdomainDict = nil
    LOG.VLog(4).Debugf("DocHandler. DumpCrawlDoc\n%s", utils.DumpCrawlDoc(crawlDoc))
}

// use for create instance from a string
func init() {
    registerCrawlTaskType(&DocHandler{})
}
