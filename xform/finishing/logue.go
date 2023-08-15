package finishing

import (
	"github.com/rj45/llbrew/ir"
	"github.com/rj45/llbrew/ir/op"
	"github.com/rj45/llbrew/ir/reg"
	"github.com/rj45/llbrew/ir/typ"
	"github.com/rj45/llbrew/sizes"
	"github.com/rj45/llbrew/xform"
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

	spillAreaSize := fn.SpillAreaSize()
	spillOffset := 0 // todo: will be > 0 when param area is allocated

	saveRegMap := map[reg.Reg]struct{}{}
	for sit := fn.InstrIter(); sit.HasNext(); sit.Next() {
		instr := sit.Instr()
		if reg.RA != reg.None && instr.Op.IsCall() {
			// todo: handle ArgSlots
			saveRegMap[reg.RA] = struct{}{}
		}

		// Replace all spill slot addresses with the actual address
		for a := 0; a < instr.NumArgs(); a++ {
			arg := instr.Arg(a)
			if arg.InSpillArea() {
				arg.ReplaceUsesWith(fn.ValueFor(typ.IntegerWordType(), arg.SpillAddress()+spillOffset))
			}
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

	frameSize := spillAreaSize + (len(saveRegMap) * sizes.WordSize())

	if frameSize == 0 {
		// nothing to do
		return
	}

	offset := 0
	for sit := fn.InstrIter(); sit.HasNext(); sit.Next() {
		instr := sit.Instr()

		if instr.Op == op.Alloca {
			size, _ := ir.IntValue(instr.Arg(0).Const())

			for u := 0; u < instr.Def(0).NumUses(); u++ {
				// todo: replace alloca instructions with values stored in stack slots,
				// then ensure load/stores have SP relative addressing to the stack slot,
				// then here, just replace values with constant offsets
			}

			offset += size * sizes.WordSize()
		}
	}

	// todo: account for SpillSlots when spilling is implemented
	// todo: account for called function ArgSlots

	// the add is put first because some arches only allow positive load/store offsets
	size := sizes.WordSize()
	spval := it.Insert(op.Add, typ.PointerType(typ.VoidType(), 0), reg.SP, -frameSize).Def(0)
	spval.SetReg(reg.SP)

	offset = spillAreaSize
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
			offset := spillAreaSize
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
			spval2 := it.Insert(op.Add, typ.PointerType(typ.VoidType(), 0), reg.SP, frameSize).Def(0)
			spval2.SetReg(reg.SP)
		}
	}
}
