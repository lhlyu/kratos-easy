package bootstrap

import (
	kratoszap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 新建日志
func newLogger(o *options) log.Logger {

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

	logger := log.With(
		filteredLogger,
		kvs...,
	)

	if o.setGlobal {
		log.SetLogger(logger)
	}

	return logger
}
