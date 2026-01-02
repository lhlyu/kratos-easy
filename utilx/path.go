package utilx

import (
	"os"
	"path/filepath"
)

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
