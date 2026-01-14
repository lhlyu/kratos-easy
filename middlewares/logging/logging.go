package logging

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/lhlyu/kratos-easy/utilx"
)

/************************
 * Option & Config
 ************************/

// options 定义日志中间件的可配置参数。
type options struct {
	maxRawSize int
}

// defaultMaxRawSize 定义日志中原始数据的默认最大长度（字节）。
const defaultMaxRawSize = 2 << 10 // 2KB

// Option 定义日志中间件的可选配置项。
type Option func(*options)

// WithMaxRawSize 设置日志中可输出的原始数据最大长度。
func WithMaxRawSize(size int) Option {
	return func(o *options) {
		if size > 0 {
			o.maxRawSize = size
		}
	}
}

// newOptions 创建并初始化日志配置。
func newOptions(opts ...Option) *options {
	o := &options{
		maxRawSize: defaultMaxRawSize,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

/************************
 * Public API
 ************************/

// Redactor 用于自定义日志脱敏行为。
//
// 实现该接口的对象将优先使用 Redact 方法生成日志内容。
type Redactor interface {
	Redact() string
}

// Server 返回服务端日志中间件。
func Server(logger log.Logger, opts ...Option) middleware.Middleware {
	return loggingMiddleware(logger, "server", newOptions(opts...))
}

// Client 返回客户端日志中间件。
func Client(logger log.Logger, opts ...Option) middleware.Middleware {
	return loggingMiddleware(logger, "client", newOptions(opts...))
}

var transports = map[string]func(ctx context.Context) (tr transport.Transporter, ok bool){
	"server": transport.FromServerContext,
	"client": transport.FromClientContext,
}

/************************
 * Middleware
 ************************/

// loggingMiddleware 创建统一的日志中间件实现。
func loggingMiddleware(
	logger log.Logger,
	kind string,
	opt *options,
) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			start := time.Now()
			helper := log.NewHelper(log.WithContext(ctx, logger))

			var (
				transportType string
				operation     string
				e             string
			)

			if h, ok := transports[kind]; ok {
				if info, ok := h(ctx); ok {
					transportType = info.Kind().String()
					operation = info.Operation()
				}
			}

			// 记录请求日志
			helper.Log(
				log.LevelInfo,
				"kind", kind,
				"transport", transportType,
				"operation", operation,
				"req", extractArgs(req, opt),
			)

			reply, err := handler(ctx, req)

			if se := errors.FromError(err); se != nil {
				e = "⟦" + se.Error() + "⟧"
			}

			kvs := make([]any, 0, 10)

			kvs = append(kvs,
				"kind", kind,
				"transport", transportType,
				"operation", operation,
				"resp", extractArgs(reply, opt),
			)

			level := log.LevelInfo
			if err != nil {
				level = log.LevelError
				kvs = append(kvs, "err", e)
			}

			kvs = append(kvs, "latency", time.Since(start).String())

			// 记录响应日志
			helper.Log(level, kvs...)

			return reply, err
		}
	}
}

/************************
 * Extract & Guard
 ************************/

// extractArgs 将请求或响应对象转换为可记录的字符串。
func extractArgs(v any, opt *options) string {
	// Redactor 优先处理
	if r, ok := v.(Redactor); ok {
		return r.Redact()
	}

	// 特殊类型保护
	if s, ok := guardByType(v, opt.maxRawSize); ok {
		return s
	}

	// JSON 序列化
	s := utilx.ToJson(v)

	// 长度裁剪
	return guardBySize(s, opt.maxRawSize)
}

// guardByType 对特殊类型进行保护性处理，避免日志膨胀或读取风险。
func guardByType(v any, max int) (string, bool) {
	switch t := v.(type) {
	case nil:
		return "null", true
	case []byte:
		if len(t) > max {
			return fmt.Sprintf("<bytes %d bytes omitted>", len(t)), true
		}
		return string(t), true
	case string:
		if len(t) > max {
			return fmt.Sprintf("<string %d bytes omitted>", len(t)), true
		}
		return t, true
	case io.Reader:
		return "<io.Reader omitted>", true
	case multipart.File:
		return "<multipart.File omitted>", true
	default:
		return "", false
	}
}

// guardBySize 对字符串按最大长度进行裁剪。
func guardBySize(s string, max int) string {
	total := len(s)
	if total <= max {
		return s
	}

	r := []rune(s)
	if len(r) > max {
		s = string(r[:max])
	}

	return fmt.Sprintf(
		"%s...<truncated, total=%d bytes>",
		s,
		total,
	)
}
