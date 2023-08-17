package emu

import (
	"os"
	"os/exec"

	"slices"
)

func Run(binfile string, trace bool) error {
	args := slices.Clone(emuArgs)

	args = append(args, binfile)
	if trace {
		args = append(args, "-trace")
	}

	runcmd := exec.Command(emuCmd, args...)
	runcmd.Stderr = os.Stderr
	runcmd.Stdout = os.Stdout
	runcmd.Stdin = os.Stdin
	return runcmd.Run()
}
