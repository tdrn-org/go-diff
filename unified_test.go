//
// Copyright (C) 2025 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package diff_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tdrn-org/go-diff"
)

const expectedUnifiedDiffPlain string = "--- ./l.txt\t0001-01-01 00:00:00 +0000 UTC\n+++ ./r.txt\t0001-01-01 00:00:00 +0000 UTC\n@@ -1,3 +1,13 @@\n+0\n+1\n+2\n+3\n+4\n+5\n+6\n+7\n+8\n+9\n a\n b\n c\n@@ -24,6 +34,3 @@\n x\n y\n z\n-ä\n-ö\n-ü\n"
const expectedUnifiedDiffAnsi string = "\x1b[97m--- ./l.txt\t0001-01-01 00:00:00 +0000 UTC\x1b[0m\n\x1b[97m+++ ./r.txt\t0001-01-01 00:00:00 +0000 UTC\x1b[0m\n\x1b[96m@@ -1,3 +1,13 @@\x1b[0m\n\x1b[32m+0\n\x1b[0m\x1b[32m+1\n\x1b[0m\x1b[32m+2\n\x1b[0m\x1b[32m+3\n\x1b[0m\x1b[32m+4\n\x1b[0m\x1b[32m+5\n\x1b[0m\x1b[32m+6\n\x1b[0m\x1b[32m+7\n\x1b[0m\x1b[32m+8\n\x1b[0m\x1b[32m+9\n\x1b[0m\x1b[97m a\n\x1b[0m\x1b[97m b\n\x1b[0m\x1b[97m c\n\x1b[0m\x1b[96m@@ -24,6 +34,3 @@\x1b[0m\n\x1b[97m x\n\x1b[0m\x1b[97m y\n\x1b[0m\x1b[97m z\n\x1b[0m\x1b[31m-ä\n\x1b[0m\x1b[31m-ö\n\x1b[0m\x1b[31m-ü\n\x1b[0m"

func TestUnified(t *testing.T) {
	result, err := diff.DiffFiles(leftFileName, rightFileName)
	require.NoError(t, err)

	// set names to non-existing file to force mtime 0
	result.LeftName = "./" + diff.DefaultLeftName
	result.RightName = "./" + diff.DefaultRightName

	// Plain output
	{
		output := &strings.Builder{}
		diff.NewPrinter(output, diff.WithAnsi(false), diff.WithUnifiedFormatter(diff.DefaultUnifiedContext)).Print(result)
		require.Equal(t, expectedUnifiedDiffPlain, output.String())
	}

	// Ansi outpu
	{
		output := &strings.Builder{}
		diff.NewPrinter(output, diff.WithAnsi(true), diff.WithUnifiedFormatter(diff.DefaultUnifiedContext)).Print(result)
		require.Equal(t, expectedUnifiedDiffAnsi, output.String())
	}
}
