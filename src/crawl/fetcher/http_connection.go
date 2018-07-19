package fetcher

import (
    "compress/gzip"
    "crypto/tls"
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "net"
    "net/http"
    "net/http/httputil"
    "net/url"
    "strings"
    "time"
    "galaxy_walker/src/utils"
    LOG "galaxy_walker/internal/gcodebase/log"
    pb "galaxy_walker/src/proto"
)

const (
    kBrowserUserAgent                 = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36"
    kConnectionTimeOut  time.Duration = time.Second *3
    kReadWriteTimeOut   time.Duration = time.Second * 5
    kCheckRedirectDepth               = 5
)

var GeneralHeader = map[string]string{
    "Accept":          "text/html;q=0.8, */*;q=0.5",
    "Accept-Charset":  "utf-8, gbk, gb2312, *;q=0.5",
    "Accept-Language": "zh-cn;q=0.8, *;q=0.5",
    "Accept-Encoding": "gzip",
    "Connection":      "close",
    //    "Connection":"keep-alive",
    "User-Agent": "XXSpider",
}

type FetchTimeout struct {
    connect   time.Duration
    readwrite time.Duration
}

var GeneralFetchTime = &FetchTimeout{
    connect:   kConnectionTimeOut,
    readwrite: kReadWriteTimeOut,
}

func timeoutDialer(to *FetchTimeout) func(net, addr string) (c net.Conn, err error) {
    return nil
}

// http://stackoverflow.com/questions/23297520/how-can-i-make-the-go-http-client-not-follow-redirects-automatically
func noRedirect(req *http.Request, via []*http.Request) error {
    return errors.New("No Redirect")
}
func multiCheckRedirect(req *http.Request, via []*http.Request) error {
    if len(via) >= kCheckRedirectDepth {
        return errors.New(fmt.Sprintf("stopped after %d redirects", kCheckRedirectDepth))
    }
    return nil
}

type Connection struct {
    clientGenerator  *HttpClientGenerator
    // eache connection  has one proxy manager.
    httpProxy *ProxyManager
    requestGenerator *HttpRequestGenerator
}

//  302 redirect no url or no header. -- BADHEADER

// run in goroutine
func (c *Connection) FetchOne(doc *pb.CrawlDoc, f func(*pb.CrawlDoc, *Connection)) {
    // TODO fetch doc and fill field
    // step 1. fill request info
    client := c.clientGenerator.
        WithSchema(doc.RequestUrl).
        WithRedirect(doc.CrawlParam.FollowRedirect).
        WithProxy(func(useproxy bool) *url.URL {
            if useproxy {
                u,e := c.httpProxy.GetProxyUrl()
                if e != nil {
                    LOG.Errorf("GetProxyUrl err:%v",e)
                    return nil
                }
                return u
            }
            return nil
        }(doc.CrawlParam.UseProxy)).
        GetClient()
    req := c.requestGenerator.
        WithCustomUA(doc.CrawlParam.CustomUa).
        WithReferer(doc.CrawlParam.Referer).
        NewRequest(doc.Url)
    // step 2. fetch
    resp, err := client.Do(req)
    // step 3. judge response code and fill the nessary field of crawldoc
    if resp != nil {
        dumpResp, dumpErr := httputil.DumpResponse(resp, false)
        respMsg, respErr := "Nil", "Nil"
        if dumpResp != nil {
            respMsg = string(dumpResp)
        }
        if dumpErr != nil {
            respErr = dumpErr.Error()
        }
        LOG.VLog(5).Debugf("Dump Response(Error:%s):\n%s", respErr,respMsg)
    }
    if err != nil && strings.Contains(err.Error(), "use of closed network connection") {
        c.httpProxy.MarkDeadProxy(c.clientGenerator.GetProxyUrl())
        doc.Code = pb.ReturnType_NOCONNECTION
        c.HandleOther(resp, err, doc)
    } else if err != nil && strings.Contains(err.Error(), "i/o timeout") {
        c.httpProxy.MarkDeadProxy(c.clientGenerator.GetProxyUrl())
        // read tcp 172.24.47.104:54386->220.181.112.244:443: i/o timeout
        // dial tcp: i/o timeout
        doc.Code = pb.ReturnType_TIMEOUT
        c.HandleOther(resp, err, doc)
    } else if err != nil && strings.Contains(err.Error(), "No Redirect") {
        // redirect error throw.
        doc.Code = pb.ReturnType(resp.StatusCode)
        c.Handle30X(resp, doc)
    } else if err == nil {
        doc.Code = pb.ReturnType(resp.StatusCode)
        if utils.IsCrawlSuccess(pb.ReturnType(resp.StatusCode)) {
            c.Handle200(resp, doc)
        } else {
            c.HandleOther(resp, nil, doc)
        }
    } else {
        c.httpProxy.MarkDeadProxy(c.clientGenerator.GetProxyUrl())
        // other?
        c.HandleOther(resp, err, doc)
    }
    // the last step: call the callback function.
    f(doc, c)
}

func (c *Connection) Handle200(resp *http.Response, doc *pb.CrawlDoc) {
    var reader io.ReadCloser
    switch resp.Header.Get("Content-Encoding") {
    case "gzip":
        reader, _ = gzip.NewReader(resp.Body)
        defer reader.Close()
    default:
        reader = resp.Body
    }
    if b, err := ioutil.ReadAll(reader); err == nil {
        doc.Content = string(b)
    }

    dumResp, _ := httputil.DumpResponse(resp, false)
    doc.Header = string(dumResp)
    doc.LastModify = resp.Header.Get("last-modified")
    doc.ContentType = resp.Header.Get("Content-Type")
    LOG.VLog(3).Debugf("Fetch Success, url:%s,reqtype:%d", doc.Url, doc.CrawlParam.Rtype)
}

func (c *Connection) Handle30X(resp *http.Response, doc *pb.CrawlDoc) {
    doc.LastModify = resp.Header.Get("last-modified")
    doc.ContentType = resp.Header.Get("Content-Type")
    redirectUrl := resp.Header.Get("Location")
    if !utils.IsInvalidUrl(redirectUrl) {
        doc.Code = pb.ReturnType_INVALIDREDIRECTURL
    } else {
        doc.RedirectUrl = redirectUrl
    }
    LOG.VLog(3).Debugf("Fetch 30X, url:%s, redirecturl:%s, reqtype:%d", doc.Url, doc.RedirectUrl, doc.CrawlParam.Rtype)
}
func (c *Connection) HandleOther(resp *http.Response, err error, doc *pb.CrawlDoc) {
    if err != nil {
        doc.ErrorInfo = err.Error()
    }
    LOG.VLog(3).Debugf("Fetch Code:%d, url:%s, reqtype:%d", doc.Code, doc.Url, doc.CrawlParam.Rtype)
}
func NewConnection() *Connection {
    return &Connection{
        clientGenerator: &HttpClientGenerator{
            redirect:  false,
            https:     false,
        },
        httpProxy: NewProxyManager(PROXY_SELECT_RR),
        requestGenerator: &HttpRequestGenerator{
            customUA: false,
            referer:  "",
        },
    }
}

////////////////HttpClientGenerator//////////////////////////////////////////////////////////
type HttpClientGenerator struct {
    // client with proxy
    clients map[string]*http.Client
    // clients with proxy and redirect.
    clientsWithRedirect map[string]*http.Client

    // client with httpdns
    httpdnsClient *http.Client
    // client with httpdns and redirect
    httpdnsClientWithRedirect *http.Client
    // http client
    client *http.Client
    // http client with redirect
    clientWithRedirect *http.Client

    redirect  bool
    https     bool // if use https, no proxy...
    proxyUrl  *url.URL
}

func (hg *HttpClientGenerator) reset() {
    hg.redirect = false
    hg.https = false
    hg.proxyUrl = nil
}
func (hg *HttpClientGenerator) WithSchema(_url string) *HttpClientGenerator {
    if strings.HasPrefix(_url, "https") {
        hg.https = true
    } else {
        hg.https = false
    }
    return hg
}
func (hg *HttpClientGenerator) WithRedirect(y bool) *HttpClientGenerator {
    hg.redirect = y
    return hg
}
func (hg *HttpClientGenerator) WithProxy(proxyUrl *url.URL) *HttpClientGenerator {
    hg.proxyUrl = proxyUrl
    return hg
}

func (hg *HttpClientGenerator) GetClient() *http.Client {
    if hg.https {
        if hg.redirect {
            if hg.httpdnsClientWithRedirect == nil {
                hg.httpdnsClientWithRedirect = hg.NewClient()
            }
            return hg.httpdnsClientWithRedirect
        } else {
            if hg.httpdnsClient == nil {
                hg.httpdnsClient = hg.NewClient()
            }
            return hg.httpdnsClient
        }
    } else {
        if hg.redirect {
            if hg.proxyUrl != nil {
                uStr := hg.proxyUrl.String()
                if _,ok := hg.clientsWithRedirect[uStr];!ok {
                    hg.clientsWithRedirect[uStr]=hg.NewClient()
                }
                return hg.clientsWithRedirect[uStr]
            } else {
                if hg.clientWithRedirect == nil {
                    hg.clientWithRedirect = hg.NewClient()
                }
                return hg.clientWithRedirect
            }
        } else {
            if hg.proxyUrl != nil {
                uStr := hg.proxyUrl.String()
                if _,ok := hg.clients[uStr];!ok {
                    hg.clients[uStr]=hg.NewClient()
                }
                return hg.clients[uStr]
            } else {
                if hg.client == nil {
                    hg.client = hg.NewClient()
                }
                return hg.client
            }
        }
    }
}
func (hg *HttpClientGenerator) NewClient() *http.Client {
    // TODO. add cache for new http client.
    LOG.VLog(4).Debugf("NewClient:https:%t,proxy:%v,redirect:%t", hg.https, hg.proxyUrl, hg.redirect)
    var client *http.Client
    ckRedirect := noRedirect
    if hg.redirect == true {
        ckRedirect = multiCheckRedirect
    }
    var tlsClientConfig *tls.Config = nil
    if hg.https {
        tlsClientConfig = &tls.Config{InsecureSkipVerify: true}
    }
    var clientProxy func(*http.Request) (*url.URL, error) = nil
    if hg.proxyUrl != nil && hg.https == false { // only http request use proxy
        clientProxy = http.ProxyURL(hg.proxyUrl)
    }
    client = &http.Client{
        CheckRedirect: ckRedirect,
        Transport: &http.Transport{
            Dial:            timeoutDialer(GeneralFetchTime),
            Proxy:           clientProxy,
            TLSClientConfig: tlsClientConfig,
        },
    }
    hg.reset()
    return client
}
func (hg *HttpClientGenerator) GetProxyUrl() *url.URL{
    return hg.proxyUrl
}

/////////////HttpRequestGenerator/////////////////////////////////////////////////////////////
// HttpRequestGenerator TODO. Add cookie and basic auth support...
type HttpRequestGenerator struct {
    customUA bool
    referer  string
}

func (rg *HttpRequestGenerator) WithCustomUA(y bool) *HttpRequestGenerator {
    rg.customUA = y
    return rg
}
func (rg *HttpRequestGenerator) WithReferer(referer string) *HttpRequestGenerator {
    rg.referer = referer
    return rg
}
func (rg *HttpRequestGenerator) NewRequest(_url string) *http.Request {
    // TODO. support POST method.... FetchHint.post_data
    req, _ := http.NewRequest("GET", _url, nil)

    for k, v := range GeneralHeader {
        req.Header.Set(k, v)
    }
    if rg.referer != "" {
        req.Header.Set("Referer", rg.referer)
    }
    if rg.customUA {
        req.Header.Set("User-Agent", kBrowserUserAgent)
    }
    dumpReq, _ := httputil.DumpRequest(req, true)
    LOG.VLog(4).Debugf("DumpRequest:\n%s", string(dumpReq))
    return req
}

//////////////////////////////////////////////////////////////////////////
