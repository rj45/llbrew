package elaboration

import (
	"github.com/rj45/llbrew/ir"
	"github.com/rj45/llbrew/ir/op"
	"github.com/rj45/llbrew/ir/reg"
	"github.com/rj45/llbrew/ir/typ"
	"github.com/rj45/llbrew/xform"
)

var _ = xform.Register(gpRelGlobals,
	xform.OnlyPass(xform.Elaboration),
)

func gpRelGlobals(it ir.Iter) {
	instr := it.Instr()

	if instr.Op == op.Add && instr.Arg(0).Reg() == reg.GP {
		// already done
		return
	}

	// insert add instructions on all global references.
	for a := 0; a < instr.NumArgs(); a++ {
		arg := instr.Arg(a)
		if global, ok := ir.GlobalValue(arg.Const()); ok {
			gp := instr.Func().ValueFor(typ.VoidPointer(), reg.GP)
			add := it.Insert(op.Add, arg.Type, gp, global)
			instr.ReplaceArg(a, add.Def(0))
		}
	}
}
