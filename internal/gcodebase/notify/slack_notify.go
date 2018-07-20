package notify

import (
	SLACK "galaxy_walker/internal/gcodebase/slack"
	"galaxy_walker/internal/github.com/nlopes/slack"
	LOG "gcodebase/log"
	"fmt"
	"time"
	"strings"
)

const (
	kSlackMsgChanSize    = 80
	kSlackMsgChanCap     = 100
	KSlackSYNCCALLPREFIX = "SYNC-"
)

type SlackNotification struct {
	// topic is no need, u can use channel name to category
	topic    string
	reminder *SLACK.SlackReminder

	localip string
	msgChan chan *SlackObj
}
type SlackObj struct {
	pretext string
	msg     string
}

func (s *SlackNotification) Init(param *NotifyParam) error {
	if s == nil {
		return fmt.Errorf("Init Obj is nil")
	}
	if param.SlackToken == "" || param.SlackChannelName == "" || param.SlackTopic == "" {
		return fmt.Errorf("SlackToken or SlackChannelName or SlackTopic is Empty")
	}
	s.topic = param.SlackTopic
	s.reminder = SLACK.GetSlackReminder(param.SlackToken, param.SlackChannelName)
	s.msgChan = make(chan *SlackObj, kSlackMsgChanCap)
	s.localip = param.Address
	go s.notifyLoop()
	LOG.Infof("SlackNotification Init Success use %s", param.String())
	return nil
}
func (s *SlackNotification) Notify(source, msg string) error {
	if s == nil || s.reminder == nil {
		LOG.Errorf("Send Slack Err,client is nil")
		return nil
	}
	if strings.HasPrefix(source, KSlackSYNCCALLPREFIX) {
		s.notifyInterval(source, msg)
		return nil
	}
	if len(s.msgChan) > kSlackMsgChanSize {
		err := fmt.Errorf("SlackNotify, MsgChan buf full,%d/%d", len(s.msgChan), kSlackMsgChanSize)
		LOG.Error(err)
		return err
	}
	s.msgChan <- &SlackObj{
		pretext: source,
		msg:     msg,
	}
	return nil
}

func (s *SlackNotification) String() string {
	if s == nil {
		return "nil"
	}
	return fmt.Sprintf("Chan:%d/%d", len(s.msgChan), kSlackMsgChanSize)
}
func (s *SlackNotification) notifyLoop() {
	for {
		select {
		case msg := <-s.msgChan:
			s.notifyInterval(msg.pretext, msg.msg)
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}
func (s *SlackNotification) notifyInterval(source, msg string) {
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		//Pretext: source,
		Text: msg,
	}
	params.Attachments = []slack.Attachment{attachment}
	topic := fmt.Sprintf("[%s,%s]", s.localip, source)
	channelID, timestamp, err := s.reminder.Client.PostMessage(s.reminder.SlackChannelName, topic, params)
	if err != nil {
		errmsg := fmt.Errorf("Slack:%v", err)
		LOG.Error(errmsg)
		return
	}
	LOG.VLog(5).DebugTag("SlackNotify", "Source:%s,Msg:%s sent to channel %s at %s", source, msg, channelID, timestamp)
}
