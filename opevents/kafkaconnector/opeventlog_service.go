package kafkaconnector

import (
	opevent "github.com/shanluzhineng/fwpkg/opevents/pkg"

	"go.uber.org/zap/zapcore"
)

type opEventLogService struct {
	kafkaPusher *opEventLogKafkaPusher
}

var _ opevent.IOpEventLogService = (*opEventLogService)(nil)

// new一个实例
func newOpEventLogService() *opEventLogService {
	return &opEventLogService{
		kafkaPusher: newOpEventLogKafkaPusher(),
	}
}

// 插入一条记录
func (service *opEventLogService) Save(item *opevent.OpEventLog) error {
	return service.kafkaPusher.PushOpEvents(item)
}

// 保存一条调试级别的任务日志
func (service *opEventLogService) SaveDebugOpEventLog(message string, opts ...opevent.EventLogOption) (err error) {
	return service.SaveOpEventLog(message, zapcore.DebugLevel, opts...)
}

// 保存一条警告级别的任务日志
func (service *opEventLogService) SaveWarnOpEventLog(message string, opts ...opevent.EventLogOption) (err error) {
	return service.SaveOpEventLog(message, zapcore.WarnLevel, opts...)
}

// 保存一条错误级别的任务日志
func (service *opEventLogService) SaveErrorOpEventLog(message string, opts ...opevent.EventLogOption) (err error) {
	return service.SaveOpEventLog(message, zapcore.ErrorLevel, opts...)
}

// 保存一条事件日志
func (service *opEventLogService) SaveOpEventLog(message string, logLevel zapcore.Level, opts ...opevent.EventLogOption) (err error) {

	newTaskEvent := opevent.NewDefaultOpEventLog(logLevel, message)
	for _, eachOpt := range opts {
		//调用回调
		eachOpt(newTaskEvent)
	}
	return service.Save(newTaskEvent)
}
