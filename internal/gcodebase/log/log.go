package log

import (
        "galaxy_walker/internal/gcodebase/conf"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
        "syscall"
	"strings"
	"sync"
        "net/http"
)

/*
* Debug(Log Level) - Info - Warning - Error -
 */

// combination version, it hide log.Logger method.
type logS struct {
	logI  *log.Logger
	level int
}

var _log logS
var _error_log logS

// level0 -- 10
var levelLog []*logS

type StringObjItf interface {
        String() string
}

func (l *logS) DebugTag(TAG, format string, v ...interface{}) {
	if l.level <= *conf.Conf.LogV {
		l.logI.SetPrefix("[Debug]")
		l.logI.Output(2, fmt.Sprintf("[%s]%s", TAG, fmt.Sprintf(format, v...)))
	}
}
func (l *logS) DebugTagStrObj(TAG string, s StringObjItf, format string, v ...interface{}) {
        if l.level <= *conf.Conf.LogV {
                l.logI.SetPrefix("[Debug]")
                l.logI.Output(2, fmt.Sprintf("[%s]%s %s", TAG,fmt.Sprintf(format, v...),s.String()))
        }
}
func (l *logS) DebugTagStrObjI(TAG string, s StringObjItf) {
        if l.level <= *conf.Conf.LogV {
                l.logI.SetPrefix("[Debug]")
                l.logI.Output(2, fmt.Sprintf("[%s]%s", TAG, s.String()))
        }
}

func (l *logS) Debugf(format string, v ...interface{}) {
	if l.level <= *conf.Conf.LogV {
		l.logI.SetPrefix("[Debug]")
		l.logI.Output(2, fmt.Sprintf(format, v...))
	}
}
func (l *logS) DebugfStrObj(s StringObjItf,format string, v ...interface{}) {
        if l.level <= *conf.Conf.LogV {
                l.logI.SetPrefix("[Debug]")
                l.logI.Output(2, fmt.Sprintf(format, v...) + s.String())
        }
}

func (l *logS) Debug(v ...interface{}) {
	if l.level <= *conf.Conf.LogV {
		l.logI.SetPrefix("[Debug]")
		l.logI.Output(2, fmt.Sprintln(v...))
	}
}
func (l *logS) DebugStrObj(s StringObjItf,v ...interface{}) {
        if l.level <= *conf.Conf.LogV {
                l.logI.SetPrefix("[Debug]")
                l.logI.Output(2, fmt.Sprintln(v...) + s.String())
        }
}

func Info(v ...interface{}) {
	_log.logI.SetPrefix("[Info]")
	_log.logI.Output(2, fmt.Sprintln(v...))
}

func Infof(format string, v ...interface{}) {
	_log.logI.SetPrefix("[Info]")
	_log.logI.Output(2, fmt.Sprintf(format, v...))
}

func Warning(v ...interface{}) {
	_log.logI.SetPrefix("[Warning]")
	_log.logI.Output(2, fmt.Sprintln(v...))
}

func Warningf(format string, v ...interface{}) {
	_log.logI.SetPrefix("[Warning]")
	_log.logI.Output(2, fmt.Sprintf(format, v...))
}

func Error(v ...interface{}) {
	_error_log.logI.SetPrefix("[Error]")
	_error_log.logI.Output(2, fmt.Sprintln(v...))
}

func Errorf(format string, v ...interface{}) {
	_error_log.logI.SetPrefix("[Error]")
	_error_log.logI.Output(2, fmt.Sprintf(format, v...))
}

func Fatal(v ...interface{}) {
	_error_log.logI.SetPrefix("[Fatal]")
	_error_log.logI.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	_error_log.logI.SetPrefix("[Fatal]")
	_error_log.logI.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// VLog user pairly with Debug
// chain function call
func VLog(level int) *logS {
	if level < 0 || level >= len(levelLog) {
		_log.level = level
		return &_log
	}
	return levelLog[level]
}

var _log_instance *log.Logger = nil
var _log_init_ctx sync.Once

func NewLogger() *log.Logger {
	_log_init_ctx.Do(func() {
		var writer io.Writer
		if logFile := *conf.Conf.LogFile; logFile != "" {
			f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				log.Fatalf("error opening file: %v", err)
			}
			writer = f
		} else {
			writer = os.Stdout
		}
		if *conf.Conf.Stdout && writer != os.Stdout {
			writer = io.MultiWriter(writer, os.Stdout)
		}
		_log_instance = log.New(writer, "", log.LstdFlags|log.Lshortfile)
	})
	return _log_instance
}

var _error_log_instance *log.Logger = nil
var _error_log_init_ctx sync.Once

func NewErrorLogger() *log.Logger {
	_error_log_init_ctx.Do(func() {
		var writer io.Writer
		if logFile := *conf.Conf.ErrorLogFile; logFile != "" {
			f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				log.Fatalf("error opening file: %v", err)
			}
			writer = f
                        // redirect stderr
                        err = syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd()))
                        if err != nil {
                                log.Fatalf("Failed to redirect stderr to file: %v", err)
                        }
		} else {
			writer = os.Stdout
		}
		if *conf.Conf.Stdout && writer != os.Stdout {
			writer = io.MultiWriter(writer, os.Stdout)
		}
		_error_log_instance = log.New(writer, "", log.LstdFlags|log.Lshortfile)
	})
	return _error_log_instance
}
func init() {
	_log.logI = NewLogger()
	_error_log.logI = NewErrorLogger()
	levelLog = append(levelLog, &logS{logI: NewLogger(), level: 0})
	levelLog = append(levelLog, &logS{logI: NewLogger(), level: 1})
	levelLog = append(levelLog, &logS{logI: NewLogger(), level: 2})
	levelLog = append(levelLog, &logS{logI: NewLogger(), level: 3})
	levelLog = append(levelLog, &logS{logI: NewLogger(), level: 4})
	levelLog = append(levelLog, &logS{logI: NewLogger(), level: 5})
	levelLog = append(levelLog, &logS{logI: NewLogger(), level: 6})
	levelLog = append(levelLog, &logS{logI: NewLogger(), level: 7})
	levelLog = append(levelLog, &logS{logI: NewLogger(), level: 8})
	levelLog = append(levelLog, &logS{logI: NewLogger(), level: 9})
	levelLog = append(levelLog, &logS{logI: NewLogger(), level: 10})
	// dump flags in log, because of dependency cycle
	//DumpFlags()
}

func escapeUsage(s string) string {
	return strings.Replace(s, "\n", "\n    # ", -1)
}

func quoteValue(v string) string {
	if !strings.ContainsAny(v, "\n#;") && strings.TrimSpace(v) == v {
		return v
	}
	v = strings.Replace(v, "\\", "\\\\", -1)
	v = strings.Replace(v, "\n", "\\n", -1)
	v = strings.Replace(v, "\"", "\\\"", -1)
	return fmt.Sprintf("\"%s\"", v)
}
func DumpFlags() {
	Info("=================Dump Flags=========================================")
	flag.VisitAll(func(f *flag.Flag) {
		if f.Name != "config" && f.Name != "dumpflags" {
			Infof("%s = %s # %s\n", f.Name, quoteValue(f.Value.String()), escapeUsage(f.Usage))
		}
	})
	Info("=================Dump Flags Finish===================================")
}
func HandleFlagsFunc(w http.ResponseWriter, r *http.Request) {
        info := ""
        flag.VisitAll(func(f *flag.Flag) {
                if f.Name != "config" && f.Name != "dumpflags" {
                        info += fmt.Sprintf("%s = %s # %s\n", f.Name, quoteValue(f.Value.String()), escapeUsage(f.Usage))
                }
        })

        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        w.Write([]byte(info))
}

// test
/*
func main() {
    LOG.Info("hello world")
    LOG.VLog(1).Debug("debug message1")
    LOG.VLog(2).Debug("debug message2")
    LOG.VLog(3).Debug("debug message3")
}
//*/