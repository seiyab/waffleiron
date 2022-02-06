package waffleiron

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type FuncParser[T any] func(r *Reader) (T, error)

func (p FuncParser[T]) Parse(r *Reader) (T, error) {
	return p(r)
}

func Rune(rn rune) Parser[rune] {
	return runeParser{rn}
}

type runeParser struct {
	rn rune
}

func (p runeParser) Parse(r *Reader) (rune, error) {
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
	return stringParser{str}
}

type stringParser struct {
	str string
}

func (p stringParser) Parse(r *Reader) (string, error) {
	overrun := int64(len(p.str)) > int64(len(r.str))-r.idx
	if overrun || !strings.HasPrefix(r.RemainingString(), p.str) {
		return "", errors.Errorf("expected %q, but not found at %s", p.str, r.pos)
	}
	r.SkipBytes(len(p.str))
	return p.str, nil
}

func Regexp(re *regexp.Regexp) Parser[string] {
	if !strings.HasPrefix(re.String(), "^") {
		return regexpParser{
			re: regexp.MustCompile("^" + re.String()),
		}
	}
	return regexpParser{re}
}

type regexpParser struct {
	re *regexp.Regexp
}

func (p regexpParser) Parse(r *Reader) (string, error) {
	str := r.RemainingString()
	loc := p.re.FindStringIndex(str)
	if len(loc) == 0 {
		return "", errors.Errorf("expected to match %q at %s", p.re, r.pos)
	}
	if loc[0] != 0 {
		panic("regex matched on loc[0] != 0. it might be bug. please submit an issue.")
	}
	r.SkipBytes(loc[1])
	return str[0:loc[1]], nil
}

func RegexpStr(str string) Parser[string] {
	return Regexp(regexp.MustCompile(str))
}

func Pure[T any](value T) Parser[T] {
	return pureParser[T]{value}
}

type pureParser[T any] struct {
	value T
}

func (p pureParser[T]) Parse(r *Reader) (T, error) {
	return p.value, nil
}
