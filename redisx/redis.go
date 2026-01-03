package redisx

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
)

// NewClient 创建 Redis 客户端
//
// source 格式示例：redis://:password@127.0.0.1:6379/0
func NewClient(logger log.Logger, source string, opts ...Option) (*redis.Client, func(), error) {
	l := log.NewHelper(logger)

	opt, err := redis.ParseURL(source)
	if err != nil {
		l.Errorw("parse redis url failed", "error", err)
		return nil, nil, err
	}

	o := &options{
		PoolSize:           10,
		MinIdleConns:       5,
		MaxIdleConns:       10,
		MaxActiveConns:     0,
		ConnMaxLifetime:    30 * time.Minute,
		ConnMaxIdleTime:    10 * time.Minute,
		PoolTimeout:        5 * time.Second,
		MaxConcurrentDials: 10,
	}

	for _, apply := range opts {
		apply(o)
	}

	// 应用配置到 redis.Options
	opt.PoolSize = o.PoolSize
	opt.MinIdleConns = o.MinIdleConns
	opt.MaxIdleConns = o.MaxIdleConns
	opt.MaxActiveConns = o.MaxActiveConns
	opt.ConnMaxLifetime = o.ConnMaxLifetime
	opt.ConnMaxIdleTime = o.ConnMaxIdleTime
	opt.PoolTimeout = o.PoolTimeout
	opt.MaxConcurrentDials = o.MaxConcurrentDials

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		l.Errorw("ping redis failed", "error", err)
		return nil, nil, err
	}

	l.Infow("redis connected")

	cleanup := func() {
		if err := client.Close(); err != nil {
			l.Errorw("close redis failed", "error", err)
		}
	}

	return client, cleanup, nil
}
