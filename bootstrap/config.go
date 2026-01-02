package bootstrap

import (
	"flag"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/joho/godotenv"
	"github.com/lhlyu/kratos-easy/utilx"
)

var FlagConf string

func init() {
	root := utilx.FindProjectRoot()
	_ = godotenv.Load(filepath.Join(root, ".env"))
	flag.StringVar(&FlagConf, "conf", filepath.Join(root, "configs"), "config path")
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
