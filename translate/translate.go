package translate

import (
	"log"

	"github.com/rj45/llir2asm/ir"
	"github.com/rj45/llir2asm/ir/typ"
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
		nglob := trans.pkg.NewGlobal(glob.Name(), translateType(glob.Type()).Pointer().Element)
		if glob.OperandsCount() > 0 {
			value := glob.Operand(0)
			switch nglob.Type.Kind() {
			case typ.IntegerKind:
				nglob.Value = ir.ConstFor(value.SExtValue())
			case typ.StructKind:
				if !value.IsAConstantAggregateZero().IsNil() {
					nglob.Value = ir.ConstFor(nglob.Type.ZeroValue())
				} else {
					value.Dump()
					panic(" -- some other struct constant")
				}
			default:
				log.Panicf("unknown kind %d", value.Type().TypeKind())
			}
		}

		// todo: set global value?
	}
}
