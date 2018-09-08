// Copyright 2018 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package example

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//start:NewReader:
//imports[bufio,os]
type Reader struct {
	rdr  *bufio.Reader
	size int
}

func NewReader(size int) *Reader {
	return &Reader{
		bufio.NewReaderSize(os.Stdin, size),
		size,
	}
}

//end:NewReader:

//start:ReadLine:
func ReadLine(r *Reader) string {
	buf := make([]byte, 0, r.size)
	for {
		line, isPrefix, err := r.rdr.ReadLine()
		if err != nil {
			panic(err)
		}
		buf = append(buf, line...)
		if !isPrefix {
			break
		}
	}
	return string(buf)
}

//end:ReadLine:

//start:ReadStrs:
//imports[strings]
//dependsOn[ReadLine]
func ReadStrs(r *Reader) []string {
	return strings.Split(ReadLine(r), " ")
}

//end:ReadStrs:

//start:ReadInt:
//imports[fmt]
func ReadInt() (i int) {
	fmt.Scan(&i)
	return
}

//end:ReadInt:

//start:ReadInts:
//imports[strconv]
//dependsOn[ReadStrs]
func ReadInts(r *Reader) []int {
	a := ReadStrs(r)
	n := make([]int, 0)
	for _, v := range a {
		if i, e := strconv.Atoi(v); e == nil {
			n = append(n, i)
		}
	}
	return n
}

//end:ReadInts:
