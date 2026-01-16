# Kratos Easy

简化 [go-kratos](https://github.com/go-kratos/kratos) 使用

## 使用

- [例子](https://github.com/lhlyu/kratos-layout)

- 获取依赖

```
go get -u github.com/lhlyu/kratos-easy@latest
```

## 文档 

### 1. bootstrap - 应用启动

#### 核心函数

```go
// Init 初始化，只执行一次
func Init(opts ...Option)

// 执行通用的启动流程
func Run[T any](cfg T, run runner[T])

// 创建 Kratos 应用实例
func NewApp(logger log.Logger, servers ...transport.Server) *kratos.App
```

#### 类型定义

```go
// runner 定义了 main 函数中 wire 注入的逻辑
type runner[T any] func(cfg T, logger log.Logger) (*kratos.App, func(), error)
```

#### 配置加载

- 自动加载 `.env` 文件
- 支持 `-conf` 命令行参数指定配置路径（默认：`configs`）
- 配置源：环境变量 + 配置文件

#### 日志选项

```go
// Option 日志配置选项函数
type Option func(*options)

// WithWriter 设置日志输出目标
func WithWriter(writer io.Writer) Option

// WithFilterLevel 设置日志过滤等级
func WithFilterLevel(level log.Level) Option

// WithJsonFormat 日志输出格式改成 json
func WithJsonFormat() Option

// WithConsoleFormat 日志输出格式改成 console
func WithConsoleFormat() Option

// WithTimeLayout 设置时间格式
func WithTimeLayout(layout string) Option

// WithDisableTrace 禁用链路追踪 trace_id
func WithDisableTrace() Option

// WithEnableSpan 启用链路追踪 span_id
func WithEnableSpan() Option

// WithDisableGlobal 禁止覆盖全局的默认日志器
func WithDisableGlobal() Option

// WithConfigDir 设置配置文件的目录路径
func WithConfigDir(elem ...string) Option

// WithLogDir 设置日志目录
func WithLogDir(dir string) Option

// WithLogMaxDays 设置日志最大保留天数
func WithLogMaxDays(days int) Option

// WithDisableFileLog 禁用文件日志输出
func WithDisableFileLog() Option
```

---

### 2. constants - 环境常量

#### 环境变量键

```go
const AppEnv = "APP_ENV"
const ProjectName = "PROJECT_NAME"
const ProjectRef = "PROJECT_REF"
const ProjectSha = "PROJECT_SHA"

// AppEnv 的枚举值如下
const (
    EnvProduction  = "production"
    EnvStaging     = "staging"
    EnvDevelopment = "development"
    EnvLocal       = "local"
)
```

#### 环境枚举值

```go
const (
    EnvProduction  = "production"
    EnvStaging     = "staging"
    EnvDevelopment = "development"
    EnvLocal       = "local"
)
```

#### 环境判断函数

```go
// CurrentEnv 获取当前环境
func CurrentEnv() string

// IsProduction 是否是正式环境
func IsProduction() bool

// IsStaging 是否是预发布环境
func IsStaging() bool

// IsDevelopment 是否是测试环境
func IsDevelopment() bool

// IsLocal 是否是本地调试
func IsLocal() bool
```

---

### 3. httpx - HTTP 响应

#### 响应结构

```go
type response struct {
    Code int32  // 0-正常; 1~1000-HTTP状态码; 1000+-业务状态码
    Msg  string // 友好的前端提示信息
    Data any    // 数据
}
```

#### 编码函数

```go
// 将 handler 的返回值包装成统一格式并写入 HTTP
func EncodeResponse(w http.ResponseWriter, r *http.Request, v any) error

// 将 error 转成统一的 HTTP 响应
func EncodeError(w http.ResponseWriter, r *http.Request, err error)
```

---

### 4. middlewares - 中间件

#### header - 管理 HTTP 响应头的中间件

```go
// 管理 HTTP 响应头中间件
func Header(opts ...Option) middleware.Middleware

// 自定义 Request ID 的 Header 名
func WithRequestIdHeader(name string) Option

// 禁用 Request Id 响应头
func DisableRequestId() Option

// WithStaticHeader 添加一个静态响应头
func WithStaticHeader(key, value string) Option
```

#### logging - 日志中间件

```go
// 服务端日志中间件
func Server(logger log.Logger, opts ...Option) middleware.Middleware

// 客户端日志中间件
func Client(logger log.Logger, opts ...Option) middleware.Middleware

// 自定义日志脱敏接口
type Redactor interface {
    Redact() string
}

// 配置选项
func WithMaxRawSize(size int) Option
```

#### validate - 验证中间件

```go
// Protobuf 消息验证中间件
func ProtoValidate(opts ...Option) middleware.Middleware

// 配置选项
func WithFriendlyMsg(fn func(err error) string) Option
```

---

### 5. mysqlx - MySQL 客户端

#### 核心函数

```go
// 新建 mysql 客户端
// source 格式：username:password@tcp(domain.com:3306)/dbname?parseTime=True&loc=Local
func NewClient(logger log.Logger, source string, opts ...Option) (*sql.DB, func(), error)
```

#### 配置选项

```go
// Option 表示 mysql 配置项的函数式选项
type Option func(*options)

// WithMaxLifetime 设置单个连接的最大存活时间（d <= 0 表示不限制）
func WithMaxLifetime(d time.Duration) Option

// WithMaxIdleTime 设置单个连接允许的最大空闲时间（d <= 0 表示不限制）
func WithMaxIdleTime(d time.Duration) Option

// WithMaxOpenConns 设置数据库允许的最大打开连接数（n <= 0 表示不限制）
func WithMaxOpenConns(n int) Option

// WithMaxIdleConns 设置连接池中允许保留的最大空闲连接数（n <= 0 表示不保留）
func WithMaxIdleConns(n int) Option

// 预设配置
// WithSmallConfig 适用于小型服务（管理后台、内部工具、低并发任务）
func WithSmallConfig() Option

// WithMediumConfig 适用于中等并发服务（常规 API 服务、业务中台服务）
func WithMediumConfig() Option

// WithHighConcurrencyConfig 适用于高并发服务（核心业务服务、高 QPS 接口）
func WithHighConcurrencyConfig() Option
```

---

### 6. redisx - Redis 客户端

#### 核心函数

```go
// 创建 Redis 客户端
// source 格式：redis://:password@127.0.0.1:6379/0
func NewClient(logger log.Logger, source string, opts ...Option) (*redis.Client, func(), error)
```

#### 配置选项

```go
// Option redis 配置项函数
type Option func(*options)

// WithPoolSize 设置连接池基础连接数
func WithPoolSize(n int) Option

// WithMaxConcurrentDials 设置最大并发创建连接数
func WithMaxConcurrentDials(n int) Option

// WithPoolTimeout 设置获取连接的最大等待时间
func WithPoolTimeout(d time.Duration) Option

// WithMinIdleConns 设置最小空闲连接数
func WithMinIdleConns(n int) Option

// WithMaxIdleConns 设置最大空闲连接数
func WithMaxIdleConns(n int) Option

// WithMaxActiveConns 设置最大活跃连接数（0 表示无限制）
func WithMaxActiveConns(n int) Option

// WithConnMaxIdleTime 设置单个连接最大空闲时间（<=0 表示不关闭）
func WithConnMaxIdleTime(d time.Duration) Option

// WithConnMaxLifetime 设置单个连接最大存活时间（<=0 表示不关闭）
func WithConnMaxLifetime(d time.Duration) Option

// 预设配置
// WithSmallConfig 小型服务默认配置
func WithSmallConfig() Option

// WithMediumConfig 中等并发默认配置
func WithMediumConfig() Option

// WithHighConcurrencyConfig 高并发默认配置
func WithHighConcurrencyConfig() Option

// WithFullCaller 启用打印详细的调用路径
func WithFullCaller() Option
```

---

### 7. utilx - 工具函数

#### json - JSON 工具

```go
// ToJson 任意类型转 json 字符串
func ToJson(v any) string

// ToJsonBytes 任意类型转 []byte
func ToJsonBytes(v any) []byte

// ToObj 将字符串转对象
func ToObj(b string, v any)

// BytesToObj 将 []byte 转对象
func BytesToObj(b []byte, v any)
```

#### string - 字符串工具

```go
// string 转 []byte（零拷贝，只读）
func StringToBytes(s string) []byte

// []byte 转 string（零拷贝，只读）
func BytesToString(b []byte) string
```

#### set - 集合工具

```go
// Set 是一个通用的集合类型，支持任何可比较类型 T
type Set[T comparable] struct

// NewSet 创建一个新的 Set，可以选择性传入初始元素
func NewSet[T comparable](items ...T) *Set[T]

// Add 向集合中添加元素
func (s *Set[T]) Add(item T)

// Remove 从集合中删除元素
func (s *Set[T]) Remove(item T)

// Contains 检查集合中是否存在元素
func (s *Set[T]) Contains(item T) bool

// Len 返回集合中元素的数量
func (s *Set[T]) Len() int

// IsEmpty 判断集合是否为空。
func (s *Set[T]) IsEmpty() bool

// Clear 清空集合
func (s *Set[T]) Clear()

// ToSlice 将集合转换为切片
func (s *Set[T]) ToSlice() []T

// Clone 返回集合的一个副本。
func (s *Set[T]) Clone() *Set[T]

// Union 返回两个集合的并集
func (s *Set[T]) Union(other *Set[T]) *Set[T]

// Intersect 返回两个集合的交集
func (s *Set[T]) Intersect(other *Set[T]) *Set[T]

// Difference 返回集合 s 相对于 other 的差集
func (s *Set[T]) Difference(other *Set[T]) *Set[T]

// All 返回一个迭代器，用于遍历集合中的所有元素。
// 支持 Go 1.23+ 的 range-over-func 特性。
func (s *Set[T]) All() iter.Seq[T]
```

#### time – 时间工具

```go
// 返回当前时间戳（毫秒）
func Now() int64
```

#### retry - 重试工具

```go
// 对一个函数进行重试
func Retry(ctx context.Context, maxRetries int, delay time.Duration, fn func() error) error
```

#### slice - 切片工具

```go
// 将切片按 batchSize 条一组顺序执行 fn
func Process[T any](data []T, batchSize int, fn func([]T) error) error

// 从切片中取前 maxLen 个元素
func Take[T any](s []T, maxLen int) []T

// Prepend 在切片 s 前面插入 elems，并返回新切片。
// 不会修改原切片。
// 当 elems 为空时，直接返回 s。
func Prepend[T any](s []T, elems ...T) []T
```

#### md5 - MD5 工具

```go
// MD5 编码，支持多个 salt
func Md5Encode(val string, salt ...string) string
```

#### go - 并发工具

```go
// 并发执行所有处理函数，返回第一个错误
func Parallel(handlers ...func() error) error

// 并发执行所有处理函数，忽略所有错误和 panic
func ParallelSafe(handlers ...func())
```

#### path – 路径工具

```go
// 找到项目根目录（查找 go.mod）
func FindProjectRoot() string

// 分隔符 + 任意基本类型变参
func Join(sep string, vals ...any) string
```

