package simplification

import (
	"github.com/rj45/llbrew/ir"
	"github.com/rj45/llbrew/ir/op"
	"github.com/rj45/llbrew/ir/typ"
	"github.com/rj45/llbrew/xform"
)

var _ = xform.Register(loadOffset,
	xform.OnlyPass(xform.Simplification),
	xform.OnOp(op.Load),
	xform.Tags(xform.LoadStoreOffset),
)

func loadOffset(it ir.Iter) {
	instr := it.Instr()
	if instr.NumArgs() > 1 {
		return
	}

	add := offset(instr)
	if add == nil {
		instr.InsertArg(1, instr.Func().ValueFor(typ.IntegerWordType(), 0))
		return
	}

	// combine the add with the load
	instr.ReplaceArg(0, add.Arg(0))
	instr.InsertArg(-1, add.Arg(1))
	it.RemoveInstr(add)
}

var _ = xform.Register(storeOffset,
	xform.OnlyPass(xform.Simplification),
	xform.OnOp(op.Store),
	xform.Tags(xform.LoadStoreOffset),
)

func storeOffset(it ir.Iter) {
	instr := it.Instr()

	if instr.NumArgs() > 2 {
		return
	}

	add := offset(instr)
	if add == nil {
		instr.InsertArg(1, instr.Func().ValueFor(typ.IntegerWordType(), 0))
		return
	}

	// combine the add with the store
	instr.ReplaceArg(0, add.Arg(0))
	instr.InsertArg(1, add.Arg(1))
	it.RemoveInstr(add)
}

func offset(instr *ir.Instr) *ir.Instr {
	if instr.Arg(0).IsConst() {
		return nil
	}

	add := instr.Arg(0).Def().Instr()
	if add.Op == op.Add && (add.Arg(1).IsConst() || add.Arg(1).InSpillArea()) && add.Def(0).NumUses() == 1 {
		return add
	}
	return nil
}
