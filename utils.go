package waffleiron

// Tuple2 is a tuple with 2 elements
type Tuple2[T, U any] struct {
	v0 T
	v1 U
}

// NewTuple2 makes Tuple2
func NewTuple2[T, U any](v0 T, v1 U) Tuple2[T, U] {
	return Tuple2[T, U]{v0, v1}
}

// Get0 returns 0-th element
func (t Tuple2[T, U]) Get0() T {
	return t.v0
}

// Get1 returns 1-st element
func (t Tuple2[T, U]) Get1() U {
	return t.v1
}

// Tuple3 is a tuple with 3 elements
type Tuple3[T, U, V any] struct {
	Tuple2[T, U]
	v2 V
}

// NewTuple3 makes Tuple3
func NewTuple3[T, U, V any](v0 T, v1 U, v2 V) Tuple3[T, U, V] {
	return Tuple3[T, U, V]{NewTuple2(v0, v1), v2}
}

// Get2 returns 2-nd element
func (t Tuple3[T, U, V]) Get2() V {
	return t.v2
}
