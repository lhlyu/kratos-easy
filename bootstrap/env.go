package bootstrap

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// env 文件加载优先级（只加载第一个存在的）
var envCandidates = []string{
	".env",
	"local.env",
	"dev.env",
	"develop.env",
	"development.env",
	"staging.env",
	"prod.env",
	"production.env",
}

func initEnv(root string) {
	for _, name := range envCandidates {
		path := filepath.Join(root, name)
		if _, err := os.Stat(path); err == nil {
			_ = godotenv.Load(path)
			return
		}
	}
}
