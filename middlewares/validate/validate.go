package validate

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"

	"buf.build/go/protovalidate"
	"google.golang.org/protobuf/proto"
)

// validator 定义旧版自定义验证器接口。
type validator interface {
	Validate() error
}

/************************
 * Option & Config
 ************************/

// Options 定义 ProtoValidate 中间件可配置项。
type Options struct {
	friendlyMsg func(error) string
}

// Option 定义配置选项类型。
type Option func(*Options)

// WithFriendlyMsg 设置自定义前端友好提示函数。
func WithFriendlyMsg(fn func(err error) string) Option {
	return func(o *Options) {
		o.friendlyMsg = fn
	}
}

// WithDetailedMsg 提供更详细的错误信息选项
func WithDetailedMsg() Option {
	return func(o *Options) {
		o.friendlyMsg = func(err error) string {
			// 返回具体的字段验证错误
			return err.Error()
		}
	}
}

// newOptions 初始化 Options。
func newOptions(opts ...Option) *Options {
	o := &Options{
		friendlyMsg: func(err error) string {
			// 默认直接返回错误信息
			return "参数不合法"
		},
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

/************************
 * Middleware
 ************************/

// ProtoValidate 返回一个中间件，用于对请求的 Protobuf 消息进行验证。
// 支持 protovalidate 和旧版 Validate() 验证器。
// 可通过 Option 自定义前端友好提示信息。
func ProtoValidate(opts ...Option) middleware.Middleware {
	o := newOptions(opts...)

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			// 验证 protovalidate 消息
			if msg, ok := req.(proto.Message); ok {
				if err := protovalidate.Validate(msg); err != nil {
					return nil, handleErr(err, o)
				}
			}

			// 验证旧版 validator 消息
			if v, ok := req.(validator); ok {
				if err := v.Validate(); err != nil {
					return nil, handleErr(err, o)
				}
			}

			return handler(ctx, req)
		}
	}
}

/************************
 * Internal
 ************************/

// handleErr 统一处理验证错误，生成 BadRequest 错误。
func handleErr(err error, o *Options) error {
	return errors.BadRequest("VALIDATOR", o.friendlyMsg(err)).WithCause(err)
}
