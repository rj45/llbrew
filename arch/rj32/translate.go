package rj32

import (
	"log"

	"github.com/rj45/llir2asm/ir"
	"github.com/rj45/llir2asm/ir/op"
)

var directTranslate = map[op.Op]Opcode{
	op.Ret:  Return,
	op.Jump: Jump,
}

var twoOperandTranslations = map[op.Op]Opcode{
	op.Add:  Add,
	op.Sub:  Sub,
	op.And:  And,
	op.Or:   Or,
	op.Xor:  Xor,
	op.Shl:  Shl,
	op.LShr: Shr,
	op.AShr: Asr,
}

// var oneOperandTranslations = map[op.Op]Opcode{
// 	op.: Not,
// 	op.Negate: Neg,
// }

var branches = map[op.Op]Opcode{
	op.Equal:         IfEq,
	op.NotEqual:      IfNe,
	op.Less:          IfLt,
	op.LessEqual:     IfLe,
	op.Greater:       IfGt,
	op.GreaterEqual:  IfGe,
	op.ULess:         IfUlt,
	op.ULessEqual:    IfUle,
	op.UGreater:      IfUgt,
	op.UGreaterEqual: IfUge,
}

func translate(it ir.Iter) {
	instr := it.Instr()
	originalOp := instr.Op
	switch instr.Op {
	case op.Copy:
		// copy is done in the finishing stage, after register allocation
	case op.Ret, op.Jump:
		it.Update(directTranslate[instr.Op.(op.Op)], 0, instr.Args())
	case op.Add, op.Sub, op.And, op.Or, op.Xor, op.Shl, op.AShr, op.LShr:
		it.Update(twoOperandTranslations[instr.Op.(op.Op)], 0, instr.Args())
	// case op.Not, op.Negate:
	// 	it.Update(oneOperandTranslations[instr.Op.(op.Op)], nil, instr.Args())
	case op.Equal, op.NotEqual, op.Less, op.LessEqual, op.Greater, op.GreaterEqual:
		def := instr.Def(0)
		if def.NumUses() > 1 || def.Use(0).Instr().Op != op.If {
			log.Panicf("Lone comparison not tied to If %s", instr.LongString())
		}
	case op.If:
		compare := instr.Arg(0).Def().Instr()
		if !compare.IsCompare() {
			log.Panicf("expecting if to have compare, but instead had: %s", compare.LongString())
		}

		branchOp := branches[compare.Op.(op.Op)]
		if branchOp == 0 {
			log.Panicf("failed to translate compare %s", compare.Op.(op.Op))
		}
		it.Update(branchOp, 0, compare.Args())
		if compare.Def(0).NumUses() == 0 {
			it.RemoveInstr(compare)
		}
	case op.Load, op.Store:
		translateLoadStore(it)
	case op.Call:
		instr.Op = Call
		it.Changed()
	// case op.Panic:
	// 	instr.Op = Error
	// 	it.Changed()
	// case op.InlineAsm:
	// 	// handled elsewhere
	default:
		// if _, ok := instr.Op.(op.Op); ok {
		// 	log.Panicf("Unknown instruction: %s", instr.LongString())
		// }
	}
	if it.Instr() == nil {
		log.Panicf("translating %s from %s left iter in bad state", originalOp, instr.LongString())
	}
}

// translateLoadStore mainly translates the function pro-/epilogue
func translateLoadStore(it ir.Iter) {
	instr := it.Instr()
	switch instr.Op {
	case op.Add:
		if instr.NumArgs() == 2 && instr.Def(0).Reg() == instr.Arg(0).Reg() {
			it.Update(Add, 0, instr.Args())
		}
	case op.Load:
		it.Update(Load, instr.Def(0).Type, instr.Args())
	case op.Store:
		it.Update(Store, 0, instr.Arg(0), instr.Arg(1), instr.Arg(2))
	}
}

func translateCopies(it ir.Iter) {
	instr := it.Instr()

	it.Update(Move, instr.Def(0).Type, instr.Args())
}
