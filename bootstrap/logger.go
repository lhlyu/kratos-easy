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

	var core zapcore.Core
	var fileWriter *dailyRotateWriter

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(globalOption.timeLayout)

	if !constants.IsLocal() {
		// 非本地环境：默认只向终端(Console)写入，若开启了 enableFile 则同时写入文件(JSON)

		// 1. 终端打印采用 console 格式
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel)

		if globalOption.enableFile {
			// 2. 程序运行前检查并删除过期日志
			cleanOldLogs(globalOption.logDir, globalOption.logMaxDays)

			// 3. 文件采用 JSON 格式，按日期分割
			fileWriter = &dailyRotateWriter{dir: globalOption.logDir}
			fileCore := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(fileWriter), zap.DebugLevel)

			core = zapcore.NewTee(consoleCore, fileCore)
		} else {
			core = consoleCore
		}
	} else {
		// 本地环境
		var encoder zapcore.Encoder
		if globalOption.format == "json" {
			encoder = zapcore.NewJSONEncoder(encoderConfig)
		} else {
			// 终端默认采用 console 格式
			encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		}
		core = zapcore.NewCore(encoder, zapcore.AddSync(globalOption.writer), zap.DebugLevel)
	}

	baseLogger := kratoszap.NewLogger(zap.New(core))

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
			_ = fileWriter.Close() // 关闭日志文件
		}
	}

	return logger, cleanup
}
