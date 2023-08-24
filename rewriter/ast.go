package rewriter

import (
	"go/ast"

	"github.com/rj45/llbrew/ir"
	"github.com/rj45/llbrew/ir/typ"
)

type rule struct {
	match *instruction
	cond  ast.Expr
	repl  value
	vars  map[string]*variable

	loc string
	pkg string

	// state
	instrs int
	values int
}

type instruction struct {
	opstr string
	op    ir.Op
	typ   typ.Type
	defs  []value
	args  []value
}

func (instr *instruction) kind() vkind {
	return instrKind
}

func (instr *instruction) String() string {
	str := "("
	if len(instr.defs) > 0 {
		str += "["
		for i, def := range instr.defs {
			if i != 0 {
				str += " "
			}
			str += def.String()
		}
		str += "] "
	}
	str += instr.opstr
	for _, arg := range instr.args {
		str += " " + arg.String()
	}

	str += ")"
	return str
}

type vkind uint8

const (
	varKind vkind = iota
	instrKind
	constKind
)

type value interface {
	kind() vkind
	String() string
}

type variable struct {
	name  string
	count int
}

var underscore *variable = &variable{name: "_"}
var elipsis *variable = &variable{name: "..."}

func (v *variable) kind() vkind {
	return varKind
}

func (v *variable) String() string {
	return v.name
}

type constant string

func (c constant) kind() vkind {
	return constKind
}

func (c constant) String() string {
	return string(c)
}
