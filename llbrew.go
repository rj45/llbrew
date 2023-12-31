package main

import (
	"flag"
	"log"
	"os"

	"github.com/rj45/llbrew/arch"
	"github.com/rj45/llbrew/compile"
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
	var dumpssa = flag.String("dumpssa", "", "Dump ssa.html for specified function")
	var binfile = flag.String("bin", "", "Output binary using customasm")
	var run = flag.Bool("run", false, "Run the program using the target emulator")
	var trace = flag.Bool("trace", false, "Run the program in the emulator with tracing turned on")

	flag.Parse()

	arch.SetArch("rj32")

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

	c.DumpLL = *llfile
	c.DumpIR = *irfile
	c.DumpSSA = *dumpssa
	c.OutFile = *outfile
	c.BinFile = *binfile
	c.Run = *run
	c.RunTrace = *trace

	err := c.Compile(flag.Args()[0])
	if err != nil {
		log.Fatal(err)
	}
}
