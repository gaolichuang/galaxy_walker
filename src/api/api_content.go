package api

import (
    "galaxy_walker/internal/github.com/gorilla/mux"
    "net/http"
)

const (
    // contentdb
    kDBOneDoc = kEndPointPreFix + "/doc"
    kDBOneDocContent = kEndPointPreFix + "/doc/content"
    kContentScanDoc = kEndPointPreFix + "/docs/{id}"
)
func (s *APIService) getOneDoc(c *APIContext, w http.ResponseWriter, r *http.Request) []byte {
    // 默认不包含content，需要content要额外处理
    // /doc?url=xxx&task=xxx
    // /doc?url=xxx
    // /doc?id=xxx&task=xxx
    // /doc?id=xxx
    vars := mux.Vars(r)
    task := vars["task"]

    return []byte(task)
}

func (s *APIService) serveDBContent(router *mux.Router) {
    router.Handle(kDBOneDoc, CommonHandlerWrapper(s.getOneDoc)).Methods("GET")
}
