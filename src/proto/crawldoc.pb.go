// Code generated by protoc-gen-go.
// source: crawldoc.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	crawldoc.proto

It has these top-level messages:
	ConnectionInfo
	OutLink
	FetchHint
	CrawlRecord
	CrawlParam
	CrawlChain
	CrawlDoc
	CrawlDocs
	CrawlRequest
	CrawlResponse
*/
package proto

import proto1 "galaxy_walker/internal/github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "galaxy_walker/internal/golang.org/x/net/context"
	grpc "galaxy_walker/internal/google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto1.ProtoPackageIsVersion1

type ReturnType int32

const (
	ReturnType_UNKNOWN ReturnType = 0
	//    NODNS                   = 1;
	ReturnType_NOCONNECTION ReturnType = 2
	//    FORBIDDENROBOTS         = 3;
	ReturnType_TIMEOUT   ReturnType = 4
	ReturnType_BADTYPE   ReturnType = 5
	ReturnType_TOOBIG    ReturnType = 6
	ReturnType_BADHEADER ReturnType = 7
	//    NETWORKERROR            = 8;
	ReturnType_SITEQUEUEFULLFETCHER ReturnType = 9
	// url format is not avaliable  crawl/base/utils.IsInvalidUrl()
	ReturnType_INVALIDURL         ReturnType = 10
	ReturnType_INVALIDREDIRECTURL ReturnType = 11
	//    META_REDIRECT           = 12;
	//    JS_REDIRECT             = 13;
	//    IP_BLACKLISTED          = 14;
	//    BADCONTENT              = 15;
	//    URL_BLACKLISTED         = 16;
	ReturnType_SITEQUEUEFULLDISPATCHER ReturnType = 17
	ReturnType_STATUS100               ReturnType = 100
	ReturnType_STATUS101               ReturnType = 101
	ReturnType_STATUS200               ReturnType = 200
	ReturnType_STATUS201               ReturnType = 201
	ReturnType_STATUS202               ReturnType = 202
	ReturnType_STATUS203               ReturnType = 203
	ReturnType_STATUS204               ReturnType = 204
	ReturnType_STATUS205               ReturnType = 205
	ReturnType_STATUS206               ReturnType = 206
	ReturnType_STATUS300               ReturnType = 300
	ReturnType_STATUS301               ReturnType = 301
	ReturnType_STATUS302               ReturnType = 302
	ReturnType_STATUS303               ReturnType = 303
	ReturnType_STATUS304               ReturnType = 304
	ReturnType_STATUS305               ReturnType = 305
	ReturnType_STATUS306               ReturnType = 306
	ReturnType_STATUS307               ReturnType = 307
	ReturnType_STATUS400               ReturnType = 400
	ReturnType_STATUS401               ReturnType = 401
	ReturnType_STATUS402               ReturnType = 402
	ReturnType_STATUS403               ReturnType = 403
	ReturnType_STATUS404               ReturnType = 404
	ReturnType_STATUS405               ReturnType = 405
	ReturnType_STATUS406               ReturnType = 406
	ReturnType_STATUS407               ReturnType = 407
	ReturnType_STATUS408               ReturnType = 408
	ReturnType_STATUS409               ReturnType = 409
	ReturnType_STATUS410               ReturnType = 410
	ReturnType_STATUS411               ReturnType = 411
	ReturnType_STATUS412               ReturnType = 412
	ReturnType_STATUS413               ReturnType = 413
	ReturnType_STATUS414               ReturnType = 414
	ReturnType_STATUS415               ReturnType = 415
	ReturnType_STATUS416               ReturnType = 416
	ReturnType_STATUS417               ReturnType = 417
	ReturnType_STATUS500               ReturnType = 500
	ReturnType_STATUS501               ReturnType = 501
	ReturnType_STATUS502               ReturnType = 502
	ReturnType_STATUS503               ReturnType = 503
	ReturnType_STATUS504               ReturnType = 504
	ReturnType_STATUS505               ReturnType = 505
	ReturnType_STATUS509               ReturnType = 509
	ReturnType_STATUS510               ReturnType = 510
)

var ReturnType_name = map[int32]string{
	0:   "UNKNOWN",
	2:   "NOCONNECTION",
	4:   "TIMEOUT",
	5:   "BADTYPE",
	6:   "TOOBIG",
	7:   "BADHEADER",
	9:   "SITEQUEUEFULLFETCHER",
	10:  "INVALIDURL",
	11:  "INVALIDREDIRECTURL",
	17:  "SITEQUEUEFULLDISPATCHER",
	100: "STATUS100",
	101: "STATUS101",
	200: "STATUS200",
	201: "STATUS201",
	202: "STATUS202",
	203: "STATUS203",
	204: "STATUS204",
	205: "STATUS205",
	206: "STATUS206",
	300: "STATUS300",
	301: "STATUS301",
	302: "STATUS302",
	303: "STATUS303",
	304: "STATUS304",
	305: "STATUS305",
	306: "STATUS306",
	307: "STATUS307",
	400: "STATUS400",
	401: "STATUS401",
	402: "STATUS402",
	403: "STATUS403",
	404: "STATUS404",
	405: "STATUS405",
	406: "STATUS406",
	407: "STATUS407",
	408: "STATUS408",
	409: "STATUS409",
	410: "STATUS410",
	411: "STATUS411",
	412: "STATUS412",
	413: "STATUS413",
	414: "STATUS414",
	415: "STATUS415",
	416: "STATUS416",
	417: "STATUS417",
	500: "STATUS500",
	501: "STATUS501",
	502: "STATUS502",
	503: "STATUS503",
	504: "STATUS504",
	505: "STATUS505",
	509: "STATUS509",
	510: "STATUS510",
}
var ReturnType_value = map[string]int32{
	"UNKNOWN":                 0,
	"NOCONNECTION":            2,
	"TIMEOUT":                 4,
	"BADTYPE":                 5,
	"TOOBIG":                  6,
	"BADHEADER":               7,
	"SITEQUEUEFULLFETCHER":    9,
	"INVALIDURL":              10,
	"INVALIDREDIRECTURL":      11,
	"SITEQUEUEFULLDISPATCHER": 17,
	"STATUS100":               100,
	"STATUS101":               101,
	"STATUS200":               200,
	"STATUS201":               201,
	"STATUS202":               202,
	"STATUS203":               203,
	"STATUS204":               204,
	"STATUS205":               205,
	"STATUS206":               206,
	"STATUS300":               300,
	"STATUS301":               301,
	"STATUS302":               302,
	"STATUS303":               303,
	"STATUS304":               304,
	"STATUS305":               305,
	"STATUS306":               306,
	"STATUS307":               307,
	"STATUS400":               400,
	"STATUS401":               401,
	"STATUS402":               402,
	"STATUS403":               403,
	"STATUS404":               404,
	"STATUS405":               405,
	"STATUS406":               406,
	"STATUS407":               407,
	"STATUS408":               408,
	"STATUS409":               409,
	"STATUS410":               410,
	"STATUS411":               411,
	"STATUS412":               412,
	"STATUS413":               413,
	"STATUS414":               414,
	"STATUS415":               415,
	"STATUS416":               416,
	"STATUS417":               417,
	"STATUS500":               500,
	"STATUS501":               501,
	"STATUS502":               502,
	"STATUS503":               503,
	"STATUS504":               504,
	"STATUS505":               505,
	"STATUS509":               509,
	"STATUS510":               510,
}

func (x ReturnType) String() string {
	return proto1.EnumName(ReturnType_name, int32(x))
}
func (ReturnType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// tag for source. you can custom
// like primary_tag and second_tag
type RequestType int32

const (
	RequestType_TESTING        RequestType = 0
	RequestType_WEB_StartUp    RequestType = 1
	RequestType_WEB_MAIN       RequestType = 2
	RequestType_WEB_HUB        RequestType = 3
	RequestType_WEB_CONTENT    RequestType = 4
	RequestType_WEB_SUBCONTENT RequestType = 5
	RequestType_WEB_DETAIL     RequestType = 6
)

var RequestType_name = map[int32]string{
	0: "TESTING",
	1: "WEB_StartUp",
	2: "WEB_MAIN",
	3: "WEB_HUB",
	4: "WEB_CONTENT",
	5: "WEB_SUBCONTENT",
	6: "WEB_DETAIL",
}
var RequestType_value = map[string]int32{
	"TESTING":        0,
	"WEB_StartUp":    1,
	"WEB_MAIN":       2,
	"WEB_HUB":        3,
	"WEB_CONTENT":    4,
	"WEB_SUBCONTENT": 5,
	"WEB_DETAIL":     6,
}

func (x RequestType) String() string {
	return proto1.EnumName(RequestType_name, int32(x))
}
func (RequestType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type Priority int32

const (
	Priority_NORMAL Priority = 0
	Priority_URGENT Priority = 1
)

var Priority_name = map[int32]string{
	0: "NORMAL",
	1: "URGENT",
}
var Priority_value = map[string]int32{
	"NORMAL": 0,
	"URGENT": 1,
}

func (x Priority) String() string {
	return proto1.EnumName(Priority_name, int32(x))
}
func (Priority) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type ConnectionInfo struct {
	Host string `protobuf:"bytes,1,opt,name=host" json:"host,omitempty"`
	Port int32  `protobuf:"varint,2,opt,name=port" json:"port,omitempty"`
}

func (m *ConnectionInfo) Reset()                    { *m = ConnectionInfo{} }
func (m *ConnectionInfo) String() string            { return proto1.CompactTextString(m) }
func (*ConnectionInfo) ProtoMessage()               {}
func (*ConnectionInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type OutLink struct {
	Url  string `protobuf:"bytes,1,opt,name=url" json:"url,omitempty"`
	Text string `protobuf:"bytes,2,opt,name=text" json:"text,omitempty"`
	Num  int32  `protobuf:"varint,3,opt,name=num" json:"num,omitempty"`
}

func (m *OutLink) Reset()                    { *m = OutLink{} }
func (m *OutLink) String() string            { return proto1.CompactTextString(m) }
func (*OutLink) ProtoMessage()               {}
func (*OutLink) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type FetchHint struct {
	Host string `protobuf:"bytes,1,opt,name=host" json:"host,omitempty"`
	Path string `protobuf:"bytes,2,opt,name=path" json:"path,omitempty"`
}

func (m *FetchHint) Reset()                    { *m = FetchHint{} }
func (m *FetchHint) String() string            { return proto1.CompactTextString(m) }
func (*FetchHint) ProtoMessage()               {}
func (*FetchHint) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type CrawlRecord struct {
	RequestTime int64           `protobuf:"varint,1,opt,name=request_time" json:"request_time,omitempty"`
	Fetcher     *ConnectionInfo `protobuf:"bytes,2,opt,name=fetcher" json:"fetcher,omitempty"`
	// time the doc fetched
	FetchTime    int64  `protobuf:"varint,3,opt,name=fetch_time" json:"fetch_time,omitempty"`
	FetchUseInms int64  `protobuf:"varint,4,opt,name=fetch_use_inms" json:"fetch_use_inms,omitempty"`
	ParentDocid  string `protobuf:"bytes,10,opt,name=parent_docid" json:"parent_docid,omitempty"`
}

func (m *CrawlRecord) Reset()                    { *m = CrawlRecord{} }
func (m *CrawlRecord) String() string            { return proto1.CompactTextString(m) }
func (*CrawlRecord) ProtoMessage()               {}
func (*CrawlRecord) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *CrawlRecord) GetFetcher() *ConnectionInfo {
	if m != nil {
		return m.Fetcher
	}
	return nil
}

type CrawlParam struct {
	Pri            Priority          `protobuf:"varint,1,opt,name=pri,enum=proto.Priority" json:"pri,omitempty"`
	Hostload       int32             `protobuf:"varint,2,opt,name=hostload" json:"hostload,omitempty"`
	RandomHostload int32             `protobuf:"varint,3,opt,name=random_hostload" json:"random_hostload,omitempty"`
	FetcherCount   int32             `protobuf:"varint,4,opt,name=fetcher_count" json:"fetcher_count,omitempty"`
	DropContent    bool              `protobuf:"varint,5,opt,name=drop_content" json:"drop_content,omitempty"`
	Nofollow       bool              `protobuf:"varint,6,opt,name=nofollow" json:"nofollow,omitempty"`
	NoExtractLink  bool              `protobuf:"varint,7,opt,name=noExtractLink" json:"noExtractLink,omitempty"`
	FakeHost       string            `protobuf:"bytes,10,opt,name=fake_host" json:"fake_host,omitempty"`
	FetchHint      *FetchHint        `protobuf:"bytes,11,opt,name=fetch_hint" json:"fetch_hint,omitempty"`
	Receivers      []*ConnectionInfo `protobuf:"bytes,12,rep,name=receivers" json:"receivers,omitempty"`
	Referer        string            `protobuf:"bytes,20,opt,name=referer" json:"referer,omitempty"`
	CustomUa       bool              `protobuf:"varint,21,opt,name=custom_ua" json:"custom_ua,omitempty"`
	UseProxy       bool              `protobuf:"varint,22,opt,name=use_proxy" json:"use_proxy,omitempty"`
	FollowRedirect bool              `protobuf:"varint,23,opt,name=follow_redirect" json:"follow_redirect,omitempty"`
	Taskid         string            `protobuf:"bytes,100,opt,name=taskid" json:"taskid,omitempty"`
	Rtype          RequestType       `protobuf:"varint,101,opt,name=rtype,enum=proto.RequestType" json:"rtype,omitempty"`
	PrimaryTag     string            `protobuf:"bytes,102,opt,name=primary_tag" json:"primary_tag,omitempty"`
	SecondaryTag   []string          `protobuf:"bytes,103,rep,name=secondary_tag" json:"secondary_tag,omitempty"`
}

func (m *CrawlParam) Reset()                    { *m = CrawlParam{} }
func (m *CrawlParam) String() string            { return proto1.CompactTextString(m) }
func (*CrawlParam) ProtoMessage()               {}
func (*CrawlParam) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *CrawlParam) GetFetchHint() *FetchHint {
	if m != nil {
		return m.FetchHint
	}
	return nil
}

func (m *CrawlParam) GetReceivers() []*ConnectionInfo {
	if m != nil {
		return m.Receivers
	}
	return nil
}

// for follow redirect.
type CrawlChain struct {
	// TODO.
	Docid      uint32 `protobuf:"varint,1,opt,name=docid" json:"docid,omitempty"`
	RequestUrl string `protobuf:"bytes,2,opt,name=request_url" json:"request_url,omitempty"`
	// url use to crawl. it's generated base on request url
	Url string `protobuf:"bytes,3,opt,name=url" json:"url,omitempty"`
	// fill at fetcher, http response information
	RedirectUrl string     `protobuf:"bytes,4,opt,name=redirect_url" json:"redirect_url,omitempty"`
	Code        ReturnType `protobuf:"varint,5,opt,name=code,enum=proto.ReturnType" json:"code,omitempty"`
	// record error.Errors.
	ErrorInfo     string `protobuf:"bytes,6,opt,name=error_info" json:"error_info,omitempty"`
	Content       string `protobuf:"bytes,7,opt,name=content" json:"content,omitempty"`
	ContentLength int64  `protobuf:"varint,8,opt,name=content_length" json:"content_length,omitempty"`
	// response header.
	Header string `protobuf:"bytes,10,opt,name=header" json:"header,omitempty"`
}

func (m *CrawlChain) Reset()                    { *m = CrawlChain{} }
func (m *CrawlChain) String() string            { return proto1.CompactTextString(m) }
func (*CrawlChain) ProtoMessage()               {}
func (*CrawlChain) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type CrawlDoc struct {
	Docid      uint32 `protobuf:"varint,1,opt,name=docid" json:"docid,omitempty"`
	RequestUrl string `protobuf:"bytes,2,opt,name=request_url" json:"request_url,omitempty"`
	// url use to crawl. it's generated base on request url
	Url string `protobuf:"bytes,3,opt,name=url" json:"url,omitempty"`
	// fill at fetcher, http response information
	RedirectUrl string     `protobuf:"bytes,4,opt,name=redirect_url" json:"redirect_url,omitempty"`
	Code        ReturnType `protobuf:"varint,5,opt,name=code,enum=proto.ReturnType" json:"code,omitempty"`
	// record error.Errors.
	ErrorInfo     string `protobuf:"bytes,6,opt,name=error_info" json:"error_info,omitempty"`
	Content       string `protobuf:"bytes,7,opt,name=content" json:"content,omitempty"`
	ContentLength int64  `protobuf:"varint,8,opt,name=content_length" json:"content_length,omitempty"`
	// compress at storage handler.
	ContentCompressed bool `protobuf:"varint,9,opt,name=content_compressed" json:"content_compressed,omitempty"`
	// response header.
	Header     string `protobuf:"bytes,10,opt,name=header" json:"header,omitempty"`
	LastModify string `protobuf:"bytes,11,opt,name=last_modify" json:"last_modify,omitempty"`
	// content type of the page, eg text/html
	ContentType       string     `protobuf:"bytes,12,opt,name=content_type" json:"content_type,omitempty"`
	IndomainOutlinks  []*OutLink `protobuf:"bytes,30,rep,name=indomain_outlinks" json:"indomain_outlinks,omitempty"`
	OutdomainOutlinks []*OutLink `protobuf:"bytes,31,rep,name=outdomain_outlinks" json:"outdomain_outlinks,omitempty"`
	// content hash(128 bit)
	// int64 chash_0   =   32;
	// int64 chash_1   =   33;
	// original encoding which is deteched by the page content
	OrigEncoding string `protobuf:"bytes,34,opt,name=orig_encoding" json:"orig_encoding,omitempty"`
	// encoding after convert to utf8
	// the same with orig_encoding if convert fail
	// utf8 if convert success
	ConvEncoding string       `protobuf:"bytes,35,opt,name=conv_encoding" json:"conv_encoding,omitempty"`
	Reservation  string       `protobuf:"bytes,40,opt,name=reservation" json:"reservation,omitempty"`
	CrawlParam   *CrawlParam  `protobuf:"bytes,50,opt,name=crawl_param" json:"crawl_param,omitempty"`
	CrawlRecord  *CrawlRecord `protobuf:"bytes,60,opt,name=crawl_record" json:"crawl_record,omitempty"`
	// http connection will follow redirect by CrawlParams.follow_redirect. so this field is not nessary and useful
	RedirectChain []*CrawlChain `protobuf:"bytes,70,rep,name=redirectChain" json:"redirectChain,omitempty"`
}

func (m *CrawlDoc) Reset()                    { *m = CrawlDoc{} }
func (m *CrawlDoc) String() string            { return proto1.CompactTextString(m) }
func (*CrawlDoc) ProtoMessage()               {}
func (*CrawlDoc) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *CrawlDoc) GetIndomainOutlinks() []*OutLink {
	if m != nil {
		return m.IndomainOutlinks
	}
	return nil
}

func (m *CrawlDoc) GetOutdomainOutlinks() []*OutLink {
	if m != nil {
		return m.OutdomainOutlinks
	}
	return nil
}

func (m *CrawlDoc) GetCrawlParam() *CrawlParam {
	if m != nil {
		return m.CrawlParam
	}
	return nil
}

func (m *CrawlDoc) GetCrawlRecord() *CrawlRecord {
	if m != nil {
		return m.CrawlRecord
	}
	return nil
}

func (m *CrawlDoc) GetRedirectChain() []*CrawlChain {
	if m != nil {
		return m.RedirectChain
	}
	return nil
}

type CrawlDocs struct {
	Docs []*CrawlDoc `protobuf:"bytes,1,rep,name=docs" json:"docs,omitempty"`
}

func (m *CrawlDocs) Reset()                    { *m = CrawlDocs{} }
func (m *CrawlDocs) String() string            { return proto1.CompactTextString(m) }
func (*CrawlDocs) ProtoMessage()               {}
func (*CrawlDocs) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *CrawlDocs) GetDocs() []*CrawlDoc {
	if m != nil {
		return m.Docs
	}
	return nil
}

type CrawlRequest struct {
	Request string `protobuf:"bytes,1,opt,name=request" json:"request,omitempty"`
}

func (m *CrawlRequest) Reset()                    { *m = CrawlRequest{} }
func (m *CrawlRequest) String() string            { return proto1.CompactTextString(m) }
func (*CrawlRequest) ProtoMessage()               {}
func (*CrawlRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

type CrawlResponse struct {
	Ok  bool  `protobuf:"varint,1,opt,name=ok" json:"ok,omitempty"`
	Ret int64 `protobuf:"varint,2,opt,name=ret" json:"ret,omitempty"`
}

func (m *CrawlResponse) Reset()                    { *m = CrawlResponse{} }
func (m *CrawlResponse) String() string            { return proto1.CompactTextString(m) }
func (*CrawlResponse) ProtoMessage()               {}
func (*CrawlResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func init() {
	proto1.RegisterType((*ConnectionInfo)(nil), "proto.ConnectionInfo")
	proto1.RegisterType((*OutLink)(nil), "proto.OutLink")
	proto1.RegisterType((*FetchHint)(nil), "proto.FetchHint")
	proto1.RegisterType((*CrawlRecord)(nil), "proto.CrawlRecord")
	proto1.RegisterType((*CrawlParam)(nil), "proto.CrawlParam")
	proto1.RegisterType((*CrawlChain)(nil), "proto.CrawlChain")
	proto1.RegisterType((*CrawlDoc)(nil), "proto.CrawlDoc")
	proto1.RegisterType((*CrawlDocs)(nil), "proto.CrawlDocs")
	proto1.RegisterType((*CrawlRequest)(nil), "proto.CrawlRequest")
	proto1.RegisterType((*CrawlResponse)(nil), "proto.CrawlResponse")
	proto1.RegisterEnum("proto.ReturnType", ReturnType_name, ReturnType_value)
	proto1.RegisterEnum("proto.RequestType", RequestType_name, RequestType_value)
	proto1.RegisterEnum("proto.Priority", Priority_name, Priority_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Client API for CrawlService service

type CrawlServiceClient interface {
	Feed(ctx context.Context, in *CrawlDocs, opts ...grpc.CallOption) (*CrawlResponse, error)
	IsHealthy(ctx context.Context, in *CrawlRequest, opts ...grpc.CallOption) (*CrawlResponse, error)
}

type crawlServiceClient struct {
	cc *grpc.ClientConn
}

func NewCrawlServiceClient(cc *grpc.ClientConn) CrawlServiceClient {
	return &crawlServiceClient{cc}
}

func (c *crawlServiceClient) Feed(ctx context.Context, in *CrawlDocs, opts ...grpc.CallOption) (*CrawlResponse, error) {
	out := new(CrawlResponse)
	err := grpc.Invoke(ctx, "/proto.CrawlService/Feed", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *crawlServiceClient) IsHealthy(ctx context.Context, in *CrawlRequest, opts ...grpc.CallOption) (*CrawlResponse, error) {
	out := new(CrawlResponse)
	err := grpc.Invoke(ctx, "/proto.CrawlService/IsHealthy", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for CrawlService service

type CrawlServiceServer interface {
	Feed(context.Context, *CrawlDocs) (*CrawlResponse, error)
	IsHealthy(context.Context, *CrawlRequest) (*CrawlResponse, error)
}

func RegisterCrawlServiceServer(s *grpc.Server, srv CrawlServiceServer) {
	s.RegisterService(&_CrawlService_serviceDesc, srv)
}

func _CrawlService_Feed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(CrawlDocs)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(CrawlServiceServer).Feed(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _CrawlService_IsHealthy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(CrawlRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(CrawlServiceServer).IsHealthy(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _CrawlService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.CrawlService",
	HandlerType: (*CrawlServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Feed",
			Handler:    _CrawlService_Feed_Handler,
		},
		{
			MethodName: "IsHealthy",
			Handler:    _CrawlService_IsHealthy_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

var fileDescriptor0 = []byte{
	// 1520 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe4, 0x56, 0x5b, 0x73, 0x1a, 0x47,
	0x16, 0xb6, 0x84, 0x40, 0xd0, 0x20, 0xd4, 0xea, 0x95, 0xed, 0x29, 0x6f, 0xed, 0xda, 0x8b, 0x77,
	0x6b, 0x55, 0xae, 0x5a, 0xaf, 0x00, 0x5d, 0xec, 0xda, 0xcd, 0x03, 0x42, 0x23, 0x8b, 0x0a, 0x02,
	0x65, 0x04, 0x71, 0xe5, 0x89, 0x9a, 0x40, 0x23, 0x4d, 0x09, 0x66, 0x48, 0xcf, 0x60, 0x5b, 0xc9,
	0x7b, 0x9e, 0x73, 0x71, 0xe2, 0xdc, 0x2f, 0xef, 0xb9, 0xe7, 0x4f, 0xc4, 0xb9, 0xfd, 0x93, 0xbc,
	0xe5, 0xfa, 0x90, 0x54, 0x4e, 0x9f, 0xee, 0x41, 0xb4, 0x2b, 0xfa, 0x05, 0x79, 0xa2, 0xbf, 0xef,
	0x9c, 0xd3, 0x7d, 0xfa, 0xf4, 0xf9, 0xce, 0x40, 0xf2, 0x5d, 0xe1, 0xde, 0x1d, 0xf4, 0x82, 0xee,
	0xf5, 0x91, 0x08, 0xa2, 0x80, 0x25, 0xf1, 0xa7, 0x70, 0x83, 0xe4, 0xab, 0x81, 0xef, 0xf3, 0x6e,
	0xe4, 0x05, 0x7e, 0xcd, 0xef, 0x07, 0x8c, 0x91, 0xb9, 0xa3, 0x20, 0x8c, 0xac, 0x99, 0x2b, 0x33,
	0x2b, 0x19, 0x07, 0xd7, 0x92, 0x1b, 0x05, 0x22, 0xb2, 0x66, 0x81, 0x4b, 0x3a, 0xb8, 0x2e, 0x54,
	0xc8, 0x7c, 0x73, 0x1c, 0xd5, 0x3d, 0xff, 0x98, 0x51, 0x92, 0x18, 0x8b, 0x81, 0x8e, 0x90, 0x4b,
	0x19, 0x10, 0xf1, 0x7b, 0x2a, 0x00, 0x36, 0x91, 0x6b, 0xe9, 0xe5, 0x8f, 0x87, 0x56, 0x02, 0xf7,
	0x90, 0xcb, 0x42, 0x99, 0x64, 0x76, 0x78, 0xd4, 0x3d, 0xda, 0xf5, 0xfc, 0xe8, 0xcc, 0x73, 0xdd,
	0xe8, 0x28, 0xde, 0x46, 0xae, 0x0b, 0x0f, 0x67, 0x48, 0xb6, 0x2a, 0xef, 0xe2, 0xf0, 0x6e, 0x20,
	0x7a, 0xec, 0x1f, 0x24, 0x27, 0xf8, 0x33, 0x63, 0x1e, 0x46, 0x9d, 0xc8, 0x1b, 0x72, 0x8c, 0x4f,
	0x38, 0x59, 0xcd, 0xb5, 0x80, 0x62, 0xff, 0x25, 0xf3, 0x7d, 0x79, 0x0e, 0x17, 0xb8, 0x53, 0xb6,
	0x74, 0x5e, 0x15, 0xe1, 0xba, 0x79, 0x75, 0x27, 0xf6, 0x62, 0x7f, 0x23, 0x04, 0x97, 0x6a, 0xc7,
	0x04, 0xee, 0x98, 0x41, 0x06, 0xf7, 0xfb, 0x27, 0xc9, 0x2b, 0xf3, 0x38, 0xe4, 0x1d, 0xcf, 0x1f,
	0x86, 0xd6, 0x1c, 0xba, 0xe4, 0x90, 0x6d, 0x87, 0xbc, 0x06, 0x9c, 0x4c, 0x6c, 0xe4, 0x0a, 0xee,
	0x47, 0x1d, 0xa8, 0xba, 0xd7, 0xb3, 0x08, 0x5e, 0x22, 0xab, 0xb8, 0x6d, 0x49, 0x15, 0x9e, 0x4f,
	0x12, 0x82, 0x77, 0xd9, 0x77, 0x85, 0x3b, 0x84, 0x88, 0xc4, 0x48, 0x78, 0x78, 0x83, 0x7c, 0x69,
	0x51, 0xe7, 0xb8, 0x2f, 0xbc, 0x40, 0x78, 0xd1, 0x89, 0x23, 0x6d, 0xec, 0x12, 0x49, 0xcb, 0xca,
	0x0c, 0x02, 0xb7, 0xa7, 0x5f, 0x63, 0x82, 0xd9, 0xbf, 0xc9, 0xa2, 0x70, 0xfd, 0x5e, 0x30, 0xec,
	0x4c, 0x5c, 0x54, 0xb1, 0xf3, 0x8a, 0xde, 0x8d, 0x1d, 0xaf, 0x92, 0x05, 0x7d, 0xd3, 0x4e, 0x37,
	0x18, 0xfb, 0x11, 0xa6, 0x9f, 0xd4, 0xe9, 0x73, 0x51, 0x95, 0x9c, 0x4c, 0xbf, 0x27, 0x82, 0x11,
	0x78, 0xf8, 0x11, 0x24, 0x6c, 0x25, 0xc1, 0x27, 0xed, 0x64, 0x25, 0x57, 0x55, 0x94, 0x4c, 0xc6,
	0x0f, 0xfa, 0xc1, 0x60, 0x10, 0xdc, 0xb5, 0x52, 0x68, 0x9e, 0x60, 0xa8, 0xd1, 0x82, 0x1f, 0xd8,
	0xf7, 0x22, 0xe1, 0x76, 0xb1, 0x49, 0xac, 0x79, 0x74, 0x30, 0x49, 0xf6, 0x57, 0x92, 0xe9, 0xbb,
	0xc7, 0x1c, 0x13, 0xd6, 0x05, 0x4a, 0x4b, 0x42, 0xa6, 0x0a, 0xcf, 0xa6, 0x5f, 0xe1, 0x08, 0xfa,
	0xc3, 0xca, 0xe2, 0xcb, 0x51, 0x5d, 0x95, 0x49, 0xdf, 0xe8, 0x77, 0xc1, 0x16, 0x82, 0x7e, 0x12,
	0xbc, 0xcb, 0xbd, 0x3b, 0x5c, 0x84, 0x56, 0xee, 0x4a, 0xe2, 0xec, 0x97, 0x3e, 0xf5, 0x63, 0x16,
	0x99, 0x17, 0xbc, 0xcf, 0x05, 0x34, 0xc7, 0x32, 0x26, 0x10, 0x43, 0x99, 0x5c, 0x77, 0x1c, 0x46,
	0x50, 0xcf, 0xb1, 0x6b, 0x9d, 0x57, 0xf7, 0x53, 0x44, 0xdb, 0x95, 0x46, 0xf9, 0xfa, 0xb0, 0xfb,
	0xbd, 0x13, 0xeb, 0x82, 0x32, 0x02, 0xb1, 0x2f, 0xb1, 0x7c, 0x09, 0x55, 0x86, 0x8e, 0xe0, 0x3d,
	0x0f, 0xce, 0x8a, 0xac, 0x8b, 0xe8, 0x92, 0x57, 0xb4, 0xa3, 0x59, 0x76, 0x81, 0xa4, 0x22, 0x37,
	0x3c, 0x86, 0xee, 0xe8, 0xe1, 0xd9, 0x1a, 0xb1, 0x15, 0x92, 0x14, 0xd1, 0xc9, 0x88, 0x5b, 0x1c,
	0x7b, 0x81, 0xe9, 0x5b, 0x38, 0xba, 0xa9, 0xc1, 0xe2, 0x28, 0x07, 0x76, 0x99, 0x64, 0xa1, 0x2f,
	0x86, 0xae, 0x38, 0xe9, 0x44, 0xee, 0xa1, 0xd5, 0xc7, 0x6d, 0x88, 0xa6, 0x5a, 0xee, 0xa1, 0x7c,
	0xec, 0x10, 0x94, 0xe2, 0xf7, 0x62, 0x97, 0x43, 0x28, 0x4c, 0xc6, 0xc9, 0x4d, 0x48, 0x70, 0x2a,
	0x3c, 0x98, 0xd5, 0x8d, 0x58, 0x3d, 0x72, 0x3d, 0x9f, 0x2d, 0x93, 0xa4, 0xea, 0x59, 0xd9, 0x8a,
	0x0b, 0x8e, 0x02, 0xf2, 0xa8, 0x58, 0x69, 0x52, 0xee, 0x4a, 0x94, 0x44, 0x53, 0x6d, 0x50, 0xbd,
	0x9e, 0x03, 0x89, 0xd3, 0x39, 0x80, 0xe2, 0x54, 0x77, 0xc5, 0x98, 0x39, 0xa5, 0x81, 0x98, 0x93,
	0x41, 0xff, 0x22, 0x73, 0xdd, 0xa0, 0xc7, 0xb1, 0xbf, 0xf2, 0xa5, 0xa5, 0xc9, 0x4d, 0xa3, 0xb1,
	0xf0, 0xf1, 0xa2, 0x68, 0x96, 0x92, 0xe4, 0x42, 0x04, 0x02, 0xf4, 0xd6, 0x0f, 0xb0, 0xdb, 0x32,
	0x4e, 0x06, 0x19, 0x9c, 0x5a, 0xf0, 0x8a, 0x71, 0xa3, 0xce, 0xab, 0x57, 0xd4, 0x10, 0xf6, 0xcf,
	0xeb, 0x65, 0x67, 0xc0, 0xfd, 0x43, 0x98, 0x26, 0x69, 0x14, 0xeb, 0x82, 0x66, 0xeb, 0x48, 0xca,
	0x97, 0x38, 0xe2, 0x6e, 0x0f, 0xba, 0x40, 0xb5, 0xa1, 0x46, 0x85, 0xfb, 0x29, 0x92, 0xc6, 0xca,
	0x80, 0x62, 0xff, 0x8c, 0x75, 0xf9, 0x0f, 0x61, 0xb1, 0x5b, 0x37, 0x18, 0x8e, 0x04, 0x0f, 0x43,
	0xde, 0xb3, 0x32, 0xd8, 0xcd, 0x4b, 0xda, 0x52, 0x9d, 0x18, 0xce, 0x2a, 0xa3, 0xac, 0xd1, 0xc0,
	0x85, 0x02, 0x0d, 0x83, 0x9e, 0xd7, 0x3f, 0x41, 0x31, 0x43, 0x8d, 0x24, 0xb5, 0x87, 0x8c, 0xac,
	0x48, 0x7c, 0x0e, 0x36, 0x7e, 0x4e, 0x55, 0x44, 0x73, 0xf2, 0xc2, 0xec, 0x7f, 0x64, 0xc9, 0x93,
	0x73, 0x0c, 0x3a, 0xb4, 0x13, 0x8c, 0xa3, 0x01, 0x0c, 0x90, 0xd0, 0xfa, 0x3b, 0xca, 0x3c, 0xaf,
	0xcb, 0xa3, 0xbf, 0x48, 0x0e, 0x8d, 0x1d, 0x9b, 0xda, 0x8f, 0x3d, 0x46, 0x18, 0xc4, 0x3c, 0x1a,
	0x7d, 0xf9, 0x0f, 0xa3, 0x97, 0x26, 0x9e, 0x93, 0x70, 0x50, 0x11, 0x4c, 0xe1, 0xc3, 0x0e, 0xf7,
	0xa1, 0xea, 0x9e, 0x7f, 0x68, 0x15, 0x30, 0xbf, 0x9c, 0x24, 0x6d, 0xcd, 0x49, 0x27, 0xc8, 0xf7,
	0xce, 0xa9, 0xd3, 0x55, 0xe5, 0x24, 0xc9, 0x89, 0xd3, 0x15, 0xd9, 0x2d, 0x21, 0x17, 0x77, 0x5c,
	0x39, 0x8d, 0xac, 0x95, 0xf8, 0xe5, 0x27, 0x14, 0x2b, 0x91, 0x2c, 0x7e, 0xac, 0x3b, 0x23, 0xf9,
	0x55, 0xb0, 0x4a, 0x38, 0xf8, 0xe2, 0x06, 0x38, 0xfd, 0x5c, 0x38, 0xa4, 0x7b, 0xfa, 0xe9, 0x58,
	0x87, 0xf2, 0x61, 0x8c, 0xc0, 0xaf, 0xa2, 0xf5, 0x7f, 0x0c, 0x62, 0xd3, 0x41, 0xea, 0x7b, 0xe9,
	0xa8, 0xbd, 0xf5, 0xc7, 0x73, 0x93, 0x2c, 0xc4, 0x3d, 0x87, 0xca, 0xb7, 0x76, 0xb0, 0x20, 0xc6,
	0x61, 0x68, 0x70, 0x4c, 0xbf, 0xc2, 0x2a, 0xc9, 0xc4, 0xaa, 0x90, 0xc5, 0x99, 0x03, 0x25, 0x84,
	0xa0, 0x0a, 0x19, 0xbc, 0x38, 0x1d, 0x0c, 0x76, 0x07, 0x8d, 0x85, 0x15, 0x92, 0xd3, 0x69, 0xa0,
	0x2e, 0xd4, 0xdc, 0xc5, 0xa5, 0xfe, 0xe4, 0xc7, 0xb0, 0x50, 0x24, 0x0b, 0xda, 0x33, 0x1c, 0x05,
	0x7e, 0xc8, 0x59, 0x9e, 0xcc, 0x06, 0xc7, 0xe8, 0x95, 0x76, 0x60, 0x25, 0xf5, 0x24, 0xb8, 0xfa,
	0x73, 0x91, 0x70, 0xe4, 0xf2, 0xda, 0x77, 0x29, 0x42, 0x4e, 0xa5, 0xc1, 0xb2, 0x64, 0xbe, 0xdd,
	0x78, 0xbc, 0xd1, 0xbc, 0xdd, 0xa0, 0xe7, 0xc0, 0x3b, 0xd7, 0x68, 0x56, 0x9b, 0x8d, 0x86, 0x5d,
	0x6d, 0xd5, 0x9a, 0x0d, 0x3a, 0x2b, 0xcd, 0xad, 0xda, 0x9e, 0xdd, 0x6c, 0xb7, 0xe8, 0x9c, 0x04,
	0x5b, 0x95, 0xed, 0xd6, 0x53, 0xfb, 0x36, 0x4d, 0x32, 0x42, 0x52, 0xad, 0x66, 0x73, 0xab, 0x76,
	0x8b, 0xa6, 0xd8, 0x02, 0xc9, 0x80, 0x61, 0xd7, 0xae, 0x6c, 0xdb, 0x0e, 0x9d, 0x87, 0x7c, 0x97,
	0x0f, 0x6a, 0x2d, 0xfb, 0x89, 0xb6, 0xdd, 0xb6, 0x77, 0xda, 0xf5, 0xfa, 0x8e, 0xdd, 0xaa, 0xee,
	0x82, 0x25, 0x03, 0xe9, 0x91, 0x5a, 0xe3, 0xc9, 0x4a, 0xbd, 0xb6, 0xdd, 0x76, 0xea, 0x94, 0x80,
	0x06, 0x98, 0xc6, 0x8e, 0xbd, 0x5d, 0x73, 0xe0, 0x58, 0xc9, 0x67, 0xe1, 0x93, 0x71, 0xd1, 0xd8,
	0x61, 0xbb, 0x76, 0xb0, 0x5f, 0x51, 0x9b, 0x2c, 0xc9, 0xd3, 0x0e, 0x5a, 0x95, 0x56, 0xfb, 0xa0,
	0xb8, 0xba, 0x4a, 0x7b, 0xd3, 0xb0, 0x48, 0x65, 0x05, 0x34, 0x2c, 0x81, 0xf5, 0x8b, 0x99, 0x69,
	0x5c, 0xa4, 0x0f, 0x0d, 0x5c, 0xa2, 0x5f, 0x1a, 0xb8, 0x4c, 0xbf, 0x32, 0xf0, 0x1a, 0xfd, 0xda,
	0xc0, 0xeb, 0xf4, 0x1b, 0x03, 0x6f, 0xd0, 0x6f, 0xa7, 0x70, 0x19, 0xce, 0xfb, 0x60, 0x76, 0x1a,
	0x17, 0xe9, 0x87, 0x06, 0x2e, 0xd1, 0x8f, 0x0c, 0x5c, 0xa6, 0x1f, 0x1b, 0x78, 0x8d, 0x7e, 0x62,
	0xe0, 0x75, 0xfa, 0xa9, 0x81, 0x37, 0xe8, 0x67, 0x06, 0xde, 0xa4, 0x9f, 0x4f, 0xe1, 0x35, 0x38,
	0xff, 0x85, 0xc4, 0x34, 0x2e, 0xd2, 0x17, 0x0d, 0x5c, 0xa2, 0x2f, 0x19, 0xb8, 0x4c, 0x5f, 0x36,
	0xf0, 0x1a, 0xbd, 0x6f, 0xe0, 0x75, 0xfa, 0x8a, 0x81, 0x37, 0xe8, 0xab, 0x06, 0xde, 0xa4, 0x0f,
	0x0c, 0x7c, 0x83, 0xbe, 0x66, 0xe0, 0x9b, 0xf4, 0xf5, 0x69, 0x5c, 0x5c, 0xa5, 0x6f, 0x18, 0xb8,
	0x48, 0xdf, 0x34, 0x70, 0x89, 0xbe, 0x65, 0xe0, 0x32, 0x7d, 0xdb, 0xc0, 0x6b, 0xf4, 0x1d, 0x03,
	0xaf, 0xd3, 0x77, 0x0d, 0xbc, 0x41, 0xdf, 0x33, 0xf0, 0x26, 0x7d, 0x7f, 0x0a, 0xaf, 0x43, 0x7d,
	0xbe, 0x37, 0x70, 0x91, 0xfe, 0x60, 0xe0, 0x12, 0xfd, 0xd1, 0xc0, 0x65, 0xfa, 0x93, 0x81, 0xd7,
	0xe8, 0xcf, 0x06, 0x5e, 0xa7, 0xbf, 0x18, 0xf8, 0x26, 0xfd, 0x75, 0x1a, 0xc3, 0x7d, 0x7f, 0x4b,
	0x5c, 0x7b, 0x8e, 0x64, 0xa7, 0xfe, 0x8b, 0xa0, 0xa0, 0xec, 0x83, 0x56, 0xad, 0x71, 0x0b, 0xf4,
	0xb6, 0x48, 0xb2, 0xb7, 0xed, 0xad, 0xce, 0x41, 0xe4, 0x8a, 0xa8, 0x3d, 0xa2, 0x33, 0x2c, 0x47,
	0xd2, 0x92, 0xd8, 0xab, 0xd4, 0xb4, 0xf8, 0x24, 0xda, 0x6d, 0x6f, 0xd1, 0x44, 0xec, 0x0b, 0xea,
	0x6c, 0xd9, 0x0d, 0xa9, 0x46, 0x46, 0xf2, 0x18, 0xdc, 0xde, 0x8a, 0xb9, 0xa4, 0xd4, 0x97, 0xe4,
	0xb6, 0xed, 0x56, 0xa5, 0x56, 0xa7, 0xa9, 0x6b, 0x05, 0x92, 0x8e, 0xff, 0x14, 0x4b, 0xc1, 0x36,
	0x9a, 0xce, 0x5e, 0xa5, 0x0e, 0x07, 0xc3, 0xba, 0xed, 0xdc, 0x92, 0x31, 0x33, 0xa5, 0x67, 0xf5,
	0xb4, 0x39, 0x80, 0xa9, 0xea, 0x75, 0x39, 0x5b, 0x25, 0x73, 0x3b, 0x1c, 0xbe, 0x4f, 0xf4, 0x91,
	0xe1, 0x14, 0x5e, 0x5a, 0x36, 0x67, 0xa4, 0x1a, 0x39, 0x85, 0x73, 0xec, 0x06, 0xc9, 0xd4, 0xc2,
	0x5d, 0xee, 0x0e, 0xa2, 0xa3, 0x13, 0xf6, 0x17, 0xd3, 0x09, 0x6f, 0x7e, 0x56, 0xe4, 0xd3, 0x29,
	0xa4, 0xcb, 0xbf, 0x07, 0x00, 0x00, 0xff, 0xff, 0x44, 0x0c, 0x58, 0x97, 0x73, 0x0d, 0x00, 0x00,
}
