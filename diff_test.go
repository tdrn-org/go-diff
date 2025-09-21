//
// Copyright (C) 2025 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package diff_test

import (
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tdrn-org/go-diff"
)

const emptyFileName string = "testdata/empty.txt"
const leftFileName string = "testdata/l.txt"
const rightFileName string = "testdata/r.txt"

func TestDiffBothEmpty(t *testing.T) {
	result := testDiff(t, emptyFileName, emptyFileName)
	require.Empty(t, result.Diffs)
}

func TestDiffLeftEmpty(t *testing.T) {
	result := testDiff(t, emptyFileName, rightFileName)
	require.Len(t, result.Diffs, 36)
}

func TestDiffRightEmpty(t *testing.T) {
	result := testDiff(t, leftFileName, emptyFileName)
	require.Len(t, result.Diffs, 29)
}

func TestDiffEqual(t *testing.T) {
	result := testDiff(t, leftFileName, leftFileName)
	require.Len(t, result.Diffs, 29)
}

func TestDiff(t *testing.T) {
	result := testDiff(t, leftFileName, rightFileName)
	require.Len(t, result.Diffs, 39)
}

func testDiff(t *testing.T, leftName string, rightName string) *diff.Result {
	fileResult := testDiffFiles(t, leftName, rightName)
	readersResult := testDiffReaders(t, leftName, rightName)
	require.Equal(t, fileResult.Diffs, readersResult.Diffs)
	linesResult := testDiffLines(t, leftName, rightName)
	require.Equal(t, fileResult.Diffs, linesResult.Diffs)
	fileResult.Print(os.Stdout)
	return fileResult
}

func testDiffFiles(t *testing.T, leftName string, rightName string) *diff.Result {
	result, err := diff.DiffFiles(leftName, rightName)
	require.NoError(t, err)
	require.Equal(t, leftName, result.LeftName)
	require.Equal(t, rightName, result.RightName)
	return result
}

func testDiffReaders(t *testing.T, leftName string, rightName string) *diff.Result {
	left, err := os.Open(leftName)
	require.NoError(t, err)
	defer left.Close()
	right, err := os.Open(rightName)
	require.NoError(t, err)
	defer right.Close()
	result, err := diff.Diff(left, right)
	require.NoError(t, err)
	require.Equal(t, diff.DefaultLeftName, result.LeftName)
	require.Equal(t, diff.DefaultRightName, result.RightName)
	return result
}

func testDiffLines(t *testing.T, leftName string, rightName string) *diff.Result {
	leftData, err := os.ReadFile(leftName)
	require.NoError(t, err)
	left := slices.AppendSeq([]string{}, strings.Lines(string(leftData)))
	rightData, err := os.ReadFile(rightName)
	require.NoError(t, err)
	right := slices.AppendSeq([]string{}, strings.Lines(string(rightData)))
	result := diff.DiffLines(left, right)
	require.Equal(t, diff.DefaultLeftName, result.LeftName)
	require.Equal(t, diff.DefaultRightName, result.RightName)
	return result
}
