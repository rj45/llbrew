package elaboration

import (
	"github.com/rj45/llbrew/ir"
	"github.com/rj45/llbrew/ir/op"
	"github.com/rj45/llbrew/xform"
)

var _ = xform.Register(ifNonCompare,
	xform.OnlyPass(xform.Elaboration),
	xform.OnOp(op.If),
)

// ifNonCompare fixes any if instructions without a corresponding
// comparison
func ifNonCompare(it ir.Iter) {
	instr := it.Instr()
	arg := instr.Arg(0)
	if instr.Arg(0).Def().Instr().IsCompare() {
		// if already a compare, do nothing
		return
	}

	compare := it.Insert(op.Equal, arg.Type, arg, 1)
	instr.ReplaceArg(0, compare.Def(0))
	it.Changed()
}
