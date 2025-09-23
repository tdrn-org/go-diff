//
// Copyright (C) 2025 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package diff

// Op colors
const ansiEql = "\x1b[97m"
const ansiAdd = "\x1b[32m"
const ansiDel = "\x1b[31m"

// Extra colors
const ansiHdr = "\x1b[97m"
const ansiLbl = "\x1b[96m"

// Reset
const ansiRst = "\x1b[0m"

type Colors struct {
	Eql string
	Add string
	Del string
	Hdr string
	Lbl string
	Rst string
}

var noColors *Colors = &Colors{}

var defaultColors *Colors = &Colors{
	Eql: ansiEql,
	Add: ansiAdd,
	Del: ansiDel,
	Hdr: ansiHdr,
	Lbl: ansiLbl,
	Rst: ansiRst,
}
