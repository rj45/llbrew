package asm

import (
	"fmt"
	"io"
	"strings"

	"github.com/rj45/llir2asm/ir"
	"github.com/rj45/llir2asm/sizes"
)

type Section string

const (
	Code Section = "code"
	Data Section = "data"
	Bss  Section = "bss"
)

type Formatter interface {
	Section(s Section) string
	GlobalLabel(global *ir.Global) string
	PCRelAddress(offsetWords int) string
	Word(val string) string
	String(val string) string
	Reserve(bytes int) string
	Comment(comment string) string
	BlockLabel(id string) string
}

type Emitter struct {
	out   io.Writer
	fmter Formatter

	emittedGlobals map[*ir.Global]bool
	emittedFuncs   map[*ir.Func]bool

	section Section
	indent  string
}

func NewEmitter(out io.Writer, fmter Formatter) *Emitter {
	return &Emitter{
		out:            out,
		fmter:          fmter,
		emittedGlobals: make(map[*ir.Global]bool),
		emittedFuncs:   make(map[*ir.Func]bool),
	}
}

func Emit(out io.Writer, fmter Formatter, prog *ir.Program) {
	emitter := NewEmitter(out, fmter)
	emitter.Program(prog)
}

func (emit *Emitter) Program(prog *ir.Program) {
	mainpkg := prog.Package("main")
	emit.assemble(mainpkg.Func("main"))
}

func (emit *Emitter) assemble(fn *ir.Func) {
	var funcs []*ir.Func
	var globals []*ir.Global

	seenFunc := map[*ir.Func]bool{fn: true}

	todo := []*ir.Func{fn}
	for len(todo) > 0 {
		fn := todo[0]
		todo = todo[1:]

		funcs, globals = emit.scan(fn, funcs, globals)

		for _, f := range funcs {
			if !seenFunc[f] && !emit.emittedFuncs[f] {
				seenFunc[f] = true
				todo = append(todo, f)
			}
		}
		funcs = funcs[:]

		for _, glob := range globals {
			if !emit.emittedGlobals[glob] {
				emit.emittedGlobals[glob] = true
				emit.global(glob)
			}
		}
		globals = globals[:]

		emit.emittedFuncs[fn] = true
		emit.fn(fn)
	}
}

func (emit *Emitter) fn(fn *ir.Func) {
	emit.ensureSection(Code)
	params := fn.Sig.Function().Params
	pstrs := make([]string, len(params))

	for i, param := range params {
		pstrs[i] = param.String()
	}

	res := fn.Sig.Function().Results
	resstr := ""
	if len(res) == 1 {
		resstr = " " + res[0].String()
	} else if len(res) > 1 {
		resstr = " ("
		for i, r := range res {
			if i != 0 {
				resstr += ","
			}
			resstr += r.String()
		}
		resstr += ")"
	}

	emit.comment("func %s(%s)%s", fn.FullName, strings.Join(pstrs, ", "), resstr)
	emit.line("%s:", fn.FullName)

	for b := 0; b < fn.NumBlocks(); b++ {
		blk := fn.Block(b)

		emit.line(emit.fmter.BlockLabel(blk.IDString()) + ":")
		emit.indent = "    "

		for i := 0; i < blk.NumInstrs(); i++ {
			instr := blk.Instr(i)

			// todo: get inline asm working
			// if instr.Op == op.InlineAsm {
			// 	asm, ok := ir.StringValue(instr.Arg(0).Const())
			// 	if !ok {
			// 		panic("expected string on InlineAsm instruction")
			// 	}
			// 	for _, line := range strings.Split(asm, "\n") {
			// 		if line == "" {
			// 			continue
			// 		}
			// 		if strings.HasSuffix(line, ":") {
			// 			emit.indent = ""
			// 			emit.line("%s", line)
			// 			emit.indent = "    "
			// 		} else {
			// 			emit.line("%s", line)
			// 		}
			// 	}

			// 	continue
			// }

			defs := make([]string, 0, instr.NumDefs())
			for d := 0; d < instr.NumDefs(); d++ {
				def := instr.Def(d)
				str := ""
				switch {
				case def.InReg():
					str = def.Reg().String()
				default:
					str = def.String()
				}
				defs = append(defs, str)
			}

			args := make([]string, 0, instr.NumArgs())
			for a := 0; a < instr.NumArgs(); a++ {
				arg := instr.Arg(a)
				str := ""
				switch {
				case arg.InReg():
					str = arg.Reg().String()
				case arg.IsConst() && arg.Const().Kind() == ir.BoolConst:
					// todo: should probably be an xform
					if b, _ := ir.BoolValue(arg.Const()); b {
						str = "1"
					} else {
						str = "0"
					}
				default:
					str = arg.String()
				}
				args = append(args, str)
			}

			if i == blk.NumInstrs()-1 {
				if blk.NumSuccs() == 1 { // jump
					if b < fn.NumBlocks()-1 && blk.Succ(0) == fn.Block(b+1) {
						// block falls through
						continue
					}
				}
				if blk.NumSuccs() > 0 {
					args = append(args, emit.fmter.BlockLabel(blk.Succ(0).IDString()))
				}
			}

			emit.line("%s", arch.Asm(instr.Op, defs, args))
		}
		emit.indent = ""
	}

	emit.line("")
}

func (emit *Emitter) global(glob *ir.Global) {
	if glob.Value != nil {
		emit.ensureSection(Data)
	} else {
		emit.ensureSection(Bss)
	}
	emit.line("%s:", emit.fmter.GlobalLabel(glob))
	if glob.Value == nil {
		size := glob.Type.Integer().Bits() / sizes.MinAddressableBits()
		emit.line("%s", emit.fmter.Reserve(int(size)))
	} else if str, ok := ir.StringValue(glob.Value); ok {
		emit.line("%s", emit.fmter.Word(emit.fmter.PCRelAddress(int(sizes.WordSize()*2))))

		emit.line("%s", emit.fmter.Word(fmt.Sprintf("%d", len(str))))
		emit.line("%s", emit.fmter.String(str))
	} else if val, ok := ir.IntValue(glob.Value); ok {
		emit.line("%s", emit.fmter.Word(fmt.Sprintf("%d", val)))
	} else {
		panic("todo: implement more types")
	}
	emit.line("")
}

func (emit *Emitter) scan(fn *ir.Func, funcs []*ir.Func, globals []*ir.Global) ([]*ir.Func, []*ir.Global) {
	for b := 0; b < fn.NumBlocks(); b++ {
		blk := fn.Block(b)
		for i := 0; i < blk.NumInstrs(); i++ {
			instr := blk.Instr(i)

			for a := 0; a < instr.NumArgs(); a++ {
				arg := instr.Arg(a)

				if arg.IsConst() {
					constant := arg.Const()
					if fnc, ok := ir.FuncValue(constant); ok {
						funcs = append(funcs, fnc)
					} else if glob, ok := ir.GlobalValue(constant); ok {
						globals = append(globals, glob)
					}
				}
			}
		}
	}

	return funcs, globals
}

func (emit *Emitter) ensureSection(section Section) {
	if emit.section != section {
		emit.line(emit.fmter.Section(section))
		emit.section = section
	}
}

func (emit *Emitter) line(fmtstr string, args ...interface{}) {
	output := fmt.Sprintf(emit.indent+fmtstr, args...)
	fmt.Fprintln(emit.out, output)
}

func (emit *Emitter) comment(fmtstr string, args ...interface{}) {
	output := emit.fmter.Comment(fmt.Sprintf(fmtstr, args...))
	fmt.Fprintln(emit.out, emit.indent+output)
}
