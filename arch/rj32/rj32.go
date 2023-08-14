package rj32

import (
	"github.com/rj45/llbrew/arch"
)

type cpuArch struct{}

var _ = arch.Register(cpuArch{})

func (cpuArch) Name() string {
	return "rj32"
}
