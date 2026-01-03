package utilx

import (
	"crypto/md5"
	"encoding/hex"
)

// Md5Encode 对传入的字符串进行 MD5 编码，并可选地追加一个或多个 salt。
//
// 该方法会先写入 val，再按顺序写入所有 salt，最终返回 32 位小写十六进制字符串。
// 如果 val 为空字符串，则直接返回空字符串，而不会计算 MD5("") 的结果。
//
// 示例：
//
//	Md5Encode("hello")                  // 无 salt
//	Md5Encode("hello", "a", "b", "c")    // 多个 salt
func Md5Encode(val string, salt ...string) string {
	if val == "" {
		return ""
	}

	h := md5.New()

	// hash.Hash.Write 的 error 永远为 nil，这里显式忽略返回值
	_, _ = h.Write(StringToBytes(val))
	for _, s := range salt {
		_, _ = h.Write(StringToBytes(s))
	}

	return hex.EncodeToString(h.Sum(nil))
}
