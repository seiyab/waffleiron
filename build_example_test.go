package waffleiron_test

import (
	"fmt"

	wi "github.com/seiyab/waffleiron"
)

func ExampleBegin3() {
	ws := wi.Untype(wi.RegexpStr(` *`))
	p := wi.Begin3(func(a int, op rune, b int) int {
		if op == '-' {
			return a - b
		}
		if op == '*' {
			return a * b
		}
		return a + b
	}).
		Skip(ws).
		Then(wi.Int()). // this will be bound as `a`
		Skip(ws).
		Then(wi.Choice( // this will be bound as `op`
			wi.Rune('+'),
			wi.Rune('-'),
			wi.Rune('*'),
		)).
		Skip(ws).
		Then(wi.Int()). // this will be bound as `b`
		Skip(ws).
		End()

	x, _ := wi.Parse("1 + 3", p)
	fmt.Println(x)
	x, _ = wi.Parse("6 - 1", p)
	fmt.Println(x)
	x, _ = wi.Parse("   3     *   2  ", p)
	fmt.Println(x)
	// Output:
	// 4
	// 5
	// 6
}
