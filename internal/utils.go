package internal

import "cmp"

// Ternary returns `a` if `cond` is true, otherwise returns `b`.
// This is a generic utility function that mimics the ternary operator found in other languages.
// The type parameter T must be cmp.Ordered, which allows comparison operations.
//
// Example usage:
//
//	max := Ternary(x > y, x, y)
//
// Parameters:
//
//	cond - the condition to evaluate
//	a    - the value to return if cond is true
//	b    - the value to return if cond is false
func Ternary[T cmp.Ordered](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}
