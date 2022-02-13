package waffleiron_test

import (
	"fmt"

	wi "github.com/seiyab/waffleiron"
)

func Example() {
	p := wi.Between(
		wi.Rune('['),
		wi.SepBy(wi.Int(), wi.Rune(',')),
		wi.Rune(']'),
	)

	result, err := wi.Parse("[0,1,2,3]", p)

	fmt.Println(err)
	fmt.Println(fmt.Sprintf("%#v", result))

	// Output:
	// <nil>
	// []int{0, 1, 2, 3}
}

func Example_recursive() {
	type obj map[string]obj

	str := wi.Between(
		wi.Rune('"'),
		wi.RegexpStr(`[^"]*`),
		wi.Rune('"'),
	)
	ws := wi.Untype(wi.Maybe(wi.RegexpStr(`[ \t\n]+`)))

	var p wi.Parser[obj]

	p = wi.Build1(func(entries []wi.Tuple2[string, obj]) obj {
		o := make(obj, 0)
		for _, e := range entries {
			o[e.Get0()] = e.Get1()
		}
		return o
	}).
		Skip(ws).
		Skip(wi.Untype(wi.Rune('{'))).
		Skip(ws).
		Accept(wi.SepBy(
			wi.Build2(wi.NewTuple2[string, obj]).
				Accept(str).
				Skip(ws).
				Skip(wi.Untype(wi.Rune(':'))).
				Skip(ws).
				Accept(wi.Ref(&p)).
				End(),
			wi.Untype(wi.Between(
				ws,
				wi.Rune(','),
				ws,
			)),
		)).
		Skip(ws).
		Skip(wi.Untype(wi.Rune('}'))).
		Skip(ws).
		End()

	result, err := wi.Parse(`{
		"a": {
			"b": {
				"c": {}
			},
			"d": {
				"e": {},
				"f": {}
			}
		}
	}`, p)

	fmt.Println(err)
	fmt.Printf("%v", result)
	// Output:
	// <nil>
	// map[a:map[b:map[c:map[]] d:map[e:map[] f:map[]]]]
}
