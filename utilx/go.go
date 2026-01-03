package utilx

import (
	"fmt"
	"runtime"
	"sync"
)

const panicBufSize = 64 << 10 // 64KB

// Parallel 并发执行所有传入的处理函数，并等待它们全部完成。
//
// 如果任意处理函数返回非 nil 的 error，或在执行过程中发生 panic，
// Parallel 会返回第一个捕获到的错误。
// 当存在多个错误或 panic 时，仅返回最先发生的一个。
//
// 如果未传入任何处理函数，则直接返回 nil。
//
// 注意：
//   - 所有处理函数都会被执行并等待完成
//   - panic 会被 recover，并转换为包含堆栈信息的 error 返回
func Parallel(handlers ...func() error) error {
	if len(handlers) == 0 {
		return nil
	}

	var (
		wg   sync.WaitGroup
		once sync.Once
		err  error
	)

	wg.Add(len(handlers))

	for _, h := range handlers {
		go func(handler func() error) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					buf := make([]byte, panicBufSize)
					n := runtime.Stack(buf, false)
					once.Do(func() {
						err = fmt.Errorf(
							"panic in parallel handler: %v\n%s",
							r,
							buf[:n],
						)
					})
				}
			}()

			if e := handler(); e != nil {
				once.Do(func() {
					err = e
				})
			}
		}(h)
	}

	wg.Wait()
	return err
}

// ParallelSafe 并发执行所有传入的处理函数，并等待它们全部完成。
//
// 与 Parallel 不同，ParallelSafe 会忽略所有处理函数返回的 error，
// 并且会吞掉执行过程中发生的 panic，保证不会向外传播任何错误。
//
// 该方法适用于“尽力而为”的场景，例如：
//   - 资源清理
//   - 日志上报
//   - 监控/埋点
//   - 后台通知
//
// 如果未传入任何处理函数，该方法会直接返回。
func ParallelSafe(handlers ...func()) {
	if len(handlers) == 0 {
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(handlers))

	for _, h := range handlers {
		go func(handler func()) {
			defer wg.Done()
			defer func() {
				_ = recover() // 明确忽略 panic
			}()
			handler()
		}(h)
	}

	wg.Wait()
}
