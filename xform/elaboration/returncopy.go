package elaboration

import (
	"github.com/rj45/llir2asm/ir"
	"github.com/rj45/llir2asm/ir/op"
	"github.com/rj45/llir2asm/ir/reg"
	"github.com/rj45/llir2asm/xform"
)

var _ = xform.Register(returnCopy,
	xform.OnlyPass(xform.Elaboration),
	xform.OnOp(op.Ret),
	xform.Once(),
)

func returnCopy(it ir.Iter) {
	ret := it.Instr()
	if ret.NumArgs() < 1 {
		return
	}

	cp := it.Insert(op.Copy, 0, ret.Args())

	for i := 0; i < ret.NumArgs(); i++ {
		cp.AddDef(cp.Func().NewValue(ret.Arg(i).Type))
		if i < len(reg.ArgRegs) {
			cp.Def(i).SetReg(reg.ArgRegs[i])
		} else {
			cp.Def(i).SetArgSlot(i - len(reg.ArgRegs))
		}
		ret.ReplaceArg(i, cp.Def(i))
	}
}
