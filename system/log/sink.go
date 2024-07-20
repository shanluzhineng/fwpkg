package log

import (
	"os"

	"go.uber.org/zap/zapcore"
)

// 日志写入sink
type ILoggerSink interface {
	CreateZapCore() zapcore.Core
}

// 创建一个新的ILoggerSink实例
func NewLoggerSink(createFunc func() zapcore.Core) ILoggerSink {
	return createZapCore(createFunc)
}

type createZapCore func() zapcore.Core

func (fn createZapCore) CreateZapCore() zapcore.Core {
	return fn()
}

// 输出到控制台的IWriterSink实现
func withConsoleSink(logLevel zapcore.Level, encoder zapcore.Encoder) ILoggerSink {
	return NewLoggerSink(func() zapcore.Core {
		return zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), logLevel)
	})
}

// 创建log的输出
func createZapCoreList() []zapcore.Core {
	if len(_registedSinks) <= 0 {
		return nil
	}
	var zapCoreList []zapcore.Core = make([]zapcore.Core, 0)
	for _, eachSinkWrapper := range _registedSinks {
		if len(eachSinkWrapper.sinks) <= 0 {
			continue
		}
		for _, eachSink := range eachSinkWrapper.sinks {
			currentZapCore := eachSink.CreateZapCore()
			if currentZapCore == nil {
				// 无法创建，则继续
				continue
			}
			zapCoreList = append(zapCoreList, currentZapCore)
		}
	}
	return zapCoreList
}
