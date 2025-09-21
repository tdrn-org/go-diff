//
// Copyright (C) 2025 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package diff_test

import (
	"io"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tdrn-org/go-diff"
)

func TestPrinter(t *testing.T) {
	setupPrinters := []func(io.Writer) (*diff.Printer, string){
		testDefaultPrinter,
		testPlainPrinter,
		testAnsiPrinter,
	}
	for _, setupPrinter := range setupPrinters {
		t.Run(runtime.FuncForPC(reflect.ValueOf(setupPrinter).Pointer()).Name(), func(t *testing.T) {
			output := &strings.Builder{}
			printer, expected := setupPrinter(output)
			left := []string{
				"removed line\n",
				"unchanged line\n",
			}
			right := []string{
				"unchanged line\n",
				"added line\n",
			}
			printer.Print(diff.DiffLines(left, right))
			require.Equal(t, expected, output.String())
		})
	}
}

func testDefaultPrinter(w io.Writer) (*diff.Printer, string) {
	printer := diff.NewPrinter(w)
	expected := "> removed line\n= unchanged line\n< added line\n"
	return printer, expected
}

func testPlainPrinter(w io.Writer) (*diff.Printer, string) {
	printer := diff.NewPrinter(w, diff.WithAnsi(false))
	expected := "> removed line\n= unchanged line\n< added line\n"
	return printer, expected
}

func testAnsiPrinter(w io.Writer) (*diff.Printer, string) {
	printer := diff.NewPrinter(w, diff.WithAnsi(true))
	expected := "\x1b[31m> removed line\n\x1b[0m\x1b[37m= unchanged line\n\x1b[0m\x1b[32m< added line\n\x1b[0m"
	return printer, expected
}
