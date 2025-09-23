//
// Copyright (C) 2025 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package diff

import (
	"fmt"
	"os"
	"time"
)

const DefaultUnifiedContext int = 3

func WithUnifiedFormatter(context int) PrinterOption {
	checkedContext := context
	if checkedContext < 0 {
		checkedContext = DefaultUnifiedContext
	}
	return PrinterOptionFunc(func(p *Printer) {
		p.formatter = &unifiedFormatter{Context: checkedContext}
	})
}

type unifiedFormatter struct {
	Context        int
	leftLine       int
	rightLine      int
	hunk           bool
	hunkIndex      int
	eqlRun         int
	hunkStartLeft  int
	hunkStartRight int
}

func (f *unifiedFormatter) Format(p *Printer, r *Result) {
	f.formatHeader(p, r)
	f.leftLine = 0
	f.rightLine = 0
	f.hunk = false
	f.eqlRun = 0
	for index, diff := range r.Diffs {
		if f.hunk {
			f.evalDiffInsideHunk(diff)
			if !f.hunk {
				f.formatHunk(p, r.Diffs)
			}
		} else {
			f.evalDiffOutsideHunk(diff)
			if f.hunk {
				f.hunkIndex = max(index-f.Context, 0)
			}
		}
	}
	if f.hunk {
		f.formatHunk(p, r.Diffs)
	}
}

func (f *unifiedFormatter) formatHeader(p *Printer, r *Result) {
	if p.Ansi() {
		colors := p.Colors()
		fmt.Fprintf(p, "%s--- %s\t%s%s\n", colors.Hdr, r.LeftName, f.modificationTime(r.LeftName), colors.Rst)
		fmt.Fprintf(p, "%s+++ %s\t%s%s\n", colors.Hdr, r.RightName, f.modificationTime(r.RightName), colors.Rst)
	} else {
		fmt.Fprintf(p, "--- %s\t%s\n", r.LeftName, f.modificationTime(r.LeftName))
		fmt.Fprintf(p, "+++ %s\t%s\n", r.RightName, f.modificationTime(r.RightName))
	}
}

func (f *unifiedFormatter) modificationTime(name string) string {
	now := time.Now()
	if name == DefaultLeftName || name == DefaultRightName {
		return now.String()
	}
	stat, err := os.Stat(name)
	if err != nil {
		return time.Time{}.UTC().String()
	}
	return stat.ModTime().String()
}

func (f *unifiedFormatter) evalDiffInsideHunk(diff LineDiff) {
	switch diff.Op {
	case EqlOp:
		f.eqlRun++
		f.hunk = f.eqlRun <= 2*f.Context
		f.leftLine++
		f.rightLine++
	case AddOp:
		f.eqlRun = 0
		f.rightLine++
	case DelOp:
		f.eqlRun = 0
		f.leftLine++
	}
}

func (f *unifiedFormatter) evalDiffOutsideHunk(diff LineDiff) {
	switch diff.Op {
	case EqlOp:
		f.leftLine++
		f.rightLine++
	case AddOp:
		f.hunk = true
		f.hunkStartLeft = max(f.leftLine-f.Context, 0)
		f.hunkStartRight = max(f.rightLine-f.Context, 0)
		f.rightLine++
	case DelOp:
		f.hunk = true
		f.hunkStartLeft = max(f.leftLine-f.Context, 0)
		f.hunkStartRight = max(f.rightLine-f.Context, 0)
		f.leftLine++
	}
}

func (f *unifiedFormatter) formatHunk(p *Printer, diffs []LineDiff) {
	hunkExtentLeft := f.leftLine - f.hunkStartLeft
	hunkExtentRight := f.rightLine - f.hunkStartRight
	if !f.hunk {
		hunkExtentLeft -= f.Context + 1
		hunkExtentRight -= f.Context + 1
	}
	f.formatRange(p, f.hunkStartLeft, hunkExtentLeft, f.hunkStartRight, hunkExtentRight)
	hunkLineLeft := f.hunkStartLeft
	hunkRemainingLeft := hunkExtentLeft
	hunkLineRight := f.hunkStartRight
	hunkRemainingRight := hunkExtentRight
	for _, diff := range diffs[f.hunkIndex:] {
		f.formatDiff(p, diff)
		switch diff.Op {
		case EqlOp:
			hunkLineLeft++
			hunkRemainingLeft--
			hunkLineRight++
			hunkRemainingRight--
		case AddOp:
			hunkLineRight++
			hunkRemainingRight--
		case DelOp:
			hunkLineLeft++
			hunkRemainingLeft--
		}
		if hunkRemainingLeft == 0 && hunkRemainingRight == 0 {
			break
		}
	}
}

func (f *unifiedFormatter) formatRange(p *Printer, startLeft int, extentLeft int, startRight int, extentRight int) {
	if p.Ansi() {
		fmt.Fprintf(p, "%s@@ -%d,%d +%d,%d @@%s\n", ansiLbl, startLeft+1, extentLeft, startRight+1, extentRight, ansiRst)
	} else {
		fmt.Fprintf(p, "@@ -%d,%d +%d,%d @@\n", startLeft+1, extentLeft, startRight+1, extentRight)
	}
}

func (f *unifiedFormatter) formatDiff(p *Printer, diff LineDiff) {
	if p.Ansi() {
		set, rst := p.OpColor(diff.Op)
		fmt.Fprintf(p.w, "%s%s %s%s", set, diff.Op, diff.Line, rst)
	} else {
		fmt.Fprintf(p.w, "%s %s", diff.Op, diff.Line)
	}
}
