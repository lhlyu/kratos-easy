package bootstrap

import (
	"sync"

	"github.com/lhlyu/kratos-easy/utilx"
)

var once sync.Once

// Init 初始化，只执行一次
func Init(opts ...Option) {
	once.Do(func() {
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
