package rj32

import (
	"github.com/rj45/llbrew/ir"
	"github.com/rj45/llbrew/ir/op"
	"github.com/rj45/llbrew/xform"
)

func (cpuArch) XformTags2() []xform.Tag {
	return nil
}

//go:generate go run github.com/rj45/llbrew/cmd/rewritegen -o rewrite.go -pkg rj32 -fn rewrite rewrite.rules

func (cpuArch) RegisterXforms() {
	xform.Register(rewrite,
		xform.OnlyPass(xform.Lowering))
	xform.Register(translateCopies,
		xform.OnlyPass(xform.Finishing),
		xform.OnOp(op.Copy))
	xform.Register(rewrite,
		xform.OnlyPass(xform.Finishing))
}

func translateCopies(it ir.Iter) {
	instr := it.Instr()

	if instr.NumArgs() == 1 {
		it.Update(Move, instr.Def(0).Type, instr.Args())
	} else if instr.NumArgs() == 2 {
		instr.Op = Swap
		it.Changed()
	} else {
		panic("parallel copy left!")
	}
}
