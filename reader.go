package waffleiron

import (
	"io"
	"strings"
)

type reader struct {
	str string
	idx int // current reading index
	loc Locator
	trc *trace
}

// newReader creates a reader for string s
func newReader(s string) *reader {
	return &reader{
		str: s,
		loc: NewLocator(s),
	}
}

func (r *reader) readRune() (rune, int, error) {
	sr := strings.NewReader(r.str)
	sr.Seek(int64(r.idx), io.SeekStart)
	c, s, err := sr.ReadRune()
	if err != nil {
		return c, s, err
	}
	r.consumeBytes(s)
	return c, s, nil
}

// remainingString returns remaining string that is not consumed yet
func (r *reader) remainingString() string {
	return r.str[r.idx:]
}

// consumeBytes consumes s bytes
func (r *reader) consumeBytes(s int) (string, error) {
	consumed := r.str[r.idx : r.idx+s]
	sr := strings.NewReader(consumed)
	for {
		_, _, err := sr.ReadRune()
		if err == io.EOF {
			r.idx += s
			return consumed, nil
		}
		if err != nil {
			return "", err
		}
	}
}

// locate returns current Position.String() of r.idx
// it returns `unknown` if loc.Locate returns error
func (r *reader) locateAndString() string {
	p, e := r.loc.Locate(r.idx)
	if e != nil {
		return "unknown"
	}
	return p.String()
}

// try runs f
// If f returns an error, r doesn't consume any bytes
func (r *reader) try(f func() error) error {
	idx := r.idx
	err := f()
	if err != nil {
		r.idx = idx
		return err
	}
	return nil
}

// withTrace adds trace information
func (r *reader) withTrace(name string, f func()) {
	r.trc = &trace{
		name:   name,
		parent: r.trc,
	}
	f()
	r.trc = r.trc.parent
}

type trace struct {
	name   string
	parent *trace
}
