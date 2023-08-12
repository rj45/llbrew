package customasm

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
)

var ErrCustomASMErred = errors.New("customasm had errors")
var ErrCustomASMNotInstalled = errors.New("customasm not found in PATH; please install rust and `cargo install customasm`")
var ErrHelperFilesNotFound = errors.New("could not find, set the LLBREW environment variable to override")

func Assemble(src []byte, binfile string) error {
	_, err := exec.LookPath("customasm")
	if err != nil {
		return ErrCustomASMNotInstalled
	}

	cpudef, err := os.CreateTemp("", "llir2asm_cpudef_*.asm")
	if err != nil {
		return fmt.Errorf("failed to create temp cpudef file for customasm: %w", err)
	}
	defer os.Remove(cpudef.Name())
	_, err = cpudef.WriteString(arch.CustomAsmCPUDef())
	cpudef.Close()
	if err != nil {
		return fmt.Errorf("failed to write cpudef: %w", err)
	}

	runasm, err := os.CreateTemp("", "llir2asm_run_*.asm")
	if err != nil {
		return fmt.Errorf("failed to create temp run.asm file for customasm: %w", err)
	}
	defer os.Remove(runasm.Name())
	_, err = runasm.WriteString(arch.CustomAsmRunAsm())
	runasm.Close()
	if err != nil {
		return fmt.Errorf("failed to write runasm: %w", err)
	}

	asmtemp, err := os.CreateTemp("", "llir2asm_asm_*.asm")
	if err != nil {
		log.Fatalln("failed to create temp asm file for customasm:", err)
	}
	defer os.Remove(asmtemp.Name())
	_, err = asmtemp.Write(src)
	asmtemp.Close()
	if err != nil {
		return fmt.Errorf("failed to write asm: %w", err)
	}

	if binfile == "" {
		bintemp, err := os.CreateTemp("", "llir2asm_*.bin")
		if err != nil {
			log.Fatalln("failed to create temp bin file for customasm:", err)
		}
		bintemp.Close() // customasm will write to it
		defer os.Remove(bintemp.Name())
		binfile = bintemp.Name()
	}

	asmcmd := exec.Command("customasm", "-q",
		"-f", arch.AssemblerFormat(),
		"-o", binfile,
		cpudef.Name(), runasm.Name(), asmtemp.Name())
	log.Println(asmcmd)
	asmcmd.Stderr = os.Stderr
	asmcmd.Stdout = os.Stdout

	err = asmcmd.Run()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return ErrCustomASMErred
		}
		return fmt.Errorf("failed to run customasm: %w", err)
	}

	return nil
}
