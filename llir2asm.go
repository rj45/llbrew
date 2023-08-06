package main

import (
	"log"

	"tinygo.org/x/go-llvm"
)

func init() {
	llvm.InitializeAllTargets()
	llvm.InitializeAllTargetMCs()
	llvm.InitializeAllTargetInfos()
	llvm.InitializeAllAsmParsers()
	llvm.InitializeAllAsmPrinters()
}

func main() {
	ctx := llvm.NewContext()
	defer ctx.Dispose()

	filename := "testdata/test.ll"

	buf, err := llvm.NewMemoryBufferFromFile(filename)
	if err != nil {
		log.Fatalf("unable to read %s: %s", filename, err)
	}
	//defer buf.Dispose() // throws error???

	mod, err := ctx.ParseIR(buf)
	if err != nil {
		log.Fatalf("unable to parse LLVM IR %s: %s", filename, err)
	}
	defer mod.Dispose()

	passBuilder := llvm.NewPassManagerBuilder()
	defer passBuilder.Dispose()

	passBuilder.SetOptLevel(0)
	passBuilder.SetSizeLevel(2)

	passManager := llvm.NewFunctionPassManagerForModule(mod)
	defer passManager.Dispose()

	passBuilder.PopulateFunc(passManager)

	passManager.InitializeFunc()
	for fn := mod.FirstFunction(); !fn.IsNil(); fn = llvm.NextFunction(fn) {
		passManager.RunFunc(fn)
	}
	passManager.FinalizeFunc()

	modPasses := llvm.NewPassManager()
	defer modPasses.Dispose()
	passBuilder.PopulateFunc(passManager)
	modPasses.Run(mod)

	mod.Dump()
}
