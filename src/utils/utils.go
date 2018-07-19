package utils

import (
    pb "galaxy_walker/src/proto"
    "galaxy_walker/internal/gcodebase/string_util"
    "strings"
)

const (
    kMaxValidUrlLength = 512
)

func GetHostName(doc *pb.CrawlDoc) string {
    if string_util.IsEmpty(doc.CrawlParam.FakeHost) {
        return doc.CrawlParam.FetchHint.Host
    }
    return doc.CrawlParam.FakeHost
}
func DumpCrawlDoc(doc *pb.CrawlDoc) string {
    docContent := doc.Content
    doc.Content = "..."
    dumpString := pb.FromProtoToString(doc)
    doc.Content = docContent
    return dumpString
}

// TODO call this function in where???
func IsInvalidUrl(_url string) bool {
    /*
        1. start with http or https
        2. url len should less then kMaxValidUrlLength
    */
    if !(strings.HasPrefix(_url, "http://") || strings.HasPrefix(_url, "https://")) {
        return false
    }
    if len(_url) > kMaxValidUrlLength {
        return false
    }
    return true
}
func IsCrawlSuccess(t pb.ReturnType) bool {
    return t == pb.ReturnType_STATUS200 || t == pb.ReturnType_STATUS201
}
func IsPermanentRedirect(t pb.ReturnType) bool {
    return t == pb.ReturnType_STATUS301
}
func IsTemporaryRedirect(t pb.ReturnType) bool {
    return t == pb.ReturnType_STATUS300 ||
        t == pb.ReturnType_STATUS302 ||
        t == pb.ReturnType_STATUS305 ||
        t == pb.ReturnType_STATUS307
}

func GetDomainFromHost(host string) string {
    hostSplit := strings.Split(host, ".")
    if len(hostSplit) <= 2 {
        return host
    }
    return strings.Join(hostSplit[1:], ".")
}
func IsSameDomain(domain1, domain2 string) bool {
    d1, d2 := strings.Split(domain1, "."), strings.Split(domain2, ".")
    if len(d1) == len(d2) {
        return domain1 == domain2
    }
    if len(d1) <= 1 || len(d2) <= 1 {
        return false
    }
    minLen := len(d1)
    if len(d2) < minLen {
        minLen = len(d2)
    }
    for i := 1; i <= minLen; i++ {
        if d1[len(d1)-i] != d2[len(d2)-i] {
            return false
        }
    }
    return true
}
