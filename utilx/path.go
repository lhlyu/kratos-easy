package utilx

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cast"
)

// FindProjectRoot 找到项目根目录
func FindProjectRoot() string {
	dir, _ := os.Getwd()

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// 已经到根目录了
			return ""
		}

		dir = parent
	}
}

// Join 分隔符 + 任意基本类型变参
func Join(sep string, vals ...any) string {
	if len(vals) == 0 {
		return ""
	}
	b := make([]string, len(vals))
	for i, v := range vals {
		if v == "" {
			continue
		}
		b[i] = cast.ToString(v)
	}
	return strings.Join(b, sep)
}
