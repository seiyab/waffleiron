package waffleiron

import (
	"fmt"
	"io"
	"strings"
)

// Reader is used to read string by parser
// You don't need to know it well
// if you don't write your own primitive parser or parser combinator
type Reader struct {
	str string
	idx int64 // current reading index
	pos Position
	trc *trace
}

// NewReader creates a reader for string s
func NewReader(s string) *Reader {
	return &Reader{
		str: s,
	}
}

// ReadRune implements the io.RuneReader interface.
func (r *Reader) ReadRune() (rune, int, error) {
	sr := strings.NewReader(r.str)
	sr.Seek(r.idx, io.SeekStart)
	c, s, err := sr.ReadRune()
	if err != nil {
		return c, s, err
	}
	r.ConsumeBytes(s)
	return c, s, nil
}

// RemainingString returns remaining string that is not consumed yet
func (r *Reader) RemainingString() string {
	return r.str[r.idx:]
}

// ConsumeBytes consumes s bytes
func (r *Reader) ConsumeBytes(s int) (string, error) {
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

// Try runs f
// If f returns an error, r doesn't consume any bytes
func (r *Reader) Try(f func() error) error {
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

// WithTrace adds trace information
func (r *Reader) WithTrace(name string, f func()) {
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
