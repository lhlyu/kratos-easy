package bootstrap

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/joho/godotenv"
	"github.com/lhlyu/kratos-easy/utilx"
)

var FlagConf string

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

func init() {
	root := utilx.FindProjectRoot()
	loadFirstEnv(root)
	flag.StringVar(&FlagConf, "conf", filepath.Join(root, "configs"), "config path")
}

// loadFirstEnv 按优先级加载第一个存在的 env 文件
func loadFirstEnv(root string) {
	for _, name := range envCandidates {
		path := filepath.Join(root, name)
		if _, err := os.Stat(path); err == nil {
			_ = godotenv.Load(path)
			return
		}
	}
}

func loadConfig(bc any) error {
	flag.Parse()

	c := config.New(
		config.WithSource(
			env.NewSource(),
			file.NewSource(FlagConf),
		),
	)

	defer func() {
		_ = c.Close()
	}()

	if err := c.Load(); err != nil {
		return err
	}

	if err := c.Scan(bc); err != nil {
		return err
	}

	return nil
}
