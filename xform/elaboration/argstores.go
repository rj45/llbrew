package elaboration

import (
	"github.com/rj45/llbrew/ir"
	"github.com/rj45/llbrew/ir/op"
	"github.com/rj45/llbrew/ir/reg"
	"github.com/rj45/llbrew/xform"
)

var _ = xform.Register(callArgStores,
	xform.OnlyPass(xform.Elaboration),
	xform.OnOp(op.Copy),
)

// callArgStores adds stores for args that are passed on the stack
func callArgStores(it ir.Iter) {
	instr := it.Instr()

	for a := 0; a < instr.NumArgs(); a++ {
		arg := instr.Arg(a)
		def := instr.Def(a)

		if def.OnStack() {
			// remove arg & def from the copy
			instr.RemoveArg(arg)
			instr.RemoveDef(def)

			it.Next()
			ptrtyp := instr.Func().Types().PointerType(def.Type, 0)
			add := it.Insert(op.Add, ptrtyp, reg.SP, def)
			it.Insert(op.Store, def.Type, add.Def(0), arg)
			it.Prev()
			it.Prev()
		}
	}
}
