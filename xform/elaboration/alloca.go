package elaboration

import (
	"github.com/rj45/llbrew/ir"
	"github.com/rj45/llbrew/ir/op"
	"github.com/rj45/llbrew/ir/reg"
	"github.com/rj45/llbrew/xform"
)

var _ = xform.Register(alloca,
	xform.OnlyPass(xform.Elaboration),
	xform.OnOp(op.Alloca),
)

func alloca(it ir.Iter) {
	instr := it.Instr()

	val := instr.Def(0)

	if val.InSpillArea() {
		return
	}

	num, ok := ir.IntValue(instr.Arg(0).Const())
	if !ok {
		panic("expecting first arg of alloca to be an int const")
	}
	size := val.Type.SizeOf()

	addr := instr.Func().AllocSpillStorage(num * size)

	val.SetSpillAddress(addr)

	uses := make([]*ir.Instr, val.NumUses())
	for u := 0; u < val.NumUses(); u++ {
		uses[u] = val.Use(u).Instr()
	}
	for _, uinstr := range uses {
		ublk := uinstr.Block()
		add := uinstr.Func().NewInstr(op.Add, val.Type, reg.SP, val)
		ublk.InsertInstr(uinstr.Index(), add)
		uinstr.ReplaceArg(uinstr.ArgIndex(val), add.Def(0))
	}
	it.RemoveInstr(instr)
	it.Changed()
}
