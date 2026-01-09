package bootstrap

import (
	kratoszap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 新建日志
func newLogger() log.Logger {

	var baseLogger log.Logger
	if globalOption.format == "json" {
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(globalOption.timeLayout)
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(globalOption.writer),
			zap.DebugLevel,
		)
		baseLogger = kratoszap.NewLogger(zap.New(core))
	} else {
		baseLogger = log.NewStdLogger(globalOption.writer)
	}

	filteredLogger := log.NewFilter(baseLogger, log.FilterLevel(globalOption.level))

	kvs := []any{
		"ts", log.Timestamp(globalOption.timeLayout),
		"caller", log.DefaultCaller,
	}

	if globalOption.enableTrace {
		kvs = append(kvs,
			"trace_id", tracing.TraceID(),
		)
	}

	if globalOption.enableSpan {
		kvs = append(kvs,
			"span_id", tracing.SpanID(),
		)
	}

	logger := log.With(
		filteredLogger,
		kvs...,
	)

	if globalOption.setGlobal {
		log.SetLogger(logger)
	}

	return logger
}
