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
