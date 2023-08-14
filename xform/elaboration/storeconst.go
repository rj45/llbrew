package elaboration

import (
	"github.com/rj45/llir2asm/ir"
	"github.com/rj45/llir2asm/ir/op"
	"github.com/rj45/llir2asm/xform"
)

var _ = xform.Register(addStoreConstCopies,
	xform.OnlyPass(xform.Elaboration),
	xform.OnOp(op.Store),
)

func addStoreConstCopies(it ir.Iter) {
	instr := it.Instr()

	arg := instr.Arg(1)
	if arg.IsConst() {
		cp := it.Insert(op.Copy, arg.Type, arg)
		instr.ReplaceArg(1, cp.Def(0))
	}
}
