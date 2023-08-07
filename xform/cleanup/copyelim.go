package cleanup

import (
	"github.com/rj45/llir2asm/ir"
	"github.com/rj45/llir2asm/ir/op"
	"github.com/rj45/llir2asm/xform"
)

var _ = xform.Register(copyElim,
	xform.OnlyPass(xform.CleanUp),
	xform.OnOp(op.Copy),
)

// copyElim eliminates any copies to the same register.
// Note: this destroys SSA, so make sure it's no longer needed
// when this runs.
func copyElim(it ir.Iter) {
	instr := it.Instr()

	for i := 0; i < instr.NumDefs(); i++ {
		def := instr.Def(i)
		arg := instr.Arg(i)

		if def.Reg() == arg.Reg() {
			def.ReplaceUsesWith(arg)
			instr.RemoveArg(arg)
			instr.RemoveDef(def)
			i--
			it.Changed()
		}
	}

	if instr.NumArgs() == 0 {
		it.Remove()
		if it.Instr() == nil {
			panic("broke iter")
		}
	}
}
