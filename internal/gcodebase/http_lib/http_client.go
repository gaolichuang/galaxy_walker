package http_lib

import (
	"net/http"
	"compress/gzip"
	"io"
        "time"
        "net"
	"net/http/httputil"
	"io/ioutil"
	"strings"
	"errors"
	"fmt"
        LOG "galaxy_walker/internal/gcodebase/log"
        "galaxy_walker/internal/gcodebase/file"
        "os"
        "encoding/json"
        "crypto/tls"
        urllib "net/url"
)
const (
        kHttpReadWriteTimeOut = 60
        kHttpConnectionTimeOut = 10
        kDefaultRetryNum = 3
        kDefaultRetryIntervalInMiSec = 100
        KHTTPDownLoadSubFix = ".dl"
)
type APIResponse struct {
        RequestId string `json:"requestId,omitempty"`
        Code   int
        Reason string      `json:"reason,omitempty"`
}
// POST, / BODY
func GetUrl(method, url, body string) (responseBody string, responseErr error){
        return GetUrlWithHeader(method,url,body,nil)
}
func GetUrlWithTimeOutAndRetry(method, url, body string,connTimeOut,rwTimeOut int,retryNum int) (string,error){
        return GetUrlWithHeaderAndTimeOutAndRetry(method,url,body,nil,connTimeOut,rwTimeOut,retryNum,kDefaultRetryIntervalInMiSec)
}
func GetUrlWithTimeOut(method, url, body string,connTimeOut,rwTimeOut int) (string,error){
        return GetUrlWithHeaderAndTimeOutAndRetry(method,url,body,nil,connTimeOut,rwTimeOut,kDefaultRetryNum,kDefaultRetryIntervalInMiSec)
}
func GetUrlWithHeader(method, url, body string, header map[string]string) (string, error){
        return GetUrlWithHeaderAndTimeOutAndRetry(method,url,body,header,0,0,kDefaultRetryNum,kDefaultRetryIntervalInMiSec)
}
func GetUrlWithHeaderAndTimeOutAndRetry(method, url, body string, header map[string]string,connTimeOut,rwTimeOut int,retryNum int,retryInterval int) (string, error){
        var resp string
        var err error
        for i := 0;i < retryNum;i++ {
                resp,err = getUrlWithHeaderAndTimeOutInternal(method,url,body,header,connTimeOut,rwTimeOut)
                if err == nil {
                        return resp,err
                }
                LOG.VLog(3).DebugTag("GetUrl","Fail %d : %v; %s %s %s",i,err,method,url,body)
                time.Sleep(time.Millisecond * time.Duration(retryInterval))
        }
        return resp,err
}
func getUrlWithHeaderAndTimeOutInternal(method, url, body string, header map[string]string,connTimeOut,rwTimeOut int) (responseBody string, responseErr error){
        if connTimeOut <= 0 {
                connTimeOut = kHttpConnectionTimeOut
        }
        if rwTimeOut <= 0 {
                rwTimeOut = kHttpReadWriteTimeOut
        }
        urlParserObj,err := urllib.Parse(url)
        if err != nil {
                return "",err
        }

        method = strings.ToUpper(method)
        client := http.Client{
                Transport: &http.Transport{
                        Dial: func(netw, addr string) (net.Conn, error) {
                                c, err := net.DialTimeout(netw, addr, time.Second * time.Duration(connTimeOut)) //设置建立连接超时
                                if err != nil {
                                        return nil, err
                                }
                                c.SetDeadline(time.Now().Add(time.Duration(rwTimeOut) * time.Second)) //设置发送接收数据超时
                                return c, nil
                        },
                        TLSClientConfig: func(schema string) *tls.Config {
                                if schema == "https" {
                                        return &tls.Config{InsecureSkipVerify: true}
                                }
                                return nil
                        }(urlParserObj.Scheme),
                },
        }
	request, _ := http.NewRequest(method, url, nil)
	if method == "POST" {
		request, _ = http.NewRequest("POST", url, strings.NewReader(body))

	}
	request.Header.Set("Accept-Encoding", "gzip,deflate")
	request.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")
	request.Header.Set("Connection", "keep-alive")
        if header != nil {
                for k,v := range header {
                        if v != "" {
                                request.Header.Set(k,v)
                        }
                }
        }
        if request.Header.Get("Content-Type") == "" {
                request.Header.Add("Content-Type","application/json")
        }

	dumpReq,dumpErr := httputil.DumpRequest(request, body != "")
	LOG.VLog(6).Debugf("Dump Request(Err:%v): url: %s \n%s", dumpErr, url, string(dumpReq))
	response, err := client.Do(request)
	if response != nil {
		dumpResp, dumpErr := httputil.DumpResponse(response, true)
		respMsg, respErr := "Nil", "Nil"
		if dumpResp != nil {
			respMsg = string(dumpResp)
		}
		if dumpErr != nil {
                        // TODO......
                        responseErr = dumpErr
			respErr = dumpErr.Error()
		}
                LOG.VLog(6).Debugf("Dump Response(Err:%s):\n%s", respErr, respMsg)
	}
	if err != nil {
                LOG.VLog(5).Debugf("%s %s Return:%s", method,url,err.Error())
		responseErr = err
	} else {
		if response.StatusCode >= 200 && response.StatusCode < 300 {
			// Check that the server actually sent compressed data
			var reader io.ReadCloser
			switch response.Header.Get("Content-Encoding") {
			case "gzip":
				reader, _ = gzip.NewReader(response.Body)
				defer reader.Close()
			default:
				reader = response.Body
			}
			if b, err := ioutil.ReadAll(reader); err == nil {
                                apiResponse := new(APIResponse)
                                err := json.Unmarshal(b,&apiResponse)
                                if err == nil {
                                        if apiResponse.Code != 200 && apiResponse.Code > 0 && apiResponse.Reason != ""{
                                                responseErr = fmt.Errorf("%s",apiResponse.Reason)
                                        }
                                }
				responseBody = string(b)
			}
                        LOG.VLog(6).Debugf(responseBody)
		} else {
                        LOG.VLog(6).Debugf("%s %s Return:%d", method,url,response.StatusCode)
			responseErr = errors.New(fmt.Sprintf("Code:%d", response.StatusCode))
		}
	}
	return
}

func DownLoadToFile(url string, localfile string) error {
        content,err := GetUrl("GET",url,"")
        if err != nil {
                return err
        }
        localtmpfile := fmt.Sprintf("%s%s",localfile,KHTTPDownLoadSubFix)
        err = file.WriteStringToFile([]byte(content),localtmpfile)
        if err != nil {
                return err
        }
        return os.Rename(localtmpfile,strings.TrimSuffix(localtmpfile,KHTTPDownLoadSubFix))
}
func DownLoadToFileWithTimeOut(url string, localfile string,ctimeout,rwtimeout int) error {
        content,err := GetUrlWithTimeOut("GET",url,"",ctimeout,rwtimeout)
        if err != nil {
                return err
        }
        localtmpfile := fmt.Sprintf("%s%s",localfile,KHTTPDownLoadSubFix)
        err = file.WriteStringToFile([]byte(content),localtmpfile)
        if err != nil {
                return err
        }
        return os.Rename(localtmpfile,strings.TrimSuffix(localtmpfile,KHTTPDownLoadSubFix))
}