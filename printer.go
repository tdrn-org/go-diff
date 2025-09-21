//
// Copyright (C) 2025 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package diff

import (
	"fmt"
	"io"
	"os"

	"github.com/mattn/go-isatty"
)

// Printer type supports configurable formatting and printing of diff results.
type Printer struct {
	w         io.Writer
	ansi      bool
	ansiEql   string
	ansiAdd   string
	ansiDel   string
	ansiRst   string
	formatter Formatter
}

const ansiEql = "\x1b[37m"
const ansiAdd = "\x1b[32m"
const ansiDel = "\x1b[31m"
const ansiRst = "\x1b[0m"

// Write as defined by [io.Writer]
func (p *Printer) Write(b []byte) (int, error) {
	return p.w.Write(b)
}

// OpAnsi returns the ansi color sequence configured for the given diff operation.
//
// Beside the color sequence, also the reset sequence is returned to reset coloring
// after all outputs have been printed. If coloring is disabled, both sequences
// are empty.
func (p *Printer) OpAnsi(op Op) (string, string) {
	if !p.ansi {
		return "", ""
	}
	switch op {
	case EqlOp:
		return p.ansiEql, p.ansiRst
	case AddOp:
		return p.ansiAdd, p.ansiRst
	case DelOp:
		return p.ansiDel, p.ansiRst
	}
	return "", ""
}

// Print prints the given diff result according to the Printer's configuration.
func (p *Printer) Print(r *Result) {
	p.formatter.Format(p, r)
}

func (p *Printer) defaultPrint(r *Result) {
	if p.ansi {
		for _, diff := range r.Diffs {
			esc, rst := p.OpAnsi(diff.Op)
			fmt.Fprintf(p.w, "%s%s %s%s", esc, diff.Op, diff.Line, rst)
		}
	} else {
		for _, diff := range r.Diffs {
			fmt.Fprintf(p.w, "%s %s", diff.Op, diff.Line)
		}
	}
}

// Formatter interface is used to format a diff result.
type Formatter interface {
	// Format is called to format the given diff result using
	// the given Printer instance.
	Format(p *Printer, r *Result)
}

// FormatterFunc typed functions are used to format diff results.
type FormatterFunc func(*Printer, *Result)

// Format formats the given diff result using
// the given Printer instance.
func (f FormatterFunc) Format(p *Printer, r *Result) {
	f(p, r)
}

// PrinterOption interface is used to configure a Printer instance.
type PrinterOption interface {
	// Apply applies the options represented by this instance
	// to the given Printer instance.
	Apply(p *Printer)
}

// PrinterOptionFunc typed functions are used to configure a Printer instance.
type PrinterOptionFunc func(*Printer)

// Apply applies options to the given Printer instance.
func (f PrinterOptionFunc) Apply(p *Printer) {
	f(p)
}

// WithAnsi explicitly enables or disables ansi color
// color output for a Printer instance.
//
// Per default the Printer instance checks the capabilities
// of the [io.Writer] provided during creation, to check
// whehter color output is suppored.
func WithAnsi(ansi bool) PrinterOption {
	return PrinterOptionFunc(func(p *Printer) {
		p.ansi = ansi
	})
}

// WithColors sets the ansi sequences to use for coloring
// the diff result.
//
// If coloring is disabled, these sequences are not used.
func WithColors(eql string, add string, del string, rst string) PrinterOption {
	return PrinterOptionFunc(func(p *Printer) {
		p.ansiEql = eql
		p.ansiAdd = add
		p.ansiDel = del
		p.ansiRst = rst
	})
}

// NewPrinter creates a new Printer instance using the given
// [io.Writer] and printer options.
func NewPrinter(w io.Writer, opts ...PrinterOption) *Printer {
	file, ok := w.(*os.File)
	ansi := ok && (isatty.IsTerminal(file.Fd()) || isatty.IsCygwinTerminal(file.Fd()))
	printer := &Printer{
		w:       w,
		ansi:    ansi,
		ansiEql: ansiEql,
		ansiAdd: ansiAdd,
		ansiDel: ansiDel,
		ansiRst: ansiRst,
		formatter: FormatterFunc(func(p *Printer, r *Result) {
			p.defaultPrint(r)
		}),
	}
	for _, opt := range opts {
		opt.Apply(printer)
	}
	return printer
}
