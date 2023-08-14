package rj32

import (
	"github.com/rj45/llbrew/ir/op"
	"github.com/rj45/llbrew/xform"
)

func (cpuArch) XformTags2() []xform.Tag {
	return []xform.Tag{xform.LoadStoreOffset}
}

func (cpuArch) RegisterXforms() {
	xform.Register(translate,
		xform.OnlyPass(xform.Lowering))
	xform.Register(translateCopies,
		xform.OnlyPass(xform.Finishing),
		xform.OnOp(op.Copy))
	xform.Register(translateLoadStore,
		xform.OnlyPass(xform.Finishing))
}
