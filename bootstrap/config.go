package bootstrap

import (
	"flag"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
)

func loadConfig(bc any) error {

	Init()

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
