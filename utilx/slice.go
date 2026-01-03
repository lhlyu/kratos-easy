package utilx

// Process 将切片按 batchSize 条一组顺序执行 fn。
// 如果 batchSize <= 0 则不执行。
// 任一 fn 返回 error 则立即终止并把错误返回。
func Process[T any](data []T, batchSize int, fn func([]T) error) error {
	if batchSize <= 0 {
		return nil
	}
	for i := 0; i < len(data); i += batchSize {
		end := min(i+batchSize, len(data))
		if err := fn(data[i:end]); err != nil {
			return err
		}
	}
	return nil
}

// Take 从切片中取前 maxLen 个元素。
// 如果切片长度小于或等于 maxLen，则返回整个切片。
func Take[T any](s []T, maxLen int) []T {
	if maxLen < 0 {
		return nil
	}
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
