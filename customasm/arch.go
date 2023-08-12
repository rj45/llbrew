package customasm

type Arch interface {
	CustomAsmCPUDef() string
	CustomAsmRunAsm() string
	AssemblerFormat() string
}

var arch Arch

func SetArch(a Arch) {
	arch = a
}
