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
	`Usage of 'golinecov':

A coverage profile is produced by 'go test':
	go test -coverprofile=go.cover

Golinecov will look for the 'go.cover' file in the current directory. If
it's not found, it'll look for it in parent directories. The profile
can be set explicitly with the -profile flag.

Display source of a function or method matching a regular expression,
annotated with coverage per-line:
	golinecov -src -func '^Builder\.String$'

Display source of a file within a package, annotated with coverage per-line:
	golinecov -src strings/builder.go

Display summary of coverage for all files in the profile:
	golinecov

Display summary of coverage for functions matching a regular expression:
	golinecov -func '^Builder\.'
`

func usage() {
	fmt.Fprintln(os.Stderr, usageMessage)
	fmt.Fprintln(os.Stderr, "Flags:")
	flag.PrintDefaults()
	os.Exit(2)
}

var (
	profile    = flag.String("profile", "", "the coverage profile to read (searches for go.cover by default)")
	funcOut    = flag.Bool("func", false, "output coverage profile information for a function or method")
	norm       = flag.Bool("norm", false, "normalize count to [0, 10]")
	showSource = flag.Bool("src", false, "show source code annotated with coverage per line")
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("golinecov: ")

	flag.Usage = usage
	flag.Parse()

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	if *profile == "" {
		var err error
		*profile, err = findProfile()
		if err != nil {
			log.Fatalf("could not find profile: %v. Generate profile with 'go test -coverprofile=go.cover'.", err)
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
