package waffleiron_test

import (
	"regexp"
	"testing"

	wi "github.com/seiyab/waffleiron"
)

func TestRegexp(t *testing.T) {
	t.Run("^[a-z,A-Z]+", func(t *testing.T) {
		p := wi.Regexp(regexp.MustCompile("^[a-z,A-Z]+"))

		t.Run("ok", func(t *testing.T) {
			testCases := []string{
				"abc",
				"XYZ",
				"T",
				"WaffleIron",
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
				"1abc",
				"abc1",
				"waffle iron",
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
