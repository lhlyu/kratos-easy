package httpx

import (
	"encoding/json"
	netHttp "net/http"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// EncodeResponse 将 handler 的返回值包装成统一格式并写入 HTTP。
func EncodeResponse(w http.ResponseWriter, r *http.Request, v any) error {
	// 支持重定向
	if rd, ok := v.(http.Redirector); ok {
		url, code := rd.Redirect()
		netHttp.Redirect(w, r, url, code)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	return json.NewEncoder(w).Encode(response{
		Code: 0,
		Msg:  "",
		Data: v,
	})
}

// EncodeError 将 error 转成统一的 HTTP 响应。
func EncodeError(w http.ResponseWriter, r *http.Request, err error) {
	se := errors.FromError(err)
	if se == nil {
		se = errors.New(errors.UnknownCode, errors.UnknownReason, "服务器异常")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(int(se.Code))

	_ = json.NewEncoder(w).Encode(response{
		Code: se.Code,
		Msg:  se.Message,
		Data: nil,
	})
}
