// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Golinecov is a program for analyzing the coverage profiles generated by
'go test -coverprofile=cover.out'.

Unlike the 'go tool cover' command, golinecov was designed to explore
coverage profiles from the command line or a text editor, instead of
relying a web browser. Golinecov can display source code annotated with
coverage per-line.  It computes the coverage for a line by taking the
minimum of all coverage profile blocks that contains the line. This is
usually good enough for most workflows, but if it's not, the user should
fallback to 'go tool cover -html'.

For usage information, please see:
	golinecov -help
*/
package main
