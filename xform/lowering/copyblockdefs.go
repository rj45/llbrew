package lowering

import (
	"github.com/rj45/llbrew/ir"
	"github.com/rj45/llbrew/ir/op"
	"github.com/rj45/llbrew/ir/reg"
	"github.com/rj45/llbrew/xform"
)

var _ = xform.Register(copyBlockDefs,
	xform.OnlyPass(xform.Lowering),
)

// copyBlockDefs inserts a copy at the beginning of each block with
// parameters. This is to aide the register allocator in lowering out
// of SSA. This will produce a parallel copy which can later be lowered
// into sequential copies. This is important so that there is no
// artificial constraints imposed on which registers can be picked due
// to the order of the sequential copies.
func copyBlockDefs(it ir.Iter) {
	// if not at the first instruction of the block, skip
	if it.InstrIndex() != 0 {
		return
	}

	blk := it.Block()

	if blk.NumDefs() == 0 {
		return
	}

	numDefs := blk.NumDefs()
	if blk.Func().Block(0) == blk && blk.NumDefs() > len(reg.ArgRegs) {
		numDefs = len(reg.ArgRegs)
	}

	allCopied := true
	for d := 0; d < numDefs; d++ {
		def := blk.Def(d)

		// check if the sole use of the def is not a copy in this block
		if def.NumUses() != 1 ||
			def.Use(0).Instr() == nil ||
			def.Use(0).Instr().Op != op.Copy ||
			def.Use(0).Instr().Block() != blk {
			allCopied = false
			break
		}
	}

	if allCopied {
		// already done
		return
	}

	instr := it.Insert(op.Copy, nil)

	for d := 0; d < numDefs; d++ {
		def := blk.Def(d)

		ndef := blk.Func().NewValue(def.Type)
		instr.AddDef(ndef)

		def.ReplaceUsesWith(ndef)
		instr.InsertArg(-1, def)
	}
}
