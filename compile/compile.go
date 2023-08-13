package compile

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/rj45/llir2asm/arch"
	"github.com/rj45/llir2asm/asm"
	"github.com/rj45/llir2asm/customasm"
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

	DumpLL  string
	DumpIR  string
	DumpSSA string
	OutFile string
	BinFile string

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

func (c *Compiler) Compile(filename string) error {

	c.Filename = filename

	err := c.loadIR()
	if err != nil {
		return err
	}
	defer c.dispose()

	// some function attributes prevent optimizations from running
	c.fixFunctionAttributes()

	// run the LLVM optimizer
	c.optimize()

	// re-split critical edges merged in optimization
	c.splitCriticalEdges()

	if c.DumpLL != "" {
		dump := createFile(c.DumpLL)
		defer dump.Close()

		fmt.Fprint(dump, c.mod.String())
	}

	// convert the LLVM program to our own IR
	c.prog = translate.Translate(c.mod)

	// convert the program into assembly
	err = c.transformProgram()
	if err != nil {
		return err
	}

	if c.DumpIR != "" {
		dump := createFile(c.DumpIR)
		defer dump.Close()

		c.prog.Emit(dump, ir.SSAString{})
	}

	out := createFile(c.OutFile)
	defer out.Close()
	outwr := io.Writer(out)

	var outbuf bytes.Buffer
	if c.BinFile != "" {
		outwr = io.MultiWriter(out, &outbuf)
	}

	asm.Emit(outwr, asm.CustomASM{}, c.prog)

	if c.BinFile != "" {
		err := customasm.Assemble(outbuf.Bytes(), c.BinFile)
		if err != nil {
			return err
		}
	}

	return nil
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
			err := ra.CheckInput()
			if err != nil {
				return fmt.Errorf("register allocation pre-check failed: %w", err)
			}
			err = ra.Allocate()
			// if *debug {
			// 	regalloc.WriteGraphvizCFG(ra)
			// 	regalloc.DumpLivenessChart(ra)
			// 	regalloc.WriteGraphvizInterferenceGraph(ra)
			// 	regalloc.WriteGraphvizLivenessGraph(ra)
			// }
			if err != nil {
				return fmt.Errorf("register allocation failed: %w", err)
			}
			w.WritePhase("regalloc", "regalloc")
			errs := verify.Verify(fn)
			for _, err := range errs {
				log.Printf("verification error: %s\n", err)
			}
			if len(errs) > 0 {
				w.Close()
				log.Fatal("verification failed")
			}

			xform.Transform(xform.CleanUp, fn)
			w.WritePhase("cleanup", "cleanup")

			xform.Transform(xform.Finishing, fn)
			w.WritePhase("finishing", "finishing")
		}
	}
	return nil
}
