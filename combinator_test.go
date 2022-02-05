package waffleiron_test

import (
	"fmt"
	"regexp"
	"testing"

	wi "github.com/seiyab/waffleiron"
)

func TestAnd(t *testing.T) {
	t.Run("And(Rune, Rune)", func(t *testing.T) {
		p := wi.And(wi.Rune('a'), wi.Rune('b'))

		type testCase struct {
			str string
			err bool
			a   rune
			b   rune
		}

		testCases := []testCase{
			{
				str: "ab",
				err: false,
				a:   'a',
				b:   'b',
			},
			{
				str: "abb",
				err: true,
			},
			{
				str: "ba",
				err: true,
			},
			{
				str: "",
				err: true,
			},
		}

		for _, tt := range testCases {
			t.Run(tt.str, func(t *testing.T) {
				v, err := wi.Parse(tt.str, p)
				if tt.err {
					if err == nil {
						t.Errorf("err = nil")
					}
				} else {
					if err != nil {
						t.Fatal(err)
					}
					if v.A != tt.a {
						t.Errorf("v.A = %q, want %q", v.A, tt.a)
					}
					if v.B != tt.b {
						t.Errorf("v.B = %q, want %q", v.B, tt.b)
					}
				}
			})
		}
	})

	t.Run("And(String, Rune)", func(t *testing.T) {
		p := wi.And(wi.String("waffle"), wi.Rune('i'))

		type testCase struct {
			str string
			err bool
			a   string
			b   rune
		}

		testCases := []testCase{
			{
				str: "wafflei",
				err: false,
				a:   "waffle",
				b:   'i',
			},
			{
				str: "waffle",
				err: true,
			},
			{
				str: "w",
				err: true,
			},
			{
				str: "i",
				err: true,
			},
			{
				str: "waffleiron",
				err: true,
			},
			{
				str: "",
				err: true,
			},
		}

		for _, tt := range testCases {
			t.Run(tt.str, func(t *testing.T) {
				v, err := wi.Parse(tt.str, p)
				if tt.err {
					if err == nil {
						t.Errorf("err = nil")
					}
				} else {
					if err != nil {
						t.Fatal(err)
					}
					if v.A != tt.a {
						t.Errorf("v.A = %q, want %q", v.A, tt.a)
					}
					if v.B != tt.b {
						t.Errorf("v.B = %q, want %q", v.B, tt.b)
					}
				}
			})
		}
	})
}

func TestTrace(t *testing.T) {
	p := wi.And(
		wi.Trace("A", wi.String("waffle")),
		wi.Trace("B", wi.And(
			wi.Trace("C", wi.Rune('i')),
			wi.Trace("D", wi.String("ron")),
		)),
	)

	type testCase struct {
		str    string
		traces []string
	}

	testCases := []testCase{
		{
			str:    "wafxleiron",
			traces: []string{"A"},
		},
		{
			str:    "wafflexron",
			traces: []string{"B", "C"},
		},
		{
			str:    "waffleixon",
			traces: []string{"B", "D"},
		},
		{
			str:    "",
			traces: []string{"A"},
		},
		{
			str:    "waffleironx",
			traces: []string{},
		},
	}

	rgx := regexp.MustCompile("at \"[^\"]+\"")

	for _, tt := range testCases {
		t.Run(tt.str, func(t *testing.T) {
			_, err := wi.Parse(tt.str, p)

			if err == nil {
				t.Fatal("expected error")
			}

			actualTraces := rgx.FindAllString(err.Error(), -1)

			if len(actualTraces) != len(tt.traces) {
				t.Fatalf("traces doesn't match: %v", err)
			}

			for i, tr := range tt.traces {
				expectedTrace := fmt.Sprintf("at %q", tr)
				if actualTraces[i] != expectedTrace {
					t.Errorf("actualTraces[%d] = %q, want %q", i, actualTraces[i], expectedTrace)
				}
			}
		})
	}
}
