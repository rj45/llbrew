package emu

import (
	"os"
	"os/exec"

	"slices"
)

func Run(binfile string, trace bool) error {
	args := slices.Clone(emuArgs)

	if trace {
		args = append(args, "-trace")
	}
	args = append(args, binfile)

	runcmd := exec.Command(emuCmd, args...)
	runcmd.Stderr = os.Stderr
	runcmd.Stdout = os.Stdout
	runcmd.Stdin = os.Stdin
	return runcmd.Run()
}
