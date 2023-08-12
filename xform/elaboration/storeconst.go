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

	if instr.Arg(0).IsConst() {
		cp := it.Insert(op.Copy, instr.Arg(0).Type, instr.Arg(0))
		instr.ReplaceArg(0, cp.Def(0))
	}
}
