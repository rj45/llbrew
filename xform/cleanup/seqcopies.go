package cleanup

import (
	"log"
	"slices"

	"github.com/rj45/llbrew/ir"
	"github.com/rj45/llbrew/ir/op"
	"github.com/rj45/llbrew/ir/reg"
	"github.com/rj45/llbrew/xform"
)

var _ = xform.Register(sequentializeCopies,
	xform.OnlyPass(xform.CleanUp),
	xform.OnOp(op.Copy),
)

// sequentializeCopies takes a copy with more than one arg and figures out
// how to do the same thing in multiple copy instructions
// Based on Algorithm 13 from Benoit Boisinot's thesis, with fixes and
// extensions by Paul Sokolovsky.
// https://github.com/pfalcon/parcopy/blob/master/parcopy1.py
func sequentializeCopies(it ir.Iter) {
	instr := it.Instr()

	if instr.NumArgs() < 2 {
		return
	}

	if !instr.Op.IsCopy() {
		log.Panicf("called with non-copy! %s", instr.LongString())
	}

	var ready []reg.Reg
	var todo []reg.Reg
	pred := make(map[reg.Reg]reg.Reg)
	loc := make(map[reg.Reg]reg.Reg)

	srcs := make(map[reg.Reg]*ir.Value)
	dests := make(map[reg.Reg]*ir.Value)

	var copied [][2]*ir.Value

	emit := func(def, arg *ir.Value) {
		if def == nil {
			panic("nil dest")
		}
		if arg == nil {
			panic("nil src")
		}

		cp := it.Insert(op.Copy, def.Type, arg)
		cpdef := cp.Def(0)
		cpdef.SetReg(def.Reg())
		def.ReplaceUsesWith(cpdef)
		copied = append(copied, [2]*ir.Value{def, arg})
		it.Changed()
	}

	findFreeReg := func() reg.Reg {
		fn := instr.Func()
		unused := slices.Clone(reg.SavedRegs)
		hasCall := false
		for it := fn.InstrIter(); it.HasNext(); it.Next() {
			instr := it.Instr()
			for a := 0; a < instr.NumArgs(); a++ {
				arg := instr.Arg(a)
				if arg.InReg() {
					slices.DeleteFunc(unused, func(e reg.Reg) bool {
						return e == arg.Reg()
					})
				}
			}
			if instr.IsCall() {
				hasCall = true
			}
		}
		if len(unused) < 1 {
			if hasCall {
				// RA should be saved to the stack and available for use
				return reg.RA
			}

			panic("todo: no temp register available")
		}
		return unused[0]
	}

	for i := 0; i < instr.NumDefs(); i++ {
		def := instr.Def(i)
		arg := instr.Arg(i)

		b := def.Reg()
		a := arg.Reg()

		if b == a {
			// wait for copy elimination first
			return
		}

		if arg.IsConst() {
			emit(def, arg)
			continue
		}

		srcs[a] = arg
		dests[b] = def

		loc[a] = a
		pred[b] = a

		for _, todob := range todo {
			if todob == b {
				panic("double destination assignment")
			}
		}

		todo = append(todo, b)
	}

	for i := 0; i < instr.NumDefs(); i++ {
		def := instr.Def(i)
		if instr.Arg(i).IsConst() {
			continue
		}

		b := def.Reg()

		if _, found := loc[b]; !found {
			ready = append(ready, b)
		}
	}

	for len(todo) > 0 {
		for len(ready) > 0 {
			b := ready[len(ready)-1]
			ready = ready[:len(ready)-1]

			a, found := pred[b]
			if !found {
				continue
			}
			c := loc[a]

			// fmt.Println("copy", b, "<-", c)
			emit(dests[b], srcs[c])

			for i, td := range todo {
				if td == c {
					// remove c from todo
					todo[i] = todo[len(todo)-1]
					todo = todo[:len(todo)-1]
				}
			}

			loc[a] = b

			if a == c {
				ready = append(ready, a)
			}
		}

		if len(todo) == 0 {
			break
		}

		if len(todo) == 2 {
			// can be emitted with a swap
			break
		}

		b := todo[len(todo)-1]
		todo = todo[:len(todo)-1]

		if b != loc[pred[b]] {
			// need test program to verify this
			free := findFreeReg()
			freedef := instr.Func().ValueFor(dests[b].Type, free)
			loc[b] = free
			srcs[b] = freedef
			emit(freedef, srcs[b])
			ready = append(ready, b)

			panic("todo: temp needed, or figure out swap chain")
		}
	}

	for _, pair := range copied {
		def, arg := pair[0], pair[1]
		instr.RemoveArg(arg)
		instr.RemoveDef(def)
	}
	if instr.NumDefs() == 0 {
		it.Remove()
	}
}
