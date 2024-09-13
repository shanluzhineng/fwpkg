package pkg

import (
	"encoding/json"
	"time"

	"github.com/shanluzhineng/fwpkg/system/lang"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap/zapcore"
)

// 执行事件
type OpEventLog struct {
	Id       uuid.UUID `json:"id"`
	TenantId string    `json:"tenantId"`
	//报告这个事件日志的时间戳
	ReportTimestamp int64     `json:"reportTimestamp"`
	ReportTime      time.Time `json:"reportTime"`
	CreationTime    time.Time `json:"creationTime"`
	//创建时间时间戳
	OnTimestamp int64 `json:"onTimestamp"`
	//创建人
	CreatorId string        `json:"creatorId,omitempty"`
	AccountId int64         `json:"accountId,omitempty"`
	LogLevel  zapcore.Level `json:"logLevel"`
	//上报此日志的ip地址
	IPAddress      string         `json:"ipAddress,omitempty"`
	AppName        string         `json:"appName"`
	AndroidId      string         `json:"androidId"`
	DeviceMobileNo string         `json:"deviceMobileNo"`
	EventMessage   string         `json:"eventMessage,omitempty"`
	Source         EventLogSource `json:"source,omitempty"`
	//操作行为名称，如add,delete
	OPAction string `json:"opAction"`
	//全链路追踪时使用
	CorrelationId string `json:"correlationId,omitempty"`

	//任务id
	// TaskId   string      `json:"taskId"`
	TaskInfo interface{} `json:"taskInfo,omitempty"`
	// //后台管理任务id
	// SceneTaskId   string      `json:"sceneTaskId"`
	SceneTaskInfo interface{} `json:"sceneTaskInfo,omitempty"`
}

// 创建时设置对象的基本信息
func (entity *OpEventLog) BeforeCreate() (err error) {
	if entity.Id == uuid.Nil {
		//创建一个新的id
		entity.Id = uuid.NewV4()
	}
	if entity.OnTimestamp <= 0 {
		entity.CreationTime = time.Now()
		entity.OnTimestamp = entity.CreationTime.UnixNano()
	}
	return
}

func (entity *OpEventLog) WithAccountId(accountId *int64) *OpEventLog {
	entity.AccountId = lang.IfValue(accountId != nil, func() int64 {
		return *accountId
	}, 0)
	return entity
}

// 设置事件的设备信息
func (entity *OpEventLog) WithDeviceNo(deviceMobileNo string) *OpEventLog {
	entity.DeviceMobileNo = deviceMobileNo
	return entity
}

// 设置日志级别
func (entity *OpEventLog) WithLogLvel(logLevel zapcore.Level) *OpEventLog {
	entity.LogLevel = logLevel
	return entity
}

func (log *OpEventLog) WithAppSource() *OpEventLog {
	log.Source = EventLogSource_App
	return log
}

func (log *OpEventLog) WithApiServerSource() *OpEventLog {
	log.Source = EventLogSource_ApiServer
	return log
}

// json序列化
func (e *OpEventLog) Bytes() []byte {
	b, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return b
}

// 使用默认参数创建一条debug记录的事件日志
func NewDebugOpEventLog(message string, opts ...EventLogOption) *OpEventLog {
	return NewDefaultOpEventLog(zapcore.DebugLevel, message, opts...)
}

// 使用默认参数创建一条warn记录的事件日志
func NewWarnOpEventLog(message string, opts ...EventLogOption) *OpEventLog {
	return NewDefaultOpEventLog(zapcore.WarnLevel, message, opts...)
}

// 使用默认参数创建一条error级别的事件日志
func NewErrorOpEventLog(message string, opts ...EventLogOption) *OpEventLog {
	return NewDefaultOpEventLog(zapcore.ErrorLevel, message, opts...)
}

// 使用默认参数创建一条事件日志记录
func NewDefaultOpEventLog(logLevel zapcore.Level, message string, opts ...EventLogOption) *OpEventLog {
	return NewOpEventLog(func(taskEvent *OpEventLog) {
		taskEvent.LogLevel = logLevel
		taskEvent.EventMessage = message
		taskEvent.Source = EventLogSource_App
	})
}

func NewOpEventLog(opts ...EventLogOption) *OpEventLog {
	var taskEvent = &OpEventLog{}
	taskEvent.Id = uuid.NewV4()
	now := time.Now()
	taskEvent.ReportTime = now
	taskEvent.ReportTimestamp = now.UnixNano()
	taskEvent.CreationTime = now
	taskEvent.OnTimestamp = now.UnixNano()
	if len(opts) > 0 {
		for _, eachOpt := range opts {
			eachOpt(taskEvent)
		}
	}
	return taskEvent
}
