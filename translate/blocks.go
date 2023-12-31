package translate

import (
	"fmt"

	"github.com/rj45/llbrew/ir"
	"tinygo.org/x/go-llvm"
)

func (trans *translator) translateBlocks(fn llvm.Value) {
	first := true
	for blk := fn.FirstBasicBlock(); !blk.IsNil(); blk = llvm.NextBasicBlock(blk) {
		nblk := trans.fn.NewBlock()
		trans.blkmap[blk] = nblk

		if first {
			first = false
			for param := fn.FirstParam(); !param.IsNil(); param = llvm.NextParam(param) {
				pval := trans.fn.NewValue(trans.translateType(param.Type()))
				nblk.AddDef(pval)
				trans.valuemap[param] = pval
			}
			nblk.SetCallRegisters(false, ir.ParamSlot)
		}

		trans.fn.InsertBlock(-1, nblk)
	}

	for blk := fn.FirstBasicBlock(); !blk.IsNil(); blk = llvm.NextBasicBlock(blk) {
		// handle preds & succs
		nblk := trans.blkmap[blk]

		br := blk.LastInstruction()
		if br.IsNil() {
			continue
		}
		if br.InstructionOpcode() != llvm.Br {
			if br.InstructionOpcode() != llvm.Ret {
				fmt.Println(opcodeMap[br.InstructionOpcode()].String())
				panic("assumed it was ret, was not")
			}
			continue // return
		}

		for i := br.OperandsCount() - 1; i >= 0; i-- {
			operand := br.Operand(i)
			if operand.Type().TypeKind() == llvm.LabelTypeKind {
				for b := fn.FirstBasicBlock(); !b.IsNil(); b = llvm.NextBasicBlock(b) {
					if b.AsValue() == operand {
						bblk := trans.blkmap[b]
						nblk.AddSucc(bblk)
						bblk.AddPred(nblk)
					}
				}
			}
		}
	}
}

func (trans *translator) translateBlockArgs(blk llvm.BasicBlock) {
	nblk := trans.blkmap[blk]
	for s := 0; s < nblk.NumSuccs(); s++ {
		nsucc := nblk.Succ(s)

		var succ llvm.BasicBlock
		for b := blk.Parent().FirstBasicBlock(); !b.IsNil(); b = llvm.NextBasicBlock(b) {
			if trans.blkmap[b] == nsucc {
				succ = b
				break
			}
		}

		for instr := succ.FirstInstruction(); !instr.IsNil(); instr = llvm.NextInstruction(instr) {
			if instr.InstructionOpcode() != llvm.PHI {
				break
			}

			predIndex := 0
			for i := 0; i < instr.IncomingCount(); i++ {
				if instr.IncomingBlock(i) == blk {
					predIndex = i
					break
				}
			}

			val := instr.Operand(predIndex)
			arg := trans.valuemap[val]

			if val.IsConstant() {
				t := val.Type()
				switch t.TypeKind() {
				case llvm.IntegerTypeKind:
					arg = nblk.Func().ValueFor(trans.translateType(t), val.SExtValue())
				default:
					instr.Dump()
					panic(" unimpl constant kind")
				}
			}

			nblk.InsertArg(-1, arg)
		}
	}
}

// func reverseSSASuccessorSort(block llvm.BasicBlock, list []llvm.BasicBlock, visited map[llvm.BasicBlock]bool) []llvm.BasicBlock {
// 	visited[block] = true

// 	br := block.LastInstruction()
// 	if br.InstructionOpcode() != llvm.Br {
// 		return append(list, block)
// 	}

// 	for i := br.OperandsCount() - 1; i >= 0; i-- {

// 		succ := br.Operand(i).AsBasicBlock()

// 		if succ.IsNil() {
// 			continue
// 		}

// 		if !visited[succ] {
// 			list = reverseSSASuccessorSort(succ, list, visited)
// 		}
// 	}

// 	return append(list, block)
// }
