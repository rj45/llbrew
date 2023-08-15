package xform

import (
	"log"
	"reflect"
	"runtime"

	"github.com/rj45/llbrew/ir"
)

type Pass int

const (
	Elaboration Pass = iota
	Simplification
	Lowering
	Legalization
	CleanUp
	Finishing

	NumPasses
)

type desc struct {
	name     string
	passes   []Pass
	tags     []Tag
	op       ir.Op
	once     bool
	disabled bool
	fn       func(ir.Iter)
}

type Option func(d *desc)

func OnlyPass(p Pass) Option {
	return func(d *desc) {
		d.passes = []Pass{p}
	}
}

func Passes(p ...Pass) Option {
	return func(d *desc) {
		d.passes = p
	}
}

func Tags(t ...Tag) Option {
	return func(d *desc) {
		d.tags = t
	}
}

func OnOp(op ir.Op) Option {
	return func(d *desc) {
		d.op = op
	}
}

func Once() Option {
	return func(d *desc) {
		d.once = true
	}
}

var xformers []desc

// Register an xform function
func Register(fn func(ir.Iter), options ...Option) int {
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	xformers = append(xformers, desc{
		name: name,
		fn:   fn,
	})
	d := &xformers[len(xformers)-1]

	for _, option := range options {
		option(d)
	}

	return 0
}

func Transform(pass Pass, fn *ir.Func) {
	active, opXforms, anyOnceXforms, otherXforms := activeXforms(pass, fn)
	tries := 0

	// do the transforms operating on any op and only once first
	for _, xform := range anyOnceXforms {
		it := fn.InstrIter()
		var iter ir.Iter = it
		perform(xform, iter)
	}

	for {
		it := fn.InstrIter()
		var iter ir.Iter = it

		for ; it.HasNext(); it.Next() {
			// run the xforms specific to the current op
			op := it.Instr().Op
			for _, xform := range opXforms[op] {
				perform(xform, iter)

				if iter.Instr() == nil {
					log.Panicf("xform %s in pass %v left iter in nil state", xform.name, pass)
				}

				if iter.Instr().Op != op {
					break
				}
			}

			// run the xforms that always run
			for _, xform := range otherXforms {
				perform(xform, iter)
				if iter.Instr() == nil {
					log.Panicf("xform %s in pass %v left iter in nil state", xform.name, pass)
				}
			}
		}

		if !it.HasChanged() {
			break
		}

		tries++
		if tries > 1000 {
			log.Panicf("transforms do not terminate: pass: %d active: %v", pass, active)
		}
	}
}

func perform(xform *desc, it ir.Iter) {
	if xform.disabled {
		return
	}
	xform.fn(it)
	if xform.once {
		xform.disabled = true
	}
}

// activeXforms determines the active xform functions for the current pass and tags
func activeXforms(pass Pass, fn *ir.Func) ([]string, map[ir.Op][]*desc, []*desc, []*desc) {
	var active []string
	opXforms := make(map[ir.Op][]*desc)
	var anyOnceXforms []*desc
	var otherXforms []*desc

next:
	for i, xf := range xformers {
		inPass := false
		for _, p := range xf.passes {
			if p == pass {
				inPass = true
				break
			}
		}
		if !inPass {
			continue
		}

		for _, tag := range xf.tags {
			if !activeTags[tag] {
				continue next
			}
		}

		// make a copy -- avoids global mutable state
		xform := xformers[i]

		if xf.op != nil {
			opXforms[xf.op] = append(opXforms[xf.op], &xform)
		} else if xf.once {
			anyOnceXforms = append(anyOnceXforms, &xform)
		} else {
			otherXforms = append(otherXforms, &xform)
		}
		active = append(active, xf.name)
	}

	return active, opXforms, anyOnceXforms, otherXforms
}
