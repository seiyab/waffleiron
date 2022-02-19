package waffleiron

import (
	"github.com/pkg/errors"

	"github.com/hashicorp/go-multierror"
)

// And returns a parser that runs p0 and then runs p1
func And[T, U any](p0 Parser[T], p1 Parser[U]) Parser[Tuple2[T, U]] {
	return Parser[Tuple2[T, U]]{p: andParser[T, U]{p0, p1}}
}

type andParser[T, U any] struct {
	p0 Parser[T]
	p1 Parser[U]
}

func (p andParser[T, U]) parse(r *reader) (Tuple2[T, U], error) {
	a, err := p.p0.parse(r)
	if err != nil {
		return Tuple2[T, U]{}, err
	}
	b, err := p.p1.parse(r)
	if err != nil {
		return Tuple2[T, U]{}, err
	}
	return NewTuple2(a, b), nil
}

// And returns a parser that runs p0 and then runs p1 and then p2
func And3[T, U, V any](p0 Parser[T], p1 Parser[U], p2 Parser[V]) Parser[Tuple3[T, U, V]] {
	return Parser[Tuple3[T, U, V]]{p: and3Parser[T, U, V]{p0, p1, p2}}
}

type and3Parser[T, U, V any] struct {
	p0 Parser[T]
	p1 Parser[U]
	p2 Parser[V]
}

func (p and3Parser[T, U, V]) parse(r *reader) (Tuple3[T, U, V], error) {
	v0, err := p.p0.parse(r)
	if err != nil {
		return Tuple3[T, U, V]{}, err
	}
	v1, err := p.p1.parse(r)
	if err != nil {
		return Tuple3[T, U, V]{}, err
	}
	v2, err := p.p2.parse(r)
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
	return Parser[T]{p: choiceParser[T]{ps: parsers}}
}

type choiceParser[T any] struct {
	ps []Parser[T]
}

func (p choiceParser[T]) parse(r *reader) (T, error) {
	var totalErr error
	for _, p := range p.ps {
		var t T
		err := r.try(func() error {
			var e error
			t, e = p.parse(r)
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
	return Parser[[]T]{p: repeatParser[T]{p}}
}

type repeatParser[T any] struct {
	p Parser[T]
}

func (p repeatParser[T]) parse(r *reader) ([]T, error) {
	ts := make([]T, 0)
	for {
		err := r.try(func() error {
			t, e := p.p.parse(r)
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
	return Parser[*T]{p: maybeParser[T]{p}}
}

type maybeParser[T any] struct {
	p Parser[T]
}

func (p maybeParser[T]) parse(r *reader) (*T, error) {
	var v T
	err := r.try(func() error {
		var e error
		v, e = p.p.parse(r)
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

func Untype[T any](p Parser[T]) Parser[any] {
	return Parser[any]{p: untypeParser[T]{p}}
}

type untypeParser[T any] struct {
	p Parser[T]
}

func (p untypeParser[T]) parse(r *reader) (interface{}, error) {
	return p.p.parse(r)
}

// Ref takes a pointer of a parser and returns a parser that works same as original one
// It can be useful for recursive parser
func Ref[T any](p *Parser[T]) Parser[T] {
	return Parser[T]{p: refParser[T]{p}}
}

type refParser[T any] struct {
	p *Parser[T]
}

func (p refParser[T]) parse(r *reader) (T, error) {
	a := *p.p
	return a.parse(r)
}

// Trace adds name for debug
func Trace[T any](name string, p Parser[T]) Parser[T] {
	return Parser[T]{p: traceParser[T]{name, p}}
}

type traceParser[T any] struct {
	name string
	p    Parser[T]
}

func (p traceParser[T]) parse(r *reader) (T, error) {
	var t T
	var err error
	r.withTrace(p.name, func() {
		t, err = p.p.parse(r)
	})
	if err != nil {
		return t, errors.Wrapf(err, "%q >", p.name)
	}
	return t, nil
}
