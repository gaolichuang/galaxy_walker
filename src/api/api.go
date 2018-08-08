package api

import (
    "net/http"
    "fmt"
    "strconv"
    "io/ioutil"
    "bytes"
    "runtime"
    "reflect"
    "encoding/json"
    "strings"
    "log"
    "html"
    "galaxy_walker/internal/gcodebase/babysitter"
    LOG "galaxy_walker/internal/gcodebase/log"
    "galaxy_walker/internal/gcodebase"
    "galaxy_walker/internal/gcodebase/file"
    "galaxy_walker/internal/gcodebase/conf"
    "galaxy_walker/internal/gcodebase/http_lib"
    "galaxy_walker/internal/github.com/gorilla/mux"
    "galaxy_walker/internal/github.com/asaskevich/govalidator"
    "galaxy_walker/src/db"
    "github.com/davecgh/go-spew/spew"
)
var CONF = conf.Conf

const (

    kSyncTaskFilterAPIKeyWord = "TAPI"
    kSyncTaskFilterDetailAPIKeyWord = "TDETAILAPI"

    KFullSyncMagicTokenKey = "token"
    KFullSyncMagicToken = "Z5EEs89bgiJC7e9RNT8w"

    kEndPointPreFix = "/api"

    kEndPointRequestProxy = kEndPointPreFix + "/requestProxy"

    kEndPointSample = kEndPointPreFix + "/sample"
    kEndPointSampleWithKey = kEndPointPreFix + "/sample/{name}"

    kEndPointETCFileServerPath = kEndPointPreFix + "/etc/"

    kEndPointTrackingLog = kEndPointPreFix + "/tracking"


    kEndPointSpew = kEndPointPreFix + "/spew"
)
var TrackingLog LOG.CustomLogger
func init() {
    loggerFile := *CONF.Crawler.TrackingLogFile
    var err error
    TrackingLog, err = LOG.NewCustomLoggerWithFlag(loggerFile, log.LstdFlags)
    base.CHECK(err == nil, "Query Log Err %v", err)
}

/////////////////
type APIContext struct {
    Version     int
    LastVersion int
    Detail      bool
    TakeEffect  bool

    RequestInfo string
}

func (a *APIContext) String() string {
    if a == nil {
        return "nil"
    }
    return fmt.Sprintf("Version:%d,LastV:%d,Detail:%t,TakeEff:%t", a.Version, a.LastVersion, a.Detail, a.TakeEffect)
}
/////////////////
type APIService struct {
    taskdb db.TaskDbItf
    contentdb db.ContentDBItf
    urldb db.UrlDbItf
}

func (s *APIService) Init() {
    s.Reload()
    s.taskdb = db.NewTaskItf()
    s.contentdb = db.NewContentDBItf()
    s.urldb = db.NewUrlDbItf()
}
func (s *APIService) Reload() {

}
func (s *APIService) Status() string {
    // TODO.. APISERVICE Status.
    return ""
}
func (s *APIService) MonitorReport(result *babysitter.MonitorResult) {
    info := s.Status()
    result.AddString(info)
}
func (s *APIService) MonitorReportHealthy() error {
    // TODO. healthy or not.
    return nil
}

func CommonHandlerWrapper(f func(*APIContext, http.ResponseWriter, *http.Request) []byte) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // parse info...
        versionStr := r.URL.Query().Get("version")
        lastVersionStr := r.URL.Query().Get("lastversion")
        detail := r.URL.Query().Get("detail")
        takeEffect := r.URL.Query().Get("effect")
        version, lastversion := 0, 0
        var err error
        if versionStr != "" {
            version, err = strconv.Atoi(versionStr)
            if err != nil {
                info := getReturnMsg(103, fmt.Sprintf("Parse Version Fail %v", err), "")
                w.Header().Set("Content-Type", "application/json")
                w.Write(info)
                return
            }
        }
        if lastVersionStr != "" {
            lastversion, err = strconv.Atoi(lastVersionStr)
        }
        if lastversion > 0 {
            // TODO. CHECK Samewith current version.
        }
        // TODO. set current version into header.

        url := r.URL.String()
        method := r.Method
        //  ReadCloser can only read once.
        buf, _ := ioutil.ReadAll(r.Body)
        rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
        rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
        r.Body = rdr2

        body, _ := ioutil.ReadAll(rdr1)
        funcName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
        requestInfo := fmt.Sprintf("From:%s, %s %s @ %s", r.RemoteAddr, method, url, funcName)

        LOG.VLog(2).DebugTag("HTTPREQUEST", "%s, Body:%s", requestInfo, string(body))
        TrackingLog.LogTag(kSyncTaskFilterAPIKeyWord + "REQUEST", "HTTP Request %s", requestInfo)
        TrackingLog.LogTag(kSyncTaskFilterDetailAPIKeyWord + "REQUEST", "HTTP Request %s,Body %s", requestInfo, string(body))
        respBody := f(&APIContext{Version: version, LastVersion: lastversion, Detail: detail == "true", TakeEffect: takeEffect == "true", RequestInfo:requestInfo}, w, r)
        logv := 2
        if method == "GET" {
            logv = 3
        }
        LOG.VLog(logv).DebugTag("HTTPRESPONSE", "From:%s, %s %s @ %s, Response:%s", r.RemoteAddr, method, url, funcName, string(respBody))
        TrackingLog.LogTag(kSyncTaskFilterAPIKeyWord + "RESPONSE", "HTTP Response From:%s, %s %s @ %s", r.RemoteAddr, method, url, funcName)
        TrackingLog.LogTag(kSyncTaskFilterDetailAPIKeyWord + "RESPONSE", "HTTP Response From:%s, %s %s @ %s, Response:%s", r.RemoteAddr, method, url, funcName, string(respBody))

        if respBody != nil {
            ct := w.Header().Get("Content-Type")
            if ct == "" {
                w.Header().Set("Content-Type", "application/json")
            }
            w.Write(respBody)
        }
    })
}

func (s *APIService) Serve(host string, port int, shost string, sport int) {
    router := mux.NewRouter().StrictSlash(true)

    router.Handle(kEndPointRequestProxy, CommonHandlerWrapper(s.RequestProxy)).Methods("POST")
    router.Handle(kEndPointSample, CommonHandlerWrapper(s.SamplePost)).Methods("POST")
    //router.HandleFunc(kEndPointSample, s.SamplePost).Methods("POST")
    router.HandleFunc(kEndPointSampleWithKey, s.SampleHandler).Methods("GET")

    router.Handle(kEndPointTrackingLog, babysitter.CommonSingleFileServer(
        *CONF.Crawler.TrackingLogFile)).Methods("GET")

    router.HandleFunc(kEndPointSpew, s.SampleSpew).Methods("GET")

    s.serveDBContent(router)
    s.serveTaskAPI(router)


    // just download and read
    confPrefixPath := file.GetConfFile(*CONF.ConfPathPrefix)
    LOG.Infof("FilePath Server Dir %s", confPrefixPath)
    router.PathPrefix(kEndPointETCFileServerPath).Handler(
        http.StripPrefix(kEndPointETCFileServerPath,
            http.FileServer(http.Dir(confPrefixPath))))


    // enable webservice or not.
    if *CONF.Crawler.EnableWebService {
        rootPath := file.GetConfFile(*CONF.Crawler.WebServiceRootPath)
        router.PathPrefix("/").Handler(
            http.StripPrefix("/",
                http.FileServer(http.Dir(rootPath))))
    }

    go func(host string, port int, router http.Handler) {
        serverAddr := fmt.Sprintf("%s:%d", host, port)
        LOG.Infof("Start HttpServer at %s", serverAddr)
        err := http.ListenAndServe(serverAddr, router)
        if err != nil {
            LOG.Fatalf("Http Server Start Fail %s, %s", serverAddr, err.Error())
        }
    }(host, port, router)

    if *CONF.UseTLS {
        go func(host string, port int, router http.Handler) {
            if port <= 0 {
                LOG.Info("Port less then 0, Not Start Https")
                return
            }
            serverAddr := fmt.Sprintf("%s:%d", host, port)
            LOG.Infof("Start HttpsServer at %s", serverAddr)
            certFile := file.GetConfFile(*CONF.CertFile)
            keyFile := file.GetConfFile(*CONF.KeyFile)
            err := http.ListenAndServeTLS(serverAddr, certFile, keyFile, router)
            if err != nil {
                LOG.Fatalf("Https Server Start Fail %s, %s", serverAddr, err.Error())
            }
        }(shost, sport, router)
    }
}


func getReturnMsg(code int, reason string, reqId string) []byte {
    info, _ := json.Marshal(&http_lib.APIResponse{
        Code:      code,
        Reason:    fmt.Sprintf("%s", reason),
        RequestId: reqId,
    })
    return info
}

func APITokenValidation(body []byte) (bool, string, map[string]string) {
    params := make(map[string]string)
    if err := json.Unmarshal(body, &params); err != nil {
        return false, fmt.Sprintf("Unmarsal Err %v %s", err, string(body)), nil
    }
    if params[KFullSyncMagicTokenKey] != KFullSyncMagicToken {
        return false, fmt.Sprintf("Token Err %s", string(body)), params
    }
    return true, "", params
}

func (s *APIService) SamplePost(c *APIContext, w http.ResponseWriter, r *http.Request) []byte {
    /*
       POST   /dnsapi/sample
    */
    body, err := ioutil.ReadAll(r.Body)
    info := getReturnMsg(200, fmt.Sprintf("%s", string(body)), "")
    if err != nil {
        info = getReturnMsg(101, fmt.Sprintf("%s", err.Error()), "")
    }
    return info
}

func (s *APIService) RequestProxy(c *APIContext, w http.ResponseWriter, r *http.Request) []byte {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        return getReturnMsg(101, fmt.Sprintf("%s", err.Error()), "")
    }
    type RequestProxyJsonObj struct {
        Url     string `json:"url"`
        // only support get,post
        Method  string            `json:"method"`
        Body    string            `json:"body"`
        Headers map[string]string `json:"headers"`
    }
    rpjo := new(RequestProxyJsonObj)
    err = json.Unmarshal(body, rpjo)
    if err != nil {
        return getReturnMsg(101, fmt.Sprintf("%s", err.Error()), "")
    }
    rpjo.Method = strings.ToLower(rpjo.Method)
    if rpjo.Method != "get" && rpjo.Method != "post" {
        return getReturnMsg(101, fmt.Sprintf("Only Support Method Post,Get, Not Support %s", rpjo.Method), "")
    }
    if !govalidator.IsURL(rpjo.Url) {
        return getReturnMsg(101, fmt.Sprintf("Not Valid Url %s", rpjo.Url), "")
    }
    LOG.VLog(4).DebugTag("RequestProxy", "Url:%s,Body:%s,Method:%s,Headers:%v", rpjo.Url, rpjo.Body, rpjo.Method, rpjo.Headers)
    respbody, respErr := http_lib.GetUrlWithHeader(rpjo.Method, rpjo.Url, rpjo.Body, rpjo.Headers)
    if respErr == nil {
        return []byte(respbody)
    }
    return []byte(respErr.Error())
}

func (s *APIService) SampleSpew(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprintf(w, "Hi there, %s!", r.URL.Path[1:])
    fmt.Fprintf(w, "<!--\n" + html.EscapeString(spew.Sdump(r)) + "\n-->")

}
func (s *APIService) SampleHandler(w http.ResponseWriter, r *http.Request) {
    /*
       GET   /dnsapi/sample/{name}
    */
    vars := mux.Vars(r)
    name := vars["name"]
    w.Write([]byte(fmt.Sprintf("From:%s,name:%s", r.RemoteAddr, name)))
}