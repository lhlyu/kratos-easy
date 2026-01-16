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

// Prepend 在切片 s 前面插入 elems，并返回新切片。
// 不会修改原切片。
// 当 elems 为空时，直接返回 s。
func Prepend[T any](s []T, elems ...T) []T {
	if len(elems) == 0 {
		return s
	}

	n := len(elems) + len(s)
	out := make([]T, n)

	// 拷贝前置元素
	copy(out, elems)
	// 拷贝原切片
	copy(out[len(elems):], s)

	return out
}
