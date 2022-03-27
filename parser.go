package waffleiron

import (
	"github.com/pkg/errors"
)

// Parse runs parser
// It returns an error if psr fails or psr doesn't consume all bytes
func Parse[T any](str string, psr Parser[T]) (T, error) {
	r := newReader(str)
	value, err := psr.parse(r)
	if err != nil {
		return value, err
	}
	if r.idx != len(r.str) {
		return value, errors.Errorf("parser stopped at %s", r.locateAndString())
	}
	return value, nil
}

// Parser is a parser to pass Parse function
type Parser[T any] struct {
	p parser[T]
}

func (p Parser[T]) parse(r *reader) (T, error) {
	return p.p.parse(r)
}

// parser is a interface for parser implementations
type parser[T any] interface {
	// parse consumes r and returns a result
	// It returns error if it fails to parse
	parse(r *reader) (T, error)
}

// Map applies function for result of parse
func Map[T, U any](p Parser[T], f func(t T) U) Parser[U] {
	return Parser[U]{p: mapParser[T, U]{p, f}}
}

type mapParser[T, U any] struct {
	p parser[T]
	f func(t T) U
}

func (p mapParser[T, U]) parse(r *reader) (U, error) {
	t, err := p.p.parse(r)
	if err != nil {
		return *new(U), err
	}
	return p.f(t), nil
}
