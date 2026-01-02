package constants

import "os"

// AppEnv 应用环境
const AppEnv = "APP_ENV"

// AppEnv 的枚举值如下
const (
	EnvProduction  = "production"
	EnvStaging     = "staging"
	EnvDevelopment = "development"
	EnvLocal       = "local"
)

// 环境变量
const (
	// ProjectName 项目名字
	ProjectName = "PROJECT_NAME"
	// ProjectRef 项目分支
	ProjectRef = "PROJECT_REF"
	// ProjectSha 项目版本
	ProjectSha = "PROJECT_SHA"
)

// CurrentEnv 获取当前环境
func CurrentEnv() string {
	return os.Getenv(AppEnv)
}

// IsProduction 是否是正式环境
func IsProduction() bool {
	return CurrentEnv() == EnvProduction
}

// IsStaging 是否是预发布环境
func IsStaging() bool {
	return CurrentEnv() == EnvStaging
}

// IsDevelopment 是否是测试环境
func IsDevelopment() bool {
	return CurrentEnv() == EnvDevelopment
}

var isLocal = os.Getenv(ProjectName) == ""

// IsLocal 是否是本地调试
func IsLocal() bool {
	if os.Getenv(ProjectName) == "" {
		return true
	}
	return CurrentEnv() == EnvLocal
}
