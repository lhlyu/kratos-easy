package bootstrap

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// options 配置选项，用于初始化日志与配置加载
type options struct {
	writer      io.Writer // 日志输出目标，默认 os.Stdout
	level       log.Level // 日志等级过滤，低于该等级的日志不会输出，默认 LevelInfo
	format      string    // 日志输出格式，可选 "console" 或 "json"，默认 "console"
	timeLayout  string    // 日志时间格式，默认 time.DateTime
	enableTrace bool      // 是否输出链路追踪 trace_id，默认开启
	enableSpan  bool      // 是否输出链路追踪 span_id，默认关闭
	setGlobal   bool      // 是否覆盖全局默认日志器，默认 true
	configDir   string    // 配置文件目录路径，默认 "configs"
	logDir      string    // 日志目录，默认 "logs"
	logMaxDays  int       // 日志最大保留天数，默认 7 天
	enableFile  bool      // 是否写入文件，默认 true
}

var globalOption = &options{
	writer:      os.Stdout,
	level:       log.LevelInfo,
	format:      "console",
	timeLayout:  time.DateTime,
	enableTrace: true,
	enableSpan:  false,
	setGlobal:   true,
	configDir:   "configs",
	logDir:      "logs",
	logMaxDays:  7,
	enableFile:  true,
}

type Option func(*options)

// WithWriter 修改日志写入器
func WithWriter(writer io.Writer) Option {
	return func(o *options) {
		o.writer = writer
	}
}

// WithFilterLevel 设置日志过滤等级
func WithFilterLevel(level log.Level) Option {
	return func(o *options) {
		o.level = level
	}
}

// WithJsonFormat 日志输出格式改成 json
func WithJsonFormat() Option {
	return func(o *options) {
		o.format = "json"
	}
}

// WithConsoleFormat 日志输出格式改成 console
func WithConsoleFormat() Option {
	return func(o *options) {
		o.format = "console"
	}
}

// WithTimeLayout 设置时间格式
func WithTimeLayout(layout string) Option {
	return func(o *options) {
		o.timeLayout = layout
	}
}

// WithDisableTrace 禁用链路追踪 trace_id
func WithDisableTrace() Option {
	return func(o *options) {
		o.enableTrace = false
	}
}

// WithEnableSpan 启用链路追踪 span_id
func WithEnableSpan() Option {
	return func(o *options) {
		o.enableSpan = true
	}
}

// WithDisableGlobal 禁止覆盖全局的默认日志器
func WithDisableGlobal() Option {
	return func(o *options) {
		o.setGlobal = false
	}
}

// WithConfigDir 设置配置文件的目录路径
func WithConfigDir(elem ...string) Option {
	return func(o *options) {
		o.configDir = filepath.Join(elem...)
	}
}

// WithLogDir 设置日志目录
func WithLogDir(dir string) Option {
	return func(o *options) {
		o.logDir = dir
	}
}

// WithLogMaxDays 设置日志最大保留天数
func WithLogMaxDays(days int) Option {
	return func(o *options) {
		o.logMaxDays = days
	}
}

// WithDisableFileLog 禁用文件日志输出
func WithDisableFileLog() Option {
	return func(o *options) {
		o.enableFile = false
	}
}
