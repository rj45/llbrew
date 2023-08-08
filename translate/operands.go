package translate

import (
	"log"

	"github.com/rj45/llir2asm/ir/typ"
	"tinygo.org/x/go-llvm"
)

func (trans *translator) translateOperands(fn llvm.Value) {
	for blk := fn.FirstBasicBlock(); !blk.IsNil(); blk = llvm.NextBasicBlock(blk) {
		for instr := blk.FirstInstruction(); !instr.IsNil(); instr = llvm.NextInstruction(instr) {
			ninstr := trans.instrmap[instr]
			for i := 0; i < instr.OperandsCount(); i++ {
				operand := instr.Operand(i)
				oval := trans.valuemap[operand]
				if oval == nil {
					if operand.IsConstant() {
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
							}
						default:
							log.Panicf("todo: other constant types: %s", ntyp)
						}

					} else {

						log.Panicf("todo: other operand types %d", operand.IntrinsicID())
					}
				} else {
					ninstr.InsertArg(-1, oval)
				}
			}
		}
	}
}
