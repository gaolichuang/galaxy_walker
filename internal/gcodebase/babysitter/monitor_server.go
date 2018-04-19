package babysitter

import (
    "net/http"
    LOG "galaxy_walker/internal/gcodebase/log"
    "galaxy_walker/internal/github.com/gorilla/mux"
    "fmt"
    "galaxy_walker/internal/gcodebase/string_util"
    "encoding/json"
    "strings"
    _ "net/http/pprof"
    "net"
    "io/ioutil"
    "galaxy_walker/internal/github.com/asaskevich/govalidator"
    "galaxy_walker/internal/gcodebase/file"
    "galaxy_walker/internal/gcodebase/conf"
    . "galaxy_walker/internal/gcodebase"
    "strconv"
)

var CONF = conf.Conf

const (
    CallBackPath           = "/statusi/callback"
    StatusUiPath           = "/statusi"
    StatusUiAPi            = "/statusi/api"
    StatusUiAPIPath        = "/statusi/api/machine"
    StatusUiAPIHealthyPath = "/statusi/api/healthy"
    PprofDebugPath         = "/debug/pprof/"

    kMaxLimitNum = 1000
    kMaxSize     = 1024000
)

type MonitorInterface interface {
    MonitorReport(*MonitorResult)
    MonitorReportHealthy() error
}
type MonitorResult struct {
    info             string
    kv               map[string]string
    StatisticMachine string
}

func (mr *MonitorResult) AddString(s string) {
    mr.info = s
}
func (mr *MonitorResult) AddKv(k, v string) {
    mr.kv[k] = v
}

type MonitorServer struct {
    result       *MonitorResult
    monitor      MonitorInterface
    callbackFunc func(*CallBackParams, *CallBackResponse) bool

    CallBackFunctionKit *CallBackFunctionKit
    handleFunc          map[string]http.HandlerFunc
    pathPrefix          map[string]string
}

func (m *MonitorServer) Init(version string) {
    bbVersion = version
    m.result = &MonitorResult{kv: make(map[string]string),
        StatisticMachine: machineInfoHtml(),
    }
    m.handleFunc = make(map[string]http.HandlerFunc)
    m.pathPrefix = make(map[string]string)
    m.CallBackFunctionKit = NewCallBackFunctionKit()
    m.CallBackFunctionKit.Init()
}
func (m *MonitorServer) AddMonitor(mi MonitorInterface) {
    m.monitor = mi
}
func (m *MonitorServer) RegisterCallBack(callback func(*CallBackParams, *CallBackResponse) bool) {
    m.callbackFunc = callback
}

// TODO: custom handlerfunc test.... for dispatcher and fetcher add custom api func.
func (m *MonitorServer) AddHandleFunc(path string, f http.HandlerFunc) {
    CHECK(strings.HasPrefix(path, "/"), "HandleFunc path should start with /")
    path = StatusUiAPi + path
    m.handleFunc[path] = f
}

func (m *MonitorServer) AddPathPrefix(urlpath, dir string) {
    CHECK(strings.HasPrefix(urlpath, "/"), "HandleFunc path should start with /")
    urlpath = StatusUiAPi + urlpath
    m.pathPrefix[urlpath] = dir
}
func (m *MonitorServer) StatusiApi(w http.ResponseWriter, r *http.Request) {
    machine := make(map[string]string)
    for k, v := range machineInfo() {
        machine[k] = v
    }
    for k, v := range statusInfo() {
        machine[k] = v
    }
    info, _ := json.Marshal(machine)
    w.Header().Set("Content-Type", "application/json")
    w.Write(info)
}
func (m *MonitorServer) StatusiHealthyApi(w http.ResponseWriter, r *http.Request) {
    type StatusHealthy struct {
        Healthy bool   `json:"Healthy"`
        Reason  string `json:"Reason"`
    }
    // json.Marshal only encode Uppercase field in struct...
    info, _ := json.Marshal(StatusHealthy{Healthy: true})
    err := m.monitor.MonitorReportHealthy()
    if err != nil {
        info, _ = json.Marshal(StatusHealthy{Healthy: false, Reason: err.Error()})
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(info)
}

func formateCallBackResponse(code int, msg string) []byte {
    cbr := &CallBackResponse{
        Code: code,
        Msg:  msg,
    }
    content, _ := json.Marshal(cbr)
    return content
}
func parseCallbackParams(body []byte) (*CallBackParams, error) {
    params := make(map[string]string)
    err := json.Unmarshal(body, &params)
    if err != nil {
        return nil, err
    }
    if _, e := params[KCallBackParamsTrackingId]; !e {
        return nil, fmt.Errorf("Params not contain %s", KCallBackParamsTrackingId)
    }

    if _, e := params[KCallBackParamsFuncName]; !e {
        return nil, fmt.Errorf("Params not contain %s", KCallBackParamsFuncName)
    }

    if _, e := params[KCallBackParamsReportAddr]; e {
        if !govalidator.IsURL(params[KCallBackParamsReportAddr]) {
            return nil, fmt.Errorf("%s ParseFail %v", KCallBackParamsReportAddr, err)
        }
    }
    cbp := new(CallBackParams)
    cbp.TrackingId = params[KCallBackParamsTrackingId]
    cbp.FuncName = params[KCallBackParamsFuncName]

    if _, e := params[KCallBackParamsReportAddr]; e {
        cbp.ReportAddr = params[KCallBackParamsReportAddr]
    }
    cbp.ParamsMap = make(map[string]string)
    for k, v := range params {
        if k == KCallBackParamsTrackingId || k == KCallBackParamsFuncName || k == KCallBackParamsReportAddr {
            continue
        }
        cbp.ParamsMap[k] = v
    }
    return cbp, nil

}

/*
POST
{
"requestId":"3dw2dkei",
"funcName":"ReportStatus",
"":"",
"":"",
......
"":""
}
/////////////////////////
{
"requestId":"3dw2dkei",
"versionNum":12345,
"funcName":"ASyncReportStatus",
"reportAddr":"http://127.0.0.1:9950/sample"
}
 */
func (m *MonitorServer) CallBackDispatch(w http.ResponseWriter, r *http.Request) {
    var info []byte
    body, _ := ioutil.ReadAll(r.Body)
    // parse params.
    var err error
    cbp, err := parseCallbackParams(body)
    if err != nil {
        info = formateCallBackResponse(101, fmt.Sprintf("CallBack ParseParamsFail Err%v, %s", err, string(body)))
        w.Header().Set("Content-Type", "application/json")
        w.Write(info)
        return
    }
    LOG.VLog(3).DebugTag("MonitorServer", "Parse Params %v", cbp.String())
    var ok bool
    resp := &CallBackResponse{
        TrackingId: cbp.TrackingId,
        FuncName:   cbp.FuncName,
        Code:       200,
    }
    if m.callbackFunc != nil {
        ok = m.callbackFunc(cbp, resp)
        if ok {
            info, _ = json.Marshal(resp)
            goto WriteMsg
        }
    }
    ok = m.CallBackFunctionKit.Dispatch(cbp, resp)
    if ok {
        info, _ = json.Marshal(resp)
        goto WriteMsg
    } else {
        info = formateCallBackResponse(102, fmt.Sprintf("MissFuncName Req:%s,Func:%s", cbp.TrackingId, cbp.FuncName))
        goto WriteMsg
    }
WriteMsg:
    w.Header().Set("Content-Type", "application/json")
    w.Write(info)
}
func (m *MonitorServer) Statusi(w http.ResponseWriter, r *http.Request) {
    ip, _, _ := net.SplitHostPort(r.RemoteAddr)
    LOG.VLog(3).Debugf("Statusi Request from %s", ip)
    m.monitor.MonitorReport(m.result)
    w.Header().Set("Server", "Golang Statusi Server")
    w.WriteHeader(200)
    var infos string
    infos += "<h1>Machine Info</h1>"
    infos += m.result.StatisticMachine
    infos += "<h1>Status Info</h1>"
    infos += statusInfoHtml()
    infos += "<h1>Healthy</h1>"
    healthy := fmt.Sprintf("<div>%t</div>", true)
    err := m.monitor.MonitorReportHealthy()
    if err != nil {
        healthy = fmt.Sprintf("%t(%s)", false, err.Error())
    }
    string_util.StringAppendF(&infos, "%s", healthy)
    infos += "<h1>Application Info</h1>"
    for k, v := range m.result.kv {
        string_util.StringAppendF(&infos, "<key>%s : <value>%s<br>", k, v)
    }
    string_util.StringAppendF(&infos, "<br>%s", m.result.info)
    w.Write([]byte(infos))
}
func (m *MonitorServer) Serve(httpListen string, httpsListen string) {
    r := mux.NewRouter().StrictSlash(true)
    r.HandleFunc(StatusUiPath, m.Statusi).Methods("GET")
    r.HandleFunc(StatusUiAPIPath, m.StatusiApi).Methods("GET")
    r.HandleFunc(StatusUiAPIHealthyPath, m.StatusiHealthyApi).Methods("GET")
    r.HandleFunc(CallBackPath, m.CallBackDispatch).Methods("POST")
    LOG.Infof("MonitorServer Serve path %s", StatusUiPath)
    LOG.Infof("MonitorServer Serve path %s", StatusUiAPIPath)
    LOG.Infof("MonitorServer Serve path %s", StatusUiAPIHealthyPath)
    for k, v := range m.handleFunc {
        r.HandleFunc(k, v)
        LOG.Infof("MonitorServer Serve path %s", k)
    }

    for k, v := range m.pathPrefix {
        r.PathPrefix(k).Handler(http.StripPrefix(k, http.FileServer(http.Dir(v))))
        LOG.Infof("MonitorServer Serve dir %s : %s", k, v)
    }

    h, p, e := net.SplitHostPort(httpListen)
    if e != nil {
        LOG.Fatal("%s Err %v", httpListen, e)
    }
    port, _ := strconv.Atoi(p)
    serverAddr := fmt.Sprintf("%s:%d", h, port)
    LOG.Infof("Starting Http Monitor at %s:%d", h, port)

    // AttachProfiler http://stackoverflow.com/questions/19591065/profiling-go-web-application-built-with-gorillas-mux-with-net-http-pprof
    r.PathPrefix(PprofDebugPath).Handler(http.DefaultServeMux)

    go func(serverAddr string, r http.Handler) {
        err := http.ListenAndServe(serverAddr, r)
        if err != nil {
            LOG.Fatalf("Http Server Start Fail, %s", serverAddr)
        }
    }(serverAddr, r)

    pi := 0
    _, hsp, e := net.SplitHostPort(httpsListen)
    if e != nil {
        pi, _ = strconv.Atoi(hsp)
    }
    if *CONF.UseTLS && pi > 0 {
        sh, sp, e := net.SplitHostPort(httpsListen)
        if e != nil {
            LOG.Fatal("Start Https Monitor %s Err %v", httpsListen, e)
        }
        sport, _ := strconv.Atoi(sp)
        httpsServerAddr := fmt.Sprintf("%s:%d", sh, sport)

        go func(serverAddr string, router http.Handler) {
            certFile := file.GetConfFile(*CONF.CertFile)
            keyFile := file.GetConfFile(*CONF.KeyFile)
            LOG.Infof("Starting Https Monitor at %s use %s", serverAddr, certFile)
            err := http.ListenAndServeTLS(serverAddr, certFile, keyFile, router)
            if err != nil {
                LOG.Fatalf("Https Server Start Fail %s, %v", serverAddr, err)
            }
        }(httpsServerAddr, r)
    }
}
