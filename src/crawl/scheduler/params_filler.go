package scheduler

import (
    "encoding/json"
    "galaxy_walker/internal/gcodebase"
    "galaxy_walker/internal/gcodebase/conf"
    "galaxy_walker/internal/gcodebase/file"
    "galaxy_walker/internal/gcodebase/hash"
    LOG "galaxy_walker/internal/gcodebase/log"
    pb "galaxy_walker/src/proto"
    "galaxy_walker/src/utils/url_parser"
    "reflect"
    "regexp"
    "strconv"
    "galaxy_walker/src/utils"
    "galaxy_walker/src/task"
)

var CONF = conf.Conf

/*
hostload,multifetcher,fake host,receivers
priority
tag,prime second
randomhostload
drop content
store engine
store db,table
request type
*/


var NormalJobD = task.JobDescription{
    IsUrgent:        false,
    PrimeTag:        "n",
    RandomHostLoad:  0,
    DropContent:     false,
    RequestType:     1,
    Use_proxy:       false,
    Custom_ua:       true,
    Follow_redirect: false,
}

var UrgentJobD = task.JobDescription{
    IsUrgent:       true,
    PrimeTag:       "U",
    RandomHostLoad: 0,
    DropContent:    false,
    RequestType:    1,
}

type ParamFillerMaster struct {
    fillers ParamFillerGroup
    jd      *task.JobDescription
}

func (m *ParamFillerMaster) RegisterParamFillerGroup(f ParamFillerGroup) {
    m.fillers = f
}
func (m *ParamFillerMaster) RegisterJobDescription(jd *task.JobDescription) {
    m.jd = jd
}
func (m *ParamFillerMaster) Init() {
    // package first...
    m.fillers.Package()
    for _, v := range m.fillers.Fillers() {
        v.Init()
    }
}
func (m *ParamFillerMaster) Fill(doc *pb.CrawlDoc) {
    for _, v := range m.fillers.Fillers() {
        LOG.VLog(4).Debugf("Fill %s by %s", doc.RequestUrl, reflect.Indirect(reflect.ValueOf(v)).Type().Name())
        v.Fill(m.jd, doc)
    }
}

type ParamFillerGroup interface {
    Package()
    Fillers() []ParamFiller
}

type DefaultParamFillerGroup struct {
    fillers []ParamFiller
}

func (d *DefaultParamFillerGroup) Fillers() []ParamFiller {
    return d.fillers
}
func (d *DefaultParamFillerGroup) Package() {
    // pay attention the sequence
    // FakeHostParamFiller & HostLoadParamFiller & MultiFetcherParamFiller mush use and ensure the sequence
    d.fillers = append(d.fillers, &PrepareParamFiller{})
    d.fillers = append(d.fillers, &FakeHostParamFiller{})
    d.fillers = append(d.fillers, &HostLoadParamFiller{})
    d.fillers = append(d.fillers, &MultiFetcherParamFiller{})
    d.fillers = append(d.fillers, &ReceiverParamFiller{})
    d.fillers = append(d.fillers, &TagParamFiller{})
}

type ParamFiller interface {
    Init()
    Fill(*task.JobDescription, *pb.CrawlDoc)
}

// prepareParamFiller should the first one
type PrepareParamFiller struct {
}

func (p *PrepareParamFiller) Init() {
}
func (p *PrepareParamFiller) Fill(jd *task.JobDescription, doc *pb.CrawlDoc) {
    base.CHECK(doc.RequestUrl != "", "Doc Request url not filled")
    // normalize request_url, fill url,host,path ...
    if doc.GetCrawlParam() == nil {
        doc.CrawlParam = &pb.CrawlParam{}
    }
    if doc.CrawlParam.GetFetchHint() == nil {
        doc.CrawlParam.FetchHint = &pb.FetchHint{}
    }
    if doc.GetCrawlRecord() == nil {
        doc.CrawlRecord = &pb.CrawlRecord{}
    }
    // fill url
    doc.Url = url_parser.NormalizeUrl(doc.RequestUrl)
    // Use uint32 url hash for docid. key in db
    doc.Docid = hash.FingerPrint32(doc.Url)
    doc.CrawlParam.FetchHint.Host = url_parser.GetURLObj(doc.Url).Host
    doc.CrawlParam.FetchHint.Path = url_parser.GetURLObj(doc.Url).Path
}

type FakeHostParamFiller struct {
    fakehost map[string]string
    file.ConfigLoader
}

func (f *FakeHostParamFiller) loadFakeHostConfigFile() {
    fname := file.GetConfFile(*CONF.Crawler.FakeHostConfigFile)
    result := f.LoadConfigWithTwoField("FakeHost", fname, ",")
    for k, v := range result {
        f.fakehost[k] = v
        LOG.VLog(3).Debugf("Load FakeHost %s : %s", k, v)
    }
}
func (f *FakeHostParamFiller) Init() {
    f.fakehost = make(map[string]string)
    f.loadFakeHostConfigFile()
}
func (f *FakeHostParamFiller) fill(jd *task.JobDescription, doc *pb.CrawlDoc) {
    for k, v := range f.fakehost {
        r, _ := regexp.Compile(k)
        regexRet := r.FindAllString(doc.CrawlParam.FetchHint.Host, -1)
        if len(regexRet) != 0 {
            doc.CrawlParam.FakeHost = v
        }
    }
}
func (f *FakeHostParamFiller) Fill(jd *task.JobDescription, doc *pb.CrawlDoc) {
    f.loadFakeHostConfigFile()
    f.fill(jd, doc)
}

type HostLoadParamFiller struct {
    hostload map[string]int
    file.ConfigLoader
}

func (h *HostLoadParamFiller) loadHostloadConfigFile() {
    fname := file.GetConfFile(*CONF.Crawler.HostLoadConfigFile)
    result := h.LoadConfigWithTwoField("HostLoad", fname, ",")
    for k, v := range result {
        hl, err := strconv.Atoi(v)
        if err != nil {
            LOG.Errorf("Load Config Atoi Error, %s %s:%s", fname, k, v)
            continue
        }
        h.hostload[k] = hl
        LOG.VLog(3).Debugf("Load HostLoad %s : %d", k, hl)
    }
}
func (h *HostLoadParamFiller) Init() {
    h.hostload = make(map[string]int)
    h.loadHostloadConfigFile()
}
func (h *HostLoadParamFiller) fill(jd *task.JobDescription, doc *pb.CrawlDoc) {
    host := utils.GetHostName(doc)
    hl := *CONF.Crawler.DefaultHostLoad
    thl, present := h.hostload[host]
    if present {
        hl = thl
    }
    doc.CrawlParam.Hostload = int32(hl)
}
func (h *HostLoadParamFiller) Fill(jd *task.JobDescription, doc *pb.CrawlDoc) {
    h.loadHostloadConfigFile() // reload
    h.fill(jd, doc)
}

type MultiFetcherParamFiller struct {
    multifetcher map[string]int
    file.ConfigLoader
}

func (f *MultiFetcherParamFiller) loadMultiFetcherConfigFile() {
    fname := file.GetConfFile(*CONF.Crawler.MultiFetcherConfigFile)
    result := f.LoadConfigWithTwoField("HostLoad", fname, ",")
    for k, v := range result {
        hl, err := strconv.Atoi(v)
        if err != nil {
            LOG.Errorf("Load Config Atoi Error, %s %s:%s", fname, k, v)
            continue
        }
        f.multifetcher[k] = hl
        LOG.VLog(3).Debugf("Load Multifetcher %s : %d", k, hl)
    }
}
func (f *MultiFetcherParamFiller) Init() {
    f.multifetcher = make(map[string]int)
    f.loadMultiFetcherConfigFile()
}
func (f *MultiFetcherParamFiller) fill(jd *task.JobDescription, doc *pb.CrawlDoc) {
    host := utils.GetHostName(doc)
    mf := 1
    thl, present := f.multifetcher[host]
    if present {
        mf = thl
    }
    doc.CrawlParam.FetcherCount = int32(mf)
}
func (f *MultiFetcherParamFiller) Fill(jd *task.JobDescription, doc *pb.CrawlDoc) {
    f.loadMultiFetcherConfigFile()
    f.fill(jd, doc)
}

type ReceiverParamFiller struct {
    receivers map[string]*pb.ConnectionInfo
    file.ConfigLoader
}

func (f *ReceiverParamFiller) loadReceiverConfigFile() {
    fname := file.GetConfFile(*CONF.Crawler.ReceiversConfigFile)
    result := f.LoadConfigWithTwoField("HostLoad", fname, ":")
    for k, v := range result {
        hl, err := strconv.Atoi(v)
        if err != nil {
            LOG.Errorf("Load Config Atoi Error, %s %s:%s", fname, k, v)
            continue
        }
        f.receivers[k+":"+v] = &pb.ConnectionInfo{Host: k, Port: int32(hl)}
        LOG.VLog(3).Debugf("Load receivers %s : %d", k, hl)
    }
}
func (f *ReceiverParamFiller) Init() {
    f.receivers = make(map[string]*pb.ConnectionInfo)
    f.loadReceiverConfigFile()
}
func (f *ReceiverParamFiller) fill(jd *task.JobDescription, doc *pb.CrawlDoc) {
    for _, v := range f.receivers {
        doc.CrawlParam.Receivers = append(doc.CrawlParam.Receivers, v)
    }
}

func (f *ReceiverParamFiller) Fill(jd *task.JobDescription, doc *pb.CrawlDoc) {
    f.loadReceiverConfigFile()
    f.fill(jd, doc)
}

type TagParamFiller struct {
}

func (h *TagParamFiller) Init() {
}
func (h *TagParamFiller) Fill(jd *task.JobDescription, doc *pb.CrawlDoc) {
    doc.CrawlParam.Pri = pb.Priority_NORMAL
    if jd.IsUrgent {
        doc.CrawlParam.Pri = pb.Priority_URGENT
    }
    if doc.CrawlParam.PrimaryTag == "" {
        doc.CrawlParam.PrimaryTag = jd.PrimeTag
    }
    for _, v := range jd.SecondTag {
        doc.CrawlParam.SecondaryTag = append(doc.CrawlParam.SecondaryTag, v)
    }
    if doc.CrawlParam.RandomHostload == 0 {
        doc.CrawlParam.RandomHostload = int32(jd.RandomHostLoad)
    }
    doc.CrawlParam.DropContent = jd.DropContent
    if doc.CrawlParam.Rtype == 0 {
        doc.CrawlParam.Rtype = pb.RequestType(jd.RequestType)
    }
    if doc.CrawlParam.Referer != "" {
        doc.CrawlParam.Referer = jd.Referer
    }
    doc.CrawlParam.CustomUa = jd.Custom_ua
    doc.CrawlParam.FollowRedirect = jd.Follow_redirect
    doc.CrawlParam.UseProxy = jd.Use_proxy
    doc.CrawlParam.Nofollow = jd.NoFollow
}

func GetJobDescriptionFromFile(filename string) *task.JobDescription {
    c, e := file.ReadFileToString(filename)
    base.CHECKERROR(e, "read file %s", filename)
    var jd task.JobDescription
    e = json.Unmarshal([]byte(c), &jd)
    base.CHECKERROR(e, "UnMarshal Error From %s", filename)
    LOG.Infof("Load JobDescription from %s : %+v", filename, jd)
    return &jd
}
