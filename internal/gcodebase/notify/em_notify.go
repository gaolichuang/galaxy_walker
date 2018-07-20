package notify

import (
	"fmt"
	"crypto/md5"
	"path/filepath"
	"gcodebase/file"
	LOG "gcodebase/log"
	"strings"
)

const (
	kEMNotifyCritical            = 16
	kEMNotifyError               = 4
	kEMNotifyWarn                = 2
	kWarnValueLengthLimiteInByte = 64
)

var EMLevelToString = map[int]string{
	kEMNotifyCritical: "EMCritical",
	kEMNotifyError:    "EMError",
	kEMNotifyWarn:     "EMWarn",
}

type EMNotification struct {
	writer     *EMFileWriter
	address    string
	WarnStatus int
}

func NewEMNotification(level int) *EMNotification {
	return &EMNotification{
		WarnStatus: level,
	}
}
func (em *EMNotification) Init(param *NotifyParam) error {
	if em == nil {
		return fmt.Errorf("Init Obj is nil")
	}
	if em.WarnStatus <= 0 {
		return fmt.Errorf("EmNotifyInit %s,Not Init WarnStatus", EMLevelToString[em.WarnStatus])
	}
	if param.EMFileDir == "" || param.EMFileSuffix == "" ||
		param.WarnName == "" || param.WarnPlatform == "" {
		return fmt.Errorf("EmNotifyInit %s,EMFileDir or EMFileSuffix or WarnName or WarnPlatform is Empty", EMLevelToString[em.WarnStatus])
	}
	param.EMFileDir = file.GetConfFile(param.EMFileDir)
	if !file.IsDir(param.EMFileDir) {
		return fmt.Errorf("EmNotifyInit %s, No Such Directory %s", EMLevelToString[em.WarnStatus], param.EMFileDir)
	}
	em.address = param.Address
	em.writer = &EMFileWriter{
		fileDir:       param.EMFileDir,
		fileSuffix:    param.EMFileSuffix,
		warn_code:     param.WarnCode,
		warn_name:     param.WarnName,
		Warn_status:   em.WarnStatus,
		warn_platform: param.WarnPlatform,
	}
	LOG.Infof("Notify %s Init Success use %s", EMLevelToString[em.WarnStatus], em.writer.String())
	return nil
}

func (em *EMNotification) Notify(source, msg string) error {
	source = fmt.Sprintf("[%s]%s", em.address, source)
	return em.writer.Write(source, msg)
}

type EMFileWriter struct {
	Warn_status   int
	fileDir       string
	fileSuffix    string
	warn_code     int
	warn_name     string
	warn_platform string
}

func (e *EMFileWriter) String() string {
	if e == nil {
		return "nil"
	}
	return fmt.Sprintf("<%s,%s;%d-%d>", e.fileDir, e.fileSuffix, e.warn_code, e.Warn_status)
}
func (e *EMFileWriter) Write(source, msg string) error {
	source = strings.Replace(source, " ", "#", -1)
	msg = strings.Replace(msg, " ", "#", -1)
	if len(source) > kWarnValueLengthLimiteInByte {
		LOG.Errorf("Length of warn_value is bigger than %d", kWarnValueLengthLimiteInByte)
		source = source[:kWarnValueLengthLimiteInByte]
	}
	EMFileCont := fmt.Sprintf("%d %s %s %d %s %s", e.warn_code, e.warn_name,
		msg, e.Warn_status, e.warn_platform, source)
	strmd5 := fmt.Sprintf("%x", md5.Sum([]byte(EMFileCont)))
	EMFileName := strmd5 + e.fileSuffix
	EMFileDir := file.GetConfFile(e.fileDir)
	EMfilePath := filepath.Join(EMFileDir, EMFileName)
	err := file.WriteStringToFile([]byte(EMFileCont), EMfilePath)
	if err != nil {
		LOG.Error("Write %s Into EM_File Err:%v", EMFileCont, err)
		return err
	}
	LOG.VLog(3).DebugTag("EMNotify", "%s Write:%s to %s", EMLevelToString[e.Warn_status], EMFileCont, EMfilePath)
	return nil
}
