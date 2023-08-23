// Some parts Copyright 2015 The Go Authors. All rights reserved.
// Use of these parts of source code is governed by a BSD-style
// license that can be found at github.com/golang/go.

package main

import (
	"flag"

	"github.com/rj45/llbrew/rewriter"
)

func main() {
	outfile := flag.String("o", "", "output file")
	pkg := flag.String("pkg", "", "package name")
	fn := flag.String("fn", "", "function name")

	flag.Parse()

	filename := flag.Arg(0)

	rewriter.Generate(*fn, filename, *outfile, *pkg)
}
