package babysitter

import (
    "fmt"
    LOG "galaxy_walker/internal/gcodebase/log"
    "strings"
    "galaxy_walker/internal/gcodebase/http_lib"
    "encoding/json"
    "galaxy_walker/internal/gcodebase/file"
    "path/filepath"
    "galaxy_walker/internal/gcodebase/time_util"
    "strconv"
    "sort"
    "galaxy_walker/internal/github.com/asaskevich/govalidator"
)

const (
    KCallBackParamsTrackingIdLength = 8
    KCallBackParamsTrackingId       = "trackingId"
    KCallBackParamsFuncName         = "funcName"
    KCallBackParamsReportAddr       = "reportAddr"
    kCallBackChanBufSize            = 50
    KIncSyncDirRemoteMD5FileKey     = "rmd5file"
    KIncSyncDirLocalMD5FileKey      = "lmd5file"
    KIncSyncDirLocalPathKey         = "localpath"
    KIncSyncDirRemotePathKey        = "remotepath"
    KIncSyncDirForceDelKey          = "forcedel"
    KSyncFileLocalFileKey           = "localfile"
    KSyncWriteFileContent           = "wcontent"
    KSyncFileRemoteFileKey          = "remotefile"

    kSyncHTTPConnectionTimeOut = 60
    KSyncHTTPRWTimeout         = 180
)

type CallBackParams struct {
    TrackingId string `json:"trackingId"`
    FuncName   string `json:"funcName"`
    // http addr for Report./ POST
    ReportAddr string `json:"reportAddr,omitempty"`

    ParamsMap map[string]string
}

func (c *CallBackParams) String() string {
    return fmt.Sprintf("<%s>-%s Addr:%s, Map:%v", c.TrackingId, c.FuncName, c.ReportAddr, c.ParamsMap)
}

type CallBackResponse struct {
    Code       int               `json:"code"`
    TrackingId string            `json:"trackingId"`
    FuncName   string            `json:"funcName"`
    Msg        string            `json:"msg"`
    OptValue   map[string]string `json:"opt,omitempty"`
}

func (c *CallBackResponse) String() string {
    if c == nil {
        return "nil"
    }
    return fmt.Sprintf("<%s>-%s Code:%d,%s, %v", c.TrackingId, c.FuncName, c.Code, c.Msg, c.OptValue)
}

type CallBackChanObj struct {
    Params *CallBackParams
    Resp   *CallBackResponse
}
type CallBackFunctionKit struct {
    FuncChan chan *CallBackChanObj
    FuncMap  map[string]func(*CallBackParams, *CallBackResponse)
}

func (c *CallBackFunctionKit) Init() {
    c.FuncMap = map[string]func(*CallBackParams, *CallBackResponse){
        "ReportStatus":          c.reportStatus,
        "IncSyncDirectory":      c.incrementSyncDirectory,
        "SyncFile":              c.syncFile,
        "ReadFile":              c.readFile,
        "WriteFile":             c.writeFile,
        "ASyncReportStatus":     c.dispatchASync,
        "ASyncIncSyncDirectory": c.dispatchASync,
        "ASyncSyncFile":         c.dispatchASync,
    }
    go c.run()
}
func (c *CallBackFunctionKit) run() {
    for {
        select {
        case p := <-c.FuncChan:
            // TODO. MERGE Command...
            LOG.VLog(2).DebugTag("CallBack", "Async Call %s", p.Params.String())
            c.FuncMap[p.Params.FuncName](p.Params, p.Resp)
            if p.Params.ReportAddr != "" {
                // CALL BACK to report Addr.
                if p.Params.ReportAddr != "" {
                    if govalidator.IsURL(p.Params.ReportAddr) {
                        content, _ := json.Marshal(p.Resp)
                        body, err := http_lib.GetUrlWithTimeOut("POST", p.Params.ReportAddr, string(content), kSyncHTTPConnectionTimeOut, KSyncHTTPRWTimeout)
                        if err != nil {
                            LOG.Errorf("CallBack Call Err %v %s", err, body)
                        }
                        LOG.VLog(2).DebugTag("CallBack", "%s %v %s", p.Params.String(), err, body)
                    } else {
                        LOG.Errorf("%s Parse ReportAddr Err", p.Params.String())
                    }
                }
            } else {
                LOG.VLog(2).DebugTag("CallBack", "Not Set CallBack %s", p.Params.String(), p.Resp.String())
            }
        }
    }
}

func (c *CallBackFunctionKit) Dispatch(cbp *CallBackParams, resp *CallBackResponse) bool {
    _, e := c.FuncMap[cbp.FuncName]
    if e {
        LOG.VLog(2).DebugTag("CallBack", "Sync Call %s", cbp.String())
        c.FuncMap[cbp.FuncName](cbp, resp)
        LOG.VLog(2).DebugTag("CallBack", "Sync Call Response %s", resp.String())
        return true
    }
    return false
}
func (c *CallBackFunctionKit) dispatchASync(cbp *CallBackParams, resp *CallBackResponse) {
    if len(c.FuncChan) >= int(0.8*float32(kCallBackChanBufSize)) {
        resp.Code = 101
        resp.Msg = fmt.Sprintf("Across Chan Len %d", len(c.FuncChan))
        return
    }
    cbp.FuncName = strings.TrimLeft(cbp.FuncName, "ASync")
    c.FuncChan <- &CallBackChanObj{
        Params: cbp,
        Resp:   resp,
    }
}
func (c *CallBackFunctionKit) reportStatus(cbp *CallBackParams, resp *CallBackResponse) {
    resp.Msg = machineInfoStr()
}

/*
POST
{
	"requestId":"3dw2dkei",
	"funcName":"IncSyncDirectory",
	"reportAddr":"http://127.0.0.1:9950/dnsapi/callback",
	"md5file":"http://127.0.0.1:9950/dnsapi/fileServer/etc.md5",
	"localpath":"/Users/zhujunchao/workspace/go/src/glory_dns/logs/etc",
	"remotepath":"http://127.0.0.1:9950/dnsapi/fileServer/etc"
}
*/
func (c *CallBackFunctionKit) incrementSyncDirectory(cbp *CallBackParams, resp *CallBackResponse) {
    /*
    md5file
            http://127.0.0.1/fileServer/etc.md5
                    {
                    "fc2dec6f1d5cbcc0da656c8f03c27c39	config.ini",
                    "53b7ca9ac9b66f7671e5a602805fb6f4	ipset/ipset.db",
                    "e30fa20568df3a5196323644cf503706	ipset/view_topology.json",
                    "9762c71f848dd0cb218f03278fb55dbb	schedule/dc.schedule",
                    "ee178120017c6bced13028367106b939	schedule/ip.schedule",
                    "ad4ad2f5f291c8519f67add55c5b35c3	smart/sh.json",
                    "c30f7472766d25af1dc80b3ffc9a58c7	version_file",
                    "3c05bc81816e7abd915f50e6fb7f1a24	view_topology/sh.tom.com.topology.json",
                    "3dd8a5d16e99521f872fd151911d12d4	view_topology.json",
                    "6124332f7ba85b199403766300d47333	zone/tom.com.hosts"
                    }
    localpath
            /usr/local/dndns/etc
    remotepath
            http://127.0.0.1/fileServer/etc
    forceDel
            true
    */
    /*
    return:
    */
    // parse and valid params
    t1 := time_util.GetTimeInMs()
    err := func(cbp *CallBackParams) error {
        if cbp.ParamsMap == nil {
            return fmt.Errorf("Params Map not Set")
        }
        _, e := cbp.ParamsMap[KIncSyncDirRemoteMD5FileKey]
        if !e {
            return fmt.Errorf("Not Set %s", KIncSyncDirRemoteMD5FileKey)
        }
        if !govalidator.IsURL(cbp.ParamsMap[KIncSyncDirRemoteMD5FileKey]) {
            return fmt.Errorf("%s not valid url %s", KIncSyncDirRemoteMD5FileKey, cbp.ParamsMap[KIncSyncDirRemoteMD5FileKey])
        }

        _, e = cbp.ParamsMap[KIncSyncDirRemotePathKey]
        if !e {
            return fmt.Errorf("Not Set %s", KIncSyncDirRemotePathKey)
        }
        if !govalidator.IsURL(cbp.ParamsMap[KIncSyncDirRemotePathKey]) {
            return fmt.Errorf("%s not valid url %s", KIncSyncDirRemotePathKey, cbp.ParamsMap[KIncSyncDirRemotePathKey])
        }

        _, e = cbp.ParamsMap[KIncSyncDirLocalPathKey]
        if !e {
            return fmt.Errorf("Not Set %s", KIncSyncDirLocalPathKey)
        }
        if !file.IsDir(cbp.ParamsMap[KIncSyncDirLocalPathKey]) {
            return fmt.Errorf("%s Set %s not exist dir", KIncSyncDirLocalPathKey, cbp.ParamsMap[KIncSyncDirLocalPathKey])
        }
        return nil
    }(cbp)
    if err != nil {
        resp.Code = 101
        resp.Msg = fmt.Sprintf("%v", err)
        return
    }
    localpath := cbp.ParamsMap[KIncSyncDirLocalPathKey]
    remotepath := cbp.ParamsMap[KIncSyncDirRemotePathKey]
    remoteMd5file := cbp.ParamsMap[KIncSyncDirRemoteMD5FileKey]
    localMd5file := cbp.ParamsMap[KIncSyncDirLocalMD5FileKey]
    remoteMd5Content, err := http_lib.GetUrlWithTimeOut("GET", remoteMd5file, "", kSyncHTTPConnectionTimeOut, KSyncHTTPRWTimeout)
    if err != nil {
        resp.Code = 103
        resp.Msg = fmt.Sprintf("Get Md5 Err %v from %s", err, remoteMd5file)
        return
    }
    rMd5Obj := make([]string, 0)
    err = json.Unmarshal([]byte(remoteMd5Content), &rMd5Obj)
    if err != nil {
        resp.Code = 104
        resp.Msg = fmt.Sprintf("Parse Md5Content Err %v from %s", err, remoteMd5file)
        return
    }
    rMd5Map := make(map[string]string)
    for _, v := range rMd5Obj {
        fileds := strings.Split(v, "\t")
        if len(fileds) != 2 {
            continue
        }
        rMd5Map[fileds[1]] = fileds[0]
    }
    if len(rMd5Map) == 0 {
        resp.Code = 105
        resp.Msg = fmt.Sprintf("Nothing Sync from %s", remoteMd5file)
        return
    }
    errs, ops := syncDirectory(localpath, localMd5file, rMd5Map, remotepath)
    if len(errs) > 0 {
        resp.Code = 103
        errMsg, _ := json.Marshal(&errs)
        resp.Msg = string(errMsg)
    } else {
        resp.Msg = fmt.Sprintf("Success Sync Use %d ms", time_util.GetTimeInMs()-t1)
        resp.OptValue = make(map[string]string)
        msg, _ := json.Marshal(&ops)
        resp.OptValue["operation"] = string(msg)
        resp.OptValue["costInMs"] = strconv.Itoa(int(time_util.GetTimeInMs() - t1))
    }
    _, e := cbp.ParamsMap[KIncSyncDirForceDelKey]
    if e && cbp.ParamsMap[KIncSyncDirForceDelKey] == "true" {
        err, num := file.DeletePath(localpath, file.KSoftDeleteSubFix)
        if err != nil {
            LOG.Error("Delete %s %s Err %v", localpath, file.KSoftDeleteSubFix, err)
        }
        LOG.VLog(2).DebugTag("SyncDir", "Delete File %s %s %d", localpath, file.KSoftDeleteSubFix, num)
        err, num = file.DeletePath(localpath, http_lib.KHTTPDownLoadSubFix)
        if err != nil {
            LOG.Error("Delete %s %s Err %v", localpath, http_lib.KHTTPDownLoadSubFix, err)
        }
        LOG.VLog(2).DebugTag("SyncDir", "Delete File %s %s %d", localpath, http_lib.KHTTPDownLoadSubFix, num)
    }
}
func syncDirectory(localpath string, localmd5file string, remotefileMd5 map[string]string, remotepath string) ([]string, []string) {
    // input localpath, remote md5 file and remotepath
    // return errors, and operations
    // fileMd5 key: filename value: md5
    allErr := make([]string, 0)
    operations := make([]string, 0)
    multiChan := make([]chan error, 0)
    lfs := file.CalDirectoryMd5Sum(localpath, "", file.KSoftDeleteSubFix)

    for k, v := range lfs {
        _, e := remotefileMd5[k]
        if !e {
            // soft delete
            localf := filepath.Join(localpath, k)
            LOG.VLog(2).DebugTag("SyncDir", "DELETE file %s", localf)
            file.SoftDeleteFile(localf)
            operations = append(operations, fmt.Sprintf("DEL\t%s", k))
        } else {
            // update
            if v != remotefileMd5[k] {
                remotef := strings.Join([]string{remotepath, k}, "/")
                localf := filepath.Join(localpath, k)
                ec := make(chan error)
                multiChan = append(multiChan, ec)
                go func(remotef, localf string, ec chan error) {
                    LOG.VLog(2).DebugTag("SyncDir", "UPDATE file From %s to %s", remotef, localf)
                    err := http_lib.DownLoadToFileWithTimeOut(remotef, localf, kSyncHTTPConnectionTimeOut, KSyncHTTPRWTimeout)
                    ec <- err
                }(remotef, localf, ec)
                operations = append(operations, fmt.Sprintf("UPD\t%s", k))
            }
        }
    }
    for k, _ := range remotefileMd5 {
        _, e := lfs[k]
        if !e {
            kdir := filepath.Dir(k)
            kdir = filepath.Join(localpath, kdir)
            if !file.Exist(kdir) {
                err := file.MkDirAll(kdir)
                if err != nil {
                    allErr = append(allErr, err.Error())
                    continue
                }
            }
            // new file
            remotef := strings.Join([]string{remotepath, k}, "/")
            localf := filepath.Join(localpath, k)
            LOG.VLog(2).DebugTag("SyncDir", "ADD file From %s to %s", remotef, localf)
            ec := make(chan error)
            multiChan = append(multiChan, ec)
            go func(remotef, localf string, ec chan error) {
                err := http_lib.DownLoadToFileWithTimeOut(remotef, localf, kSyncHTTPConnectionTimeOut, KSyncHTTPRWTimeout)
                ec <- err
            }(remotef, localf, ec)
            operations = append(operations, fmt.Sprintf("ADD\t%s", k))
        }
    }
    for _, c := range multiChan {
        err := <-c
        if err != nil {
            allErr = append(allErr, err.Error())
        }
    }
    // save lfs to localpath.md5
    newlfs := file.CalDirectoryMd5Sum(localpath, "", file.KSoftDeleteSubFix)
    retj := make([]string, 0)
    for k, v := range newlfs {
        retj = append(retj, fmt.Sprintf("%s\t%s", v, k))
    }
    sort.Strings(retj)
    content, _ := json.Marshal(&retj)
    // TODO. modify local md5 file...
    err := file.WriteStringToFile(content, localmd5file)
    if err != nil {
        allErr = append(allErr, err.Error())
    }
    return allErr, operations
}
func (c *CallBackFunctionKit) readFile(cbp *CallBackParams, resp *CallBackResponse) {
    /*
    localfile
    */
    err := func(cbp *CallBackParams) error {
        if cbp.ParamsMap == nil {
            return fmt.Errorf("Params Map not Set")
        }
        _, e := cbp.ParamsMap[KSyncFileLocalFileKey]
        if !e {
            return fmt.Errorf("Not Set %s", KSyncFileLocalFileKey)
        }
        return nil
    }(cbp)
    if err != nil {
        resp.Code = 101
        resp.Msg = fmt.Sprintf("%v", err)
        return
    }
    localfile := cbp.ParamsMap[KSyncFileLocalFileKey]
    if !file.IsRegular(localfile) {
        resp.Code = 102
        resp.Msg = fmt.Sprintf("%s not regular file.", localfile)
        return
    }

    content, err := file.ReadFileToString(localfile)
    if err != nil || len(content) > 10*1024 {
        resp.Code = 103
        resp.Msg = fmt.Sprintf("read %s Err:%v or too big", localfile, err)
        return
    }
    resp.Msg = fmt.Sprintf("%s", string(content))
}

func (c *CallBackFunctionKit) writeFile(cbp *CallBackParams, resp *CallBackResponse) {
    /*
    localfile
    */
    err := func(cbp *CallBackParams) error {
        if cbp.ParamsMap == nil {
            return fmt.Errorf("Params Map not Set")
        }

        if _, e := cbp.ParamsMap[KSyncFileLocalFileKey]; !e {
            return fmt.Errorf("Not Set %s", KSyncFileLocalFileKey)
        }
        if _, e := cbp.ParamsMap[KSyncWriteFileContent]; !e {
            return fmt.Errorf("Not Set %s", KSyncWriteFileContent)
        }
        return nil
    }(cbp)
    if err != nil {
        resp.Code = 101
        resp.Msg = fmt.Sprintf("%v", err)
        return
    }
    localfile := cbp.ParamsMap[KSyncFileLocalFileKey]
    content := cbp.ParamsMap[KSyncWriteFileContent]

    t1 := time_util.GetTimeInMs()
    err = file.WriteStringToFile([]byte(content), localfile)
    if err != nil {
        resp.Code = 103
        resp.Msg = fmt.Sprintf("Write %s to %s Err:%v", content, localfile, err)
        return
    }
    resp.Msg = fmt.Sprintf("Success Write %s to %s use %d ms", content, localfile, time_util.GetTimeInMs()-t1)
}
func (c *CallBackFunctionKit) syncFile(cbp *CallBackParams, resp *CallBackResponse) {
    /*
    remotefile:
    localfile:
    */
    t1 := time_util.GetTimeInMs()
    err := func(cbp *CallBackParams) error {
        if cbp.ParamsMap == nil {
            return fmt.Errorf("Params Map not Set")
        }
        _, e := cbp.ParamsMap[KSyncFileLocalFileKey]
        if !e {
            return fmt.Errorf("Not Set %s", KSyncFileLocalFileKey)
        }
        if !govalidator.IsURL(cbp.ParamsMap[KSyncFileRemoteFileKey]) {
            return fmt.Errorf("%s not valid url %s", KSyncFileRemoteFileKey, cbp.ParamsMap[KSyncFileRemoteFileKey])
        }
        return nil
    }(cbp)
    if err != nil {
        resp.Code = 101
        resp.Msg = fmt.Sprintf("%v", err)
        return
    }
    err = http_lib.DownLoadToFileWithTimeOut(cbp.ParamsMap[KSyncFileRemoteFileKey], cbp.ParamsMap[KSyncFileLocalFileKey], kSyncHTTPConnectionTimeOut, KSyncHTTPRWTimeout)
    if err != nil {
        resp.Code = 102
        resp.Msg = fmt.Sprintf("%v", err)
        return
    }
    resp.Msg = fmt.Sprintf("Success %s use %d ms", cbp.ParamsMap[KSyncFileLocalFileKey], time_util.GetTimeInMs()-t1)
}

func NewCallBackFunctionKit() *CallBackFunctionKit {
    return &CallBackFunctionKit{
        FuncChan: make(chan *CallBackChanObj, kCallBackChanBufSize),
    }
}
