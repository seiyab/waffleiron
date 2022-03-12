package waffleiron_test

import (
	"testing"

	wi "github.com/seiyab/waffleiron"
)

func TestParseJsonLike(t *testing.T) {
	var p wi.Parser[JsonLike]
	ws := wi.Untype(wi.RegexpStr(`[ \t\n]*`))
	num := wi.Map(
		wi.Int(),
		func(n int) JsonLike {
			return Number(n)
		},
	)
	str := wi.Between(
		wi.Rune('"'),
		wi.RegexpStr(`[^"]*`),
		wi.Rune('"'),
	)

	obj := wi.Begin1(func(entries []wi.Tuple2[string, JsonLike]) JsonLike {
		o := make(Object)
		for _, e := range entries {
			o[e.Get0()] = e.Get1()
		}
		return o
	}).
		Skip(ws).
		Skip(wi.Untype(wi.Rune('{'))).
		Then(wi.SepBy(
			wi.Begin2(wi.NewTuple2[string, JsonLike]).
				Skip(ws).
				Then(str).
				Skip(ws).
				Skip(wi.Untype(wi.Rune(':'))).
				Skip(ws).
				Then(wi.Ref(&p)).
				Skip(ws).
				End(),
			wi.Rune(','),
		)).
		Skip(wi.Untype(wi.Rune('}'))).
		Skip(ws).
		End()

	p = wi.Choice(obj, num)

	text := `{
			"abc": 29,
			"x": {
				"yz": 100
			}
		}`

	result, err := wi.Parse(text, p)
	if err != nil {
		t.Fatal(err)
	}

	if result.AsMap()["abc"].AsNumber() != 29 {
		t.Errorf(`result.AsMap()["abc"].AsNumber() = %d, want 29`, result.AsMap()["abc"].AsNumber())
	}

	if result.AsMap()["x"].AsMap()["yz"].AsNumber() != 100 {
		t.Errorf(
			`result.AsMap()["x"].AsMap()["yz"].AsNumber() = %d, want 100`,
			result.AsMap()["x"].AsMap()["yz"].AsNumber(),
		)
	}
}

type JsonLike interface {
	AsNumber() int
	// AsString() string
	AsMap() map[string]JsonLike
	// AsArray() JsonLike[]
}

type Object map[string]JsonLike

func (o Object) AsMap() map[string]JsonLike {
	return o
}

func (o Object) AsNumber() int {
	panic("Object.AsNumber")
}

type Number int

func (n Number) AsNumber() int {
	return int(n)
}

func (n Number) AsMap() map[string]JsonLike {
	panic("Number.AsMap")
}
