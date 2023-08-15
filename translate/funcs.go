package translate

import (
	"github.com/rj45/llbrew/ir"
	"tinygo.org/x/go-llvm"
)

func (trans *translator) translateFunctions() {
	for fn := trans.mod.FirstFunction(); !fn.IsNil(); fn = llvm.NextFunction(fn) {
		trans.translateFunction(fn)
	}
}

func (trans *translator) translateFunction(fn llvm.Value) {
	trans.fn = trans.pkg.NewFunc(fn.Name(), trans.translateFuncType(fn))

	trans.fn.Referenced = true

	trans.blkmap = make(map[llvm.BasicBlock]*ir.Block)
	trans.valuemap = make(map[llvm.Value]*ir.Value)
	trans.instrmap = make(map[llvm.Value]*ir.Instr)

	trans.translateBlocks(fn)
	trans.translateInstructions(fn)
	trans.translateAllOperands(fn)
}
