package elaboration

import (
	"github.com/rj45/llbrew/ir"
	"github.com/rj45/llbrew/ir/op"
	"github.com/rj45/llbrew/ir/reg"
	"github.com/rj45/llbrew/xform"
)

var _ = xform.Register(paramSlots,
	xform.OnlyPass(xform.Elaboration),
	xform.Once(),
)

func paramSlots(it ir.Iter) {
	blk := it.Block()
	types := blk.Func().Types()

	for i := len(reg.ArgRegs); i < blk.NumDefs(); i++ {
		param := blk.Def(i)
		add := it.Insert(op.Add, types.PointerType(param.Type, 0), reg.SP)
		load := it.Insert(op.Load, param.Type, add.Def(0))
		param.ReplaceUsesWith(load.Def(0))
		add.InsertArg(-1, param)
	}
}
