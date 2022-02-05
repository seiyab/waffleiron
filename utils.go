package waffleiron

type Empty struct{}

type Pair[T, U any] struct {
	A T
	B U
}
