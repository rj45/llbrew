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
	types := instr.Func().Types()

	addr := instr.Arg(0)
	ptrtyp := addr.Type
	elemtyp := addr.Type
	offset := 0

	for a := 1; a < instr.NumArgs(); a++ {
		index := instr.Arg(a)
		iindex := 0
		if index.IsConst() {
			iindex, _ = ir.IntValue(index.Const())
		}
		et, elemoffset := elementOffsetType(elemtyp, iindex)

		elemtyp = et
		ptrtyp = types.PointerType(et, 0)

		if index.IsConst() {
			offset += elemoffset
		} else {
			if et.SizeOf() > 1 {
				mul := it.Insert(op.Mul, ptrtyp, index, elemtyp.SizeOf())
				index = mul.Def(0)
			}
			if offset > 0 {
				// clear accumulated offset
				add := it.Insert(op.Add, ptrtyp, addr, offset)
				addr = add.Def(0)
				offset = 0
			}
			add := it.Insert(op.Add, ptrtyp, addr, index)
			addr = add.Def(0)
		}
	}

	if offset == 0 {
		instr.Def(0).ReplaceUsesWith(addr)
		it.RemoveInstr(instr)
		return
	}

	it.Update(op.Add, ptrtyp, addr, offset)
}

func elementOffsetType(t typ.Type, index int) (et typ.Type, offset int) {
	switch t := t.(type) {
	case *typ.Pointer:
		et = t.Element
		offset = et.SizeOf() * index
	case *typ.Array:
		et = t.Element
		offset = et.SizeOf() * index
	case *typ.Struct:
		et = t.Elements[index]
		offset = t.OffsetOf(index)
	default:
		offset = 0
		et = t
	}
	return
}
