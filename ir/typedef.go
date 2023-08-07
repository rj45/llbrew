package ir

import "github.com/rj45/llir2asm/ir/typ"

// TypeDef is a type definition
type TypeDef struct {
	pkg *Package

	Name       string
	Referenced bool

	Type typ.Type
}
