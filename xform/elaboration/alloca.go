package elaboration

import (
	"github.com/rj45/llbrew/ir"
	"github.com/rj45/llbrew/ir/op"
	"github.com/rj45/llbrew/ir/reg"
	"github.com/rj45/llbrew/ir/typ"
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

	t := val.Type
	if t.Kind() == typ.PointerKind {
		t = t.(*typ.Pointer).Element
	}

	size := t.SizeOf()

	if num == 0 {
		num++
	}

	addr := instr.Func().AllocSpillStorage(num * size)

	val.SetSpillAddress(addr)

	uses := make([]*ir.User, val.NumUses())
	for u := 0; u < val.NumUses(); u++ {
		uses[u] = val.Use(u)
	}
	for _, use := range uses {
		if use.IsInstr() {
			uinstr := use.Instr()
			ublk := uinstr.Block()
			if ublk == nil {
				// todo: investigate why this happens, it shouldn't
				continue
			}
			add := uinstr.Func().NewInstr(op.Add, val.Type, reg.SP, val)
			ublk.InsertInstr(uinstr.Index(), add)
			uinstr.ReplaceArg(uinstr.ArgIndex(val), add.Def(0))
		} else {
			panic("other use")
		}
	}
	it.RemoveInstr(instr)
	it.Changed()
}
