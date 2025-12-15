package iter

import (
	"reflect"
	"testing"
)

func TestFromSlice(t *testing.T) {
	slice := []int{1, 2, 3}
	iter := FromSlice(slice)

	for i := 0; i < len(slice); i++ {
		val, ok := iter.next()
		if !ok {
			t.Fatalf("Expected next to return true, but got false at index %d", i)
		}
		if val != slice[i] {
			t.Errorf("Expected value %d, but got %d at index %d", slice[i], val, i)
		}
	}

	_, ok := iter.next()
	if ok {
		t.Error("Expected next to return false at the end of the slice")
	}
}

func TestToSlice(t *testing.T) {
	slice := []int{1, 2, 3}
	iter := FromSlice(slice)
	result := ToSlice(iter)

	if !reflect.DeepEqual(slice, result) {
		t.Errorf("Expected slice %v, but got %v", slice, result)
	}

	// Test with empty slice
	emptySlice := []int{}
	iter = FromSlice(emptySlice)
	result = ToSlice(iter)
	if len(result) != 0 {
		t.Errorf("Expected empty slice, but got %v", result)
	}
}
