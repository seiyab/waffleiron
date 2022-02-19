package waffleiron

import (
	"github.com/pkg/errors"
)

// Parse runs parser
// It returns an error if psr fails or psr doesn't consume all bytes
func Parse[T any](str string, psr Parser[T]) (T, error) {
	r := newReader(str)
	value, err := psr.Parse(r)
	if err != nil {
		return value, err
	}
	if r.idx != int64(len(r.str)) {
		return value, errors.Errorf("parser stopped at %s", r.pos)
	}
	return value, nil
}

// Parser is a parser
type Parser[T any] interface {
	// Parse consumes r and returns a result
	// It returns error if it fails to parse
	Parse(r *reader) (T, error)
}

// Map applies function for result of parse
func Map[T, U any](p Parser[T], f func(t T) U) Parser[U] {
	return mapParser[T, U]{p, f}
}

type mapParser[T, U any] struct {
	p Parser[T]
	f func(t T) U
}

// Parse implements Parser interface
func (p mapParser[T, U]) Parse(r *reader) (U, error) {
	t, err := p.p.Parse(r)
	if err != nil {
		return *new(U), err
	}
	return p.f(t), nil
}
