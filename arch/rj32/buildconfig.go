package rj32

import _ "embed"

func (cpuArch) AssemblerFormat() string {
	return "logisim16"
}

func (cpuArch) EmulatorCmd() string {
	return "emurj"
}

func (cpuArch) EmulatorArgs() []string {
	return []string{"-run"}
}

//go:embed customasm/cpudef.asm
var cpudef string

func (cpuArch) CustomAsmCPUDef() string {
	return cpudef
}

//go:embed customasm/run.asm
var runasm string

func (cpuArch) CustomAsmRunAsm() string {
	return runasm
}
