package bootstrap

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/lhlyu/kratos-easy/constants"
)

// dailyRotateWriter 按日期分割的日志写入器
type dailyRotateWriter struct {
	dir     string
	maxDays int
	mu      sync.Mutex
	file    *os.File
	lastDay string
	closed  bool
}

func (w *dailyRotateWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.closed {
		return 0, nil // 闭合后静默丢弃，避免影响 shutdown 流程
	}

	day := time.Now().Format("2006-01-02")

	// 1. 日期切换或文件未打开
	if day != w.lastDay || w.file == nil {
		isNewDay := day != w.lastDay
		if err = w.rotate(day); err != nil {
			return 0, err
		}
		if isNewDay {
			go cleanOldLogs(w.dir, w.maxDays)
		}
	}

	// 2. 尝试写入
	n, err = w.file.Write(p)
	if err != nil {
		// 3. 写入失败尝试重试一次（可能是文件被删或其它 IO 异常）
		_, _ = os.Stderr.WriteString("write log failed, trying to reopen: " + err.Error() + "\n")
		if err = w.reopen(); err != nil {
			return 0, err
		}
		return w.file.Write(p)
	}

	return n, nil
}

func (w *dailyRotateWriter) rotate(day string) error {
	if w.file != nil {
		_, _ = fmt.Fprintf(w.file, "\n%s [ROTATE] Log rotated from %s to %s\n", time.Now().Format("2006-01-02 15:04:05"), w.lastDay, day)
		_ = w.file.Close()
		w.file = nil
	}

	return w.open(day)
}

func (w *dailyRotateWriter) open(day string) error {
	if err := os.MkdirAll(w.dir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory %s: %w", w.dir, err)
	}
	filename := filepath.Join(w.dir, day+".log")
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file %s: %w", filename, err)
	}
	w.file = f
	w.lastDay = day
	return nil
}

func (w *dailyRotateWriter) reopen() error {
	day := time.Now().Format("2006-01-02")
	if w.file != nil {
		_ = w.file.Close()
		w.file = nil
	}
	return w.open(day)
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
	if w.closed {
		return nil
	}
	w.closed = true
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

	// 基础日志器，默认输出到终端
	baseLogger := log.NewStdLogger(globalOption.writer)

	// 非本地环境且开启了文件日志，则同时写入文件(JSON)
	if !constants.IsLocal() && globalOption.enableFile {
		// 程序运行前检查并删除过期日志
		cleanOldLogs(globalOption.logDir, globalOption.logMaxDays)

		fileWriter = &dailyRotateWriter{
			dir:     globalOption.logDir,
			maxDays: globalOption.logMaxDays,
		}

		kzLogger := log.NewStdLogger(fileWriter)

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
		if fileWriter != nil {
			_ = fileWriter.Sync()
			_ = fileWriter.Close()
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
