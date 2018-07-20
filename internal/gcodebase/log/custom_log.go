package log

import (
        "log"
        "os"
        "io"
        "fmt"
        "galaxy_walker/internal/gcodebase/time_util"
        "strings"
        "sync"
        "time"
)

type CustomLogger interface {
        Log(...interface{})
        Logf(string, ...interface{})
        LogStr(string)
        LogTag(string, string, ...interface{})
        LogFile() string
}
type CustomLoggerSingle struct {
        logfile string
        Logger  *log.Logger
}

func (c *CustomLoggerSingle)LogFile() string {
        return c.logfile
}
func (c *CustomLoggerSingle)Log(v ...interface{}) {
        c.Logger.Output(2, fmt.Sprintln(v...))
}
func (c *CustomLoggerSingle)LogStr(str string) {
        c.Logger.Output(2, str)
}
func (c *CustomLoggerSingle)Logf(format string, v ...interface{}) {
        c.Logger.Output(2, fmt.Sprintf(format, v...))
}
func (c *CustomLoggerSingle)LogTag(TAG, format string, v ...interface{}) {
        c.Logger.Output(2, "["+TAG + "]"+fmt.Sprintf(format,v...))
        //c.Logger.Output(2, fmt.Sprintf("[%s]%s", TAG, fmt.Sprintf(format, v...)))
}

func NewCustomLogger(logFile string) (*CustomLoggerSingle, error) {
        // TODO add date to query log.
        return NewCustomLoggerWithFlag(logFile, log.LstdFlags)
}
func NewCustomLoggerWithFlag(logFile string, flag int) (*CustomLoggerSingle, error) {
        var writer io.Writer
        if logFile != "" {
                f, err := os.OpenFile(logFile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
                if err != nil {
                        log.Printf("error opening file %s: %v\n", logFile, err)
                        return nil, err
                }
                writer = f
        } else {
                writer = os.Stdout
        }
        return &CustomLoggerSingle{
                logfile:logFile,
                Logger:log.New(writer, "", flag),
        }, nil
}
//////////////////////////////////////////////////////
type CustomLoggerWithBuf struct {
        logfile       string
        Logger        *log.Logger

        Buf           []string
        FlushSize     int
        FlushTimeOut  int64

        BufChan       chan string

        LastFlushTime int64
        sync.RWMutex
}

func (c *CustomLoggerWithBuf)LogFile() string {
        return c.logfile
}
func (c *CustomLoggerWithBuf)flushProcessor() {
        for {
                logf := <-c.BufChan
                c.Buf = append(c.Buf, logf)
                if len(c.Buf) > c.FlushSize || (time_util.GetCurrentTimeStamp() - c.LastFlushTime) > c.FlushTimeOut {
                        c.Logger.Output(2, strings.Join(c.Buf, "\n"))
                        c.LastFlushTime = time_util.GetCurrentTimeStamp()
                        c.Buf = nil
                }
        }
}
func (c *CustomLoggerWithBuf)Logf(format string, v ...interface{}) {
        c.Lock()
        defer c.Unlock()
        str := fmt.Sprintf("%s %s", time.Now().Format("15:04:05"), fmt.Sprintf(format, v...))
        c.BufChan <- str

}
func (c *CustomLoggerWithBuf)LogStr(str string) {
        c.Lock()
        defer c.Unlock()
        c.BufChan <- str
}
func (c *CustomLoggerWithBuf)Log(v ...interface{}) {
        c.Lock()
        defer c.Unlock()
        str := fmt.Sprintf("%s %s", time.Now().Format("15:04:05"), fmt.Sprintln(v...))
        c.BufChan <- str

}
func (c *CustomLoggerWithBuf)LogTag(TAG, format string, v ...interface{}) {
        c.Lock()
        defer c.Unlock()
        str := fmt.Sprintf("%s [%s]%s", time.Now().Format("15:04:05"), TAG, fmt.Sprintf(format, v...))
        c.BufChan <- str
}

func NewCustomLoggerWithBuf(logFile string, num int, timelen int64) (*CustomLoggerWithBuf, error) {
        var writer io.Writer
        if logFile != "" {
                f, err := os.OpenFile(logFile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
                if err != nil {
                        log.Printf("error opening file %s: %v\n", logFile, err)
                        return nil, err
                }
                writer = f
        } else {
                writer = os.Stdout
        }
        customL := &CustomLoggerWithBuf{
                logfile:logFile,
                Logger:log.New(writer, "", log.Ltime),
                Buf:make([]string, 0),
                FlushSize:num,
                FlushTimeOut:timelen,
                BufChan:make(chan string, num * 4),
        }
        go customL.flushProcessor()
        return customL, nil
}
//////////////////////////////////////////////////////