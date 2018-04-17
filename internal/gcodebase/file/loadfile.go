package file

import (
	. "galaxy_walker/internal/gcodebase"
        LOG "galaxy_walker/internal/gcodebase/log"
        "galaxy_walker/internal/gcodebase/time_util"
	"fmt"
	"strconv"
	"strings"
	"io/ioutil"
        "sync"
)

const (
        kReloadVersionMin = 10
)

type ConfigLoader struct {
        Last_load_time    int64
        last_load_version int

        reload_interval   int64
        last_version      int

        // TODO. file md5...
        last_fingerprint  string

        version_file      string
}

func (c *ConfigLoader) ShouldReloadConfigWithContentChange(filename string) bool {
        if c.ShouldReloadConfig() == false {
                return false
        }
        fp := GetMD5FromFile(filename)
        if fp != c.last_fingerprint {
                c.last_fingerprint = fp
                return true
        }
        return false
}
func (c *ConfigLoader) ShouldReloadConfig() bool {
        // only call once...
        if time_util.GetCurrentTimeStamp() - c.Last_load_time > c.reload_interval {
                c.Last_load_time = time_util.GetCurrentTimeStamp()
                newVersion := c.GetVersionFromFile()
                if newVersion < 0 {
                        return true
                } else {
                        if c.last_version < newVersion || (c.last_version != newVersion && newVersion < kReloadVersionMin) {
                                c.last_version = newVersion
                                return true
                        }
                }
        }
        return false
}
func (c *ConfigLoader) GetVersionFromFile() int {
        currentVersion := -1
        FileLineReaderSoftly(c.version_file, "#", func(line string) {
                value, e := strconv.Atoi(line)
                if e != nil {
                        LOG.Errorf("GerVersionFail %s", line)
                } else {
                        if currentVersion == -1 {
                                currentVersion = value
                        }
                }
        })
        if currentVersion < 0 {
                LOG.Errorf("GerVersion Fail %s", c.version_file)
        }
        return currentVersion
}
func (c *ConfigLoader) SetConfigVersionFile(filename string) *ConfigLoader {
        c.version_file = GetConfFile(filename)
        CHECK(Exist(c.version_file), "ConfigLoader %s not exist", c.version_file)
        c.Last_load_time = -1
        return c
}
func (c *ConfigLoader) SetReloadInterval(interval int64) *ConfigLoader {
        c.reload_interval = interval
        CHECK(c.reload_interval > 0, "interval should bigger then 0")
        return c
}
func (c *ConfigLoader) NextLoadInfo() string {
        remain := c.reload_interval - time_util.GetCurrentTimeStamp() + c.Last_load_time
        if remain < 0 {
                remain = 0
        }
        return fmt.Sprintf("Remain:%d,Version:%d", remain, c.last_version)
}
func (c *ConfigLoader) LoadConfigWithTwoField(name, filename, splitS string) map[string]string {
        c.Last_load_time = time_util.GetCurrentTimeStamp()
        filename = GetConfFile(filename)
        result := make(map[string]string)
        LOG.Infof("Load Config %s %s", name, filename)
        FileLineReader(filename, "#", func(line string) {
                addr := strings.Split(line, splitS)
                if len(addr) != 2 {
                        LOG.Errorf("%s Load Config Format Error, %s : %s", name, filename, line)
                        return
                }
                addr0 := strings.TrimSpace(addr[0])
                result[addr0] = strings.TrimSpace(addr[1])
                LOG.VLog(6).Debugf("Load %s  %s : %s", name, addr0, addr[1])
        })
        return result
}

type DirectoryConfLoader struct {
        // key: filename, value: md5, judge reload or not.
        ConfFingerPrint       map[string]string
        // TODO. set last_reload_timestamp into each filename; ShouldReload should be deleted.
        //	ConfLoadTimeStamp map[string]uint64
        ConfPathFileNum       map[string]int

        ReloadInterval        int64
        last_reload_timestamp int64

        last_version          int
        version_file          string
}

func (d *DirectoryConfLoader) SetReloadInterval(interval int64) *DirectoryConfLoader {
        d.ConfFingerPrint = make(map[string]string)
        d.ConfPathFileNum = make(map[string]int)
        d.ReloadInterval = interval
        CHECK(d.ReloadInterval > 0, "interval should bigger then 0")
        return d
}
func (d *DirectoryConfLoader) SetConfigVersionFile(filename string) *DirectoryConfLoader {
        d.version_file = GetConfFile(filename)
        CHECK(Exist(d.version_file), "ConfigLoader %s not exist", d.version_file)
        return d
}
func (d *DirectoryConfLoader) GetVersionFromFile() int {
        currentVersion := -1
        FileLineReaderSoftly(d.version_file, "#", func(line string) {
                value, e := strconv.Atoi(line)
                if e != nil {
                        LOG.Errorf("GerVersionFail %s", line)
                } else {
                        if currentVersion == -1 {
                                currentVersion = value
                        }
                }
        })
        if currentVersion < 0 {
                LOG.Errorf("GerVersion Fail %s", d.version_file)
        }
        return currentVersion
}
func (d *DirectoryConfLoader)ShouldReload() bool {
        // only call once...
        if time_util.GetCurrentTimeStamp() - d.last_reload_timestamp > d.ReloadInterval {
                d.last_reload_timestamp = time_util.GetCurrentTimeStamp()
                newVersion := d.GetVersionFromFile()
                if d.last_version != newVersion {
                        //if d.last_version < newVersion || (d.last_version != newVersion && newVersion < kReloadVersionMin) {
                        d.last_version = newVersion
                        return true
                }
        }
        return false
}
func (d *DirectoryConfLoader) ShouldReloadConfigWithContentChange(filename string) bool {
        // only check  md5
        filename = GetConfFile(filename)
        ffp := GetMD5FromFile(filename)
        fp, e := d.ConfFingerPrint[filename]
        if e && ffp == fp {
                return false
        }
        d.ConfFingerPrint[filename] = ffp
        return true
}
func (d *DirectoryConfLoader)ShouldReloadWithDirectoryChange(filepath, filesubfix string) bool {
        filepath = GetConfFile(filepath)
        if !Exist(filepath) {
                LOG.Errorf("%s Not Exist", filepath)
                return false
        }
        files, err := ioutil.ReadDir(filepath)
        if err != nil {
                return false
        }
        newFileNum := 0
        for _, f := range files {
                if !f.IsDir() && strings.HasSuffix(f.Name(), filesubfix) {
                        newFileNum += 1
                }
        }
        orgFileNum, e := d.ConfPathFileNum[filepath]
        if !e {
                orgFileNum = 0
        }
        if orgFileNum != newFileNum {
                LOG.VLog(5).Debugf("%s (%s) ShouldReload %d => %d", filepath, filesubfix, orgFileNum, newFileNum)
                d.ConfPathFileNum[filepath] = newFileNum
                return true
        }
        return false
}

const (
        FileOPTypeInsUpd = 1
        FileOPTypeDel = 2
)

type Token struct {
        Err      error
        Type     int
        // ins
        Value    interface{}

        // del
        FileName string
}

func (d *DirectoryConfLoader)TravelConfPath(filepath, filesubfix string, force bool, funcF func(string) (error, interface{})) (error, chan *Token) {
        return d.travelConfPathHelper([]string{filepath}, filesubfix, force, 1000, funcF)
}
func (d *DirectoryConfLoader)TravelConfMultiPath(filepath []string, filesubfix string, force bool, funcF func(string) (error, interface{})) (error, chan *Token) {
        return d.travelConfPathHelper(filepath, filesubfix, force, 1000, funcF)
}
func (d *DirectoryConfLoader)travelConfPathHelper(paths []string, subfix string, force bool, chansize int, funcF func(string) (error, interface{})) (error, chan *Token) {
        files := make([]string, 0)
        for _, path := range paths {
                path = GetConfFile(path)
                if !IsDir(path) {
                        LOG.Errorf("travelConfPathHelper path %s Not dir", path)
                        continue
                }
                LOG.VLog(4).DebugTag("DirectoryConfLoader","begin to load %s",path)
                fs, err := ioutil.ReadDir(path)
                if err != nil {
                        return err, nil
                }
                for _, f := range fs {
                        if !f.IsDir() && strings.HasSuffix(f.Name(), subfix) {
                                files = append(files, fmt.Sprintf("%s/%s", path, f.Name()))
                        }
                }
        }
        LOG.VLog(3).DebugTag("DirectoryConfLoader", "load %d files subfix:%s,from %v", len(files),subfix,paths)
        LOG.VLog(4).DebugTag("DirectoryConfLoader", "%v load files subfix:%s %v",paths,subfix,files)
        t := make(chan *Token, chansize)
        go d.travelConfPath(files, subfix, force, t, funcF)
        return nil, t
}
func (d *DirectoryConfLoader)travelConfPath(files []string, subfix string, force bool, tokens chan *Token, funcF func(string) (error, interface{})) {
        defer close(tokens)
        var wait sync.WaitGroup
        travelfiles := make(map[string]bool)
        for _, filename := range files {
                travelfiles[filename] = true
                // check md5
                ffp := GetMD5FromFile(filename)
                fp, e := d.ConfFingerPrint[filename]
                if force == false && e && ffp == fp {
                        // no need to reload
                        LOG.VLog(5).Debugf("%s FP(%s) NotChange", filename, fp)
                        continue
                }
                d.ConfFingerPrint[filename] = ffp
                wait.Add(1)
                go func(filename string) {
                        e, v := funcF(filename)
                        tokens <- &Token{
                                Err:e,
                                Value:v,
                                Type:FileOPTypeInsUpd,
                                FileName:filename,
                        }
                        wait.Done()
                }(filename)
        }
        wait.Wait()

        // delete file with prefix...
        delfiles := make([]string, 0)
        for k, _ := range d.ConfFingerPrint {
                if strings.HasSuffix(k, subfix) {
                        if _, ok := travelfiles[k]; !ok {
                                delfiles = append(delfiles, k)
                        }
                }
        }
        for _, k := range delfiles {
                delete(d.ConfFingerPrint, k)
                LOG.VLog(3).DebugTag("File", "delete file %s", k)
                tokens <- &Token{
                        Type:FileOPTypeDel,
                        FileName:k,
                }
        }
}

func (d *DirectoryConfLoader)TravelConfFile(filepath, filesubfix string, force bool, funcF func(string) error) (error, int) {
        filepath = GetConfFile(filepath)
        if !Exist(filepath) {
                return fmt.Errorf("%s Not Exist", filepath), 0
        }
        files, err := ioutil.ReadDir(filepath)
        if err != nil {
                return err, 0
        }
        doneChannels := make([]chan error, 0)
        for _, f := range files {
                if !f.IsDir() && strings.HasSuffix(f.Name(), filesubfix) {
                        filename := fmt.Sprintf("%s/%s", filepath, f.Name())
                        // check md5
                        ffp := GetMD5FromFile(filename)
                        fp, e := d.ConfFingerPrint[filename]
                        if force == false && e && ffp == fp {
                                // no need to reload
                                LOG.VLog(5).Debugf("%s FP(%s) NotChange", filename, fp)
                                continue
                        }
                        d.ConfFingerPrint[filename] = ffp
                        /*
			err := funcF(filename)
                        if err != nil {
                                return err,pNum
                        }
                        */
                        done := make(chan error)
                        doneChannels = append(doneChannels, done)
                        go func(d chan error, filename string) {
                                err := funcF(filename)
                                d <- err
                        }(done, filename)
                }
        }
        // wait for all error...
        var retErr error
        for _, done := range doneChannels {
                err := <-done
                if err != nil {
                        retErr = err
                }
        }
        return retErr, len(doneChannels)
}