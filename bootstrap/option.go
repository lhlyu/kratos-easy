package bootstrap

import (
	"io"

	"github.com/go-kratos/kratos/v2/log"
)

// options 配置选项
type options struct {
	writer      io.Writer
	level       log.Level
	format      string
	timeLayout  string
	enableTrace bool
	enableSpan  bool
	setGlobal   bool
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
