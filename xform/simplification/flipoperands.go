package simplification

import (
	"github.com/rj45/llbrew/ir"
	"github.com/rj45/llbrew/xform"
)

var _ = xform.Register(swapOperands,
	xform.OnlyPass(xform.Simplification),
)

func swapOperands(it ir.Iter) {
	instr := it.Instr()
	if !instr.Op.IsCommutative() {
		return
	}

	a := followUp(instr.Arg(0))
	b := followUp(instr.Arg(1))
	dest := followDown(instr.Def(0))

	flip := false
	if a.IsConst() && !b.IsConst() {
		flip = true
	}

	if b.InReg() && !a.InReg() {
		flip = true
	}

	if dest.InReg() && b.InReg() && dest.Reg() == b.Reg() {
		flip = true
	}

	if flip {
		arg := instr.Arg(0)
		instr.RemoveArg(arg)
		instr.InsertArg(-1, arg)
		it.Changed()
	}
}

func followUp(val *ir.Value) *ir.Value {
	def := val.Def()
	if def != nil && def.IsInstr() {
		instr := def.Instr()
		if !instr.Op.IsCopy() {
			return val
		}
		for d := 0; d < instr.NumDefs(); d++ {
			if instr.Def(d) == val {
				return followUp(instr.Arg(d))
			}
		}
	}
	return val
}

func followDown(val *ir.Value) *ir.Value {
	if val.NumUses() > 1 || val.NumUses() == 0 {
		return val
	}
	use := val.Use(0)
	if use.IsInstr() {
		instr := use.Instr()
		if instr.Op.IsCopy() {
			idx := instr.ArgIndex(val)
			return followDown(instr.Def(idx))
		}
	}
	return val
}
