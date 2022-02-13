package waffleiron_test

import (
	"testing"

	wi "github.com/seiyab/waffleiron"
)

func TestBuild0(t *testing.T) {
	p := wi.Build1(func(x int) int {
		return x * 3
	}).
		Skip(wi.Untype(wi.Rune('['))).
		Accept(wi.Int()).
		Skip(wi.Untype(wi.Rune(']'))).
		End()

	t.Run("ok", func(t *testing.T) {
		v, err := wi.Parse("[100]", p)
		if err != nil {
			t.Fatal(err)
		}
		if v != 300 {
			t.Errorf("v = %d, want %d", v, 300)
		}
	})

	t.Run("error", func(t *testing.T) {
		testCases := []string{
			"100]",
			"[]",
			"[100",
			"@100",
		}
		for _, tt := range testCases {
			t.Run(tt, func(t *testing.T) {
				_, err := wi.Parse(tt, p)
				if err == nil {
					t.Error("expected error")
				}
			})
		}
	})
}
