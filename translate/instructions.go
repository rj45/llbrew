package translate

import (
	"github.com/rj45/llir2asm/ir"
	"github.com/rj45/llir2asm/ir/op"
	"tinygo.org/x/go-llvm"
)

func (trans *translator) translateInstructions(fn llvm.Value) {
	for blk := fn.FirstBasicBlock(); !blk.IsNil(); blk = llvm.NextBasicBlock(blk) {
		nblk := trans.blkmap[blk]

		for instr := blk.FirstInstruction(); !instr.IsNil(); instr = llvm.NextInstruction(instr) {
			if instr.InstructionOpcode() == llvm.PHI {
				def := trans.fn.NewValue(translateType(instr.Type()))
				nblk.AddDef(def)
				trans.valuemap[instr] = def

				continue
			}

			opcode := translateOpcode(instr.InstructionOpcode())
			if opcode == op.Invalid {
				switch instr.InstructionOpcode() {
				case llvm.ICmp:
					opcode = intPredicateMap[instr.IntPredicate()]
				case llvm.Br:
					if instr.OperandsCount() == 1 {
						opcode = op.Jump
					} else if instr.OperandsCount() == 3 {
						opcode = op.If
					} else {
						instr.Dump()
						panic(" unknown branch format")
					}
				}
			}

			ninstr := trans.fn.NewInstr(opcode, translateType(instr.Type()))
			nblk.InsertInstr(-1, ninstr)
			trans.instrmap[instr] = ninstr
			if ninstr.NumDefs() > 0 {
				trans.valuemap[instr] = ninstr.Def(0)
			}
		}
	}
}

func translateOpcode(opcode llvm.Opcode) ir.Op {
	return opcodeMap[opcode]
}

var intPredicateMap = [...]op.Op{
	llvm.IntEQ:  op.Equal,
	llvm.IntNE:  op.NotEqual,
	llvm.IntSLT: op.Less,
	llvm.IntSLE: op.LessEqual,
	llvm.IntSGE: op.GreaterEqual,
	llvm.IntSGT: op.Greater,
	llvm.IntULT: op.ULess,
	llvm.IntULE: op.ULessEqual,
	llvm.IntUGE: op.UGreaterEqual,
	llvm.IntUGT: op.UGreater,
}

var opcodeMap = [...]op.Op{
	llvm.Ret:         op.Ret,
	llvm.Switch:      op.Switch,
	llvm.IndirectBr:  op.IndirectBr,
	llvm.Invoke:      op.Invoke,
	llvm.Unreachable: op.Unreachable,

	// Standard Binary Operators
	llvm.Add:  op.Add,
	llvm.FAdd: op.FAdd,
	llvm.Sub:  op.Sub,
	llvm.FSub: op.FSub,
	llvm.Mul:  op.Mul,
	llvm.FMul: op.FMul,
	llvm.UDiv: op.UDiv,
	llvm.SDiv: op.SDiv,
	llvm.FDiv: op.FDiv,
	llvm.URem: op.URem,
	llvm.SRem: op.SRem,
	llvm.FRem: op.FRem,

	// Logical Operators
	llvm.Shl:  op.Shl,
	llvm.LShr: op.LShr,
	llvm.AShr: op.AShr,
	llvm.And:  op.And,
	llvm.Or:   op.Or,
	llvm.Xor:  op.Xor,

	// Memory Operators
	llvm.Alloca:        op.Alloca,
	llvm.Load:          op.Load,
	llvm.Store:         op.Store,
	llvm.GetElementPtr: op.GetElementPtr,

	// Cast Operators
	llvm.Trunc:    op.Trunc,
	llvm.ZExt:     op.ZExt,
	llvm.SExt:     op.SExt,
	llvm.FPToUI:   op.FPToUI,
	llvm.FPToSI:   op.FPToSI,
	llvm.UIToFP:   op.UIToFP,
	llvm.SIToFP:   op.SIToFP,
	llvm.FPTrunc:  op.FPTrunc,
	llvm.FPExt:    op.FPExt,
	llvm.PtrToInt: op.PtrToInt,
	llvm.IntToPtr: op.IntToPtr,
	llvm.BitCast:  op.BitCast,

	// Other Operators
	llvm.FCmp:   op.FCmp,
	llvm.PHI:    op.PHI,
	llvm.Call:   op.Call,
	llvm.Select: op.Select,

	// UserOp1
	// UserOp2
	llvm.VAArg:          op.VAArg,
	llvm.ExtractElement: op.ExtractElement,
	llvm.InsertElement:  op.InsertElement,
	llvm.ShuffleVector:  op.ShuffleVector,
	llvm.ExtractValue:   op.ExtractValue,
	llvm.InsertValue:    op.InsertValue,
}
