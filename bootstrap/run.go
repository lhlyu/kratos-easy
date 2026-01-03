package bootstrap

import (
	"time"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/lhlyu/kratos-easy/constants"
)

// runner 定义了 main 函数中 wire 注入的逻辑
// T 是每个项目自定义的 Bootstrap 配置结构体类型
type runner[T any] func(cfg T, logger log.Logger) (*kratos.App, func(), error)

// Run 执行通用的启动流程
func Run[T any](cfg T, run runner[T], opts ...Option) {

	if constants.IsLocal() {
		opts = append(
			opts,
			WithTimeLayout(time.TimeOnly),
			WithFilterLevel(log.LevelDebug),
		)
	}
	if constants.IsDevelopment() {
		opts = append(
			opts,
			WithFilterLevel(log.LevelDebug),
		)
	}

	logger := newLogger(opts...)

	// 加载配置到传入的泛型结构体
	if err := loadConfig(cfg); err != nil {
		panic(err)
	}

	// 执行业务注入逻辑 (调用 main 里的 wireApp)
	app, cleanup, err := run(cfg, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// 运行服务
	if err := app.Run(); err != nil {
		panic(err)
	}
}
