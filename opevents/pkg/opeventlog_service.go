package pkg

import (
	"fmt"

	"github.com/shanluzhineng/fwpkg/system/log"
	jsonUtil "github.com/shanluzhineng/fwpkg/utils/json"
	"go.uber.org/zap/zapcore"
)

type IOpEventLogStoreService interface {
	Save(item *OpEventLog) error
}

type IOpEventLogService interface {
	// GetPageList(input OpEventLogSearch) (list []OpEventLog, total int64, err error)

	IOpEventLogStoreService
	SaveDebugOpEventLog(message string, opts ...EventLogOption) (err error)
	SaveWarnOpEventLog(message string, opts ...EventLogOption) (err error)
	SaveErrorOpEventLog(message string, opts ...EventLogOption) (err error)
	SaveOpEventLog(message string, logLevel zapcore.Level, opts ...EventLogOption) (err error)
}

type defaultOpEventLogService struct {
}

var _ IOpEventLogService = (*defaultOpEventLogService)(nil)

func newDefaultOpEventLogService() IOpEventLogService {
	return &defaultOpEventLogService{}
}

func (s *defaultOpEventLogService) Save(item *OpEventLog) error {
	msg := fmt.Sprintf("保存一条opevent数据:%s", jsonUtil.ObjectToJson(item))
	if item.LogLevel == zapcore.DebugLevel {
		log.Logger.Debug(msg)
	} else if item.LogLevel == zapcore.InfoLevel {
		log.Logger.Info(msg)
	} else if item.LogLevel == zapcore.WarnLevel {
		log.Logger.Warn(msg)
	} else if item.LogLevel >= zapcore.ErrorLevel {
		log.Logger.Error(msg)
	}
	return nil
}

func (s *defaultOpEventLogService) SaveDebugOpEventLog(message string, opts ...EventLogOption) (err error) {
	return s.SaveOpEventLog(message, zapcore.DebugLevel, opts...)
}

func (s *defaultOpEventLogService) SaveWarnOpEventLog(message string, opts ...EventLogOption) (err error) {
	return s.SaveOpEventLog(message, zapcore.WarnLevel, opts...)
}

func (s *defaultOpEventLogService) SaveErrorOpEventLog(message string, opts ...EventLogOption) (err error) {
	return s.SaveOpEventLog(message, zapcore.ErrorLevel, opts...)
}

func (s *defaultOpEventLogService) SaveOpEventLog(message string, logLevel zapcore.Level, opts ...EventLogOption) (err error) {
	newTaskEvent := NewDefaultOpEventLog(logLevel, message)
	for _, eachCallback := range opts {
		//调用回调
		eachCallback(newTaskEvent)
	}
	return s.Save(newTaskEvent)
}
