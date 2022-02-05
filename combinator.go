package waffleiron

import (
	"github.com/pkg/errors"
)

func And[T, U any](p1 Parser[T], p2 Parser[U]) Parser[Pair[T, U]] {
	return FuncParser[Pair[T, U]](func(rd *Reader) (Pair[T, U], error) {
		a, err := p1.Parse(rd)
		if err != nil {
			return Pair[T, U]{}, err
		}
		b, err := p2.Parse(rd)
		if err != nil {
			return Pair[T, U]{}, err
		}
		return Pair[T, U]{A: a, B: b}, nil
	})
}

func Trace[T any](name string, p Parser[T]) Parser[T] {
	return FuncParser[T](func(rd *Reader) (T, error) {
		var t T
		var err error
		rd.WithTrace(name, func() {
			t, err = p.Parse(rd)
		})
		if err != nil {
			return t, errors.Wrapf(err, "at %q", name)
		}
		return t, nil
	})
}
