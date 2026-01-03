package redisx

import "time"

// options redis 配置选项
type options struct {
	// PoolSize 表示连接池的基础连接数。
	// 当连接不足时，会动态分配更多连接。
	// 默认：10 * runtime.GOMAXPROCS(0)
	PoolSize int

	// MaxConcurrentDials 表示同时创建连接的最大 goroutine 数。
	// 如果 <=0，则默认等于 PoolSize。
	// 如果大于 PoolSize，则上限为 PoolSize。
	MaxConcurrentDials int

	// PoolTimeout 表示获取连接时等待的最大时间，如果超时返回错误。
	// 默认：ReadTimeout + 1 秒
	PoolTimeout time.Duration

	// MinIdleConns 表示连接池中最小空闲连接数。
	// 默认：0
	MinIdleConns int

	// MaxIdleConns 表示连接池中最大空闲连接数。
	// 默认：0
	MaxIdleConns int

	// MaxActiveConns 表示连接池中最大活跃连接数。
	// 为 0 表示无限制。
	MaxActiveConns int

	// ConnMaxIdleTime 表示单个连接允许的最大空闲时间。
	// <=0 表示不关闭空闲连接
	// 默认：30 分钟
	ConnMaxIdleTime time.Duration

	// ConnMaxLifetime 表示单个连接允许存在的最大时间。
	// <=0 表示不关闭
	ConnMaxLifetime time.Duration
}

// Option redis 配置项函数
type Option func(*options)

// WithPoolSize 设置连接池基础连接数
func WithPoolSize(n int) Option {
	return func(o *options) {
		if n > 0 {
			o.PoolSize = n
		}
	}
}

// WithMaxConcurrentDials 设置最大并发创建连接数
func WithMaxConcurrentDials(n int) Option {
	return func(o *options) {
		if n > 0 {
			o.MaxConcurrentDials = n
		}
	}
}

// WithPoolTimeout 设置获取连接的最大等待时间
func WithPoolTimeout(d time.Duration) Option {
	return func(o *options) {
		if d > 0 {
			o.PoolTimeout = d
		}
	}
}

// WithMinIdleConns 设置最小空闲连接数
func WithMinIdleConns(n int) Option {
	return func(o *options) {
		if n >= 0 {
			o.MinIdleConns = n
		}
	}
}

// WithMaxIdleConns 设置最大空闲连接数
func WithMaxIdleConns(n int) Option {
	return func(o *options) {
		if n >= 0 {
			o.MaxIdleConns = n
		}
	}
}

// WithMaxActiveConns 设置最大活跃连接数
func WithMaxActiveConns(n int) Option {
	return func(o *options) {
		if n > 0 {
			o.MaxActiveConns = n
		}
	}
}

// WithConnMaxIdleTime 设置单个连接最大空闲时间
func WithConnMaxIdleTime(d time.Duration) Option {
	return func(o *options) {
		if d > 0 {
			o.ConnMaxIdleTime = d
		}
	}
}

// WithConnMaxLifetime 设置单个连接最大存活时间
func WithConnMaxLifetime(d time.Duration) Option {
	return func(o *options) {
		if d > 0 {
			o.ConnMaxLifetime = d
		}
	}
}

// WithSmallConfig 小型服务默认配置
func WithSmallConfig() Option {
	return func(o *options) {
		o.PoolSize = 10
		o.MinIdleConns = 5
		o.MaxIdleConns = 10
		o.MaxActiveConns = 0
		o.ConnMaxLifetime = 30 * time.Minute
		o.ConnMaxIdleTime = 10 * time.Minute
		o.PoolTimeout = 5 * time.Second
		o.MaxConcurrentDials = 10
	}
}

// WithMediumConfig 中等并发默认配置
func WithMediumConfig() Option {
	return func(o *options) {
		o.PoolSize = 50
		o.MinIdleConns = 25
		o.MaxIdleConns = 50
		o.MaxActiveConns = 0
		o.ConnMaxLifetime = 30 * time.Minute
		o.ConnMaxIdleTime = 5 * time.Minute
		o.PoolTimeout = 3 * time.Second
		o.MaxConcurrentDials = 50
	}
}

// WithHighConcurrencyConfig 高并发默认配置
func WithHighConcurrencyConfig() Option {
	return func(o *options) {
		o.PoolSize = 100
		o.MinIdleConns = 50
		o.MaxIdleConns = 100
		o.MaxActiveConns = 0
		o.ConnMaxLifetime = 15 * time.Minute
		o.ConnMaxIdleTime = 3 * time.Minute
		o.PoolTimeout = 2 * time.Second
		o.MaxConcurrentDials = 100
	}
}
