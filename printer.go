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

func (p *Printer) Write(b []byte) (int, error) {
	return p.w.Write(b)
}

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

type Formatter interface {
	Format(p *Printer, r *Result)
}

type FormatterFunc func(*Printer, *Result)

func (f FormatterFunc) Format(p *Printer, r *Result) {
	f(p, r)
}

type PrinterOption interface {
	Apply(p *Printer)
}

type PrinterOptionFunc func(*Printer)

func (f PrinterOptionFunc) Apply(p *Printer) {
	f(p)
}

func WithAnsi(ansi bool) PrinterOption {
	return PrinterOptionFunc(func(p *Printer) {
		p.ansi = ansi
	})
}

func WithColors(eql string, add string, del string, rst string) PrinterOption {
	return PrinterOptionFunc(func(p *Printer) {
		p.ansiEql = eql
		p.ansiAdd = add
		p.ansiDel = del
		p.ansiRst = rst
	})
}

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
