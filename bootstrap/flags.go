package bootstrap

import (
	"flag"
	"os"
	"path/filepath"
)

var FlagConf string

// 默认配置文件优先级
var defaultConfCandidates = []string{
	"config.test.yaml",
	"config.local.yaml",
	"config.yaml",
}

// findDefaultConf 返回第一个存在的配置文件，如果都不存在就返回 configs 目录
func findDefaultConf(root string) string {
	dir := globalOption.configDir

	for _, name := range defaultConfCandidates {
		path := filepath.Join(root, dir, name)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	// 都没找到就用默认路径
	return filepath.Join(root, dir)
}

// initFlags 初始化 flag
func initFlags(root string) {
	flag.StringVar(
		&FlagConf,
		"conf",
		findDefaultConf(root),
		"config path",
	)
}
