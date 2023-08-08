package translate

import (
	"github.com/rj45/llir2asm/ir"
	"tinygo.org/x/go-llvm"
)

type translator struct {
	mod llvm.Module

	prog *ir.Program
	pkg  *ir.Package

	fn       *ir.Func
	blkmap   map[llvm.BasicBlock]*ir.Block
	valuemap map[llvm.Value]*ir.Value
	instrmap map[llvm.Value]*ir.Instr
}

func Translate(mod llvm.Module) *ir.Program {
	t := &translator{mod: mod}

	t.initProgram()
	t.translateGlobals()
	t.translateFunctions()
	return t.prog
}

func (trans *translator) initProgram() {
	trans.prog = &ir.Program{}
	trans.pkg = &ir.Package{
		Name: "main",
	}
	trans.prog.AddPackage(trans.pkg)
}

func (trans *translator) translateGlobals() {
	for glob := trans.mod.FirstGlobal(); !glob.IsNil(); glob = llvm.NextGlobal(glob) {
		trans.pkg.NewGlobal(glob.Name(), translateType(glob.Type()))
		// todo: set global value?
	}
}
