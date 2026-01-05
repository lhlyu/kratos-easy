package header

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

// Header 返回一个用于统一管理 HTTP 响应头的中间件。
//
// 默认行为：
//   - 返回 x-request-id（基于 OpenTelemetry TraceID）
//
// 可通过 Option 自定义或禁用相关行为。
func Header(opts ...Option) middleware.Middleware {
	opt := newOptions(opts...)

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			if info, ok := transport.FromServerContext(ctx); ok {
				// 当前仅作用于 HTTP Server
				if info.Kind().String() == "http" {
					applyHeaders(ctx, info, opt)
				}
			}
			return handler(ctx, req)
		}
	}
}
