package waffleiron

import (
	"strings"

	"github.com/pkg/errors"
)

type FuncParser[T any] func(r *Reader) (T, error)

func (p FuncParser[T]) Parse(r *Reader) (T, error) {
	return p(r)
}

func Rune(rn rune) Parser[rune] {
	return RuneParser{rn}
}

type RuneParser struct {
	rn rune
}

func (p RuneParser) Parse(r *Reader) (rune, error) {
	ch, _, err := r.ReadRune()
	if err != nil {
		return 0, errors.Wrapf(err, "at %s", r.pos)
	}
	if ch != p.rn {
		return 0, errors.Errorf("expected %q, found %q at %s", p.rn, ch, r.pos)
	}
	return ch, nil
}

func String(str string) Parser[string] {
	return StringParser{str}
}

type StringParser struct {
	str string
}

func (p StringParser) Parse(r *Reader) (string, error) {
	overrun := int64(len(p.str)) > int64(len(r.str))-r.idx
	if overrun || !strings.HasPrefix(r.str[r.idx:], p.str) {
		return "", errors.Errorf("expected %q, but not found at %s", p.str, r.pos)
	}
	r.SkipBytes(len(p.str))
	return p.str, nil
}
