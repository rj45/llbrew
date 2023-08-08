package compile

import (
	"fmt"

	"tinygo.org/x/go-llvm"
)

func (c *Compiler) loadIR() error {
	c.ctx = llvm.NewContext()
	c.initLevel = 1

	buf, err := llvm.NewMemoryBufferFromFile(c.Filename)
	if err != nil {
		c.dispose()
		return fmt.Errorf("unable to read %s: %w", c.Filename, err)
	}
	//defer buf.Dispose() // throws error???

	mod, err := c.ctx.ParseIR(buf)
	if err != nil {
		c.dispose()
		return fmt.Errorf("unable to parse LLVM IR %s: %w", c.Filename, err)
	}
	c.mod = mod
	c.initLevel = 2

	return nil
}

func (c *Compiler) fixFunctionAttributes() {
	for fn := c.mod.FirstFunction(); !fn.IsNil(); fn = llvm.NextFunction(fn) {
		// remove disabling optimizations
		fn.RemoveEnumFunctionAttribute(llvm.AttributeKindID("optnone"))

		// set any size optimizations
		if c.OptSize >= 1 {
			fn.AddFunctionAttr(c.ctx.CreateEnumAttribute(llvm.AttributeKindID("optsize"), 0))
		}
		if c.OptSize >= 2 {
			fn.AddFunctionAttr(c.ctx.CreateEnumAttribute(llvm.AttributeKindID("minsize"), 0))
		}
	}
}

func (c *Compiler) optimize() {
	// {
	// 	pm := llvm.NewPassManager()
	// 	defer pm.Dispose()
	// 	pm.AddGlobalDCEPass()
	// 	pm.AddGlobalOptimizerPass()
	// 	pm.AddIPSCCPPass()
	// 	pm.AddInstructionCombiningPass()
	// 	pm.AddAggressiveDCEPass()
	// 	pm.AddFunctionAttrsPass()
	// 	pm.AddGlobalDCEPass()
	// 	if !pm.Run(mod) {
	// 		fmt.Println("initial pass did nothing!")
	// 	}
	// }

	passBuilder := llvm.NewPassManagerBuilder()
	defer passBuilder.Dispose()

	passBuilder.SetOptLevel(c.OptSpeed)
	passBuilder.SetSizeLevel(c.OptSize)

	{
		passManager := llvm.NewFunctionPassManagerForModule(c.mod)
		defer passManager.Dispose()
		passBuilder.PopulateFunc(passManager)

		passManager.InitializeFunc()
		for fn := c.mod.FirstFunction(); !fn.IsNil(); fn = llvm.NextFunction(fn) {
			passManager.RunFunc(fn)
		}
		passManager.FinalizeFunc()
	}

	{
		modPasses := llvm.NewPassManager()
		defer modPasses.Dispose()
		passBuilder.Populate(modPasses)
		modPasses.Run(c.mod)
	}

	// err = llvm.VerifyModule(mod, llvm.PrintMessageAction)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
