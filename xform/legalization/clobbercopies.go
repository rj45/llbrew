package legalization

import (
	"fmt"

	"github.com/rj45/llir2asm/ir"
	"github.com/rj45/llir2asm/ir/op"
	"github.com/rj45/llir2asm/xform"
)

var _ = xform.Register(addClobberCopies,
	xform.OnlyPass(xform.Legalization),
)

// addClobberCopies adds copies for operands that get clobbered
// on two-operand architectures
func addClobberCopies(it ir.Iter) {
	instr := it.Instr()
	if !instr.Op.ClobbersArg() {
		return
	}

	if instr.NumArgs() < 1 {
		return
	}

	def := instr.Arg(0).Def()
	cand := def.Instr()
	if cand.Op != nil && cand.Op.IsCopy() && cand.NumDefs() == 1 && cand.Def(0).NumUses() == 1 && cand.Block() == instr.Block() {
		// already added the copy
		fmt.Println("copy already added for: ", instr.LongString(), "copy:", cand.LongString())
		return
	}

	cp := it.Insert(op.Copy, instr.Arg(0).Type, instr.Arg(0))
	fmt.Println("inserting copy for: ", instr.LongString(), "copy:", cp.LongString())
	instr.ReplaceArg(0, cp.Def(0))
}
