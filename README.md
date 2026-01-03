# Kratos Easy

简化 [go-kratos](https://github.com/go-kratos/kratos) 使用

## 使用

[例子](https://github.com/lhlyu/kratos-layout)

- 获取依赖

```
go get -u github.com/lhlyu/kratos-easy@latest
```

- `cmd/main.go` 使用

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
	return bootstrap.NewApp(
		logger,
		gs,
		hs,
	)
}

func main() {
	var bc conf.Conf

	bootstrap.Run(&bc, wireApp)
}
```