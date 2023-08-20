package finishing

import (
	"github.com/rj45/llbrew/ir"
	"github.com/rj45/llbrew/ir/op"
	"github.com/rj45/llbrew/ir/reg"
	"github.com/rj45/llbrew/xform"
)

var _ = xform.Register(logue,
	xform.OnlyPass(xform.Finishing),
	xform.Once(),
)

func logue(it ir.Iter) {
	if !it.HasNext() {
		// todo: remove this when external assembly is properly implemented
		return // empty fn
	}
	fn := it.Block().Func()
	types := fn.Types()
	frame := &fn.Frame

	saveRegMap := map[reg.Reg]struct{}{}
	for sit := fn.InstrIter(); sit.HasNext(); sit.Next() {
		instr := sit.Instr()
		if reg.RA != reg.None && instr.Op.IsCall() {
			saveRegMap[reg.RA] = struct{}{}
		}

		for d := 0; d < instr.NumDefs(); d++ {
			reg := instr.Def(d).Reg()
			if reg.IsSavedReg() {
				saveRegMap[reg] = struct{}{}
			}
		}

		for a := 0; a < instr.NumArgs(); a++ {
			reg := instr.Arg(a).Reg()
			if reg.IsSavedReg() {
				saveRegMap[reg] = struct{}{}
			}
		}
	}

	index := 0
	for _, sreg := range reg.SavedRegs {
		if _, found := saveRegMap[sreg]; found {
			slot := frame.SlotID(ir.SavedSlot, index)
			it.Insert(op.Store, types.IntegerWordType(), reg.SP, slot, sreg)
			index++
		}
	}
	if _, found := saveRegMap[reg.RA]; found {
		slot := frame.SlotID(ir.SavedSlot, index)
		it.Insert(op.Store, types.PointerType(types.VoidType(), 0), reg.SP, slot, reg.RA)
	}

	var rets []*ir.Instr

	for it := fn.InstrIter(); it.HasNext(); it.Next() {
		instr := it.Instr()
		if instr.Op.IsReturn() {
			rets = append(rets, instr)
			index := 0
			for _, sreg := range reg.SavedRegs {
				if _, found := saveRegMap[sreg]; found {
					slot := frame.SlotID(ir.SavedSlot, index)
					load := it.Insert(op.Load, types.IntegerWordType(), reg.SP, slot)
					load.Def(0).SetReg(sreg)
					index++
				}
			}
			if _, found := saveRegMap[reg.RA]; found {
				slot := frame.SlotID(ir.SavedSlot, index)
				load := it.Insert(op.Load, types.PointerType(types.VoidType(), 0), reg.SP, slot)
				load.Def(0).SetReg(reg.RA)
			}
		}
	}

	frame.Scan()
	if frame.FrameSize() > 0 {
		it := fn.InstrIter()
		it.Insert(op.Store, types.PointerType(types.VoidType(), 0), reg.SP, 0, reg.SP)
		add := it.Insert(op.Add, types.PointerType(types.VoidType(), 0), reg.SP, -frame.FrameSize())
		add.Def(0).SetReg(reg.SP)

		for _, ret := range rets {
			add := fn.NewInstr(op.Add, types.PointerType(types.VoidType(), 0), reg.SP, frame.FrameSize())
			add.Def(0).SetReg(reg.SP)

			ret.Block().InsertInstr(ret.Index(), add)
		}
	}

	frame.ReplaceOffsets()
	it.Changed()
}
