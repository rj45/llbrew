package emu

type Arch interface {
	EmulatorArgs() []string
	EmulatorCmd() string
}

var emuArgs []string
var emuCmd string

func SetArch(a Arch) {
	emuArgs = a.EmulatorArgs()
	emuCmd = a.EmulatorCmd()
}
