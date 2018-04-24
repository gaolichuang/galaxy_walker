package task

type JobDescription struct {
    IsUrgent        bool     `json:"isUrgent,omitempty"`
    PrimeTag        string   `json:"primeTag,omitempty"`
    SecondTag       []string `json:"secondTag,omitempty"`
    RandomHostLoad  int      `json:"randomHostLoad,omitempty"`
    DropContent     bool     `json:"dropContent,omitempty"`
    StoreEngine     string   `json:"storeEngine,omitempty"`
    StoreDb         string   `json:"storeDb,omitempty"`
    StoreTable      string   `json:"storeTable,omitempty"`
    RequestType     int      `json:"requestType,omitempty"`
    Referer         string   `json:"referer,omitempty"`
    Custom_ua       bool     `json:"custom_ua,omitempty"`
    Follow_redirect bool     `json:"follow_redirect,omitempty"`
    Use_proxy       bool     `json:"use_proxy,omitempty"`
    // if true, nofollow href will extract too.
    NoFollow bool `json:"nofollow,omitempty"`
}
