package waffleiron

import (
	"github.com/pkg/errors"

	"github.com/hashicorp/go-multierror"
)

func And[T, U any](p0 Parser[T], p1 Parser[U]) Parser[Pair[T, U]] {
	return andParser[T, U]{p0, p1}
}

type andParser[T, U any] struct {
	p0 Parser[T]
	p1 Parser[U]
}

func (p andParser[T, U]) Parse(r *Reader) (Pair[T, U], error) {
	a, err := p.p0.Parse(r)
	if err != nil {
		return Pair[T, U]{}, err
	}
	b, err := p.p1.Parse(r)
	if err != nil {
		return Pair[T, U]{}, err
	}
	return Pair[T, U]{A: a, B: b}, nil
}

func Or[T any](p0 Parser[T], ps ...Parser[T]) Parser[T] {
	parsers := make([]Parser[T], len(ps)+1)
	parsers[0] = p0
	for i, p := range ps {
		parsers[i+1] = p
	}
	return orParser[T]{ps: parsers}
}

type orParser[T any] struct {
	ps []Parser[T]
}

func (p orParser[T]) Parse(r *Reader) (T, error) {
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
	return traceParser[T]{name, p}
}

type traceParser[T any] struct {
	name string
	p    Parser[T]
}

func (p traceParser[T]) Parse(r *Reader) (T, error) {
	var t T
	var err error
	r.WithTrace(p.name, func() {
		t, err = p.p.Parse(r)
	})
	if err != nil {
		return t, errors.Wrapf(err, "at %q", p.name)
	}
	return t, nil
}
