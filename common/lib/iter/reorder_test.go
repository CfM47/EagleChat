package iter

import (
	"reflect"
	"strconv"
	"testing"
)

func TestZip(t *testing.T) {
	a := FromSlice([]int{1, 2, 3})
	b := FromSlice([]string{"a", "b", "c"})
	zipped := Zip(a, b)
	result := ToSlice(zipped)

	expected := []Tuple[int, string]{
		{1, "a"},
		{2, "b"},
		{3, "c"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Zip result = %v; want %v", result, expected)
	}

	// Test with different lengths
	a = FromSlice([]int{1, 2})
	b = FromSlice([]string{"a", "b", "c"})
	zipped = Zip(a, b)
	result = ToSlice(zipped)
	expected = []Tuple[int, string]{
		{1, "a"},
		{2, "b"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Zip with different lengths: result = %v; want %v", result, expected)
	}
}

func TestZipWith(t *testing.T) {
	a := []int{1, 2, 3}
	b := []string{"a", "b", "c"}
	result := ZipWith(a, b, func(x int, y string) string {
		return strconv.Itoa(x) + y
	})
	expected := []string{"1a", "2b", "3c"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ZipWith result = %v; want %v", result, expected)
	}
}

func TestUnzip(t *testing.T) {
	tuples := []Tuple[int, string]{
		{1, "a"},
		{2, "b"},
		{3, "c"},
	}
	a, b := Unzip(tuples)

	expectedA := []int{1, 2, 3}
	expectedB := []string{"a", "b", "c"}

	if !reflect.DeepEqual(a, expectedA) {
		t.Errorf("Unzip first slice = %v; want %v", a, expectedA)
	}
	if !reflect.DeepEqual(b, expectedB) {
		t.Errorf("Unzip second slice = %v; want %v", b, expectedB)
	}
}

func TestPartition(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6}
	isEven := func(n int) bool { return n%2 == 0 }
	evens, odds := Partition(slice, isEven)

	expectedEvens := []int{2, 4, 6}
	expectedOdds := []int{1, 3, 5}

	if !reflect.DeepEqual(evens, expectedEvens) {
		t.Errorf("Partition evens = %v; want %v", evens, expectedEvens)
	}
	if !reflect.DeepEqual(odds, expectedOdds) {
		t.Errorf("Partition odds = %v; want %v", odds, expectedOdds)
	}
}
