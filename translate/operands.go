package translate

import (
	"fmt"
	"log"

	"github.com/rj45/llbrew/ir"
	"github.com/rj45/llbrew/ir/op"
	"github.com/rj45/llbrew/ir/typ"
	"tinygo.org/x/go-llvm"
)

func (trans *translator) translateAllOperands(fn llvm.Value) {
	for blk := fn.FirstBasicBlock(); !blk.IsNil(); blk = llvm.NextBasicBlock(blk) {
		for instr := blk.FirstInstruction(); !instr.IsNil(); instr = llvm.NextInstruction(instr) {
			if instr.InstructionOpcode() == llvm.PHI {
				// translated separately
				continue
			}

			ninstr := trans.instrmap[instr]
			trans.translateOperands(instr, ninstr)
		}

		trans.translateBlockArgs(blk)
	}
}

func (trans *translator) translateOperands(instr llvm.Value, ninstr *ir.Instr) {
	for i := 0; i < instr.OperandsCount(); i++ {
		operand := instr.Operand(i)
		oval := trans.valuemap[operand]
		if oval != nil {
			ninstr.InsertArg(-1, oval)
			continue
		}

		opertyp := operand.Type()
		ntyp := translateType(opertyp)
		switch opertyp.TypeKind() {
		case llvm.IntegerTypeKind:
			ninstr.InsertArg(-1, trans.fn.ValueFor(ntyp, operand.SExtValue()))
		case llvm.FunctionTypeKind:
			panic(operand.Name())
		case llvm.PointerTypeKind:
			if ntyp.Pointer().Element.Kind() == typ.FunctionKind {
				otherfn := trans.pkg.Func(operand.Name())
				ninstr.InsertArg(0, trans.fn.ValueFor(ntyp.Pointer().Element, otherfn))
			} else if !operand.IsAGlobalVariable().IsNil() {
				globalName := operand.Name()
				glob := trans.pkg.Global(globalName)
				glob.Referenced = true
				val := trans.fn.ValueFor(translateType(operand.Type()), glob)
				argidx := -1
				if ninstr.Op == op.Store {
					argidx = 0
				}
				ninstr.InsertArg(argidx, val)
			} else if operand.Opcode() == llvm.GetElementPtr {
				gep := ninstr.Func().NewInstr(op.GetElementPtr, translateType(operand.Type()))
				ninstr.Block().InsertInstr(ninstr.Index(), gep)
				trans.translateOperands(operand, gep)
				ninstr.InsertArg(0, gep.Def(0))
			} else {
				fmt.Println(operand.Opcode())
				fmt.Println(operand.Opcode())
				fmt.Println(ntyp)
				instr.Dump()
				panic(" other constant pointer")
			}
		case llvm.LabelTypeKind:
			if ninstr.Op == op.If || ninstr.Op == op.Jump {
				// branch labels handled in block translation
				continue
			}
			instr.Dump()
			fmt.Println(" encountered")

		default:
			log.Panicf("todo: other operand types: %s", ntyp.String())
		}

	}
	trans.fixupInstruction(ninstr)
}

func (trans *translator) fixupInstruction(instr *ir.Instr) {
	if instr.Op == op.Store {
		swap := false
		// operands can be backwards sometimes
		if instr.Arg(1).Type.Kind() == typ.PointerKind {
			def := instr.Arg(1).Def().Instr()
			if def != nil && def.Op == op.Alloca {
				swap = true
			}
		}
		// todo: add other situations

		if swap {
			arg := instr.Arg(0)
			instr.RemoveArg(arg)
			instr.InsertArg(-1, arg)
		}
	}
}
