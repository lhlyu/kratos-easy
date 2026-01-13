package utilx

import "iter"

// Set 是一个通用的集合类型，支持任何可比较类型 T。
// 内部基于 map[T]struct{} 实现，支持常用集合操作。
type Set[T comparable] struct {
	data map[T]struct{}
}

// NewSet 创建一个新的 Set，可以选择性传入初始元素。
func NewSet[T comparable](items ...T) *Set[T] {
	s := &Set[T]{data: make(map[T]struct{}, len(items))}
	for _, item := range items {
		s.data[item] = struct{}{}
	}
	return s
}

// Add 向集合中添加一个或多个元素。
func (s *Set[T]) Add(items ...T) {
	for _, item := range items {
		s.data[item] = struct{}{}
	}
}

// Remove 从集合中删除一个或多个元素。
func (s *Set[T]) Remove(items ...T) {
	for _, item := range items {
		delete(s.data, item)
	}
}

// Contains 检查集合中是否存在元素。
func (s *Set[T]) Contains(item T) bool {
	_, ok := s.data[item]
	return ok
}

// Len 返回集合中元素的数量。
func (s *Set[T]) Len() int {
	return len(s.data)
}

// IsEmpty 判断集合是否为空。
func (s *Set[T]) IsEmpty() bool {
	return len(s.data) == 0
}

// Clear 清空集合。
func (s *Set[T]) Clear() {
	s.data = make(map[T]struct{})
}

// ToSlice 将集合转换为切片。
func (s *Set[T]) ToSlice() []T {
	result := make([]T, 0, len(s.data))
	for k := range s.data {
		result = append(result, k)
	}
	return result
}

// Clone 返回集合的一个副本。
func (s *Set[T]) Clone() *Set[T] {
	result := NewSet[T]()
	for k := range s.data {
		result.data[k] = struct{}{}
	}
	return result
}

// Union 返回两个集合的并集。
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for k := range s.data {
		result.Add(k)
	}
	for k := range other.data {
		result.Add(k)
	}
	return result
}

// Intersect 返回两个集合的交集。
func (s *Set[T]) Intersect(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for k := range s.data {
		if other.Contains(k) {
			result.Add(k)
		}
	}
	return result
}

// Difference 返回集合 s 相对于 other 的差集。
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for k := range s.data {
		if !other.Contains(k) {
			result.Add(k)
		}
	}
	return result
}

// All 返回一个迭代器，用于遍历集合中的所有元素。
// 支持 Go 1.23+ 的 range-over-func 特性。
func (s *Set[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for k := range s.data {
			if !yield(k) {
				return
			}
		}
	}
}
