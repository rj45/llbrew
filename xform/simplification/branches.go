package simplification

import (
	"github.com/rj45/llbrew/ir"
	"github.com/rj45/llbrew/ir/op"
	"github.com/rj45/llbrew/xform"
)

var _ = xform.Register(swapIfBranches,
	xform.OnlyPass(xform.Simplification),
	xform.OnOp(op.If),
)

func swapIfBranches(it ir.Iter) {
	instr := it.Instr()

	compare := instr.Arg(0).Def().ID.InstrIn(instr.Func())

	// if the false branch of the `if` is not the very next block, but the true branch is
	if instr.Block().Succ(1).Index() != instr.Block().Index()+1 {
		if instr.Block().Succ(0).Index() == instr.Block().Index()+1 {
			if opper, ok := compare.Op.(interface{ Opposite() op.Op }); ok {
				compare.Update(opper.Opposite(), nil, compare.Args())
				it.Changed()
			} else {
				// not := it.Insert(op.Not, compare.Def(0).Type, compare.Def(0))
				// instr.ReplaceArg(0, not.Def(0))
				panic("todo: figure out how to do logical not")
			}

			instr.Block().SwapSuccs()
		} else {
			panic("not able to legalize branch")
		}
	}
}
