package header

import (
	"context"

	"github.com/go-kratos/kratos/v2/transport"
	"go.opentelemetry.io/otel/trace"
)

// applyRequestId 向响应头写入 Request Id
func applyRequestId(
	ctx context.Context,
	h transport.Header,
	opt *options,
) {
	if !opt.enableRequestId {
		return
	}

	sc := trace.SpanContextFromContext(ctx)
	if sc.IsValid() {
		h.Set(opt.requestIdHeader, sc.TraceID().String())
	}
}
