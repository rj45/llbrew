package ir

import "github.com/rj45/llbrew/ir/typ"

// Global is a global variable or literal stored in memory
type Global struct {
	pkg *Package

	Name       string
	FullName   string
	Type       typ.Type
	Referenced bool

	// initial value
	Value Const
}

func (glob *Global) String() string {
	return glob.FullName
}

func (glob *Global) Package() *Package {
	return glob.pkg
}
