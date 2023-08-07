package rj32

import (
	"go/types"

	"github.com/rj45/llir2asm/arch"
)

type cpuArch struct{}

var _ = arch.Register(cpuArch{})

func (cpuArch) Name() string {
	return "rj32"
}

func isUnsigned(typ types.Type) bool {
	basic, ok := typ.(*types.Basic)
	if !ok {
		return false
	}
	switch basic.Kind() {
	case types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64, types.Uintptr:
		return true
	}
	return false
}
