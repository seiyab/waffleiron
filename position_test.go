package waffleiron_test

import (
	"fmt"
	"testing"

	wi "github.com/seiyab/waffleiron"
)

func TestLocator(t *testing.T) {
	t.Run("ASCII only", func(t *testing.T) {
		s := `012345678
0123456789012345678
012345678

123456789`

		type testCase struct {
			index    int
			position wi.Position
		}

		testCases := []testCase{
			{6, wi.Position{0, 6}},
			{9, wi.Position{0, 9}},
			{10, wi.Position{1, 0}},
			{11, wi.Position{1, 1}},
			{14, wi.Position{1, 4}},
			{20, wi.Position{1, 10}},
			{29, wi.Position{1, 19}},
			{30, wi.Position{2, 0}},
			{39, wi.Position{2, 9}},
			{40, wi.Position{3, 0}},
			{41, wi.Position{4, 0}},
			{42, wi.Position{4, 1}},
		}

		run := func(t *testing.T, tt testCase, l wi.Locator) {
			p, err := l.Locate(tt.index)
			if err != nil {
				t.Fatal(err)
			}
			if p != tt.position {
				t.Errorf("p = %v, want %v", p, tt.position)
			}
		}

		t.Run(fmt.Sprintf("no cache"), func(t *testing.T) {
			l := wi.NewLocator(s)
			for i, tt := range testCases {
				t.Run(fmt.Sprintf("case[%d]", i), func(t *testing.T) {
					run(t, tt, l)
				})
			}
		})

		t.Run(fmt.Sprintf("cached"), func(t *testing.T) {
			l := wi.NewLocator(s)
			_, err := l.Locate(len(s) - 1)
			if err != nil {
				t.Fatal(err)
			}
			for i, tt := range testCases {
				t.Run(fmt.Sprintf("case[%d]", i), func(t *testing.T) {
					run(t, tt, l)
				})
			}
		})
	})
}
