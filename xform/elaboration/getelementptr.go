package elaboration

import (
	"github.com/rj45/llbrew/ir"
	"github.com/rj45/llbrew/ir/op"
	"github.com/rj45/llbrew/ir/typ"
	"github.com/rj45/llbrew/xform"
)

var _ = xform.Register(getElementPtr,
	xform.OnlyPass(xform.Elaboration),
	xform.OnOp(op.GetElementPtr),
)

func getElementPtr(it ir.Iter) {
	instr := it.Instr()

	base := instr.Arg(0)
	index, ok := ir.IntValue(instr.Arg(1).Const())
	if !ok {
		// todo: handle if this is in a register
		panic("bad assumption, expecting index to be int const")
	}
	element, ok := ir.IntValue(instr.Arg(2).Const())
	if !ok {
		// todo: handle if this is in a register
		panic("bad assumption, expecting element to be int const")
	}

	pointedto := base.Type.(*typ.Pointer).Element

	offset := pointedto.SizeOf()*index +
		pointedto.(*typ.Struct).OffsetOf(element)

	elemtyp := pointedto.(*typ.Struct).Elements[element]
	elemtypptr := instr.Func().Types().PointerType(elemtyp, 0)

	if offset == 0 {
		instr.Def(0).ReplaceUsesWith(base)
		it.RemoveInstr(instr)
		return
	}

	it.Update(op.Add, elemtypptr, base, offset)
}
