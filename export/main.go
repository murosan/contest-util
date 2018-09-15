// Copyright 2018 murosan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	MainFilePath = "main.go"
	OutputPath   = "output/main.go"
	UtilDir      = "cutil"
)

var (
	deleteRegex    = regexp.MustCompile(`"github\.com/murosan/contest/cutil"|cutil\.`)
	headRegex      = regexp.MustCompile(`(?s)package.*import.?(\(.*?\)|".*?")`)
	quotedRegex    = regexp.MustCompile(`"([A-Za-z]*)"`)
	libRegex       = regexp.MustCompile(`cutil\.([A-Za-z]*)`)
	importRegex    = regexp.MustCompile(`imports\[(.*)]`)
	dependsOnRegex = regexp.MustCompile(`dependsOn\[(.*)]`)
	tokenRegex     = regexp.MustCompile(`//(imports|dependsOn|start:|end:).*`)

	ignoreFileRegexp = regexp.MustCompile(`.*_test.go`)
	utils            = listNames(UtilDir)
	emp              = []byte("")
	delm             = []byte(",")
)

func main() {
	abs, err := filepath.Abs(MainFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	// load main.go
	content, err := ioutil.ReadFile(abs)
	if err != nil {
		log.Fatalln("failed to load. path="+abs, err)
	}

	// extract method names
	arr := libRegex.FindAllStringSubmatch(string(content), -1)
	funcNames := make([]string, len(arr))
	for i, v := range arr {
		funcNames[i] = v[1]
	}

	// delete cutil from main
	clean := deleteRegex.ReplaceAll(content, emp)

	used := map[string]bool{}
	imports := map[string]bool{}
	var codes [][]byte

	// extract import from main
	head := headRegex.Find(clean)
	importMatches := quotedRegex.FindAllSubmatch(head, -1)
	for _, m := range importMatches {
		s := string(m[1])
		imports[s] = true
	}

	// load cutils
	var lib []byte
	for _, path := range utils {
		utl, err := ioutil.ReadFile(filepath.Clean(UtilDir + "/" + path))
		if err != nil {
			log.Fatal("failed to load file. path="+path, err)
		}
		lib = append(lib, headRegex.ReplaceAll(utl, emp)...)
	}

	// get funcs from cutil
	var searchRec func(string)
	searchRec = func(name string) {
		_, ok := used[name]
		if !ok {
			used[name] = true
			r := regexp.MustCompile("(?s)//start:" + name + ":.*//end:" + name + ":")
			block := r.Find(lib)
			importCsv := importRegex.FindSubmatch(block)
			dependsOn := dependsOnRegex.FindSubmatch(block)
			if len(importCsv) > 0 {
				for _, v := range bytes.Split(importCsv[1], delm) {
					imports[string(v)] = true
				}
			}
			if len(dependsOn) > 0 {
				for _, v := range bytes.Split(dependsOn[1], delm) {
					searchRec(string(v))
				}
			}
			codes = append(codes, tokenRegex.ReplaceAll(block, emp))
		}
	}

	for _, name := range funcNames {
		searchRec(name)
	}

	out := clean
	for _, v := range codes {
		out = append(out, v...)
	}

	imps := make([]string, len(imports))
	cnt := 0
	for k := range imports {
		imps[cnt] = "\"" + k + "\""
		cnt++
	}
	joined := strings.Join(imps, "\n")
	joined = "package main\n\nimport (\n" + joined + "\n)"

	noHead := headRegex.ReplaceAll(out, emp)

	result := append([]byte(joined), noHead...)
	pretty, e := format.Source(result)
	if e != nil {
		panic(e)
	}

	dir := filepath.Dir(OutputPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}

	e = ioutil.WriteFile(filepath.Clean(OutputPath), pretty, 0644)
	if e != nil {
		panic(e)
	}

	// copy to clipboard (mac only)
	c := exec.Command("pbcopy")
	in, err := c.StdinPipe()
	if err != nil {
		panic(err)
	}
	if err := c.Start(); err != nil {
		panic(err)
	}
	if _, err := in.Write(pretty); err != nil {
		panic(err)
	}
	if err := in.Close(); err != nil {
		panic(err)
	}
	c.Wait()
}

func listNames(dirname string) []string {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
	}
	var names []string
	for _, f := range files {
		if !f.IsDir() && !ignoreFileRegexp.MatchString(f.Name()) {
			names = append(names, f.Name())
		}
	}
	fmt.Println(names)
	return names
}
