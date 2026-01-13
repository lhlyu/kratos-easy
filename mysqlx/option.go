package mysqlx

import "time"

// options mysql 配置选项
type options struct {
	// maxOpen 表示数据库允许的最大打开连接数。
	maxOpen int

	// maxIdleCount 表示连接池中允许保留的最大空闲连接数。
	maxIdleCount int

	// maxLifetime 表示单个连接的最大存活时间。
	maxLifetime time.Duration

	// maxIdleTime 表示单个连接允许的最大空闲时间。
	maxIdleTime time.Duration
}

// Option 表示 mysql 配置项的函数式选项。
type Option func(*options)

// WithMaxLifetime 设置单个连接的最大存活时间。
//
// 如果 d <= 0，则表示不限制连接的存活时间。
func WithMaxLifetime(d time.Duration) Option {
	return func(o *options) {
		if d > 0 {
			o.maxLifetime = d
		}
	}
}

// WithMaxIdleTime 设置单个连接允许的最大空闲时间。
//
// 如果 d <= 0，则表示不限制连接的空闲时间。
func WithMaxIdleTime(d time.Duration) Option {
	return func(o *options) {
		if d > 0 {
			o.maxIdleTime = d
		}
	}
}

// WithMaxOpenConns 设置数据库允许的最大打开连接数。
//
// 如果 n <= 0，则表示不限制最大连接数。
func WithMaxOpenConns(n int) Option {
	return func(o *options) {
		if n > 0 {
			o.maxOpen = n
		}
	}
}

// WithMaxIdleConns 设置连接池中允许保留的最大空闲连接数。
//
// 如果 n <= 0，则表示不保留空闲连接。
func WithMaxIdleConns(n int) Option {
	return func(o *options) {
		if n > 0 {
			o.maxIdleCount = n
		}
	}
}

// WithCustomConfig 添加更多预设
func WithCustomConfig(maxOpen, maxIdle int, lifetime, idleTime time.Duration) Option {
	return func(o *options) {
		o.maxOpen = maxOpen
		o.maxIdleCount = maxIdle
		o.maxLifetime = lifetime
		o.maxIdleTime = idleTime
	}
}

// WithSmallConfig 返回适用于小型服务的默认数据库配置。
//
// 适用场景：
//   - 管理后台
//   - 内部工具
//   - 低并发任务服务
func WithSmallConfig() Option {
	return func(o *options) {
		o.maxOpen = 20
		o.maxIdleCount = 10
		o.maxLifetime = 30 * time.Minute
		o.maxIdleTime = 10 * time.Minute
	}
}

// WithMediumConfig 返回适用于中等并发服务的默认数据库配置。
//
// 适用场景：
//   - 常规 API 服务
//   - 业务中台服务
func WithMediumConfig() Option {
	return func(o *options) {
		o.maxOpen = 50
		o.maxIdleCount = 25
		o.maxLifetime = 30 * time.Minute
		o.maxIdleTime = 5 * time.Minute
	}
}

// WithHighConcurrencyConfig 返回适用于高并发服务的默认数据库配置。
//
// 适用场景：
//   - 核心业务服务
//   - 高 QPS 接口
//   - 流量波动明显的服务
func WithHighConcurrencyConfig() Option {
	return func(o *options) {
		o.maxOpen = 100
		o.maxIdleCount = 50
		o.maxLifetime = 15 * time.Minute
		o.maxIdleTime = 3 * time.Minute
	}
}
