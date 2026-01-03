# Kratos Easy 速查表

> 简化 go-kratos 使用的工具库

## 📦 安装

```bash
go get -u github.com/lhlyu/kratos-easy@latest
```

---

## 🚀 Bootstrap - 应用启动

### 快速开始

```go
package main

import (
    "kratos-layout/internal/conf"
    "github.com/go-kratos/kratos/v2"
    "github.com/go-kratos/kratos/v2/log"
    "github.com/go-kratos/kratos/v2/transport/grpc"
    "github.com/go-kratos/kratos/v2/transport/http"
    "github.com/lhlyu/kratos-easy/bootstrap"
    _ "go.uber.org/automaxprocs"
)

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
    return bootstrap.NewApp(logger, gs, hs)
}

func main() {
    var bc conf.Conf
    bootstrap.Run(&bc, wireApp)
}
```

### 核心函数

#### `bootstrap.NewApp`
创建 Kratos 应用实例，自动设置 metadata 和日志

```go
func NewApp(logger log.Logger, servers ...transport.Server) *kratos.App
```

**特性：**
- 自动生成唯一实例 ID
- 从环境变量读取应用信息（PROJECT_NAME, PROJECT_REF, PROJECT_SHA, APP_ENV）
- 自动记录启动日志

#### `bootstrap.Run`
执行通用的启动流程

```go
func Run[T any](cfg T, run runner[T], opts ...Option)
```

**特性：**
- 自动加载配置（支持环境变量和配置文件）
- 自动初始化日志（根据环境自动调整）
- 本地环境：时间格式 `TimeOnly`，日志级别 `Debug`
- 开发环境：日志级别 `Debug`

### 日志配置选项

#### `WithWriter`
设置日志输出目标

```go
bootstrap.WithWriter(os.Stdout)
```

#### `WithFilterLevel`
设置日志过滤等级

```go
bootstrap.WithFilterLevel(log.LevelInfo)
```

#### `WithJsonFormat` / `WithConsoleFormat`
设置日志输出格式

```go
bootstrap.WithJsonFormat()      // JSON 格式
bootstrap.WithConsoleFormat()   // 控制台格式（默认）
```

#### `WithTimeLayout`
设置时间格式

```go
bootstrap.WithTimeLayout(time.DateTime)
```

#### `WithDisableTrace` / `WithEnableSpan`
控制链路追踪

```go
bootstrap.WithDisableTrace()  // 禁用 trace_id
bootstrap.WithEnableSpan()    // 启用 span_id
```

#### `WithDisableGlobal`
禁止覆盖全局日志器

```go
bootstrap.WithDisableGlobal()
```

### 配置加载

- 自动加载 `.env` 文件（从项目根目录）
- 支持命令行参数 `-conf` 指定配置路径（默认：`configs`）
- 支持环境变量和配置文件两种配置源

---

## 🌐 HTTP - HTTP 响应处理

### 统一响应格式

```go
type response struct {
    Code int32 `json:"code"`  // 0-正常; 1-1000 HTTP状态码; 1000+ 业务状态码
    Msg  string `json:"msg"`  // 友好的前端提示信息
    Data any    `json:"data"` // 数据
}
```

### 编码器

#### `EncodeResponse`
将 handler 返回值包装成统一格式

```go
func EncodeResponse(w http.ResponseWriter, r *http.Request, v any) error
```

**特性：**
- 自动包装为统一响应格式
- 支持重定向（实现 `http.Redirector` 接口）

#### `EncodeError`
将 error 转成统一的 HTTP 响应

```go
func EncodeError(w http.ResponseWriter, r *http.Request, err error)
```

**特性：**
- 自动从 Kratos errors 提取状态码和消息
- 未知错误默认返回 "服务器异常"

---

## 🔧 Middlewares - 中间件

### 日志中间件

#### `logging.Server`
服务端日志中间件

```go
func Server(logger log.Logger, opts ...Option) middleware.Middleware
```

#### `logging.Client`
客户端日志中间件

```go
func Client(logger log.Logger, opts ...Option) middleware.Middleware
```

**配置选项：**

```go
// 设置日志中可输出的原始数据最大长度（默认 2KB）
logging.WithMaxRawSize(4096)
```

**特性：**
- 自动记录请求和响应
- 自动提取 trace_id 并设置到响应头 `x-trace-id`
- 自动处理特殊类型（[]byte, string, io.Reader, multipart.File）
- 支持自定义脱敏（实现 `Redactor` 接口）

**Redactor 接口：**

```go
type Redactor interface {
    Redact() string
}
```

### 验证中间件

#### `validate.ProtoValidate`
Protobuf 消息验证中间件

```go
func ProtoValidate(opts ...Option) middleware.Middleware
```

**配置选项：**

```go
// 自定义前端友好提示函数
validate.WithFriendlyMsg(func(err error) string {
    return "参数不合法"
})
```

**特性：**
- 支持 `protovalidate` 验证
- 支持旧版 `Validate()` 方法验证
- 验证失败返回 `BadRequest` 错误

---

## 🛠️ Utilx - 工具函数

### JSON 处理

```go
// 任意类型转 JSON 字符串
utilx.ToJson(v any) string

// 任意类型转 []byte
utilx.ToJsonBytes(v any) []byte

// 字符串转对象
utilx.ToObj(b string, v any)

// []byte 转对象
utilx.BytesToObj(b []byte, v any)
```

### 字符串转换（零拷贝）

```go
// string 转 []byte（只读，零拷贝）
utilx.StringToBytes(s string) []byte

// []byte 转 string（只读，零拷贝）
utilx.BytesToString(b []byte) string
```

**⚠️ 注意：** 使用 unsafe 实现，仅适用于只读、短生命周期场景

### 时间工具

```go
// 获取当前时间戳（毫秒）
utilx.Now() int64
```

### 切片处理

```go
// 按批次处理切片
utilx.Process(data []T, batchSize int, fn func([]T) error) error

// 取前 N 个元素
utilx.Take(s []T, maxLen int) []T
```

**示例：**

```go
// 批量处理
err := utilx.Process(items, 100, func(batch []Item) error {
    return db.BatchInsert(batch)
})

// 取前 10 个
top10 := utilx.Take(items, 10)
```

### 集合操作

```go
// 创建集合
s := utilx.NewSet[int](1, 2, 3)

// 添加元素
s.Add(4)

// 删除元素
s.Remove(1)

// 检查存在
s.Contains(2) // true

// 获取长度
s.Len()

// 清空集合
s.Clear()

// 转切片
s.ToSlice()

// 并集
s1.Union(s2)

// 交集
s1.Intersect(s2)

// 差集
s1.Difference(s2)
```

### 重试机制

```go
func Retry(ctx context.Context, maxRetries int, delay time.Duration, fn func() error) error
```

**示例：**

```go
err := utilx.Retry(ctx, 3, time.Second, func() error {
    return api.Call()
})
```

**特性：**
- 支持 context 取消
- 支持固定延迟
- 返回最后一次错误

### MD5 编码

```go
func Md5Encode(val string, salt ...string) string
```

**示例：**

```go
hash := utilx.Md5Encode("hello")              // 无 salt
hash := utilx.Md5Encode("hello", "a", "b")   // 多个 salt
```

**特性：**
- 空字符串直接返回空字符串
- 支持多个 salt 追加

### 路径处理

```go
// 查找项目根目录（查找 go.mod）
utilx.FindProjectRoot() string

// 连接路径（支持任意基本类型）
utilx.Join(sep string, vals ...any) string
```

**示例：**

```go
root := utilx.FindProjectRoot()
path := utilx.Join("/", "api", "v1", "users")  // "api/v1/users"
```

### 并发处理

#### `Parallel`
并发执行并等待所有完成，返回第一个错误

```go
func Parallel(handlers ...func() error) error
```

**示例：**

```go
err := utilx.Parallel(
    func() error { return task1() },
    func() error { return task2() },
    func() error { return task3() },
)
```

**特性：**
- 捕获 panic 并转换为 error
- 返回第一个错误
- 所有任务都会执行完成

#### `ParallelSafe`
并发执行，忽略所有错误和 panic

```go
func ParallelSafe(handlers ...func())
```

**适用场景：**
- 资源清理
- 日志上报
- 监控/埋点
- 后台通知

---

## 📋 Constants - 常量

### 环境变量

```go
constants.AppEnv          // "APP_ENV"
constants.ProjectName     // "PROJECT_NAME"
constants.ProjectRef      // "PROJECT_REF"
constants.ProjectSha      // "PROJECT_SHA"
```

### 环境枚举

```go
constants.EnvProduction   // "production"
constants.EnvStaging      // "staging"
constants.EnvDevelopment  // "development"
constants.EnvLocal        // "local"
```

### 环境判断

```go
constants.CurrentEnv() string
constants.IsProduction() bool
constants.IsStaging() bool
constants.IsDevelopment() bool
constants.IsLocal() bool
```

---

## 📝 使用示例

### 完整应用启动示例

```go
package main

import (
    "kratos-layout/internal/conf"
    "github.com/go-kratos/kratos/v2"
    "github.com/go-kratos/kratos/v2/log"
    "github.com/go-kratos/kratos/v2/transport/grpc"
    "github.com/go-kratos/kratos/v2/transport/http"
    "github.com/lhlyu/kratos-easy/bootstrap"
    "github.com/lhlyu/kratos-easy/httpx"
    "github.com/lhlyu/kratos-easy/middlewares/logging"
    "github.com/lhlyu/kratos-easy/middlewares/validate"
    _ "go.uber.org/automaxprocs"
)

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
    return bootstrap.NewApp(logger, gs, hs)
}

func wireApp(cfg *conf.Conf, logger log.Logger) (*kratos.App, func(), error) {
    // 创建 HTTP 服务器
    hs := http.NewServer(
        http.Address(cfg.Server.Http.Addr),
        http.Middleware(
            logging.Server(logger),
            validate.ProtoValidate(),
        ),
        http.ResponseEncoder(httpx.EncodeResponse),
        http.ErrorEncoder(httpx.EncodeError),
    )
    
    // 创建 gRPC 服务器
    gs := grpc.NewServer(
        grpc.Address(cfg.Server.Grpc.Addr),
        grpc.Middleware(
            logging.Server(logger),
            validate.ProtoValidate(),
        ),
    )
    
    app := newApp(logger, gs, hs)
    return app, func() {}, nil
}

func main() {
    var bc conf.Conf
    bootstrap.Run(&bc, wireApp,
        bootstrap.WithJsonFormat(),
        bootstrap.WithFilterLevel(log.LevelInfo),
    )
}
```

### 工具函数使用示例

```go
// JSON 处理
data := map[string]string{"key": "value"}
jsonStr := utilx.ToJson(data)

// 集合操作
s1 := utilx.NewSet(1, 2, 3)
s2 := utilx.NewSet(3, 4, 5)
union := s1.Union(s2)        // {1, 2, 3, 4, 5}
intersect := s1.Intersect(s2) // {3}

// 重试
err := utilx.Retry(ctx, 3, time.Second, func() error {
    return db.Query()
})

// 并发执行
err := utilx.Parallel(
    func() error { return task1() },
    func() error { return task2() },
)
```

---

## 🔗 相关链接

- [go-kratos](https://github.com/go-kratos/kratos)
- [使用示例](https://github.com/lhlyu/kratos-layout)

---

## 📄 License

详见 [LICENSE](LICENSE) 文件

