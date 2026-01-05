package header

import (
	"context"

	"github.com/go-kratos/kratos/v2/transport"
)

// applyHeaders 根据配置向响应写入 Header
func applyHeaders(
	ctx context.Context,
	info transport.Transporter,
	opt *options,
) {
	h := info.ReplyHeader()

	applyRequestId(ctx, h, opt)
	applyStaticHeaders(h, opt)
}
