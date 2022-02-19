package waffleiron

import (
	"github.com/pkg/errors"

	"github.com/hashicorp/go-multierror"
)

// And returns a parser that runs p0 and then runs p1
func And[T, U any](p0 Parser[T], p1 Parser[U]) Parser[Tuple2[T, U]] {
	return andParser[T, U]{p0, p1}
}

type andParser[T, U any] struct {
	p0 Parser[T]
	p1 Parser[U]
}

// Parse implements Parser interface
func (p andParser[T, U]) Parse(r *Reader) (Tuple2[T, U], error) {
	a, err := p.p0.Parse(r)
	if err != nil {
		return Tuple2[T, U]{}, err
	}
	b, err := p.p1.Parse(r)
	if err != nil {
		return Tuple2[T, U]{}, err
	}
	return NewTuple2(a, b), nil
}

// And returns a parser that runs p0 and then runs p1 and then p2
func And3[T, U, V any](p0 Parser[T], p1 Parser[U], p2 Parser[V]) Parser[Tuple3[T, U, V]] {
	return and3Parser[T, U, V]{p0, p1, p2}
}

type and3Parser[T, U, V any] struct {
	p0 Parser[T]
	p1 Parser[U]
	p2 Parser[V]
}

// Parse implements Parser interface
func (p and3Parser[T, U, V]) Parse(r *Reader) (Tuple3[T, U, V], error) {
	v0, err := p.p0.Parse(r)
	if err != nil {
		return Tuple3[T, U, V]{}, err
	}
	v1, err := p.p1.Parse(r)
	if err != nil {
		return Tuple3[T, U, V]{}, err
	}
	v2, err := p.p2.Parse(r)
	if err != nil {
		return Tuple3[T, U, V]{}, err
	}
	return NewTuple3(v0, v1, v2), nil
}

// Choice returns a parser that tries each parsers and returns first successful result
func Choice[T any](p0 Parser[T], ps ...Parser[T]) Parser[T] {
	parsers := make([]Parser[T], len(ps)+1)
	parsers[0] = p0
	for i, p := range ps {
		parsers[i+1] = p
	}
	return choiceParser[T]{ps: parsers}
}

type choiceParser[T any] struct {
	ps []Parser[T]
}

// Parse implements Parser interface
func (p choiceParser[T]) Parse(r *Reader) (T, error) {
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

// Repeat returns a parser that repeatedly tries p until it fails
// If first trial fails, it results empty slice without any errors
func Repeat[T any](p Parser[T]) Parser[[]T] {
	return repeatParser[T]{p}
}

type repeatParser[T any] struct {
	p Parser[T]
}

// Parse implements Parser interface
func (p repeatParser[T]) Parse(r *Reader) ([]T, error) {
	ts := make([]T, 0)
	for {
		err := r.Try(func() error {
			t, e := p.p.Parse(r)
			if e != nil {
				return e
			}
			ts = append(ts, t)
			return nil
		})
		if err != nil {
			return ts, nil
		}
	}
}

// SepBy returns a parser that repeatedly parse with p separated by sep
func SepBy[T, U any](p Parser[T], sep Parser[U]) Parser[[]T] {
	return Choice(
		Map(
			And(
				Repeat(And(p, sep)),
				p,
			),
			func(v Tuple2[[]Tuple2[T, U], T]) []T {
				ts := make([]T, 0)
				for _, t := range v.Get0() {
					ts = append(ts, t.Get0())
				}
				return append(ts, v.Get1())
			},
		),
		Pure[[]T](nil),
	)
}

// Maybe returns a parser that tries to run p
// If p fails, it results nil without any errors
func Maybe[T any](p Parser[T]) Parser[*T] {
	return maybeParser[T]{p}
}

type maybeParser[T any] struct {
	p Parser[T]
}

// Parse implements Parser interface
func (p maybeParser[T]) Parse(r *Reader) (*T, error) {
	var v T
	err := r.Try(func() error {
		var e error
		v, e = p.p.Parse(r)
		return e
	})
	if err != nil {
		return nil, nil
	}
	return &v, nil
}

// Between returns a parser that runs open and then runs p and then runs close
// The results of open and close are discarded
func Between[T, U, V any](open Parser[T], p Parser[U], close Parser[V]) Parser[U] {
	return Map(
		And3(
			open,
			p,
			close,
		),
		func(t Tuple3[T, U, V]) U {
			return t.Get1()
		},
	)
}

func Untype[T any](p Parser[T]) Parser[interface{}] {
	return untypeParser[T]{p}
}

type untypeParser[T any] struct {
	p Parser[T]
}

func (p untypeParser[T]) Parse(r *Reader) (interface{}, error) {
	return p.p.Parse(r)
}

// Ref takes a pointer of a parser and returns a parser that works same as original one
// It can be useful for recursive parser
func Ref[T any](p *Parser[T]) Parser[T] {
	return refParser[T]{p}
}

type refParser[T any] struct {
	p *Parser[T]
}

func (p refParser[T]) Parse(r *Reader) (T, error) {
	a := *p.p
	return a.Parse(r)
}

// Trace adds name for debug
func Trace[T any](name string, p Parser[T]) Parser[T] {
	return traceParser[T]{name, p}
}

type traceParser[T any] struct {
	name string
	p    Parser[T]
}

// Parse implements Parser interface
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
