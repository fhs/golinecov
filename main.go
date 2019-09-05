// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const usageMessage = "" +
	`Usage of 'go tool cover':
Given a coverage profile produced by 'go test':
	go test -coverprofile=c.out

Open a web browser displaying annotated source code:
	go tool cover -html=c.out

Write out an HTML file instead of launching a web browser:
	go tool cover -html=c.out -o coverage.html

Display coverage percentages to stdout for each function:
	go tool cover -func=c.out

Finally, to generate modified source code with coverage annotations
(what go test -cover does):
	go tool cover -mode=set -var=CoverageVariableName program.go
`

func usage() {
	fmt.Fprintln(os.Stderr, usageMessage)
	fmt.Fprintln(os.Stderr, "Flags:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "\n  Only one of -html, -func, or -mode may be set.")
	os.Exit(2)
}

var (
	profile    = flag.String("profile", "", "the coverage profile to read (searches for go.cover by default)")
	funcOut    = flag.Bool("func", false, "output coverage profile information for a function")
	norm       = flag.Bool("norm", false, "normalize count to [0, 10]")
	showSource = flag.Bool("src", false, "show source code annotated with coverage per line")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	if *profile == "" {
		var err error
		*profile, err = findProfile()
		if err != nil {
			log.Fatalf("could not find profile: %v", err)
		}
	}

	if *funcOut {
		var err error
		if flag.NArg() == 0 {
			err = funcOutput(w, *profile, ".")
		} else {
			err = funcOutput(w, *profile, flag.Arg(0))
		}
		if err != nil {
			log.Fatalf("%v\n", err)
		}
		return
	}

	if flag.NArg() == 0 {
		err := textOutput(w, *profile, "")
		if err != nil {
			fmt.Fprintf(os.Stderr, "cover: %v\n", err)
			os.Exit(2)
		}
		return
	}

	for i := 0; i < flag.NArg(); i++ {
		err := textOutput(w, *profile, flag.Arg(i))
		if err != nil {
			log.Fatalf("%v\n", err)
		}
	}
}

func findProfile() (string, error) {
	d, err := os.Getwd()
	if err != nil {
		return "", err
	}
	filename := "go.cover"

	for {
		fn := filepath.Join(d, filename)
		info, err := os.Stat(fn)
		if !os.IsNotExist(err) && !info.IsDir() {
			return fn, nil
		}
		dd := filepath.Dir(d)
		if dd == d {
			break
		}
		d = dd
	}
	return "", fmt.Errorf("go.cover file not found")
}
