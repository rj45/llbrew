package translate

import (
	"fmt"
	"log"
	"strings"

	"github.com/rj45/llbrew/ir"
	"github.com/rj45/llbrew/ir/typ"
	"tinygo.org/x/go-llvm"
)

type translator struct {
	mod llvm.Module

	prog  *ir.Program
	pkg   *ir.Package
	types *typ.Types

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
	trans.prog = ir.NewProgram()
	trans.pkg = &ir.Package{
		Name: "main",
	}
	trans.prog.AddPackage(trans.pkg)
	trans.types = trans.prog.Types()
}

func fixupGlobalName(name string) string {
	name = strings.TrimPrefix(name, ".")
	name = strings.Replace(name, ".", "_", -1)
	return name
}

func (trans *translator) translateGlobals() {
	for glob := trans.mod.FirstGlobal(); !glob.IsNil(); glob = llvm.NextGlobal(glob) {
		t := trans.translateType(glob.Type())
		name := fixupGlobalName(glob.Name())
		nglob := trans.pkg.NewGlobal(name, t.(*typ.Pointer).Element)
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
			case typ.ArrayKind:
				if !value.IsAConstantAggregateZero().IsNil() {
					nglob.Value = ir.ConstFor(nglob.Type.ZeroValue())
				} else if value.IsConstantString() {
					nglob.Value = ir.ConstFor(value.ConstGetAsString())
				} else {
					fmt.Println("operands", value.OperandsCount())
					fmt.Println("name", value.Name())
					fmt.Println("opcode", value.Opcode())
					fmt.Println("metadata", value.HasMetadata())
					fmt.Println("string", value.ConstGetAsString())

					value.Dump()
					panic(" --- something else!")
				}
			default:
				log.Panicf("unknown kind %d", value.Type().TypeKind())
			}
		}

		// todo: set global value?
	}
}
