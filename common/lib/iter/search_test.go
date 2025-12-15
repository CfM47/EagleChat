package iter

import (
	"strconv"
	"testing"
)

func TestContains(t *testing.T) {
	slice := []int{1, 2, 3, 4}
	iter := FromSlice(slice)
	if !Contains(iter, 3) {
		t.Error("Contains should find 3 in the slice")
	}

	iter = FromSlice(slice)
	if Contains(iter, 5) {
		t.Error("Contains should not find 5 in the slice")
	}
}

func TestFind(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	isEven := func(n int) bool { return n%2 == 0 }

	iter := FromSlice(slice)
	index, val, found := Find(iter, isEven)
	if !found || index != 1 || val != 2 {
		t.Errorf("Find failed, got: index=%d, val=%d, found=%v", index, val, found)
	}

	isGreaterThan5 := func(n int) bool { return n > 5 }
	iter = FromSlice(slice)
	_, _, found = Find(iter, isGreaterThan5)
	if found {
		t.Error("Find should not find a value greater than 5")
	}
}

func TestFindMap(t *testing.T) {
	slice := []string{"a", "b", "123", "d"}
	iter := FromSlice(slice)
	val, found := FindMap(iter, func(s string) (int, bool) {
		if n, err := strconv.Atoi(s); err == nil {
			return n, true
		}
		return 0, false
	})

	if !found || val != 123 {
		t.Errorf("FindMap failed, got: val=%d, found=%v", val, found)
	}

	slice = []string{"a", "b", "c"}
	iter = FromSlice(slice)
	_, found = FindMap(iter, func(s string) (int, bool) {
		if n, err := strconv.Atoi(s); err == nil {
			return n, true
		}
		return 0, false
	})
	if found {
		t.Error("FindMap should not find a number in a slice of strings")
	}
}

func TestAny(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	isEven := func(n int) bool { return n%2 == 0 }
	iter := FromSlice(slice)
	if !Any(iter, isEven) {
		t.Error("Any should find an even number")
	}

	isGreaterThan5 := func(n int) bool { return n > 5 }
	iter = FromSlice(slice)
	if Any(iter, isGreaterThan5) {
		t.Error("Any should not find a number greater than 5")
	}
}

func TestAll(t *testing.T) {
	slice := []int{2, 4, 6, 8}
	isEven := func(n int) bool { return n%2 == 0 }
	iter := FromSlice(slice)
	if !All(iter, isEven) {
		t.Error("All should confirm all numbers are even")
	}

	slice = []int{2, 4, 5, 8}
	iter = FromSlice(slice)
	if All(iter, isEven) {
		t.Error("All should find that not all numbers are even")
	}
}
