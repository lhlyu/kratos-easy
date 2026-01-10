package bootstrap

import (
	"flag"
	"os"
	"path/filepath"
)

var FlagConf string

// ConfDir 配置文件所在的目录
var ConfDir string

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
	ConfDir, _ = getDirPath(FlagConf)
}

// getDirPath 根据传入路径返回目录路径
// 如果 path 是目录，则直接返回
// 如果 path 是文件，则返回文件所在目录
func getDirPath(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if info.IsDir() {
		// 是目录，直接返回
		return path, nil
	}

	// 是文件，返回它所在的目录
	return filepath.Dir(path), nil
}
