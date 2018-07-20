package notify

import (
	LOG "gcodebase/log"
	"gcodebase/file"
	"path/filepath"
	"fmt"
)

const (
	kFileNotifyName = "notify.file.log"
)

type FileNotification struct {
	logger LOG.CustomLogger
}

func (f *FileNotification) Init(param *NotifyParam) error {
	if f == nil {
		return fmt.Errorf("Init Obj is nil")
	}
	if param.FileNotifyPath == "" {
		return fmt.Errorf("File Notify Path is Empty")
	}
	fpath := file.GetConfFile(param.FileNotifyPath)
	if !file.IsDir(fpath) {
		LOG.Errorf("%s NOT Dir", fpath)
		return fmt.Errorf("%s NOT EXIST", fpath)
	}
	var err error
	fname := filepath.Join(fpath, kFileNotifyName)
	f.logger, err = LOG.NewCustomLogger(fname)
	if err != nil {
		LOG.Errorf("Create Notify File %s Err:%v", fname, err)
		return err
	}
	LOG.Infof("FileNotification Init Success use %s", fname)
	return nil
}

func (f *FileNotification) Notify(source, msg string) error {
	if f.logger != nil {
		f.logger.LogTag(source, msg)
	}
	return nil
}
