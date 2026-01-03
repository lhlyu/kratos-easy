package utilx

import (
	"encoding/json"
)

// ToJson 任意类型转 json 字符串
func ToJson(v any) string {
	b, _ := json.Marshal(v)
	return BytesToString(b)
}

// ToJsonBytes 任意类型转 []byte
func ToJsonBytes(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}

// ToObj 将字符串 转 对象
func ToObj(b string, v any) {
	_ = json.Unmarshal(StringToBytes(b), v)
}

// BytesToObj 将 []byte 转 对象
func BytesToObj(b []byte, v any) {
	_ = json.Unmarshal(b, v)
}
