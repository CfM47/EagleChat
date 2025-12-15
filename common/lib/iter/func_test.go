package iter

import (
	"strconv"
	"strings"
	"testing"
)

func TestCompose(t *testing.T) {
	add2 := func(x int) int { return x + 2 }
	intToString := func(x int) string { return strconv.Itoa(x) }

	composed := Compose(add2, intToString)

	if composed(3) != "5" {
		t.Errorf("Compose(add2, intToString)(3) = %s; want 5", composed(3))
	}
}

func TestPipe(t *testing.T) {
	add2 := func(x int) int { return x + 2 }
	multiplyBy3 := func(x int) int { return x * 3 }
	subtract1 := func(x int) int { return x - 1 }

	result := Pipe(5, add2, multiplyBy3, subtract1)

	if result != 20 {
		t.Errorf("Pipe(5, add2, multiplyBy3, subtract1) = %d; want 20", result)
	}

	// Test with a different type
	toUpper := func(s string) string { return strings.ToUpper(s) }
	addWorld := func(s string) string { return s + " WORLD" }

	strResult := Pipe("hello", toUpper, addWorld)
	if strResult != "HELLO WORLD" {
		t.Errorf("Pipe(\"hello\", toUpper, addWorld) = %s; want HELLO WORLD", strResult)
	}
}

