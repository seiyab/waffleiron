package waffleiron

import (
	"github.com/pkg/errors"
)

func Parse[T any](str string, psr Parser[T]) (T, error) {
	r := NewReader(str)
	value, err := psr.Parse(r)
	if err != nil {
		return value, err
	}
	if r.idx != int64(len(r.str)) {
		return value, errors.Errorf("parser stopped at %s", r.pos)
	}
	return value, nil
}

type Parser[T any] interface {
	Parse(r *Reader) (T, error)
}

func Map[T, U any](p Parser[T], f func(t T) U) Parser[U] {
	return mapParser[T, U]{p, f}
}

type mapParser[T, U any] struct {
	p Parser[T]
	f func(t T) U
}

func (p mapParser[T, U]) Parse(r *Reader) (U, error) {
	t, err := p.p.Parse(r)
	if err != nil {
		return *new(U), err
	}
	return p.f(t), nil
}
