package proto

import "encoding/json"

type JobDescription struct {
    IsUrgent        bool     `json:"isUrgent,omitempty"`
    PrimeTag        string   `json:"primeTag,omitempty"`
    SecondTag       []string `json:"secondTag,omitempty"`
    RandomHostLoad  int      `json:"randomHostLoad,omitempty"`
    DropContent     bool     `json:"dropContent,omitempty"`
    RequestType     int      `json:"requestType,omitempty"`
    Referer         string   `json:"referer,omitempty"`
    Custom_ua       bool     `json:"custom_ua,omitempty"`
    Follow_redirect bool     `json:"follow_redirect,omitempty"`
    Use_proxy       bool     `json:"use_proxy,omitempty"`
    // if true, nofollow href will extract too.
    NoFollow bool `json:"nofollow,omitempty"`
}
func (j *JobDescription)ToString() string {
    if j == nil {return ""}
    s,_ := json.Marshal(j)
    return string(s)
}

const (
    KTaskStatusDebug = "debug"
    KTaskStatusOnline = "online"
    KTaskStatusStarting = "starting"
)
type TaskDescription struct {
    Name string `json:"name,omitempty"`
    // Status debug, online, starting.
    Status string `json:"status,omitempty"`
    Desc *JobDescription `json:"job,omitempty"`
    ExpireAt int `json:"expireAt,omitempty"`
    CreateAt int `json:"createAt,omitempty"`
    Report  *TaskReport `json:"report,omitempty"`
}

type LevelPair struct {
    Level string `json:"level,omitempty"`
    Total int `json:"total,omitempty"`
}
type TaskReport struct {
    Total int `json:"total,omitempty"`
    Success int `json:"success,omitempty"`
    Fail1 int `json:"fail1,omitempty"`
    Fail2 int `json:"fail2,omitempty"`
    Fail3 int `json:"fail3,omitempty"`
    Levels []*LevelPair `json:"levelPair,omitempty"`
}

func IsFreshTask(t *TaskDescription) bool {
    return t.Status == KTaskStatusOnline
}