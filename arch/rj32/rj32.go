package rj32

import (
	"github.com/rj45/llir2asm/arch"
)

type cpuArch struct{}

var _ = arch.Register(cpuArch{})

func (cpuArch) Name() string {
	return "rj32"
}
