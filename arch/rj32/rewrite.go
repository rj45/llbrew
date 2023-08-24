// Code generated from rewrite.rules using 'go generate'; DO NOT EDIT.

package rj32

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
	case op.Ret:
		for {
			instr.Op = Return
			it.Changed()
			break
		}
	case op.Jump:
		for {
			instr.Op = Jump
			it.Changed()
			break
		}
	case op.Call:
		for {
			instr.Op = Call
			it.Changed()
			break
		}
	case op.Add:
		for {
			instr.Op = Add
			it.Changed()
			break
		}
	case op.Sub:
		for {
			instr.Op = Sub
			it.Changed()
			break
		}
	case op.And:
		for {
			instr.Op = And
			it.Changed()
			break
		}
	case op.Or:
		for {
			instr.Op = Or
			it.Changed()
			break
		}
	case op.Xor:
		for {
			instr.Op = Xor
			it.Changed()
			break
		}
	case op.Shl:
		for {
			instr.Op = Shl
			it.Changed()
			break
		}
	case op.LShr:
		for {
			instr.Op = Shr
			it.Changed()
			break
		}
	case op.AShr:
		for {
			instr.Op = Asr
			it.Changed()
			break
		}
	case op.SExt:
		for {
			instr.Op = Sxt
			it.Changed()
			break
		}
	case op.If:
		for {
			var a, b *ir.Value
			if instr.NumArgs() != 1 {
				break
			}
			instr_arg0 := instr.Arg(0).Def().Instr()
			if instr_arg0.Op != op.Equal {
				break
			}
			if instr_arg0.NumArgs() != 2 {
				break
			}
			a = instr_arg0.Arg(0)
			b = instr_arg0.Arg(1)
			it.Update(IfEq, instr.Type(), a, b)
			break
		}
		for {
			var a, b *ir.Value
			if instr.NumArgs() != 1 {
				break
			}
			instr_arg0 := instr.Arg(0).Def().Instr()
			if instr_arg0.Op != op.NotEqual {
				break
			}
			if instr_arg0.NumArgs() != 2 {
				break
			}
			a = instr_arg0.Arg(0)
			b = instr_arg0.Arg(1)
			it.Update(IfNe, instr.Type(), a, b)
			break
		}
		for {
			var a, b *ir.Value
			if instr.NumArgs() != 1 {
				break
			}
			instr_arg0 := instr.Arg(0).Def().Instr()
			if instr_arg0.Op != op.Less {
				break
			}
			if instr_arg0.NumArgs() != 2 {
				break
			}
			a = instr_arg0.Arg(0)
			b = instr_arg0.Arg(1)
			it.Update(IfLt, instr.Type(), a, b)
			break
		}
		for {
			var a, b *ir.Value
			if instr.NumArgs() != 1 {
				break
			}
			instr_arg0 := instr.Arg(0).Def().Instr()
			if instr_arg0.Op != op.LessEqual {
				break
			}
			if instr_arg0.NumArgs() != 2 {
				break
			}
			a = instr_arg0.Arg(0)
			b = instr_arg0.Arg(1)
			it.Update(IfLe, instr.Type(), a, b)
			break
		}
		for {
			var a, b *ir.Value
			if instr.NumArgs() != 1 {
				break
			}
			instr_arg0 := instr.Arg(0).Def().Instr()
			if instr_arg0.Op != op.Greater {
				break
			}
			if instr_arg0.NumArgs() != 2 {
				break
			}
			a = instr_arg0.Arg(0)
			b = instr_arg0.Arg(1)
			it.Update(IfGt, instr.Type(), a, b)
			break
		}
		for {
			var a, b *ir.Value
			if instr.NumArgs() != 1 {
				break
			}
			instr_arg0 := instr.Arg(0).Def().Instr()
			if instr_arg0.Op != op.GreaterEqual {
				break
			}
			if instr_arg0.NumArgs() != 2 {
				break
			}
			a = instr_arg0.Arg(0)
			b = instr_arg0.Arg(1)
			it.Update(IfGe, instr.Type(), a, b)
			break
		}
		for {
			var a, b *ir.Value
			if instr.NumArgs() != 1 {
				break
			}
			instr_arg0 := instr.Arg(0).Def().Instr()
			if instr_arg0.Op != op.ULess {
				break
			}
			if instr_arg0.NumArgs() != 2 {
				break
			}
			a = instr_arg0.Arg(0)
			b = instr_arg0.Arg(1)
			it.Update(IfUlt, instr.Type(), a, b)
			break
		}
		for {
			var a, b *ir.Value
			if instr.NumArgs() != 1 {
				break
			}
			instr_arg0 := instr.Arg(0).Def().Instr()
			if instr_arg0.Op != op.ULessEqual {
				break
			}
			if instr_arg0.NumArgs() != 2 {
				break
			}
			a = instr_arg0.Arg(0)
			b = instr_arg0.Arg(1)
			it.Update(IfUle, instr.Type(), a, b)
			break
		}
		for {
			var a, b *ir.Value
			if instr.NumArgs() != 1 {
				break
			}
			instr_arg0 := instr.Arg(0).Def().Instr()
			if instr_arg0.Op != op.UGreater {
				break
			}
			if instr_arg0.NumArgs() != 2 {
				break
			}
			a = instr_arg0.Arg(0)
			b = instr_arg0.Arg(1)
			it.Update(IfUgt, instr.Type(), a, b)
			break
		}
		for {
			var a, b *ir.Value
			if instr.NumArgs() != 1 {
				break
			}
			instr_arg0 := instr.Arg(0).Def().Instr()
			if instr_arg0.Op != op.UGreaterEqual {
				break
			}
			if instr_arg0.NumArgs() != 2 {
				break
			}
			a = instr_arg0.Arg(0)
			b = instr_arg0.Arg(1)
			it.Update(IfUge, instr.Type(), a, b)
			break
		}
	case op.Load:
		for {
			var a *ir.Value
			if instr.NumArgs() != 1 {
				break
			}
			a = instr.Arg(0)
			if !(!a.IsConst()) {
				break
			}
			it.Update(Load, instr.Type(), a, 0)
			break
		}
		for {
			var a *ir.Value
			if instr.NumArgs() != 1 {
				break
			}
			a = instr.Arg(0)
			if !(a.IsConst()) {
				break
			}
			it.Update(Load, instr.Type(), reg.GP, a)
			break
		}
		for {
			var a, b *ir.Value
			if instr.NumArgs() != 2 {
				break
			}
			a = instr.Arg(0)
			b = instr.Arg(1)
			if !(b.IsConst()) {
				break
			}
			it.Update(Load, instr.Type(), a, b)
			break
		}
		for {
			var a, b *ir.Value
			if instr.NumArgs() != 2 {
				break
			}
			a = instr.Arg(0)
			b = instr.Arg(1)
			if !(!b.IsConst()) {
				break
			}
			i0 := it.Insert(Add, a.Type, a, b)
			v0 := i0.Def(0)
			it.Update(Load, instr.Type(), v0, 0)
			break
		}
	case Load:
		for {
			var a, b *ir.Value
			if instr.NumArgs() != 2 {
				break
			}
			instr_arg0 := instr.Arg(0).Def().Instr()
			if instr_arg0.Op != Add {
				break
			}
			if instr_arg0.NumArgs() != 2 {
				break
			}
			instr_arg0_arg0 := instr_arg0.Arg(0).Def().Instr()
			if instr_arg0_arg0.Op != Move {
				break
			}
			if instr_arg0_arg0.NumArgs() != 1 {
				break
			}
			a = instr_arg0_arg0.Arg(0)
			b = instr_arg0.Arg(1)
			if !instr.Arg(1).HasConstValue(0) {
				break
			}
			if !(b.IsConst()) {
				break
			}
			it.Update(Load, instr.Type(), a, b)
			break
		}
		for {
			var a, b *ir.Value
			if instr.NumArgs() != 2 {
				break
			}
			instr_arg0 := instr.Arg(0).Def().Instr()
			if instr_arg0.Op != Add {
				break
			}
			if instr_arg0.NumArgs() != 2 {
				break
			}
			a = instr_arg0.Arg(0)
			b = instr_arg0.Arg(1)
			if !instr.Arg(1).HasConstValue(0) {
				break
			}
			if !(b.IsConst()) {
				break
			}
			it.Update(Load, instr.Type(), a, b)
			break
		}
	case op.Store:
		for {
			var a, b *ir.Value
			if instr.NumArgs() != 2 {
				break
			}
			a = instr.Arg(0)
			b = instr.Arg(1)
			if !(!a.IsConst()) {
				break
			}
			it.Update(Store, instr.Type(), a, 0, b)
			break
		}
		for {
			var a, b *ir.Value
			if instr.NumArgs() != 2 {
				break
			}
			a = instr.Arg(0)
			b = instr.Arg(1)
			if !(a.IsConst()) {
				break
			}
			it.Update(Store, instr.Type(), reg.GP, a, b)
			break
		}
		for {
			var a, b, c *ir.Value
			if instr.NumArgs() != 3 {
				break
			}
			a = instr.Arg(0)
			b = instr.Arg(1)
			c = instr.Arg(2)
			if !(b.IsConst()) {
				break
			}
			it.Update(Store, instr.Type(), a, b, c)
			break
		}
		for {
			var a, b, c *ir.Value
			if instr.NumArgs() != 3 {
				break
			}
			a = instr.Arg(0)
			b = instr.Arg(1)
			c = instr.Arg(2)
			if !(!b.IsConst()) {
				break
			}
			i0 := it.Insert(Add, a.Type, a, b)
			v0 := i0.Def(0)
			it.Update(Store, instr.Type(), v0, 0, c)
			break
		}
	case Store:
		for {
			var a, b, c *ir.Value
			if instr.NumArgs() != 3 {
				break
			}
			instr_arg0 := instr.Arg(0).Def().Instr()
			if instr_arg0.Op != Add {
				break
			}
			if instr_arg0.NumArgs() != 2 {
				break
			}
			instr_arg0_arg0 := instr_arg0.Arg(0).Def().Instr()
			if instr_arg0_arg0.Op != Move {
				break
			}
			if instr_arg0_arg0.NumArgs() != 1 {
				break
			}
			a = instr_arg0_arg0.Arg(0)
			b = instr_arg0.Arg(1)
			if !instr.Arg(1).HasConstValue(0) {
				break
			}
			c = instr.Arg(2)
			if !(b.IsConst()) {
				break
			}
			it.Update(Store, instr.Type(), a, b, c)
			break
		}
		for {
			var a, b, c *ir.Value
			if instr.NumArgs() != 3 {
				break
			}
			instr_arg0 := instr.Arg(0).Def().Instr()
			if instr_arg0.Op != Add {
				break
			}
			if instr_arg0.NumArgs() != 2 {
				break
			}
			a = instr_arg0.Arg(0)
			b = instr_arg0.Arg(1)
			if !instr.Arg(1).HasConstValue(0) {
				break
			}
			c = instr.Arg(2)
			if !(b.IsConst()) {
				break
			}
			it.Update(Store, instr.Type(), a, b, c)
			break
		}
	}
}
