package waffleiron

import (
	"github.com/pkg/errors"

	"github.com/hashicorp/go-multierror"
)

func And[T, U any](p1 Parser[T], p2 Parser[U]) Parser[Pair[T, U]] {
	return FuncParser[Pair[T, U]](func(r *Reader) (Pair[T, U], error) {
		a, err := p1.Parse(r)
		if err != nil {
			return Pair[T, U]{}, err
		}
		b, err := p2.Parse(r)
		if err != nil {
			return Pair[T, U]{}, err
		}
		return Pair[T, U]{A: a, B: b}, nil
	})
}

func Or[T any](p0 Parser[T], ps ...Parser[T]) Parser[T] {
	parsers := make([]Parser[T], len(ps)+1)
	parsers[0] = p0
	for i, p := range ps {
		parsers[i+1] = p
	}
	return OrParser[T]{ps: parsers}
}

type OrParser[T any] struct {
	ps []Parser[T]
}

func (p OrParser[T]) Parse(r *Reader) (T, error) {
	var totalErr error
	for _, p := range p.ps {
		var t T
		err := r.Try(func() error {
			var e error
			t, e = p.Parse(r)
			return e
		})
		if err == nil {
			return t, nil
		}

		totalErr = multierror.Append(totalErr, err)
	}
	return *new(T), totalErr
}

func Trace[T any](name string, p Parser[T]) Parser[T] {
	return FuncParser[T](func(r *Reader) (T, error) {
		var t T
		var err error
		r.WithTrace(name, func() {
			t, err = p.Parse(r)
		})
		if err != nil {
			return t, errors.Wrapf(err, "at %q", name)
		}
		return t, nil
	})
}
