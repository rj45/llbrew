package elaboration

import (
	"github.com/rj45/llir2asm/ir"
	"github.com/rj45/llir2asm/ir/op"
	"github.com/rj45/llir2asm/ir/reg"
	"github.com/rj45/llir2asm/xform"
)

var _ = xform.Register(loadGPRelGlobals,
	xform.OnlyPass(xform.Elaboration),
	xform.OnOp(op.Load),
)

var _ = xform.Register(storeGPRelGlobals,
	xform.OnlyPass(xform.Elaboration),
	xform.OnOp(op.Store),
)

func loadGPRelGlobals(it ir.Iter) {
	instr := it.Instr()

	if instr.Arg(0).IsConst() {
		instr.InsertArg(0, instr.Func().ValueFor(0, reg.GP))
	}
}

func storeGPRelGlobals(it ir.Iter) {
	instr := it.Instr()

	if instr.Arg(1).IsConst() {
		instr.InsertArg(1, instr.Func().ValueFor(0, reg.GP))
	}
}
