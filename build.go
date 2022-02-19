package waffleiron

/* Builder 5 is not implemented yet
type Builder5[T5, T4, T3, T2, T1, U any] struct {
	b builtParser[T5, T4, T3, T2, T1, U]
}
*/

type Builder3[T5, T4, T3, T2, T1, U any] struct {
	b builtParser[T5, T4, T3, T2, T1, U]
}

func Begin3[T3, T2, T1, U any](f func(v3 T3, v2 T2, v1 T1) U) Builder3[any, any, T3, T2, T1, U] {
	return Builder3[any, any, T3, T2, T1, U]{
		b: builtParser[any, any, T3, T2, T1, U]{
			f:     func(_, _ any, v3 T3, v2 T2, v1 T1) U { return f(v3, v2, v1) },
			pSkip: zeroParser,
			p5:    zeroParser,
			p4:    zeroParser,
			p3:    Parser[T3]{},
			p2:    Parser[T2]{},
			p1:    Parser[T1]{},
		},
	}
}

func (b Builder3[T5, T4, T3, T2, T1, U]) Skip(p Parser[any]) Builder3[T5, T4, T3, T2, T1, U] {
	b.b.p4 = Map(And(b.b.p4, p), discardRight[T4])
	return Builder3[T5, T4, T3, T2, T1, U]{b: b.b}
}

func (b Builder3[T5, T4, T3, T2, T1, U]) Then(p Parser[T3]) Builder2[T5, T4, T3, T2, T1, U] {
	b.b.p3 = p
	return Builder2[T5, T4, T3, T2, T1, U]{b: b.b}
}

type Builder2[T5, T4, T3, T2, T1, U any] struct {
	b builtParser[T5, T4, T3, T2, T1, U]
}

func Begin2[T2, T1, U any](f func(v2 T2, v1 T1) U) Builder2[any, any, any, T2, T1, U] {
	return Builder2[any, any, any, T2, T1, U]{
		b: builtParser[any, any, any, T2, T1, U]{
			f:     func(_, _, _ any, v2 T2, v1 T1) U { return f(v2, v1) },
			pSkip: zeroParser,
			p5:    zeroParser,
			p4:    zeroParser,
			p3:    zeroParser,
			p2:    Parser[T2]{},
			p1:    Parser[T1]{},
		},
	}
}

func (b Builder2[T5, T4, T3, T2, T1, U]) Skip(p Parser[any]) Builder2[T5, T4, T3, T2, T1, U] {
	b.b.p3 = Map(And(b.b.p3, p), discardRight[T3])
	return Builder2[T5, T4, T3, T2, T1, U]{b: b.b}
}

func (b Builder2[T5, T4, T3, T2, T1, U]) Then(p Parser[T2]) Builder1[T5, T4, T3, T2, T1, U] {
	b.b.p2 = p
	return Builder1[T5, T4, T3, T2, T1, U]{b: b.b}
}

type Builder1[T5, T4, T3, T2, T1, U any] struct {
	b builtParser[T5, T4, T3, T2, T1, U]
}

func Begin1[T, U any](f func(t T) U) Builder1[any, any, any, any, T, U] {
	return Builder1[any, any, any, any, T, U]{
		b: builtParser[any, any, any, any, T, U]{
			f:     func(_, _, _, _ any, t T) U { return f(t) },
			pSkip: zeroParser,
			p5:    zeroParser,
			p4:    zeroParser,
			p3:    zeroParser,
			p2:    zeroParser,
			p1:    Parser[T]{},
		},
	}
}

func (b Builder1[T5, T4, T3, T2, T1, U]) Skip(p Parser[any]) Builder1[T5, T4, T3, T2, T1, U] {
	b.b.p2 = Map(And(b.b.p2, p), discardRight[T2])
	return b
}

func (b Builder1[T5, T4, T3, T2, T1, U]) Then(p Parser[T1]) Builder0[T5, T4, T3, T2, T1, U] {
	b.b.p1 = p
	return Builder0[T5, T4, T3, T2, T1, U]{b: b.b}
}

type Builder0[T5, T4, T3, T2, T1, U any] struct {
	b builtParser[T5, T4, T3, T2, T1, U]
}

func (b Builder0[T5, T4, T3, T2, T1, U]) Skip(p Parser[any]) Builder0[T5, T4, T3, T2, T1, U] {
	b.b.p1 = Map(And(b.b.p1, p), func(a Tuple2[T1, any]) T1 {
		return a.Get0()
	})
	return b
}

func (b Builder0[T5, T4, T3, T2, T1, U]) End() Parser[U] {
	return Parser[U]{p: b.b}
}

type builtParser[T5, T4, T3, T2, T1, U any] struct {
	f     func(v5 T5, v4 T4, v3 T3, v2 T2, v1 T1) U
	pSkip Parser[any]
	p5    Parser[T5]
	p4    Parser[T4]
	p3    Parser[T3]
	p2    Parser[T2]
	p1    Parser[T1]
}

func (p builtParser[T5, T4, T3, T2, T1, U]) parse(r *reader) (U, error) {
	_, err := p.pSkip.parse(r)
	if err != nil {
		return *new(U), err
	}
	v5, err := p.p5.parse(r)
	if err != nil {
		return *new(U), err
	}
	v4, err := p.p4.parse(r)
	if err != nil {
		return *new(U), err
	}

	v3, err := p.p3.parse(r)
	if err != nil {
		return *new(U), err
	}
	v2, err := p.p2.parse(r)
	if err != nil {
		return *new(U), err
	}
	v1, err := p.p1.parse(r)
	if err != nil {
		return *new(U), err
	}
	return p.f(v5, v4, v3, v2, v1), nil
}

var zeroParser Parser[any] = Untype(Pure(0))

func discardRight[T any](a Tuple2[T, any]) T {
	return a.Get0()
}
