package compile

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/rj45/llir2asm/arch"
	"github.com/rj45/llir2asm/html"
	"github.com/rj45/llir2asm/ir"
	"github.com/rj45/llir2asm/regalloc"
	"github.com/rj45/llir2asm/regalloc/verify"
	"github.com/rj45/llir2asm/translate"
	"github.com/rj45/llir2asm/xform"
	"tinygo.org/x/go-llvm"

	_ "github.com/rj45/llir2asm/arch/rj32"

	_ "github.com/rj45/llir2asm/xform/cleanup"
	_ "github.com/rj45/llir2asm/xform/elaboration"
	_ "github.com/rj45/llir2asm/xform/finishing"
	_ "github.com/rj45/llir2asm/xform/legalization"
	_ "github.com/rj45/llir2asm/xform/lowering"
	_ "github.com/rj45/llir2asm/xform/simplification"
)

func init() {
	llvm.InitializeAllTargets()
	llvm.InitializeAllTargetMCs()
	llvm.InitializeAllTargetInfos()
	llvm.InitializeAllAsmParsers()
	llvm.InitializeAllAsmPrinters()
}

type Compiler struct {
	Filename string

	OptSpeed int
	OptSize  int

	DumpLL  io.WriteCloser
	DumpIR  io.WriteCloser
	DumpSSA string

	ctx       llvm.Context
	mod       llvm.Module
	initLevel int

	prog *ir.Program
}

func (c *Compiler) dispose() {
	// yuck :-/
	if c.initLevel > 1 {
		c.mod.Dispose()
	}

	if c.initLevel > 0 {
		c.ctx.Dispose()
	}
}

func (c *Compiler) Compile(filename string) (*ir.Program, error) {

	c.Filename = filename

	err := c.loadIR()
	if err != nil {
		return nil, err
	}
	defer c.dispose()

	// some function attributes prevent optimizations from running
	c.fixFunctionAttributes()

	// run the LLVM optimizer
	c.optimize()

	fmt.Fprint(c.DumpLL, c.mod.String())

	// convert the LLVM program to our own IR
	c.prog = translate.Translate(c.mod)

	// convert the program into assembly
	err = c.transformProgram()
	if err != nil {
		return nil, err
	}

	c.prog.Emit(c.DumpIR, ir.SSAString{})

	return c.prog, nil
}

func (c *Compiler) transformProgram() error {
	arch.SetArch("rj32")

	for _, pkg := range c.prog.Packages() {
		for _, fn := range pkg.Funcs() {
			var w dumper = nopDumper{}
			if c.DumpSSA != "" && strings.Contains(fn.FullName, c.DumpSSA) {
				w = html.NewHTMLWriter("ssa.html", fn)
			}
			defer w.Close()

			w.WritePhase("initial", "initial")

			xform.Transform(xform.Elaboration, fn)
			w.WritePhase("elaboration", "elaboration")

			xform.Transform(xform.Simplification, fn)
			w.WritePhase("simplification", "simplification")

			xform.Transform(xform.Lowering, fn)
			w.WritePhase("lowering", "lowering")

			xform.Transform(xform.Legalization, fn)
			w.WritePhase("legalization", "legalization")

			ra := regalloc.NewRegAlloc(fn)
			err := ra.Allocate()
			// if *debug {
			// 	regalloc.WriteGraphvizCFG(ra)
			// 	regalloc.DumpLivenessChart(ra)
			// 	regalloc.WriteGraphvizInterferenceGraph(ra)
			// 	regalloc.WriteGraphvizLivenessGraph(ra)
			// }
			if err != nil {
				return fmt.Errorf("register allocation failed: %w", err)
			}
			errs := verify.Verify(fn)
			for _, err := range errs {
				log.Printf("verification error: %s\n", err)
			}
			if len(errs) > 0 {
				log.Fatal("verification failed")
			}
			w.WritePhase("regalloc", "regalloc")

			xform.Transform(xform.CleanUp, fn)
			w.WritePhase("cleanup", "cleanup")

			xform.Transform(xform.Finishing, fn)
			w.WritePhase("finishing", "finishing")
		}
	}
	return nil
}
