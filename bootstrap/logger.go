package bootstrap

import (
	"io"
	"os"
	"time"

	kratoszap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// options 日志配置选项
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

// 新建日志
func newLogger(opts ...Option) log.Logger {
	// 默认配置
	o := &options{
		writer:      os.Stdout,
		level:       log.LevelInfo,
		format:      "console",
		timeLayout:  time.DateTime,
		enableTrace: true,
		enableSpan:  false,
		setGlobal:   true,
	}
	for _, opt := range opts {
		opt(o)
	}

	var baseLogger log.Logger
	if o.format == "json" {
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(o.timeLayout)
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(o.writer),
			zap.DebugLevel,
		)
		baseLogger = kratoszap.NewLogger(zap.New(core))
	} else {
		baseLogger = log.NewStdLogger(o.writer)
	}

	filteredLogger := log.NewFilter(baseLogger, log.FilterLevel(o.level))

	kvs := []any{
		"ts", log.Timestamp(o.timeLayout),
		"caller", log.DefaultCaller,
	}

	if o.enableTrace {
		kvs = append(kvs,
			"trace_id", tracing.TraceID(),
		)
	}

	if o.enableSpan {
		kvs = append(kvs,
			"span_id", tracing.SpanID(),
		)
	}

	if o.enableTrace || o.enableSpan {
		// 启用链路追踪
		tp := trace.NewTracerProvider()
		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(propagation.TraceContext{})
	}

	logger := log.With(
		filteredLogger,
		kvs...,
	)

	if o.setGlobal {
		log.SetLogger(logger)
	}

	return logger
}
