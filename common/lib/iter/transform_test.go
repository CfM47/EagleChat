package iter

import (
	"reflect"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	slice := []int{1, 2, 3}
	iter := FromSlice(slice)
	mapIter := Map(iter, func(n int) string { return strconv.Itoa(n) })
	result := ToSlice(mapIter)
	expected := []string{"1", "2", "3"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Map result = %v; want %v", result, expected)
	}
}

func TestFilter(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6}
	iter := FromSlice(slice)
	filterIter := Filter(iter, func(n int) bool { return n%2 == 0 })
	result := ToSlice(filterIter)
	expected := []int{2, 4, 6}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Filter result = %v; want %v", result, expected)
	}
}

func TestFilterMap(t *testing.T) {
	slice := []string{"1", "two", "3", "four"}
	iter := FromSlice(slice)
	filterMapIter := FilterMap(iter, func(s string) (int, bool) {
		n, err := strconv.Atoi(s)
		return n, err == nil
	})
	result := ToSlice(filterMapIter)
	expected := []int{1, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FilterMap result = %v; want %v", result, expected)
	}
}

func TestFlatMap(t *testing.T) {
	slice := [][]int{{1, 2}, {3, 4}}
	iter := FromSlice(slice)
	flatMapIter := FlatMap(iter, func(s []int) Iter[int] { return FromSlice(s) })
	result := ToSlice(flatMapIter)
	expected := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FlatMap result = %v; want %v", result, expected)
	}
}
