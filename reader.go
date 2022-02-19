package waffleiron

import (
	"fmt"
	"io"
	"strings"
)

type reader struct {
	str string
	idx int64 // current reading index
	pos Position
	trc *trace
}

// newReader creates a reader for string s
func newReader(s string) *reader {
	return &reader{
		str: s,
	}
}

func (r *reader) readRune() (rune, int, error) {
	sr := strings.NewReader(r.str)
	sr.Seek(r.idx, io.SeekStart)
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
	consumed := r.str[r.idx : r.idx+int64(s)]
	sr := strings.NewReader(consumed)
	for {
		c, s, err := sr.ReadRune()
		if err == io.EOF {
			return consumed, nil
		}
		if err != nil {
			return "", err
		}
		r.pos.column += 1
		r.idx += int64(s)
		if c == '\n' {
			r.pos.line += 1
			r.pos.column = 0
		}
	}
}

// try runs f
// If f returns an error, r doesn't consume any bytes
func (r *reader) try(f func() error) error {
	var idx int64 = r.idx
	var pos Position = r.pos
	err := f()
	if err != nil {
		r.idx = idx
		r.pos = pos
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

// Position represents a position in a string
type Position struct {
	line   int64
	column int64
}

// String implements the fmt.Stringer interface
func (p Position) String() string {
	return fmt.Sprintf("line %d, column %d", p.line, p.column)
}

type trace struct {
	name   string
	parent *trace
}
