package notify

import (
	"fmt"
	"sync"
	LOG "gcodebase/log"
	"strings"
)

const (
	KNotifyLevelCritical = 1
	KNotifyLevelError    = 2
	KNotifyLevelWarn     = 3
	KNotifyLevelInfo     = 4
)

type Notification interface {
	Notify(source, msg string) error
	Init(param *NotifyParam) error
}

type NotifyParam struct {
	FileNotifyPath string

	SlackTopic       string
	SlackToken       string
	SlackChannelName string

	//for em
	EMFileDir    string
	EMFileSuffix string
	WarnCode     int
	WarnName     string
	WarnPlatform string

	Address string
}

func (p *NotifyParam) String() string {
	if p == nil {
		return "nil"
	}
	return fmt.Sprintf("File:%s;Slack:%s,%s,%s;EM:%s,%d",
		p.FileNotifyPath, p.SlackToken, p.SlackTopic, p.SlackChannelName, p.EMFileDir, p.WarnCode)

}

type NotifyLevelManager struct {
	notifyCritical []Notification
	notifyError    []Notification
	notifyWarn     []Notification
	notifyInfo     []Notification

	//flag
	fileNotify  bool
	slackNotify bool
	eMNotify    bool

	slackNotifyObj *SlackNotification
	fileNotifyObj  *FileNotification
}

func (nm *NotifyLevelManager) String() string {
	if nm == nil {
		return "nil"
	}
	return nm.slackNotifyObj.String()
}
func (nm *NotifyLevelManager) RegisterNotifyCritical(nt Notification) {
	nm.notifyCritical = append(nm.notifyCritical, nt)
}

func (nm *NotifyLevelManager) RegisterNotifyError(nt Notification) {
	nm.notifyError = append(nm.notifyError, nt)
}

func (nm *NotifyLevelManager) RegisterNotifyWarn(nt Notification) {
	nm.notifyWarn = append(nm.notifyWarn, nt)
}

func (nm *NotifyLevelManager) RegisterNotifyInfo(nt Notification) {
	nm.notifyInfo = append(nm.notifyInfo, nt)
}

func (nm *NotifyLevelManager) RegisterNotification() {
	nm.slackNotifyObj = &SlackNotification{}
	nm.fileNotifyObj = &FileNotification{}
	nm.RegisterNotifyCritical(NewEMNotification(kEMNotifyCritical))
	nm.RegisterNotifyCritical(nm.slackNotifyObj)
	nm.RegisterNotifyCritical(nm.fileNotifyObj)

	nm.RegisterNotifyError(NewEMNotification(kEMNotifyError))
	nm.RegisterNotifyError(nm.slackNotifyObj)
	nm.RegisterNotifyError(nm.fileNotifyObj)

	nm.RegisterNotifyWarn(NewEMNotification(kEMNotifyWarn))
	nm.RegisterNotifyWarn(nm.slackNotifyObj)
	nm.RegisterNotifyWarn(nm.fileNotifyObj)

	nm.RegisterNotifyInfo(nm.slackNotifyObj)
	nm.RegisterNotifyInfo(nm.fileNotifyObj)
}

func (nm *NotifyLevelManager) initInternal(nts []Notification, param *NotifyParam) {
	for _, nt := range nts {
		switch nt.(type) {
		case *EMNotification:
			if err := nt.Init(param); err != nil {
				LOG.Errorf("EMNotification Init Err:%v", err)
				nm.eMNotify = false
			}
		case *SlackNotification:
			if err := nt.Init(param); err != nil {
				LOG.Errorf("SlackNotification Init Err:%v", err)
				nm.slackNotify = false
			}
		case *FileNotification:
			if err := nt.Init(param); err != nil {
				LOG.Errorf("FileNotification Init Err:%v", err)
				nm.fileNotify = false
			}
		}
	}
}

func (nm *NotifyLevelManager) Init(param *NotifyParam) {
	nm.eMNotify = true
	nm.slackNotify = true
	nm.fileNotify = true
	nm.initInternal(nm.notifyCritical, param)
	nm.initInternal(nm.notifyError, param)
	nm.initInternal(nm.notifyWarn, param)
	nm.initInternal(nm.notifyInfo, param)
}

func (nm *NotifyLevelManager) notifyInternal(nts []Notification, source, msg string) error {
	retErr := make([]string, 0)
	for _, nt := range nts {
		switch nt.(type) {
		case *EMNotification:
			if nm.eMNotify {
				if err := nt.Notify(source, msg); err != nil {
					retErr = append(retErr, err.Error())
				}
			}
		case *SlackNotification:
			if nm.slackNotify {
				if err := nt.Notify(source, msg); err != nil {
					retErr = append(retErr, err.Error())
				}
			}
		case *FileNotification:
			if nm.fileNotify {
				if err := nt.Notify(source, msg); err != nil {
					retErr = append(retErr, err.Error())
				}
			}
		}
	}
	if len(retErr) <= 0 {
		return nil
	}
	return fmt.Errorf(strings.Join(retErr, ";"))
}

func (nm *NotifyLevelManager) Notify(source, msg string, level int) error {
	var err error
	switch level {
	case KNotifyLevelCritical:
		err = nm.notifyInternal(nm.notifyCritical, source, msg)
	case KNotifyLevelError:
		err = nm.notifyInternal(nm.notifyError, source, msg)
	case KNotifyLevelWarn:
		err = nm.notifyInternal(nm.notifyWarn, source, msg)
	case KNotifyLevelInfo:
		err = nm.notifyInternal(nm.notifyInfo, source, msg)
	default:
		err = fmt.Errorf("Notify Level %d is InValid, please input %d or %d or %d or %d", level,
			KNotifyLevelCritical, KNotifyLevelError, KNotifyLevelWarn, KNotifyLevelInfo)
	}
	if err != nil {
		LOG.VLog(3).DebugTag("NotifyManager", "Notify Err %v", err)
	}
	return err
}

func NewNotifyLevelManager(param *NotifyParam) *NotifyLevelManager {
	_notify_manager := &NotifyLevelManager{
		notifyCritical: make([]Notification, 0),
		notifyError:    make([]Notification, 0),
		notifyWarn:     make([]Notification, 0),
		notifyInfo:     make([]Notification, 0),
	}
	_notify_manager.RegisterNotification()
	_notify_manager.Init(param)
	return _notify_manager
}

///////////// Singleton  ////////////////////////
var _nm_instance *NotifyLevelManager = nil
var _nm_init_ctx sync.Once

func GetNotifyLevelManage(param *NotifyParam) *NotifyLevelManager {
	_nm_init_ctx.Do(func() {
		_nm_instance = &NotifyLevelManager{
			notifyCritical: make([]Notification, 0),
			notifyError:    make([]Notification, 0),
			notifyWarn:     make([]Notification, 0),
			notifyInfo:     make([]Notification, 0),
		}
		_nm_instance.RegisterNotification()
		_nm_instance.Init(param)
	})
	return _nm_instance
}
