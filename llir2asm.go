package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/rj45/llir2asm/asm"
	"github.com/rj45/llir2asm/compile"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	var optz = flag.Bool("Oz", false, "Optimize for min size")
	var opts = flag.Bool("Os", false, "Optimize for size")
	var opt0 = flag.Bool("O0", false, "Disable all optimizations")
	var opt1 = flag.Bool("O1", false, "Minimal speed optimizations")
	var opt2 = flag.Bool("O2", false, "Maximal speed optimizations")
	var outfile = flag.String("o", "-", "Output assembly file")
	var llfile = flag.String("ll", "", "Dump optimized llvm IR to file")
	var irfile = flag.String("ir", "", "Dump pre-optimized llir2asm IR to file")

	flag.Parse()
	c := compile.Compiler{}

	c.OptSize = 1
	c.OptSpeed = 1

	if *opts {
		c.OptSize = 1
	} else if *optz {
		c.OptSize = 2
	}

	if *opt0 {
		c.OptSpeed = 0
		c.OptSize = 0
	} else if *opt1 {
		c.OptSpeed = 1
	} else if *opt2 {
		c.OptSpeed = 2
	}

	c.DumpLL = createFile(*llfile)
	defer c.DumpLL.Close()

	c.DumpIR = createFile(*irfile)
	defer c.DumpIR.Close()

	prog, err := c.Compile(flag.Args()[0])
	if err != nil {
		log.Fatal(err)
	}

	out := createFile(*outfile)
	defer out.Close()

	asm.Emit(out, asm.CustomASM{}, prog)
}

func createFile(filename string) io.WriteCloser {
	if filename == "-" {
		return nopWriteCloser{os.Stdout}
	}
	if filename == "" {
		return nopWriteCloser{}
	}
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

type nopWriteCloser struct{ w io.Writer }

func (nopWriteCloser) Close() error {
	return nil
}

func (n nopWriteCloser) Write(p []byte) (int, error) {
	if n.w == nil {
		return len(p), nil
	}
	return n.w.Write(p)
}
