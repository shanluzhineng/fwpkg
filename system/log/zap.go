package log

import (
	"fmt"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/shanluzhineng/fwpkg/utils/io"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogConfiguration struct {
	//日志级别
	Level string `mapstructure:"level" json:"level" yaml:"level"`
	// 日志前缀
	Prefix string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	// 格式,json,console
	Format string `mapstructure:"format" json:"format" yaml:"format"`
	// 输出在的文件夹
	Directory string `mapstructure:"directory" json:"directory"  yaml:"directory"`
	// 编码级
	EncodeLevel string `mapstructure:"encode-level" json:"encode-level" yaml:"encode-level"`
	// 栈名
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktrace-key" yaml:"stacktrace-key"`
	// 日志留存时间, 天
	MaxAge int64 `mapstructure:"max-age" json:"max-age" yaml:"max-age"`
	// 单个文件大小，单位为 bytes
	MaxSize int64 `mapstructure:"max-size" json:"max-size" yaml:"max-size"`
	// 是否显示行号
	ShowLine bool `mapstructure:"show-line" json:"show-line" yaml:"show-line"`
	//是否输出到控制台
	ToConsole bool `mapstructure:"to-console" json:"to-console" yaml:"to-console"`
}

type sinkWrapper struct {
	name  string
	sinks []ILoggerSink
}

// 日志配置
var (
	DefaultLogConfiguration *LogConfiguration = &LogConfiguration{}
	Logger                  *zap.Logger
	_registedSinks          map[string]sinkWrapper = make(map[string]sinkWrapper)
)

const (
	SinkName_Console = "console"
	SinkName_File    = "file"
)

func init() {
	//初始化默认的配置
	DefaultLogConfiguration.Directory = "log"
	DefaultLogConfiguration.Level = zap.InfoLevel.String()
	DefaultLogConfiguration.Format = "console"
	DefaultLogConfiguration.EncodeLevel = "LowercaseColorLevelEncoder"
	DefaultLogConfiguration.StacktraceKey = "stacktrace"
	DefaultLogConfiguration.MaxAge = 30
	DefaultLogConfiguration.MaxSize = 200 * 1024 * 1024 // 200m
	DefaultLogConfiguration.ToConsole = true
}

// 注册一个日志sink
// name: sink的名称
// level:sink处理的level
func RegistWriterSink(name string, sinks ...ILoggerSink) error {
	if len(name) <= 0 {
		return fmt.Errorf("name不能为空")
	}

	var sinkWrapper sinkWrapper = sinkWrapper{
		name: name,
	}
	sinkWrapper.sinks = append(sinkWrapper.sinks, sinks...)

	_registedSinks[name] = sinkWrapper
	return nil
}

// 构建默认的logger
func BuildDefaultLogger(actions ...func(logConfiguration *LogConfiguration)) {
	for _, eachAction := range actions {
		eachAction(DefaultLogConfiguration)
	}
	Logger = NewLog(DefaultLogConfiguration)
	zap.ReplaceGlobals(Logger)
}

func WithSinks(name string, sinks ...ILoggerSink) {
	RegistWriterSink(name, sinks...)
	BuildDefaultLogger()
}

// 使用固定字段
func WithField(fields ...zap.Field) {
	if Logger == nil {
		BuildDefaultLogger()
	}
	Logger = Logger.With(fields...)
}

// 创建一个日志
func NewLog(configuration *LogConfiguration, sinkOpts ...func()) *zap.Logger {
	if ok, _ := io.PathExists(configuration.Directory); !ok {
		// 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", configuration.Directory)
		_ = os.Mkdir(configuration.Directory, os.ModePerm)
	}
	//检测level配置是否合理
	level, err := zapcore.ParseLevel(configuration.Level)
	if err != nil {
		msg := fmt.Sprintf("无效的日志level参数值:%s,系统将使用默认的info值配置", configuration.Level)
		fmt.Println(msg)
		level = zap.DebugLevel
	}

	encoder := getEncoder(configuration)
	//注册默认的sink
	if configuration.ToConsole {
		RegistWriterSink(SinkName_Console, withConsoleSink(level, encoder))
	}
	RegistWriterSink(SinkName_File, withFileRotateWriteSink(level,
		encoder,
		configuration.Directory,
		rotatelogs.ForceNewFile(),
		rotatelogs.WithClock(rotatelogs.Local),
		rotatelogs.WithRotationTime(time.Hour*24),
		rotatelogs.WithRotationSize(configuration.MaxSize),
		rotatelogs.WithMaxAge(time.Duration(configuration.MaxAge)*24*time.Hour)))

	//用来注册额外的sink
	for _, eachSinkOpt := range sinkOpts {
		eachSinkOpt()
	}
	zapcoreList := createZapCoreList()
	logger := zap.New(zapcore.NewTee(zapcoreList...),
		zap.AddStacktrace(zapcore.ErrorLevel))
	if configuration.ShowLine {
		//显示行号
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}

func CreateEncoderConfig(configuration *LogConfiguration) (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  configuration.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     createTimeEncoder(configuration),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, //不输出go代码行数
	}
	switch {
	case configuration.EncodeLevel == "LowercaseLevelEncoder":
		// 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case configuration.EncodeLevel == "LowercaseColorLevelEncoder":
		// 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case configuration.EncodeLevel == "CapitalLevelEncoder":
		// 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case configuration.EncodeLevel == "CapitalColorLevelEncoder":
		// 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return config
}

func getEncoder(configuration *LogConfiguration) zapcore.Encoder {
	if configuration.Format == "json" {
		return zapcore.NewJSONEncoder(CreateEncoderConfig(configuration))
	}
	return zapcore.NewConsoleEncoder(CreateEncoderConfig(configuration))
}

// 自定义的时间编码
func createTimeEncoder(configuration *LogConfiguration) zapcore.TimeEncoder {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(configuration.Prefix + "2006/01/02-15:04:05.000"))
	}
}

// 耗时统计函数, defer 调用，defer定义时已固定定startTime
// 使用方式： defer TimeCost("RedisDelKey")()
func TimeCost(funcName string) func() {
	start := time.Now()
	return func() {
		tc := time.Since(start)
		Logger.Info(fmt.Sprintf("[%s]>>> time cost = %v", funcName, tc))
	}
}
