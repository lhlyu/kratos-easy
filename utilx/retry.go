package utilx

import (
	"context"
	"time"
)

// Retry 对一个函数进行重试
// ctx: 用于控制超时和取消
// maxRetries: 最大重试次数
// delay: 每次重试之间的固定延迟
// fn: 需要重试的函数，它返回一个 error
func Retry(ctx context.Context, maxRetries int, delay time.Duration, fn func() error) error {
	var lastErr error
	for i := range maxRetries {
		// 在每次尝试前检查 context 是否已取消
		select {
		case <-ctx.Done():
			return ctx.Err() // 如果 context 被取消，立即返回
		default:
		}

		// 执行函数
		err := fn()
		if err == nil {
			// 成功，直接返回
			return nil
		}
		lastErr = err

		// 如果不是最后一次尝试，则等待
		if i < maxRetries-1 {
			// 使用 time.NewTimer 而不是 time.Sleep，以便能响应 context 的取消
			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop() // 防止内存泄漏
				return ctx.Err()
			case <-timer.C:
				// 等待结束，继续下一次重试
			}
		}
	}
	// 所有重试都失败了，返回最后一次遇到的错误
	return lastErr
}
