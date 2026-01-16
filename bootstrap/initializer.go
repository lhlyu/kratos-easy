package bootstrap

import (
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/lhlyu/kratos-easy/constants"
	"github.com/lhlyu/kratos-easy/utilx"
)

var once sync.Once

// Init 初始化，只执行一次
func Init(opts ...Option) {
	once.Do(func() {
		if constants.IsLocal() {
			opts = utilx.Prepend(
				opts,
				WithTimeLayout(time.TimeOnly),
				WithFilterLevel(log.LevelDebug),
			)
		}
		if constants.IsDevelopment() {
			opts = utilx.Prepend(
				opts,
				WithFilterLevel(log.LevelDebug),
			)
		}
		for _, opt := range opts {
			opt(globalOption)
		}
		initInternal()
	})
}

// initInternal contains real initialization logic.
func initInternal() {
	root := utilx.FindProjectRoot()

	initEnv(root)
	initFlags(root)
}
