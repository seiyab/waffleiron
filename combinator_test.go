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
					if v.Get0() != tt.a {
						t.Errorf("v.Get0() = %q, want %q", v.Get0(), tt.a)
					}
					if v.Get1() != tt.b {
						t.Errorf("v.Get1() = %q, want %q", v.Get1(), tt.b)
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
					if v.Get0() != tt.a {
						t.Errorf("v.Get0() = %q, want %q", v.Get0(), tt.a)
					}
					if v.Get1() != tt.b {
						t.Errorf("v.Get1() = %q, want %q", v.Get1(), tt.b)
					}
				}
			})
		}
	})
}

func TestAnd3(t *testing.T) {
	p := wi.And3(
		wi.Rune('w'),
		wi.Rune('a'),
		wi.Rune('f'),
	)

	t.Run("ok", func(t *testing.T) {
		v, err := wi.Parse("waf", p)
		if err != nil {
			t.Fatal(err)
		}
		if v.Get0() != 'w' {
			t.Errorf("v.Get0() = %q, want %q", v.Get0(), 'w')
		}
		if v.Get1() != 'a' {
			t.Errorf("v.Get1() = %q, want %q", v.Get1(), 'a')
		}
		if v.Get2() != 'f' {
			t.Errorf("v.Get2() = %q, want %q", v.Get2(), 'f')
		}
	})

	t.Run("error", func(t *testing.T) {
		testCases := []string{
			"w",
			"waffle",
			"x",
		}

		for _, tt := range testCases {
			t.Run(tt, func(t *testing.T) {
				_, err := wi.Parse(tt, p)

				if err == nil {
					t.Errorf("expected error")
				}
			})
		}
	})
}

func TestChoice(t *testing.T) {
	p := wi.Choice(
		wi.String("waffleiron"),
		wi.String("waffle"),
		wi.String("iron"),
		wi.String("parser"),
	)

	t.Run("ok", func(t *testing.T) {
		testCases := []string{
			"waffleiron",
			"waffle",
			"iron",
			"parser",
		}

		for _, tt := range testCases {
			t.Run(tt, func(t *testing.T) {
				v, err := wi.Parse(tt, p)

				if err != nil {
					t.Fatal(err)
				}
				if v != tt {
					t.Errorf("v = %q, want %q", v, tt)
				}
			})
		}
	})

	t.Run("error", func(t *testing.T) {
		testCases := []string{
			"",
			"_waffleiron",
			"parsers",
			"p",
		}

		for _, tt := range testCases {
			t.Run(tt, func(t *testing.T) {
				_, err := wi.Parse(tt, p)

				if err == nil {
					t.Fatal("expected error")
				}
			})
		}

	})
}

func TestMaybe(t *testing.T) {
	p := wi.Maybe(wi.String("waffle"))

	t.Run("ok", func(t *testing.T) {
		testCases := []string{
			"waffle",
			"",
		}
		for _, tt := range testCases {
			t.Run(tt, func(t *testing.T) {
				v, err := wi.Parse(tt, p)
				if err != nil {
					t.Fatal(err)
				}
				if (v == nil) != (tt == "") {
					t.Error("unexpected v")
				}
			})
		}
	})
}

func TestRepeat(t *testing.T) {
	t.Run("Repeat(Rune)", func(t *testing.T) {
		p := wi.Repeat(wi.Rune('a'))

		t.Run("ok", func(t *testing.T) {
			testCases := []string{
				"",
				"a",
				"aaa",
				"aaaaaaaaaaa",
			}

			for _, tt := range testCases {
				t.Run(tt, func(t *testing.T) {
					v, err := wi.Parse(tt, p)
					if err != nil {
						t.Fatal(err)
					}
					if len(v) != len(tt) {
						t.Errorf("len(v) = %d, want %d", len(v), len(tt))
					}
				})
			}
		})

		t.Run("error", func(t *testing.T) {
			testCases := []string{
				"b",
				"xaaaa",
				"aaaax",
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
