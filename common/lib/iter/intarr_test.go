package iter

import "testing"

func TestSum(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := FromSlice(slice)
	sum := Sum(iter)
	if sum != 15 {
		t.Errorf("Sum = %d; want 15", sum)
	}

	fSlice := []float64{1.1, 2.2, 3.3}
	fIter := FromSlice(fSlice)
	fSum := Sum(fIter)
	if fSum != 6.6 {
		t.Errorf("Sum = %f; want 6.6", fSum)
	}
}

func TestProduct(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	iter := FromSlice(slice)
	product := Product(iter)
	if product != 120 {
		t.Errorf("Product = %d; want 120", product)
	}
}

func TestMin(t *testing.T) {
	slice := []int{3, 1, 4, 1, 5, 9}
	iter := FromSlice(slice)
	minVal, ok := Min(iter)
	if !ok || minVal != 1 {
		t.Errorf("Min = %d, ok = %v; want 1, true", minVal, ok)
	}

	fSlice := []float64{3.14, 1.618, 2.718}
	fIter := FromSlice(fSlice)
	fMinVal, ok := Min(fIter)
	if !ok || fMinVal != 1.618 {
		t.Errorf("Min = %f, ok = %v; want 1.618, true", fMinVal, ok)
	}
}

func TestMax(t *testing.T) {
	slice := []int{3, 1, 4, 1, 5, 9}
	iter := FromSlice(slice)
	maxVal, ok := Max(iter)
	if !ok || maxVal != 9 {
		t.Errorf("Max = %d, ok = %v; want 9, true", maxVal, ok)
	}

	fSlice := []float64{3.14, 1.618, 2.718}
	fIter := FromSlice(fSlice)
	fMaxVal, ok := Max(fIter)
	if !ok || fMaxVal != 3.14 {
		t.Errorf("Max = %f, ok = %v; want 3.14, true", fMaxVal, ok)
	}
}

func TestMin_Empty(t *testing.T) {
	iter := FromSlice([]int{})
	_, ok := Min(iter)
	if ok {
		t.Errorf("Min of empty slice should return ok = false")
	}
}

func TestMax_Empty(t *testing.T) {
	iter := FromSlice([]int{})
	_, ok := Max(iter)
	if ok {
		t.Errorf("Max of empty slice should return ok = false")
	}
}
