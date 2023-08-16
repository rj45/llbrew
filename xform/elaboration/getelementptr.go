package elaboration

import (
	"github.com/rj45/llbrew/ir"
	"github.com/rj45/llbrew/ir/op"
	"github.com/rj45/llbrew/ir/typ"
	"github.com/rj45/llbrew/xform"
)

// var _ = xform.Register(constGetElementPtr,
// 	xform.OnlyPass(xform.Elaboration),
// 	xform.OnOp(op.GetElementPtr),
// )

func constGetElementPtr(it ir.Iter) {
	instr := it.Instr()

	base := instr.Arg(0)
	index, ok := ir.IntValue(instr.Arg(1).Const())
	if !ok {
		// handled by getElementPtr
		return
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

var _ = xform.Register(getElementPtr,
	xform.OnlyPass(xform.Elaboration),
	xform.OnOp(op.GetElementPtr),
)

func getElementPtr(it ir.Iter) {
	instr := it.Instr()

	base := instr.Arg(0)
	index := instr.Arg(1)

	pointedto := base.Type.(*typ.Pointer).Element

	offset := 0

	if index.IsConst() {
		index, ok := ir.IntValue(instr.Arg(1).Const())
		if !ok {
			// handled by getElementPtr
			return
		}
		offset = pointedto.SizeOf() * index
	} else {
		if pointedto.SizeOf() > 1 {
			mul := it.Insert(op.Mul, base.Type, index, pointedto.SizeOf())
			index = mul.Def(0)
		}
		add := it.Insert(op.Add, base.Type, base, index)
		base = add.Def(0)
	}

	elemtyp := pointedto
	if st, ok := pointedto.(*typ.Struct); ok {
		element, ok := ir.IntValue(instr.Arg(2).Const())
		if !ok {
			// todo: handle if this is in a register
			panic("bad assumption, expecting element to be int const")
		}
		offset += st.OffsetOf(element)
		elemtyp = st.Elements[element]
	}

	elemtypptr := instr.Func().Types().PointerType(elemtyp, 0)

	if offset == 0 {
		instr.Def(0).ReplaceUsesWith(base)
		it.RemoveInstr(instr)
		return
	}

	it.Update(op.Add, elemtypptr, base, offset)
}
