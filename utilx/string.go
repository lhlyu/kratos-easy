package utilx

import (
	"unsafe"
)

// StringToBytes 将 string 转为只读 []byte，不发生内存拷贝。
// 返回的 []byte 仅用于只读场景，严禁修改。
func StringToBytes(s string) []byte {
	if s == "" {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// BytesToString 将 []byte 转换为 string，不发生内存拷贝。
//
// 该方法使用 unsafe 实现，仅适用于只读、短生命周期场景。
// 返回的 string 与原始 []byte 共享同一块底层内存，
// 如果在转换后继续修改 []byte 的内容，将导致未定义行为。
//
// 使用场景示例：
//   - 编解码
//   - 哈希计算
//   - 日志格式化
//   - 协议解析（只读）
//
// 注意：
//   - 严禁在转换后修改原始 []byte
//   - 不要长期持有返回的 string
//   - 不要在并发场景下修改 []byte
func BytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(b), len(b))
}
