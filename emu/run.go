package emu

import (
	"io"
	"os"
	"os/exec"

	"slices"
)

func Run(binfile string, trace bool, wr io.Writer) error {
	args := slices.Clone(emuArgs)

	args = append(args, binfile)
	if trace {
		args = append(args, "-trace")
	}

	runcmd := exec.Command(emuCmd, args...)
	runcmd.Stderr = os.Stderr
	runcmd.Stdout = wr
	runcmd.Stdin = os.Stdin
	return runcmd.Run()
}
