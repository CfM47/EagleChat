package iter

import (
	"reflect"
	"testing"
)

func TestFold(t *testing.T) {
	slice := []int{1, 2, 3, 4}
	iter := FromSlice(slice)
	sum := Fold(iter, 0, func(acc int, v int) int { return acc + v })
	if sum != 10 {
		t.Errorf("Fold sum = %d; want 10", sum)
	}

	iter = FromSlice(slice)
	product := Fold(iter, 1, func(acc int, v int) int { return acc * v })
	if product != 24 {
		t.Errorf("Fold product = %d; want 24", product)
	}
}

func TestReduce(t *testing.T) {
	slice := []int{1, 2, 3, 4}
	iter := FromSlice(slice)
	sum, ok := Reduce(iter, func(a, b int) int { return a + b })
	if !ok || sum != 10 {
		t.Errorf("Reduce sum = %d, ok = %v; want 10, true", sum, ok)
	}

	// Test with empty slice
	emptySlice := []int{}
	iter = FromSlice(emptySlice)
	_, ok = Reduce(iter, func(a, b int) int { return a + b })
	if ok {
		t.Error("Reduce on empty slice should return ok = false")
	}
}

func TestScan(t *testing.T) {
	slice := []int{1, 2, 3, 4}
	iter := FromSlice(slice)
	scanIter := Scan(iter, 0, func(acc int, v int) int { return acc + v })
	result := ToSlice(scanIter)
	expected := []int{1, 3, 6, 10}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Scan result = %v; want %v", result, expected)
	}
}
