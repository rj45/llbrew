// Code generated from rewrite.rules using 'go generate'; DO NOT EDIT.

package elaboration

import (
	"github.com/rj45/llbrew/ir"
	"github.com/rj45/llbrew/ir/op"
	"github.com/rj45/llbrew/ir/reg"
)

func rewrite(it ir.Iter) {
	instr := it.Instr()
	blk := instr.Block()
	fn := blk.Func()
	types := fn.Types()
	_ = types
	_ = reg.SP
	switch instr.Op {
	case op.Store:
		for {
			var a, b, c *ir.Value
			if instr.NumArgs() != 3 {
				break
			}
			a = instr.Arg(0)
			b = instr.Arg(1)
			c = instr.Arg(2)
			if !(c.IsConst()) {
				break
			}
			i0 := it.Insert(op.Copy, c.Type, c)
			v0 := i0.Def(0)
			it.Update(op.Store, instr.Type(), a, b, v0)
			break
		}
		for {
			var a, c *ir.Value
			if instr.NumArgs() != 2 {
				break
			}
			a = instr.Arg(0)
			c = instr.Arg(1)
			if !(c.IsConst()) {
				break
			}
			i0 := it.Insert(op.Copy, c.Type, c)
			v0 := i0.Def(0)
			it.Update(op.Store, instr.Type(), a, v0)
			break
		}
	case op.If:
		for {
			var a *ir.Value
			if instr.NumArgs() != 1 {
				break
			}
			a = instr.Arg(0)
			if !(!a.Op().IsCompare()) {
				break
			}
			i0 := it.Insert(op.NotEqual, a.Type, a, 0)
			v0 := i0.Def(0)
			it.Update(op.If, instr.Type(), v0)
			break
		}
	case op.IntToPtr:
		for {
			var a *ir.Value
			if instr.NumArgs() != 1 {
				break
			}
			a = instr.Arg(0)
			if !(a.Type == types.IntegerWordType()) {
				break
			}
			it.ReplaceWith(a)
			break
		}
	}
}
