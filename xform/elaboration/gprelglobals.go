package elaboration

import (
	"github.com/rj45/llir2asm/ir"
	"github.com/rj45/llir2asm/ir/op"
	"github.com/rj45/llir2asm/ir/reg"
	"github.com/rj45/llir2asm/ir/typ"
	"github.com/rj45/llir2asm/xform"
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
