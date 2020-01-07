package main

import (
	"bufio"
	"io"
)

type Config struct {
	Trailling bool
}

// when to add a comma
// between an ending char (bool, end quote, digit, ], } and a ) and a starting char (start quote, digit, bool, [, {)
// if trailling, between ending char and ending char

type Fixer struct {
	config *Config
	in     *bufio.Reader
	out    *bufio.Writer

	n int
}

func (f *Fixer) WriteByte(b byte) error {
	err := f.out.WriteByte(b)
	if err != nil {
		return err
	}
	f.n += 1
	return nil
}

func (f *Fixer) Write(b []byte) error {
	n, err := f.out.Write(b)
	if err != nil {
		return err
	}
	f.n += n
	return nil
}

// consumeString reads the entire string and writes it to out, untouched.
// It handles backslashes (\")
func (f *Fixer) consumeString() error {
	var bytes []byte
	var err error

	for len(bytes) == 0 || bytes[len(bytes)-2] == '\\' {
		bytes, err = f.in.ReadBytes('"')
		if err != nil {
			return err
		}
		if err := f.Write(bytes); err != nil {
			return err
		}

		assert(len(bytes) >= 2, "len bytes should be greater than 2, got %d in %q", len(bytes), bytes)
	}
	return nil
}

func (f *Fixer) consumeComment() error {
	var bytes []byte
	var err error

	// FIXME: better handling of different line endings
	bytes, err = f.in.ReadBytes('\n')
	if err != nil {
		return err
	}
	if err := f.Write(bytes); err != nil {
		return err
	}
	return nil

}

func (f *Fixer) Fix() error {

	for {
		// pay attention to
		// after a comment, consume the whole thing
		// within a string, consume the whole thing

		byte, err := f.in.ReadByte()
		if err != nil {
			return err
		}

		if err := f.WriteByte(byte); err != nil {
			return err
		}

		if byte == '"' {
			if err := f.consumeString(); err != nil {
				return err
			}
		} else if byte == '/' {
			next, err := f.in.Peek(1)
			if err != nil {
				return err
			}

			if next[0] == '/' {
				if err := f.consumeComment(); err != nil {
					return err
				}
			}

			// otherwise, don't do anything. We just peeked at the next
			// character, it's going to be consumed automatically
			// by something else
		}

	}
	return nil
}

func (f *Fixer) Written() int {
	return f.n
}

func (f *Fixer) Flush() error {
	return f.out.Flush()
}

// Fix writes everything from in to out, just adding commas where needed
// returns the number of bytes written, and error
func Fix(config *Config, in io.Reader, out io.Writer) (int, error) {
	f := &Fixer{
		config: config,
		in:     bufio.NewReader(in),
		out:    bufio.NewWriter(out),
	}
	err := f.Fix()
	if err == io.EOF {
		err = nil
	}
	n := f.Written()
	if err != nil {
		return n, err
	}
	err = f.Flush()
	return n, err
}
