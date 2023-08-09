package translate

import (
	"log"

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

			op := translateOpcode(instr.InstructionOpcode())
			ninstr := trans.fn.NewInstr(op, translateType(instr.Type()))
			nblk.InsertInstr(-1, ninstr)
			trans.instrmap[instr] = ninstr
			if ninstr.NumDefs() > 0 {
				trans.valuemap[instr] = ninstr.Def(0)
			}
		}
	}
}

func translateOpcode(opcode llvm.Opcode) ir.Op {
	op2 := opcodeMap[opcode]
	if op2 == op.Invalid {
		log.Panicf("bad opcode %d", opcode)
	}
	return op2
}

var opcodeMap = [...]op.Op{
	llvm.Ret:         op.Ret,
	llvm.Br:          op.Br,
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
	llvm.ICmp:   op.ICmp,
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
