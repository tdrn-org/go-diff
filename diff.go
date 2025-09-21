//
// Copyright (C) 2025 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

// Package diff provides functions for line based diffing as well as printing the results.
package diff

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
)

// Op defines the diff operation associated with a specific line.
type Op int

const (
	// EqlOp indicates, left and right lines are equal.
	EqlOp Op = 0
	// AddOp indicates, line exists only on the right and has been added.
	AddOp Op = 1
	// DelOp indicates, line exists only on the left and has been deleted.
	DelOp Op = -1
)

// String matches the Op to it's string representation (Eql: "=", Add: "<", Del: ">").
func (op Op) String() string {
	switch op {
	case EqlOp:
		return "="
	case AddOp:
		return "<"
	case DelOp:
		return ">"
	}
	return "?"
}

// LineDiff represents the diff result for a single line.
type LineDiff struct {
	// Op indicates the diff operation associated with this line.
	Op Op
	// Line contains the actual line.
	Line string
}

// DefaultLeftName is used to name the left side of a diff
// in case no specific name has been given.
const DefaultLeftName = "l.txt"

// DefaultRightName is used to name the right side of a diff
// in case no specific name has been given.
const DefaultRightName = "r.txt"

// Result contains the result of a Diff operation.
type Result struct {
	// LeftName contains the name of the left side
	// of the Diff operation (e.g. file name).
	LeftName string
	// RightName contains the name of the right side
	// of the Diff operation (e.g. file name).
	RightName string
	// Diffs contains for all compared lines the diff result.
	Diffs []LineDiff
}

// Print prints the diff result to the given writer.
func (r *Result) Print(w io.Writer) {
	for _, diff := range r.Diffs {
		fmt.Fprintf(w, "%s %s", diff.Op, diff.Line)
	}
}

func (r *Result) reverse() {
	slices.Reverse(r.Diffs)
}

func (r *Result) keepLine(line string) {
	r.Diffs = append(r.Diffs, LineDiff{Op: EqlOp, Line: line})
}

func (r *Result) addLine(line string) {
	r.Diffs = append(r.Diffs, LineDiff{Op: AddOp, Line: line})
}

func (r *Result) addLines(lines []string) {
	for _, line := range lines {
		r.addLine(line)
	}
}

func (r *Result) deleteLine(line string) {
	r.Diffs = append(r.Diffs, LineDiff{Op: DelOp, Line: line})
}

func (r *Result) deleteLines(lines []string) {
	for _, line := range lines {
		r.deleteLine(line)
	}
}

// DiffFiles runs a diff operation on the two given file names.
func DiffFiles(leftName string, rightName string) (*Result, error) {
	left, err := os.Open(leftName)
	if err != nil {
		return nil, err
	}
	defer left.Close()
	right, err := os.Open(rightName)
	if err != nil {
		return nil, err
	}
	defer right.Close()
	differ, err := differFromReaders(left, leftName, right, rightName)
	if err != nil {
		return nil, err
	}
	return differ.run(), nil
}

// DiffLines runs a diff operation on the two given string arrays.
func DiffLines(left []string, right []string) *Result {
	return differFromLines(left, DefaultLeftName, right, DefaultRightName).run()
}

// Diff runs a diff operation on the two given reader's contents.
func Diff(left io.Reader, right io.Reader) (*Result, error) {
	differ, err := differFromReaders(left, DefaultLeftName, right, DefaultRightName)
	if err != nil {
		return nil, err
	}
	return differ.run(), nil
}

type differ struct {
	Left      []string
	LeftName  string
	Right     []string
	RightName string
}

func differFromLines(left []string, leftName string, right []string, rightName string) *differ {
	return &differ{
		Left:      left,
		LeftName:  leftName,
		Right:     right,
		RightName: rightName,
	}
}

func differFromReaders(left io.Reader, leftName string, right io.Reader, rightName string) (*differ, error) {
	leftLines, err := readLines(left)
	if err != nil {
		return nil, err
	}
	rightLines, err := readLines(right)
	if err != nil {
		return nil, err
	}
	return differFromLines(leftLines, leftName, rightLines, rightName), nil
}

func (p *differ) run() *Result {
	l := len(p.Left)
	r := len(p.Right)
	max := l + r
	result := &Result{
		LeftName:  p.LeftName,
		RightName: p.RightName,
		Diffs:     make([]LineDiff, 0, max),
	}
	if p.runFast(result, l, r, max) {
		return result
	}
	p.runFull(result, l, r, max)
	return result
}

func (p *differ) runFast(result *Result, l int, r int, max int) bool {
	if max == 0 {
		return true
	}
	if l == 0 {
		result.addLines(p.Right)
		return true
	}
	if r == 0 {
		result.deleteLines(p.Left)
		return true
	}
	return false
}

func (p *differ) runFull(result *Result, l int, r int, max int) {
	v := make([]int, 2*max+1)
	trace := make([][]int, 0, max)
	for d := 0; d <= max; d++ {
		dv := make([]int, len(v))
		copy(dv, v)
		trace = append(trace, dv)
		for k := -d; k <= d; k += 2 {
			var x int
			if k == -d || (k != d && v[max+k-1] < v[max+k+1]) {
				x = v[max+k+1]
			} else {
				x = v[max+k-1] + 1
			}
			y := x - k
			for x < l && y < r && p.Left[x] == p.Right[y] {
				x++
				y++
			}
			v[max+k] = x
			if x >= l && y >= r {
				p.runFullBacktrack(result, trace, l, r, max)
				return
			}
		}
	}
	panic("unexpected")
}

func (p *differ) runFullBacktrack(result *Result, trace [][]int, l int, r int, max int) {
	x := l
	y := r
	for d := len(trace) - 1; d >= 0; d-- {
		v := trace[d]
		k := x - y
		var prevK int
		if k == -d || (k != d && v[max+k-1] < v[max+k+1]) {
			prevK = k + 1
		} else {
			prevK = k - 1
		}
		prevX := v[max+prevK]
		prevY := prevX - prevK
		for x > prevX && y > prevY {
			x--
			y--
			result.keepLine(p.Left[x])
		}
		if d > 0 {
			if prevX < x {
				x--
				result.deleteLine(p.Left[x])
			} else {
				y--
				result.addLine(p.Right[y])
			}
		}
	}
	result.reverse()
}

func readLines(r io.Reader) ([]string, error) {
	buf := bufio.NewReader(r)
	lines := make([]string, 0)
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}
	return lines, nil
}
