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

		switch opertyp.TypeKind() {
		case llvm.IntegerTypeKind:
			ntyp := trans.translateType(opertyp)
			ninstr.InsertArg(-1, trans.fn.ValueFor(ntyp, operand.SExtValue()))
		case llvm.FunctionTypeKind:
			panic(operand.Name())
		case llvm.PointerTypeKind:
			ntyp := trans.translateType(opertyp)
			if ntyp.(*typ.Pointer).Element.Kind() == typ.FunctionKind {
				otherfn := trans.pkg.Func(operand.Name())
				ninstr.InsertArg(0, trans.fn.ValueFor(ntyp.(*typ.Pointer).Element, otherfn))
			} else if !operand.IsAGlobalVariable().IsNil() {
				globalName := fixupGlobalName(operand.Name())
				glob := trans.pkg.Global(globalName)
				glob.Referenced = true
				val := trans.fn.ValueFor(trans.translateType(operand.Type()), glob)
				ninstr.InsertArg(-1, val)
			} else if operand.Opcode() == llvm.GetElementPtr {
				gep := ninstr.Func().NewInstr(op.GetElementPtr, trans.translateType(operand.Type()))
				ninstr.Block().InsertInstr(ninstr.Index(), gep)
				trans.translateOperands(operand, gep)
				ninstr.InsertArg(-1, gep.Def(0))
			} else if operand.Opcode() == llvm.IntToPtr {
				itop := ninstr.Func().NewInstr(op.IntToPtr, trans.translateType(operand.Type()))
				ninstr.Block().InsertInstr(ninstr.Index(), itop)
				trans.translateOperands(operand, itop)
				ninstr.InsertArg(-1, itop.Def(0))
			} else {
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

		case llvm.VectorTypeKind:

		default:
			ntyp := trans.translateType(opertyp)
			log.Panicf("todo: other operand types: %s %s", ntyp.String(), opertyp.TypeKind())
		}

	}
	trans.fixupInstruction(ninstr)
}

func (trans *translator) fixupInstruction(instr *ir.Instr) {
	// stores are backwards from what we expect, so swap the args
	if instr.Op == op.Store {
		arg := instr.Arg(0)
		instr.RemoveArg(arg)
		instr.InsertArg(-1, arg)
	}
}
