package arch

import (
	"log"
	"strings"

	"github.com/rj45/llir2asm/asm"
	"github.com/rj45/llir2asm/customasm"
	"github.com/rj45/llir2asm/ir/reg"
	"github.com/rj45/llir2asm/sizes"
	"github.com/rj45/llir2asm/xform"
)

const defaultArch = "rj32"

type Architecture interface {
	Name() string
	reg.Arch
	sizes.Arch
	xform.Arch
	asm.Arch
	customasm.Arch
}

var arch Architecture

var arches map[string]Architecture

func Arch() Architecture {
	return arch
}

func Register(a Architecture) int {
	if arches == nil {
		arches = make(map[string]Architecture)
	}
	name := strings.ToLower(a.Name())
	arches[name] = a
	if name == defaultArch {
		SetArch(name)
	}
	return 0
}

func SetArch(name string) {
	arch = arches[strings.ToLower(name)]
	if arch == nil {
		log.Panicf("unknown arch %s", name)
	}
	reg.SetArch(arch)
	sizes.SetArch(arch)
	xform.SetArch(arch)
	asm.SetArch(arch)
	customasm.SetArch(arch)
}
