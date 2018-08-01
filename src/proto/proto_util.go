package proto

// proto3 syntax
// https://developers.google.com/protocol-buffers/docs/proto3

import (
    LOG "galaxy_walker/internal/gcodebase/log"
    "galaxy_walker/internal/github.com/golang/protobuf/proto"
)

func Serialize(pb proto.Message) ([]byte, error) {
    r, e := proto.Marshal(pb)
    if e != nil {
        LOG.Error("MarShal failed")
    }
    return r, e
}
func Deserialize(s []byte, pb proto.Message) error {
    e := proto.Unmarshal(s, pb)
    if e != nil {
        LOG.Error("UnMarShal failed")
    }
    return e
}
func CloneCrawlDoc(doc *CrawlDoc) *CrawlDoc {
    nd := proto.Clone(doc)
    return nd.(*CrawlDoc)
}
func FromProtoToString(pb proto.Message) string {
    return proto.MarshalTextString(pb)
}
func FromStringToProto(s string, pb proto.Message) error {
    e := proto.UnmarshalText(s, pb)
    if e != nil {
        LOG.Error("UnMarShal failed")
    }
    return e
}

var RequestTypeToStringMap map[RequestType]string = map[RequestType]string{
    RequestType_TESTING:        "RequestType_TESTING",
    RequestType_WEB_StartUp:    "RequestType_WEB_StartUp",
    RequestType_WEB_MAIN:       "RequestType_WEB_MAIN",
    RequestType_WEB_HUB:        "RequestType_WEB_HUB",
    RequestType_WEB_CONTENT:    "RequestType_WEB_CONTENT",
    RequestType_WEB_SUBCONTENT: "RequestType_WEB_SUBCONTENT",
    RequestType_WEB_DETAIL:     "RequestType_WEB_DETAIL",
}

func RequestTypeToString(requestType RequestType) string {
    return RequestTypeToStringMap[requestType]
}
