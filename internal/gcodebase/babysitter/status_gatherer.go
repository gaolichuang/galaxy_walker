package babysitter

import (
	"fmt"
        "galaxy_walker/internal/gcodebase/string_util"
	"strings"
	"os"
        "galaxy_walker/internal/gcodebase/time_util"
	"encoding/json"
        "time"
)

var StartAtTime = time_util.GetReadableTimeNow()
var StartAtTimeStamp = time_util.GetCurrentTimeStamp()
var bbVersion = ""
func machineInfo() map[string]string {
	//TODO Statistic pid,cmd, cpunum,total mem, ip port hostname, uptime,
	machine := make(map[string]string)
	machine["cmd"] = strings.Join(os.Args, " ")
	machine["pid"] = fmt.Sprintf("%d", os.Getpid())
	machine["uid"] = fmt.Sprintf("%d", os.Getuid())
	machine["hostname"], _ = os.Hostname()
	machine["StartAt"] = StartAtTime
        machine["version"] = bbVersion
	return machine
}
func statusInfo() map[string]string {
	//TODO process mem,cpu,fd, load; dynamic
        status := make(map[string]string)
        dur := time_util.GetCurrentTimeStamp() - StartAtTimeStamp
        status["Duration"] = (time.Second*time.Duration(dur)).String()
	return status
}
func machineInfoStr() string {
        // just collect one time, information do not change
        machine := machineInfo()
        c,_ := json.Marshal(&machine)
        return string(c)
}
func machineInfoHtml() string {
	// just collect one time, information do not change
	machine := machineInfo()
	var info string
	for k, v := range machine {
		string_util.StringAppendF(&info, "<key>%s</key>:<value>%s</value><br>", k, v)
	}
	return info
}

// dynamic, collect when each request
func statusInfoHtml() string {
        status := statusInfo()
        var info string
        for k, v := range status {
                string_util.StringAppendF(&info, "<key>%s</key>:<value>%s</value><br>", k, v)
        }
        return info
}
