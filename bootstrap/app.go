package bootstrap

import (
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/google/uuid"
	"github.com/lhlyu/kratos-easy/constants"
)

// NewApp 创建一个 Kratos 应用实例，带默认 metadata 和日志输出
func NewApp(logger log.Logger, servers ...transport.Server) *kratos.App {
	// 生成唯一实例 ID
	id := uuid.NewString()

	// 获取应用信息
	name := os.Getenv(constants.ProjectName)
	if name == "" {
		name = "unknown-service"
	}
	ref := os.Getenv(constants.ProjectRef)
	if ref == "" {
		ref = "unknown-ref"
	}
	version := os.Getenv(constants.ProjectSha)
	if version == "" {
		version = "unknown-version"
	}
	env := os.Getenv(constants.AppEnv)
	if env == "" {
		env = "unknown-env"
	}

	// 构造 metadata
	md := map[string]string{
		"env": env,
	}

	// 启动日志
	_ = logger.Log(
		log.LevelInfo,
		"msg", "starting service",
		"service.name", name,
		"service.id", id,
		"service.ref", ref,
		"service.version", version,
		"env", env,
	)

	// 返回 Kratos App
	return kratos.New(
		kratos.ID(id),
		kratos.Name(name),
		kratos.Version(version),
		kratos.Metadata(md),
		kratos.Logger(logger),
		kratos.Server(servers...),
	)
}
