package translate

import "tinygo.org/x/go-llvm"

func (trans *translator) translateBlocks(fn llvm.Value) {
	first := true
	for blk := fn.FirstBasicBlock(); !blk.IsNil(); blk = llvm.NextBasicBlock(blk) {
		nblk := trans.fn.NewBlock()
		trans.blkmap[blk] = nblk

		if first {
			first = false
			for param := fn.FirstParam(); !param.IsNil(); param = llvm.NextParam(param) {
				pval := trans.fn.NewValue(translateType(param.Type()))
				nblk.AddDef(pval)
				trans.valuemap[param] = pval
			}
		}

		trans.fn.InsertBlock(-1, nblk)
	}

	for blk := fn.FirstBasicBlock(); !blk.IsNil(); blk = llvm.NextBasicBlock(blk) {
		// handle preds & succs
		nblk := trans.blkmap[blk]
		bval := blk.AsValue()
		for i := 0; i < bval.IncomingCount(); i++ {
			inc := bval.IncomingBlock(i)
			pred := trans.blkmap[inc]
			nblk.AddPred(pred)
			pred.AddSucc(nblk)
		}
	}
}
