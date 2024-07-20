package log

import (
	"path"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap/zapcore"
)

// 将日志文件拆分输出到文件系统中
// level 日志级别
// options 参数
// directory 日志文件存储的文件夹
func withFileRotateWriteSink(level zapcore.Level,
	encoder zapcore.Encoder,
	directory string,
	options ...rotatelogs.Option) ILoggerSink {
	return NewLoggerSink(func() zapcore.Core {
		fileWriter, err := rotatelogs.New(path.Join(directory, "%Y-%m-%d", "log.log"), options...)
		if err != nil {
			return nil
		}
		if fileWriter == nil {
			return nil
		}
		return zapcore.NewCore(encoder, zapcore.AddSync(fileWriter), level)
	})
}
