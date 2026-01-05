package header

import "github.com/go-kratos/kratos/v2/transport"

// applyStaticHeaders 写入用户自定义的静态响应头
func applyStaticHeaders(
	h transport.Header,
	opt *options,
) {
	for k, v := range opt.staticHeaders {
		h.Set(k, v)
	}
}
