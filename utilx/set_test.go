package utilx

import (
	"testing"
)

func TestSet_IsEmpty(t *testing.T) {
	s := NewSet[int]()
	if !s.IsEmpty() {
		t.Errorf("expected empty set, got %d items", s.Len())
	}
	s.Add(1)
	if s.IsEmpty() {
		t.Errorf("expected non-empty set after Add")
	}
}

func TestSet_AddVariadic(t *testing.T) {
	s := NewSet[int]()
	s.Add(1, 2, 3)
	if s.Len() != 3 {
		t.Errorf("expected 3 items, got %d", s.Len())
	}
	if !s.Contains(1) || !s.Contains(2) || !s.Contains(3) {
		t.Errorf("set missing items after variadic Add")
	}
}

func TestSet_RemoveVariadic(t *testing.T) {
	s := NewSet[int](1, 2, 3, 4)
	s.Remove(1, 3)
	if s.Len() != 2 {
		t.Errorf("expected 2 items after Remove, got %d", s.Len())
	}
	if s.Contains(1) || s.Contains(3) {
		t.Errorf("set should not contain removed items")
	}
	if !s.Contains(2) || !s.Contains(4) {
		t.Errorf("set missing remaining items")
	}
}

func TestSet_Clone(t *testing.T) {
	s1 := NewSet[int](1, 2, 3)
	s2 := s1.Clone()
	if s2.Len() != 3 {
		t.Errorf("cloned set length mismatch: expected 3, got %d", s2.Len())
	}
	if !s2.Contains(1) || !s2.Contains(2) || !s2.Contains(3) {
		t.Errorf("cloned set missing items")
	}
	s1.Add(4)
	if s2.Contains(4) {
		t.Errorf("cloned set should not be affected by original set modifications")
	}
}

func TestSet_All(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	s := NewSet[int](items...)

	count := 0
	visited := make(map[int]bool)
	for item := range s.All() {
		count++
		visited[item] = true
	}

	if count != len(items) {
		t.Errorf("expected to iterate %d items, got %d", len(items), count)
	}
	for _, item := range items {
		if !visited[item] {
			t.Errorf("item %d was not visited during iteration", item)
		}
	}

	// Test break in iteration
	count = 0
	for _ = range s.All() {
		count++
		if count == 2 {
			break
		}
	}
	if count != 2 {
		t.Errorf("expected to break after 2 items, got %d", count)
	}
}
