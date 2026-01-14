package bootstrap

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	kratoszap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/lhlyu/kratos-easy/constants"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// dailyRotateWriter 按日期分割的日志写入器
type dailyRotateWriter struct {
	dir     string
	mu      sync.Mutex
	file    *os.File
	lastDay string
}

func (w *dailyRotateWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	day := time.Now().Format("2006-01-02")
	if day != w.lastDay {
		if w.file != nil {
			_ = w.file.Close()
		}
		if err := os.MkdirAll(w.dir, 0755); err != nil {
			return 0, err
		}
		filename := filepath.Join(w.dir, day+".log")
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return 0, err
		}
		w.file = f
		w.lastDay = day
	}
	if w.file == nil {
		return 0, errors.New("log file not open")
	}
	return w.file.Write(p)
}

func (w *dailyRotateWriter) Sync() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file != nil {
		return w.file.Sync()
	}
	return nil
}

// Close 关闭日志文件
func (w *dailyRotateWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file != nil {
		return w.file.Close()
	}
	return nil
}

// cleanOldLogs 清理过期日志
func cleanOldLogs(dir string, maxDays int) {
	if maxDays <= 0 {
		return
	}
	files, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	cutoff := time.Now().AddDate(0, 0, -maxDays)
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".log") {
			continue
		}
		// 期望格式 2006-01-02.log
		name := strings.TrimSuffix(f.Name(), ".log")
		t, err := time.Parse("2006-01-02", name)
		if err != nil {
			continue
		}
		if t.Before(cutoff) {
			_ = os.Remove(filepath.Join(dir, f.Name()))
		}
	}
}

// newLogger 新建日志，返回 logger 和清理函数
func newLogger() (log.Logger, func()) {
	var fileWriter *dailyRotateWriter
	var zLogger *zap.Logger
	var kzLogger *kratoszap.Logger

	// 基础日志器，默认输出到终端
	baseLogger := log.NewStdLogger(globalOption.writer)

	// 非本地环境且开启了文件日志，则同时写入文件(JSON)
	if !constants.IsLocal() && globalOption.enableFile {
		// 1. 程序运行前检查并删除过期日志
		cleanOldLogs(globalOption.logDir, globalOption.logMaxDays)

		// 2. 文件采用 JSON 格式，按日期分割
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(globalOption.timeLayout)
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

		fileWriter = &dailyRotateWriter{dir: globalOption.logDir}
		fileCore := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(fileWriter), zap.DebugLevel)
		zLogger = zap.New(fileCore)
		kzLogger = kratoszap.NewLogger(zLogger)

		baseLogger = MultiLogger(baseLogger, kzLogger)
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

	// 返回清理函数
	cleanup := func() {
		if kzLogger != nil {
			_ = kzLogger.Close()
		}
		if zLogger != nil {
			_ = zLogger.Sync()
		}
		if fileWriter != nil {
			_ = fileWriter.Close() // 关闭日志文件
		}
	}

	return logger, cleanup
}

type multiLogger struct {
	loggers []log.Logger
}

func (l *multiLogger) Log(level log.Level, keyvals ...any) error {
	for _, logger := range l.loggers {
		if err := logger.Log(level, keyvals...); err != nil {
			return err
		}
	}
	return nil
}

// MultiLogger 组合多个日志器
func MultiLogger(loggers ...log.Logger) log.Logger {
	return &multiLogger{loggers: loggers}
}
