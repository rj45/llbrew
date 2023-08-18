package finishing

import (
	"github.com/rj45/llbrew/ir"
	"github.com/rj45/llbrew/ir/op"
	"github.com/rj45/llbrew/ir/reg"
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
+------------------------+   > beginnings of next frame
|                        |  | <-- SP
+------------------------+  |
| potential stack arg 3  | /
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
| saved ra (if needed)   | / <-- caller SP
+------------------------+ \
| callee stack arg 3     |  |
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
	types := fn.Types()

	spillAreaSize := fn.SpillAreaSize()

	numParamSlots := max(len(fn.Sig.Params)-len(reg.ArgRegs), 0)
	paramSlots := make([]int, numParamSlots)
	paramOffset := 0
	for i := 0; i < numParamSlots; i++ {
		t := fn.Sig.Params[len(reg.ArgRegs)+i]
		paramOffset += t.SizeOf()
		paramSlots[i] = paramOffset
	}

	// TODO: count arg slots and replace
	argSize := fn.NumArgSlots() * sizes.WordSize()

	if argSize > 0 {
		argSize++
	}

	spillOffset := argSize // todo: will be > 0 when arg area is allocated

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
				actualAddress := arg.SpillAddress() + spillOffset
				arg.SetConst(ir.ConstFor(actualAddress))
			} else if arg.InArgSlot() {
				actualAddress := (arg.ArgSlot() + 1) * sizes.WordSize()
				arg.SetConst(ir.ConstFor(actualAddress))
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

	frameSize := spillAreaSize + argSize + (len(saveRegMap) * sizes.WordSize())

	offset := 0
	for sit := fn.InstrIter(); sit.HasNext(); sit.Next() {
		instr := sit.Instr()

		if instr.Op == op.Alloca {
			size, _ := ir.IntValue(instr.Arg(0).Const())

			offset += size * sizes.WordSize()
		}

		for a := 0; a < instr.NumArgs(); a++ {
			arg := instr.Arg(a)
			if arg.InParamSlot() {
				arg.SetConst(ir.ConstFor(paramSlots[arg.ParamSlot()] + frameSize))
			}
		}
	}

	// the add is put first because some arches only allow positive load/store offsets
	size := sizes.WordSize()
	var spval *ir.Value
	if frameSize > 0 {
		spval = it.Insert(op.Add, types.PointerType(types.VoidType(), 0), reg.SP, -frameSize).Def(0)
		spval.SetReg(reg.SP)
	}

	offset = spillAreaSize + argSize
	for _, sreg := range reg.SavedRegs {
		if _, found := saveRegMap[sreg]; found {
			it.Insert(op.Store, types.IntegerWordType(), reg.SP, offset, sreg)
			offset += size
		}
	}
	if _, found := saveRegMap[reg.RA]; found {
		it.Insert(op.Store, types.PointerType(types.VoidType(), 0), reg.SP, offset, reg.RA)
		offset += size
	}

	for ; it.HasNext(); it.Next() {
		instr := it.Instr()
		if instr.Op.IsReturn() {
			offset := spillAreaSize + argSize
			for _, sreg := range reg.SavedRegs {
				if _, found := saveRegMap[sreg]; found {
					load := it.Insert(op.Load, types.IntegerWordType(), reg.SP, offset)
					load.Def(0).SetReg(sreg)
					offset += size
				}
			}
			if _, found := saveRegMap[reg.RA]; found {
				load := it.Insert(op.Load, types.PointerType(types.VoidType(), 0), reg.SP, offset)
				load.Def(0).SetReg(reg.RA)
				offset += size
			}
			if frameSize > 0 {
				spval2 := it.Insert(op.Add, types.PointerType(types.VoidType(), 0), reg.SP, frameSize).Def(0)
				spval2.SetReg(reg.SP)
			}
		}
	}
}
