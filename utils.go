package waffleiron

type Empty struct{}

type Tuple2[T, U any] struct {
	v0 T
	v1 U
}

func NewTuple2[T, U any](v0 T, v1 U) Tuple2[T, U] {
	return Tuple2[T, U]{v0, v1}
}

func (t Tuple2[T, U]) Get0() T {
	return t.v0
}

func (t Tuple2[T, U]) Get1() U {
	return t.v1
}

type Tuple3[T, U, V any] struct {
	Tuple2[T, U]
	v2 V
}

func NewTuple3[T, U, V any](v0 T, v1 U, v2 V) Tuple3[T, U, V] {
	return Tuple3[T, U, V]{NewTuple2(v0, v1), v2}
}

func (t Tuple3[T, U, V]) Get2() V {
	return t.v2
}
