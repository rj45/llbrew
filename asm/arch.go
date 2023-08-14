package asm

import "github.com/rj45/llbrew/ir"

type Arch interface {
	Asm(op ir.Op, defs []string, args []string, emit func(string))
}

var arch Arch

func SetArch(a Arch) {
	arch = a
}
