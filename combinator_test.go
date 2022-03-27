package waffleiron_test

import (
	"fmt"
	"regexp"
	"strings"
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

	t.Run("And(Word, Rune)", func(t *testing.T) {
		p := wi.And(wi.Word("waffle"), wi.Rune('i'))

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
		wi.Word("waffleiron"),
		wi.Word("waffle"),
		wi.Word("iron"),
		wi.Word("parser"),
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
	p := wi.Maybe(wi.Word("waffle"))

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
					for i := range v {
						if v[i] != 'a' {
							t.Errorf("v[%d] = %q, want %q", i, v[i], 'a')
						}
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

func TestSepBy(t *testing.T) {
	p := wi.SepBy(
		wi.RegexpStr("[0-9]+"),
		wi.Rune(','),
	)

	t.Run("ok", func(t *testing.T) {
		type testCase struct {
			str string
			arr []string
		}
		testCases := []testCase{
			{
				str: "",
				arr: []string{},
			},
			{
				str: "0",
				arr: []string{"0"},
			},
			{
				str: "0,123,456",
				arr: []string{"0", "123", "456"},
			},
		}
		for _, tt := range testCases {
			t.Run(tt.str, func(t *testing.T) {
				v, err := wi.Parse(tt.str, p)
				if err != nil {
					t.Fatal(err)
				}
				if len(v) != len(tt.arr) {
					t.Fatalf("len(v = %s) = %d, want %d", v, len(v), len(tt.arr))
				}
				for i := range v {
					if v[i] != tt.arr[i] {
						t.Errorf("v[%d] = %q, want %q", i, v[i], tt.arr[i])
					}
				}
			})
		}
	})
}

func TestBetween(t *testing.T) {
	p := wi.Between(
		wi.Rune('('),
		wi.RegexpStr("[0-9]+"),
		wi.Rune(')'),
	)

	t.Run("ok", func(t *testing.T) {
		testCases := []string{
			"(0)",
			"(123)",
			"(123456789)",
		}
		for _, tt := range testCases {
			t.Run(tt, func(t *testing.T) {
				_, err := wi.Parse(tt, p)
				if err != nil {
					t.Fatal(err)
				}
			})
		}
	})
}

func TestTrace(t *testing.T) {
	p := wi.And(
		wi.Trace("A", wi.Word("waffle")),
		wi.Trace("B", wi.And(
			wi.Trace("C", wi.Rune('i')),
			wi.Trace("D", wi.Word("ron")),
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

	rgx := regexp.MustCompile("\"[^\"]+\" >")

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
				expectedTrace := fmt.Sprintf("%q >", tr)
				if actualTraces[i] != expectedTrace {
					t.Errorf("actualTraces[%d] = %q, want %q", i, actualTraces[i], expectedTrace)
				}
			}
		})
	}
}

func TestErrorPosition(t *testing.T) {
	s := `- waffle
- iron
- parser ***
- combinator`

	p := wi.SepBy(wi.And(wi.Word("- "), wi.RegexpStr(`^[\w]+`)), wi.Rune('\n'))

	_, err := wi.Parse(s, p)

	if err == nil {
		t.Errorf("expected error")
	}

	expected := wi.Position{Line: 2, Column: 8}.String()
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("err.Error() = %q, want to contain %q", err.Error(), expected)
	}
}
