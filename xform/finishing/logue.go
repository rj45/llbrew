package finishing

import (
	"github.com/rj45/llir2asm/ir"
	"github.com/rj45/llir2asm/ir/op"
	"github.com/rj45/llir2asm/ir/reg"
	"github.com/rj45/llir2asm/ir/typ"
	"github.com/rj45/llir2asm/sizes"
	"github.com/rj45/llir2asm/xform"
)

var _ = xform.Register(logue,
	xform.OnlyPass(xform.Finishing),
	xform.Once(),
)

/*
SP points at next empty place on the stack

Here is what the call frames look like:

low mem addresses
+------------------------+  |
|                        |   > beginnings of next frame
+------------------------+  |
| potential stack arg 3  | /  <-- SP
+------------------------+ \
| potential stack arg 4  |  |
+------------------------+  |
| local spill 0          |  |
+------------------------+  |
| local spill 1          |  |
+------------------------+  |
| saved reg 1            |   > current callee's frame
+------------------------+  |
| saved reg 2            |  |
+------------------------+  |
| saved ra (if needed)   | /
+------------------------+ \
| callee stack arg 3     |  | <-- caller SP
+------------------------+  |
| callee stack arg 4     |  |
+------------------------+  |
| spill local 0          |  |
+------------------------+   > caller's stack frame
| spill local 1          |  |
+------------------------+  |
| saved reg 1            |  |
+------------------------+  |
| saved reg 2            |  |
+------------------------+  |
| saved RA               | /
+------------------------+
high mem addresses

Note that function parameters (arguments) are on the caller's frame,
in the area known as "ArgSlots".

So the order on the from the SP is:
	- ArgSlots for calls
	- SpillSlots for local variables on the stack
	- Saved registers
	- Params for the function incoming parameters

*/

func logue(it ir.Iter) {
	if !it.HasNext() {
		// todo: remove this when external assembly is properly implemented
		return // empty fn
	}
	fn := it.Block().Func()

	saveRegMap := map[reg.Reg]struct{}{}
	for sit := fn.InstrIter(); sit.HasNext(); sit.Next() {
		instr := sit.Instr()
		if reg.RA != reg.None && instr.Op.IsCall() {
			// todo: handle ArgSlots
			saveRegMap[reg.RA] = struct{}{}
		}

		// todo: add FP to saveRegs if it's used and the arch
		// requires it to be saved

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

	if len(saveRegMap) == 0 {
		// nothing to do
		return
	}

	// todo: account for SpillSlots when spilling is implemented
	// todo: account for called function ArgSlots

	// the add is put first because some arches only allow positive load/store offsets
	size := int(sizes.WordSize())
	spval := it.Insert(op.Add, typ.PointerType(typ.VoidType(), 0), reg.SP, -len(saveRegMap)*size).Def(0)
	spval.SetReg(reg.SP)

	offset := 0
	for _, sreg := range reg.SavedRegs {
		if _, found := saveRegMap[sreg]; found {
			it.Insert(op.Store, typ.IntegerWordType(), spval, offset, sreg)
			offset += size
		}
	}
	if _, found := saveRegMap[reg.RA]; found {
		it.Insert(op.Store, typ.PointerType(typ.VoidType(), 0), spval, offset, reg.RA)
		offset += size
	}

	for ; it.HasNext(); it.Next() {
		instr := it.Instr()
		if instr.Op.IsReturn() {
			offset := 0
			for _, sreg := range reg.SavedRegs {
				if _, found := saveRegMap[sreg]; found {
					load := it.Insert(op.Load, typ.IntegerWordType(), spval, offset)
					load.Def(0).SetReg(sreg)
					offset += size
				}
			}
			if _, found := saveRegMap[reg.RA]; found {
				load := it.Insert(op.Load, typ.PointerType(typ.VoidType(), 0), spval, offset)
				load.Def(0).SetReg(reg.RA)
				offset += size
			}
			spval2 := it.Insert(op.Add, typ.PointerType(typ.VoidType(), 0), reg.SP, len(saveRegMap)*size).Def(0)
			spval2.SetReg(reg.SP)
		}
	}
}
