package conf

import (
	"flag"
	"galaxy_walker/internal/github.com/vharitonsky/iniflags"
	"os"
	"sync"
)

type ConfType struct {
	LogFile      *string
	ErrorLogFile *string
	LogV         *int
	Stdout       *bool
	ReBorn       *bool

	UseTLS             *bool
	CertFile           *string // server
	KeyFile            *string // server
	CaFile             *string // client
	ServerHostOverride *string // client
	ConnectionTimeout  *int

	ConfPathPrefix     *string
	ReloadConfInterval *int
	FilePersistentPath *string

	ConnectionMode *string
}

var _conf = ConfType{
	LogFile:            flag.String("log_file", "", "log to file"),
	ErrorLogFile:       flag.String("error_log_file", "", "log to file"),
	LogV:               flag.Int("v", 3, "log level for debug"),
	Stdout:             flag.Bool("stdout", true, "output stdout or not"),
	ReBorn:             flag.Bool("reborn", false, "application reborn or not"),
	UseTLS:             flag.Bool("use_tls", false, "use tls or not"),
	CertFile:           flag.String("cert_file", "", "TLS cert file"),
	KeyFile:            flag.String("key_file", "", "TLS key file"),
	CaFile:             flag.String("ca_file", "", "The file containning the CA root cert file"),
	ServerHostOverride: flag.String("server_host_override", "x.a.com", "The server name use to verify the hostname returned by TLS handshake"),
	ConnectionTimeout:  flag.Int("connection_timeout", 10, "rpc connection second."),

	ConfPathPrefix:     flag.String("conf_path_prefix", "/Application/mustard", "conf common prefix"),
	ReloadConfInterval: flag.Int("reload_conf_interval", 60, "reload conf interval"),
	FilePersistentPath: flag.String("persistent_file_path", "persistent", "persistent file path, relative by conf_path"),
	ConnectionMode:     flag.String("conn_mode", "SRTT", "rr or srtt"),
}

var Conf *ConfType

func init() {
	Conf = NewConfType()
	// TODO. Set Default Config File
	DefaultConfigFile := "config.ini"
	_, err := os.Stat(DefaultConfigFile)
	if !os.IsNotExist(err) {
		iniflags.SetConfigFile(DefaultConfigFile)
	}
	iniflags.Parse()
}

var _conf_instance *ConfType = nil
var _conf_init_ctx sync.Once

func NewConfType() *ConfType {
	_conf_init_ctx.Do(func() {
		_conf_instance = &_conf
	})
	return _conf_instance
}
