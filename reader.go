package waffleiron

import (
	"fmt"
	"io"
	"strings"
)

type Reader struct {
	str string
	idx int64 // current reading index
	pos Position
	trc *trace
}

func NewReader(s string) *Reader {
	return &Reader{
		str: s,
	}
}

func (r *Reader) ReadRune() (rune, int, error) {
	sr := strings.NewReader(r.str)
	sr.Seek(r.idx, io.SeekStart)
	c, s, err := sr.ReadRune()
	if err != nil {
		return c, s, err
	}
	r.SkipBytes(s)
	return c, s, nil
}

func (r *Reader) RemainingString() string {
	return r.str[r.idx:]
}

func (r *Reader) SkipBytes(s int) error {
	sr := strings.NewReader(r.str[r.idx : r.idx+int64(s)])
	for {
		c, s, err := sr.ReadRune()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		r.pos.column += 1
		r.idx += int64(s)
		if c == '\n' {
			r.pos.line += 1
			r.pos.column = 0
		}
	}
}

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

func (r *Reader) WithTrace(name string, f func()) {
	r.trc = &trace{
		name:   name,
		parent: r.trc,
	}
	f()
	r.trc = r.trc.parent
}

type Position struct {
	line   int64
	column int64
}

func (p Position) String() string {
	return fmt.Sprintf("line %d, column %d", p.line, p.column)
}

type trace struct {
	name   string
	parent *trace
}
