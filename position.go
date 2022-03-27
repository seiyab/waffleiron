package waffleiron

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Position represents a position in a string
type Position struct {
	Line   int
	Column int
}

// String implements the fmt.Stringer interface
func (p Position) String() string {
	return fmt.Sprintf("Line %d, Column %d", p.Line, p.Column)
}

type Locator struct {
	s        string
	newlines []int
}

func NewLocator(s string) Locator {
	return Locator{
		s:        s,
		newlines: []int{-1},
	}
}

func (l *Locator) Locate(byteIndex int) (Position, error) {
	if byteIndex >= len(l.s) {
		return Position{}, errors.Errorf("index ouf of bounds (%d >= %d)", byteIndex, len(l.s))
	}
	line := upperBound(l.newlines, func(n int) bool {
		return n < byteIndex
	})
	head := l.newlines[line] + 1
	reader := strings.NewReader(l.s[head:])
	col := 0
	needle := 0
	for head+needle < byteIndex {
		r, size, err := reader.ReadRune()
		if err != nil {
			return Position{}, errors.Wrap(err, "failed to read rune")
		}
		if r == '\n' {
			l.newlines = append(l.newlines, head+needle)
			line += 1
			col = 0
			needle += size
			continue
		}
		col += 1
		needle += size
	}
	return Position{
		Line:   line,
		Column: col,
	}, nil
}

func upperBound[T any](slice []T, ok func(v T) bool) int {
	var f func(left, right int) int
	f = func(left, right int) int {
		if right-left <= 1 {
			return left
		}
		m := (left + right) / 2
		if ok(slice[m]) {
			return f(m, right)
		}
		return f(left, m)
	}
	return f(0, len(slice))
}
