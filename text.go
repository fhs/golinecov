// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

// textOutput reads the profile data from profile and generates a text
// coverage report, writing it to writer w.
func textOutput(w io.Writer, profile, gofile string) error {
	profiles, err := ParseProfiles(profile)
	if err != nil {
		return err
	}

	var d templateData

	dirs, err := findPkgs(profiles)
	if err != nil {
		return err
	}

	for _, profile := range profiles {
		fn := profile.FileName
		if gofile != "" && fn != gofile {
			continue
		}
		if profile.Mode == "set" {
			d.Set = true
		}
		file, err := findFile(dirs, fn)
		if err != nil {
			return err
		}
		src, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("can't read %q: %v", fn, err)
		}
		var buf strings.Builder
		if *showSource {
			err = textGen(&buf, src, profile.Boundaries(src))
			if err != nil {
				return err
			}
		}
		d.Files = append(d.Files, &templateFile{
			Name:     fn,
			Body:     buf.String(),
			Coverage: percentCovered(profile),
		})
	}
	if len(d.Files) == 0 {
		return fmt.Errorf("no coverage profile found for file %q", gofile)
	}

	if *showSource {
		for _, file := range d.Files {
			_, err = io.WriteString(w, file.Body)
			if err != nil {
				return err
			}
			fmt.Fprintf(w, "\n")
		}
	}
	for _, file := range d.Files {
		_, err = fmt.Fprintf(w, "%5.1f%% %v\n", file.Coverage, file.Name)
		if err != nil {
			return err
		}
	}
	return err
}

// percentCovered returns, as a percentage, the fraction of the statements in
// the profile covered by the test run.
// In effect, it reports the coverage of a given source file.
func percentCovered(p *Profile) float64 {
	var total, covered int64
	for _, b := range p.Blocks {
		total += int64(b.NumStmt)
		if b.Count > 0 {
			covered += int64(b.NumStmt)
		}
	}
	if total == 0 {
		return 0
	}
	return float64(covered) / float64(total) * 100
}

func printLine(w io.Writer, minCount int, line string) {
	if minCount == -1 {
		fmt.Fprintf(w, "%*s %s", *width, "-", line)
	} else {
		fmt.Fprintf(w, "%*d %s", *width, minCount, line)
	}
}

// textGen generates a text coverage report with the provided filename,
// source code, and tokens, and writes it to the given Writer.
func textGen(w io.Writer, src []byte, boundaries []Boundary) error {
	bcount := -1
	minCount := -1
	var line []byte
	for i := range src {
		for len(boundaries) > 0 && boundaries[0].Offset == i {
			b := boundaries[0]
			if b.Start {
				bcount = b.CountOrNorm()
			} else {
				bcount = -1
			}
			boundaries = boundaries[1:]
		}
		line = append(line, src[i])
		if src[i] == '\n' {
			printLine(w, minCount, string(line))
			line = nil
			minCount = bcount
		} else if minCount == -1 || bcount < minCount {
			minCount = bcount
		}
	}
	if len(line) > 0 {
		printLine(w, minCount, string(line))
		line = nil
	}
	return nil
}

type templateData struct {
	Files []*templateFile
	Set   bool
}

type templateFile struct {
	Name     string
	Body     string
	Coverage float64
}
