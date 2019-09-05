// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
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
	output  = flag.String("o", "", "file for output; default: stdout")
	htmlOut = flag.String("html", "", "generate HTML representation of coverage profile")
	funcOut = flag.String("func", "", "output coverage profile information for each function")
)

var profile string // The profile to read; the value of -html or -func

func main() {
	flag.Usage = usage
	flag.Parse()

	// Usage information when no arguments.
	if flag.NFlag() == 0 && flag.NArg() == 0 {
		flag.Usage()
	}

	// Output HTML or function coverage information.
	var err error
	if *htmlOut != "" {
		err = htmlOutput(profile, *output)
	} else {
		err = funcOutput(profile, *output)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "cover: %v\n", err)
		os.Exit(2)
	}
}
