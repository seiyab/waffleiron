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
