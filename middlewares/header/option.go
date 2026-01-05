package header

const defaultRequestIDHeader = "x-request-id"

// options 定义响应头中间件的配置项
type options struct {
	enableRequestId bool
	requestIdHeader string
	staticHeaders   map[string]string
}

// Option 定义配置函数
type Option func(*options)

// newOptions 初始化配置
func newOptions(opts ...Option) *options {
	o := &options{
		enableRequestId: true,
		requestIdHeader: defaultRequestIDHeader,
		staticHeaders:   make(map[string]string),
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// WithRequestIdHeader 自定义 Request ID 的 Header 名
func WithRequestIdHeader(name string) Option {
	return func(o *options) {
		if name != "" {
			o.requestIdHeader = name
		}
	}
}

// DisableRequestId 禁用 Request Id 响应头
func DisableRequestId() Option {
	return func(o *options) {
		o.enableRequestId = false
	}
}

// WithStaticHeader 添加一个静态响应头
func WithStaticHeader(key, value string) Option {
	return func(o *options) {
		if key != "" {
			o.staticHeaders[key] = value
		}
	}
}
