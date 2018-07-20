package application

import (
        "os"
        "path/filepath"
        "fmt"
        "strconv"
        "strings"
        "time"
        "galaxy_walker/internal/gcodebase/file"
        "sync"
        LOG "galaxy_walker/internal/gcodebase/log"
)

type Application struct {
        PIDPath string
        CMD string
        pidFile string
}
func (a *Application)SetPIDPath(path string) *Application {
        a.PIDPath = path
        a.pidFile = filepath.Join(path,fmt.Sprintf("%s.pid",a.CMD))
        return a
}

func (a *Application)GetLastPID() (int,error) {
        content,err := file.ReadFileToString(a.pidFile)
        if err != nil {
                return -1,err
        }
        lpid,err := strconv.Atoi(string(content))
        if err != nil {
                return -1,err
        }
        if os.Getpid() == lpid {
                return lpid,fmt.Errorf("%d is self",lpid)
        }
        LOG.VLog(2).DebugTag("Application","Get Pid:%d from %s",lpid,a.pidFile)
        return lpid,nil
}
func (a *Application)StopByPID(pid int) error {
        if os.Getpid() == pid {
                return fmt.Errorf("PID:%d is self",pid)
        }
        process,err := os.FindProcess(pid)
        if err != nil {
                LOG.VLog(2).DebugTag("Application","%d not found",pid)
                return err
        }
        for i:=0;i<10;i++ {
                if err := process.Kill(); err != nil {
                        if strings.Contains(err.Error(),"os: process already finished") {
                                LOG.VLog(2).DebugTag("Application","%d already finished",pid)
                                return nil
                        }
                }
                LOG.VLog(2).DebugTag("Application","%d wait for fninsh...",pid)
                time.Sleep(time.Millisecond*500)
        }
        return err
}
func (a *Application)SaveCurrentPID() (int,error) {
        pid := strconv.Itoa(os.Getpid())
        LOG.VLog(2).DebugTag("Application","Save pid %s to %s",pid,a.pidFile)
        return os.Getpid(),file.WriteStringToFile([]byte(pid),a.pidFile)
}


var _app_instance *Application = nil
var _app_init_ctx sync.Once

func GetApplication() *Application {
        _app_init_ctx.Do(func() {
                cmd := filepath.Base(os.Args[0])
                _app_instance = &Application{
                        CMD:cmd,
                        PIDPath:"/tmp",
                        pidFile:fmt.Sprintf("/tmp/%s.pid",cmd),
                }
        })
        return _app_instance
}

