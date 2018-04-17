package http_lib

import (
        "crypto/tls"
        "errors"
        "fmt"
        "net"
        "net/http"
        "net/http/httputil"
        "net/url"
        "strings"
        "time"
        LOG "galaxy_walker/internal/gcodebase/log"
        "io"
)

const (
        BROWSER_UA = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36"
        CONNECTION_TIMEOUT time.Duration = time.Duration(3) * time.Second
        READ_WRITE_TIMEOUT time.Duration = time.Duration(10) * time.Second
        CHECK_REDIRECT_DEPTH = 5
)

var GeneralHeader = map[string]string{
        "Accept":          "text/html;q=0.8, */*;q=0.5",
        "Accept-Charset":  "utf-8, gbk, gb2312, *;q=0.5",
        "Accept-Language": "zh-cn;q=0.8, *;q=0.5",
        "Accept-Encoding": "gzip",
        "Connection":      "close",
        //    "Connection":"keep-alive",
        "User-Agent": BROWSER_UA,
}

type FetchTimeout struct {
        connect   time.Duration
        readwrite time.Duration
}

var GeneralFetchTime = &FetchTimeout{
        connect:   CONNECTION_TIMEOUT,
        readwrite: READ_WRITE_TIMEOUT,
}

func timeoutDialer(to *FetchTimeout) func(net, addr string) (c net.Conn, err error) {
        return nil
}

// http://stackoverflow.com/questions/23297520/how-can-i-make-the-go-http-client-not-follow-redirects-automatically
func noRedirect(req *http.Request, via []*http.Request) error {
        return errors.New("No Redirect")
}
func multiCheckRedirect(req *http.Request, via []*http.Request) error {
        if len(via) >= CHECK_REDIRECT_DEPTH {
                return errors.New(fmt.Sprintf("stopped after %d redirects", CHECK_REDIRECT_DEPTH))
        }
        return nil
}



////////////////HttpClientGenerator//////////////////////////////////////////////////////////
/*

        req := c.requestGenerator.
                WithCustomUA(.CustomUa).
                WithReferer(.Referer).
                NewRequest(Url)
        client := c.clientGenerator.
                WithSchema(https).
                WithRedirect(..FollowRedirect).
                WithProxy(..UseProxy).
                Fetch(req)
*/
type HttpClientGenerator struct {
        redirect  bool
        https     bool // if use https, no proxy...

        ProxyUrls []*url.URL
        ProxyUrlIdx int
}

func (hg *HttpClientGenerator) reset() {
        hg.redirect = false
        hg.https = false
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
func (hg *HttpClientGenerator) WithProxy(proxyurls string) *HttpClientGenerator {

        return hg
}
func (p *HttpClientGenerator) rrProxyUrl() (*url.URL, error) {
        if len(p.ProxyUrls) == 0 {
                return nil, errors.New("No Alive Proxy")
        }
        id := p.ProxyUrlIdx % len(p.ProxyUrls)
        p.ProxyUrlIdx = (p.ProxyUrlIdx + 1) % len(p.ProxyUrls)
        return p.ProxyUrls[id],nil
}

func (hg *HttpClientGenerator) NewClient() *http.Client {
        LOG.VLog(4).Debugf("NewClient:https:%t,proxy:%v,redirect:%t", hg.https, hg.ProxyUrls, hg.redirect)
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
        proxyUrl,_ := hg.rrProxyUrl()
        if proxyUrl != nil && hg.https == false {
                // only http request use proxy
                clientProxy = http.ProxyURL(proxyUrl)
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
func (hg *HttpClientGenerator)Fetch(request *http.Request) (error, *http.Response) {
        // step 2. fetch
        client := hg.NewClient()
        resp, err := client.Do(request)

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
                LOG.VLog(4).Debugf("Dump Response(Error:%s):\n%s", respMsg, respErr)
        }
        if err != nil && strings.Contains(err.Error(), "use of closed network connection") {
        }
        return err,resp
}

/////////////HttpRequestGenerator/////////////////////////////////////////////////////////////
// HttpRequestGenerator TODO. Add cookie and basic auth support...
type HttpRequestGenerator struct {
        customUA bool
        referer  string
        method string
        body io.Reader
}

func (rg *HttpRequestGenerator) WithMethod(method string) *HttpRequestGenerator {
        rg.method = method
        return rg
}
func (rg *HttpRequestGenerator) WithBody(content string) *HttpRequestGenerator {
        rg.body = strings.NewReader(content)
        return rg
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
        req, _ := http.NewRequest(rg.method, _url, rg.body)

        for k, v := range GeneralHeader {
                req.Header.Set(k, v)
        }
        if rg.referer != "" {
                req.Header.Set("Referer", rg.referer)
        }
        if rg.customUA {
                req.Header.Set("User-Agent", BROWSER_UA)
        }
        dumpReq, _ := httputil.DumpRequest(req, true)
        LOG.VLog(4).Debugf("DumpRequest:\n%s", string(dumpReq))
        return req
}

//////////////////////////////////////////////////////////////////////////
