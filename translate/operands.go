package translate

import (
	"fmt"
	"log"

	"github.com/rj45/llir2asm/ir/op"
	"github.com/rj45/llir2asm/ir/typ"
	"tinygo.org/x/go-llvm"
)

func (trans *translator) translateOperands(fn llvm.Value) {
	for blk := fn.FirstBasicBlock(); !blk.IsNil(); blk = llvm.NextBasicBlock(blk) {
		for instr := blk.FirstInstruction(); !instr.IsNil(); instr = llvm.NextInstruction(instr) {
			if instr.InstructionOpcode() == llvm.PHI {
				// translated separately
				continue
			}

			ninstr := trans.instrmap[instr]
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
					} else {
						instr.Dump()
						panic(" other constant pointer")
					}
				case llvm.LabelTypeKind:
					if ninstr.Op == op.Br {
						// branch labels handled in block translation
						continue
					}
					instr.Dump()
					fmt.Println(" encountered")

				default:
					log.Panicf("todo: other operand types: %s", ntyp.String())
				}

			}
		}

		trans.translateBlockArgs(blk)
	}
}
