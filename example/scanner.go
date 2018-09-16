// Copyright 2018 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cutil

import (
	"bufio"
	"io"
	"strconv"
)

// 1. set the separator
//    ex1. SetSplitter(" ") // default
//    ex2. SetSplitterByRune(',')
// 2. run scanning

//start:NewScanner:
//imports[bufio,os,io]
type Scanner struct {
	sc *bufio.Scanner
}

func NewScanner(r io.Reader) *Scanner {
	s := &Scanner{bufio.NewScanner(r)}
	s.sc.Split(bufio.ScanWords)
	return s
}

//end:NewScanner:

//start:ScanStr:
//dependsOn[NewScanner]
func (s *Scanner) ScanStr() string {
	s.sc.Scan()
	return s.sc.Text()
}

//end:ScanStr:

//start:ScanInt:
//imports[strconv]
//dependsOn[NewScanner]
func (s *Scanner) ScanInt() int {
	s.sc.Scan()
	i, e := strconv.Atoi(s.sc.Text())
	if e != nil {
		panic(e)
	}
	return i
}

//end:ScanInt:

//start:ScanLine:
//dependsOn[NewScanner]
func (s *Scanner) ScanLine() string {
	s.sc.Scan()
	return s.sc.Text()
}

//end:ScanLine:

//start:ScanStrs:
//dependsOn[ScanStr]
func (s *Scanner) ScanStrs(len int) []string {
	a := make([]string, len)
	for i := 0; i < len; i++ {
		a[i] = s.ScanStr()
	}
	return a
}

//end:ScanStrs:

//start:ScanInts:
//dependsOn[ScanInt]
func (s *Scanner) ScanInts(len int) []int {
	a := make([]int, len)
	for i := 0; i < len; i++ {
		a[i] = s.ScanInt()
	}
	return a
}

//end:ScanInts:

//start:SetSplitter:
//dependsOn[NewScanner]
func (s *Scanner) SetSplitter(sep string) {
	switch sep {
	case "":
		s.sc.Split(bufio.ScanRunes)
	case " ":
		s.sc.Split(bufio.ScanWords)
	default:
		s.sc.Split(bufio.ScanLines)
	}
}

//end:SetSplitter:
