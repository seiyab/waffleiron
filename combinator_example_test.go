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
	type obj struct {
		name    string
		entries []obj
	}

	str := wi.Between(
		wi.Rune('"'),
		wi.RegexpStr(`[^"]*`),
		wi.Rune('"'),
	)
	ws := wi.RegexpStr(`[ \t\n]+`)

	var p wi.Parser[obj]

	p = wi.Map(
		wi.And3(
			str,
			wi.Between(wi.Maybe(ws), wi.Rune(':'), wi.Maybe(ws)),
			wi.Between(
				wi.And3(wi.Maybe(ws), wi.Rune('{'), wi.Maybe(ws)),
				wi.SepBy(
					wi.Ref(&p),
					wi.And(wi.Rune(','), wi.Maybe(ws)),
				),
				wi.And3(wi.Maybe(ws), wi.Rune('}'), wi.Maybe(ws)),
			),
		),
		func(t wi.Tuple3[string, rune, []obj]) obj {
			return obj{
				name:    t.Get0(),
				entries: t.Get2(),
			}
		},
	)

	result, err := wi.Parse(`"a": {
		"b": {
			"c": {}
		},
		"d": {
			"e": {},
			"f": {}
		}
	}`, p)

	fmt.Println(err)
	fmt.Printf("%v", result)
	// Output:
	// <nil>
	// {a [{b [{c []}]} {d [{e []} {f []}]}]}
}
