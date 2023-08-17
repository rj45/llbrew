package elaboration

import (
	"github.com/rj45/llbrew/ir"
	"github.com/rj45/llbrew/ir/op"
	"github.com/rj45/llbrew/ir/reg"
	"github.com/rj45/llbrew/ir/typ"
	"github.com/rj45/llbrew/xform"
)

var _ = xform.Register(calls,
	xform.OnlyPass(xform.Elaboration),
	xform.OnOp(op.Call),
)

func calls(it ir.Iter) {
	instr := it.Instr()
	fnType := instr.Arg(0).Type.(*typ.Function)

	allSet := true
	for a := 0; a < instr.NumArgs(); a++ {
		if instr.Arg(a).Def() == nil || instr.Arg(a).Def().Instr().Op != op.Copy {
			allSet = false
			break
		}
		if a < len(reg.ArgRegs) && instr.Arg(a).Reg() != reg.ArgRegs[a] {
			allSet = false
			break
		} else if a >= len(reg.ArgRegs) && instr.Arg(a).ArgSlot() != a-len(reg.ArgRegs) {
			allSet = false
			break
		}
	}

	for d := 0; d < instr.NumDefs(); d++ {
		if instr.Def(d).NumUses() != 1 || instr.Def(d).Use(0).Instr().Op != op.Copy {
			allSet = false
			break
		}
		if d < len(reg.ArgRegs) && instr.Def(d).Reg() != reg.ArgRegs[d] {
			allSet = false
			break
		} else if d >= len(reg.ArgRegs) && instr.Def(d).ArgSlot() != d-len(reg.ArgRegs) {
			allSet = false
			break
		}
	}

	if allSet {
		return
	}

	if instr.NumDefs() > 0 && instr.Def(0).NumUses() == 1 && instr.Def(0).Use(0).Instr().Op == op.Copy {
		return
	}

	// todo:
	// - add parallel copy for clobbered regs?

	if instr.NumArgs() > 1 {
		params := fnType.Params

		args := make([]interface{}, instr.NumArgs()-1)
		for i := 1; i < instr.NumArgs(); i++ {
			args[i-1] = instr.Arg(i)
		}

		paramCopy := it.Insert(op.Copy, nil, args...)
		for i := 0; i < paramCopy.NumArgs(); i++ {
			paramCopy.AddDef(paramCopy.Func().NewValue(params[i]))
			if i < len(reg.ArgRegs) {
				paramCopy.Def(i).SetReg(reg.ArgRegs[i])
			} else {
				paramCopy.Def(i).SetArgSlot(i - len(reg.ArgRegs))
			}
			instr.ReplaceArg(i+1, paramCopy.Def(i))
		}
	}

	if instr.NumDefs() > 0 {
		results := fnType.Results

		args := make([]interface{}, instr.NumDefs())
		for i := 0; i < instr.NumDefs(); i++ {
			args[i] = instr.Def(i)
		}

		it.Next()
		resCopy := it.Insert(op.Copy, nil, args...)
		it.Prev()
		for i := 0; i < resCopy.NumArgs(); i++ {
			resCopy.AddDef(resCopy.Func().NewValue(results[i]))
			if i < len(reg.ArgRegs) {
				resCopy.Arg(i).SetReg(reg.ArgRegs[i])
			} else {
				resCopy.Arg(i).SetArgSlot(i - len(reg.ArgRegs))
			}

			// todo: could use a version of this that doesn't
			// clobber the current instruction or something
			instr.Def(i).ReplaceUsesWith(resCopy.Def(i))

			// switch this back to what it was
			resCopy.ReplaceArg(i, instr.Def(i))
		}
	}
}
