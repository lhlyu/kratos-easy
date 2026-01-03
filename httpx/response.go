package httpx

type response struct {
	// 状态码: 0 - 正常; 1 ~ 1000 - http状态码; 1000以上 业务状态码
	Code int32 `json:"code"`
	// 友好的前端提示信息
	Msg string `json:"msg"`
	// 数据
	Data any `json:"data"`
}
