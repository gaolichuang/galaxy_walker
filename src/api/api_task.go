package api

import (
    "galaxy_walker/internal/github.com/gorilla/mux"
    pb "galaxy_walker/src/proto"
    "net/http"
    "encoding/json"
    "time"
    "fmt"
    "strconv"
    "io/ioutil"
    "galaxy_walker/internal/gcodebase/time_util"
    LOG "galaxy_walker/internal/gcodebase/log"
)

const (
    // contentdb
    kTask    = kEndPointPreFix + "/tasks"
    kOneTask = kEndPointPreFix + "/tasks/{name}"
)

func (s *APIService) createTask(c *APIContext, w http.ResponseWriter, r *http.Request) []byte {
    /*
    POST /api/tasks/xxx?expire=1w4d&mode=debug/online
    body
    {JobDescription}
    */
    parserDuration := func (str string) (error,time.Duration) {
        unitMap := map[byte]time.Duration {
            'w':time.Hour*24*7,
            'd':time.Hour*24,
        }
        var ret time.Duration
        pos := 0
        for i:=0;i<len(str);i++ {
            if str[i]>='0' && str[i]<='9' {
                continue
            }
            if _,ok := unitMap[str[i]];!ok {
                return fmt.Errorf("invalid unit %s",string(str[i])),0
            }
            n,e := strconv.Atoi(str[pos:i])
            if e != nil {
                return fmt.Errorf("parse int error %s %v",str[pos:i],e),0
            }
            ret += time.Duration(n) * unitMap[str[i]]
            pos = i+1
        }
        return nil,ret
    }
    err,expire := parserDuration(r.URL.Query().Get("expire"))
    if err != nil {
        return getReturnMsg(201,"invalid expire,sample 1w2d","")
    }
    if expire <= 0 {
        expire = time.Hour*24*30 //  one month
    }
    mode := r.URL.Query().Get("mode")
    vars := mux.Vars(r)
    name := vars["name"]
    obj := s.taskdb.Get(name)
    if obj != nil {
        return getReturnMsg(201,name + " already exist.","")
    }
    body, _ := ioutil.ReadAll(r.Body)
    job := pb.JobDescription{}
    if err := json.Unmarshal(body, &job); err != nil {
        return getReturnMsg(102, fmt.Sprintf("Unmarsal Err %v %s", err, string(body)), "")
    }
    err = s.taskdb.Put(&pb.TaskDescription{
        Name:name,
        CreateAt:int(time_util.GetCurrentTimeStamp()),
        ExpireAt:int(time_util.GetCurrentTimeStamp()+int64(expire)),
        Status: func(debug bool) string {
            if debug {
                return pb.KTaskStatusDebug
            }
            return pb.KTaskStatusOnline
        }(mode=="debug"),
        Desc:&job,
    })
    if err != nil {
        return getReturnMsg(103,fmt.Sprintf("taskdb err %v",err),"")
    }
    return getReturnMsg(200,"success " + name,"")
}
func (s *APIService) updateTask(c *APIContext, w http.ResponseWriter, r *http.Request) []byte {
    /*
    PUT /api/tasks/xxx?expire=1w4d&mode=debug/online
    body
    {JobDescription}
    */
    parserDuration := func (str string) (error,time.Duration) {
        unitMap := map[byte]time.Duration {
            'w':time.Hour*24*7,
            'd':time.Hour*24,
        }
        var ret time.Duration
        pos := 0
        for i:=0;i<len(str);i++ {
            if str[i]>='0' && str[i]<='9' {
                continue
            }
            if _,ok := unitMap[str[i]];!ok {
                return fmt.Errorf("invalid unit %s",string(str[i])),0
            }
            n,e := strconv.Atoi(str[pos:i])
            if e != nil {
                return fmt.Errorf("parse int error %s %v",str[pos:i],e),0
            }
            ret += time.Duration(n) * unitMap[str[i]]
            pos = i+1
        }
        return nil,ret
    }
    err,expire := parserDuration(r.URL.Query().Get("expire"))
    if err != nil {
        return getReturnMsg(201,"invalid expire,sample 1w2d","")
    }
    var expirTimestamp  int64
    if expire > 0 {
        expirTimestamp = time_util.GetCurrentTimeStamp() + int64(expire)
    }
    umode := ""
    mode := r.URL.Query().Get("mode")
    if mode == "debug" {
        umode = "debug"
    } else if mode == "online" {
        umode = "online"
    }
    vars := mux.Vars(r)
    name := vars["name"]
    obj := s.taskdb.Get(name)
    if obj == nil {
        return getReturnMsg(201,name + " not exist.","")
    }
    body, _ := ioutil.ReadAll(r.Body)
    var job *pb.JobDescription
    if len(body) > 0 {
        job = &pb.JobDescription{}
        if err := json.Unmarshal(body, &job); err != nil {
            return getReturnMsg(102, fmt.Sprintf("Unmarsal Err %v %s", err, string(body)), "")
        }
    }
    LOG.VLog(3).DebugTag("API","updateTask %s,%d,%s,%v",name,expirTimestamp,umode,job)
    err = s.taskdb.Update(name,expirTimestamp,umode,job)
    if err != nil {
        return getReturnMsg(103,fmt.Sprintf("taskdb err %v",err),"")
    }
    return getReturnMsg(200,"success " + name,"")
}
func (s *APIService) deleteTask(c *APIContext, w http.ResponseWriter, r *http.Request) []byte {
    /*
    DELETE /api/tasks/xxx
    */
    body, _ := ioutil.ReadAll(r.Body)
    validate, msg, _ := APITokenValidation(body)
    if !validate {
        return getReturnMsg(102, msg, "")
    }
    vars := mux.Vars(r)
    name := vars["name"]
    obj := s.taskdb.Get(name)
    if obj == nil {
        return getReturnMsg(201,name + " not exist.","")
    }
    LOG.VLog(1).DebugTag("API","Delete task %v",obj)
    err := s.taskdb.Delete(name)
    if err != nil {
        return getReturnMsg(103,fmt.Sprintf("taskdb delete err %v",err),"")
    }
    return getReturnMsg(200,"success delete " + name,"")
}
func (s *APIService) getTask(c *APIContext, w http.ResponseWriter, r *http.Request) []byte {
    vars := mux.Vars(r)
    name := vars["name"]
    obj := s.taskdb.Get(name)
    info, _ := json.Marshal(&obj)
    return info
}
func (s *APIService) listTask(c *APIContext, w http.ResponseWriter, r *http.Request) []byte {
    // /tasks?detail=true
    detail := r.URL.Query().Get("detail")
    tasks, taskObjs := s.taskdb.List()
    if detail == "true" {
        info, _ := json.Marshal(&taskObjs)
        return info
    }
    info, _ := json.Marshal(&tasks)
    return info
}

func (s *APIService) serveTaskAPI(router *mux.Router) {
    router.Handle(kTask, CommonHandlerWrapper(s.listTask)).Methods("GET")
    router.Handle(kOneTask, CommonHandlerWrapper(s.getTask)).Methods("GET")
    router.Handle(kOneTask, CommonHandlerWrapper(s.createTask)).Methods("POST")   // create
    router.Handle(kOneTask, CommonHandlerWrapper(s.updateTask)).Methods("PUT")    // update
    router.Handle(kOneTask, CommonHandlerWrapper(s.deleteTask)).Methods("DELETE") // delete
}
