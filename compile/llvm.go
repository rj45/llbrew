package compile

import (
	"fmt"
	"log"
	"strings"

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
	pbo := llvm.NewPassBuilderOptions()
	defer pbo.Dispose()

	// disable vectorization... I am not going to implement SIMD support
	pbo.SetLoopInterleaving(false)
	pbo.SetLoopVectorization(false)
	pbo.SetSLPVectorization(false)

	defaultPass := "default<O0>"

	if c.OptSize > 0 {
		switch c.OptSize {
		case 1:
			defaultPass = "default<Os>"
		case 2:
			defaultPass = "default<Oz>"
		}
	} else {
		switch c.OptSpeed {
		case 1:
			defaultPass = "default<O1>"
		case 2:
			defaultPass = "default<O2>"
		case 3:
			defaultPass = "default<O3>"
		}
	}

	targ, err := llvm.GetTargetFromTriple(c.mod.Target())
	if err != nil {
		log.Fatal(err)
	}

	tm := targ.CreateTargetMachine(c.mod.Target(), "", "",
		llvm.CodeGenLevelAggressive, llvm.RelocStatic, llvm.CodeModelSmall)

	passes := []string{
		defaultPass,

		// this pass makes it so exception invokes are converted to calls instead
		// thus making less work for us. This can be removed once invoke is supported.
		"lowerinvoke",

		// This lowers switch statements to a series of branches. This can be removed
		// when jump tables are supported.
		"lowerswitch",

		// Merge return makes it so only one return statement exists, thus the function
		// epilogue only gets generated once
		"mergereturn",

		// Simplify the control flow graph. lowerinvoke needs this to clean up the dead branches.
		"simplifycfg",

		// breaking critical edges is important to do last, since it's required by the
		// register allocator.
		"break-crit-edges",
	}

	err = c.mod.RunPasses(strings.Join(passes, ","), tm, pbo)
	if err != nil {
		log.Fatal(err)
	}

	// err = llvm.VerifyModule(mod, llvm.PrintMessageAction)
	// if err != nil {
	// 	log.Fatal(err)
	// }

}
