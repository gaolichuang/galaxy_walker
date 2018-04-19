package babysitter

import (
    "net/http"
    "strconv"
    "fmt"
    "time"
    "strings"
    "galaxy_walker/internal/gcodebase/file"
    "galaxy_walker/internal/gcodebase/cmd"
)

func CommonSingleFileServer(filename string) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        filter := r.URL.Query().Get("filter")
        limitStr := r.URL.Query().Get("limit")
        sizeStr := r.URL.Query().Get("size")
        reversed := r.URL.Query().Get("reversed")
        var err error
        line := kMaxLimitNum
        size := kMaxSize
        if limitStr != "" {
            line, err = strconv.Atoi(limitStr)
            if err != nil {
                info := []byte(fmt.Sprintf("Line Should Set integer %v %s", err, limitStr))
                w.Header().Set("Content-Type", "text/plain")
                w.Write(info)
                return
            }
        }
        if sizeStr != "" {
            size, err = strconv.Atoi(sizeStr)
            if err != nil {
                info := []byte(fmt.Sprintf("Size Should Set integer %v %s", err, sizeStr))
                w.Header().Set("Content-Type", "text/plain")
                w.Write(info)
                return
            }
        }
        if !file.Exist(filename) {
            info := []byte(fmt.Sprintf("%s not exist", filename))
            w.Header().Set("Content-Type", "text/plain")
            w.Write(info)
            return
        }

        if line > kMaxLimitNum {
            line = kMaxLimitNum
        }
        cmds := make([][]string, 0)
        if filter != "" {
            fs := strings.Split(filter, ":")
            cmds = append(cmds, []string{"grep", fs[0], filename})
            for i := 1; i < len(fs); i++ {
                cmds = append(cmds, []string{"grep", fs[i]})
            }
        } else {
            cmds = append(cmds, []string{"tail", fmt.Sprintf("-%d", line), filename})
        }
        var output []byte
        if reversed == "true" {
            // http://stackoverflow.com/questions/742466/how-can-i-reverse-the-order-of-lines-in-a-file
            cmds = append(cmds, []string{"tac"})
            //rrcmd := []string{"sed","1!G;h;$!d"}
        }
        err, _, output, _ = cmd.RunMultiCommandWithTimeOut(time.Second*10, cmds...)
        if err != nil {
            info := []byte(fmt.Sprintf("Run Cmd Err %v %s", err, string(output)))
            w.Header().Set("Content-Type", "application/json")
            w.Write(info)
            return
        }
        if len(output) > size {
            output = output[:size]
        }
        w.Header().Set("Content-Type", "text/plain")
        w.Write(output)
    })
}
