package rj32

import (
	"strings"
)

type Opcode int

//go:generate go run github.com/dmarkham/enumer -type=Opcode -transform snake

const (
	// Natively implemented instructions
	Nop Opcode = iota
	Rets
	Error
	Halt
	Rcsr
	Wcsr
	Move
	Loadc
	Jump
	Imm
	Call
	Imm2
	Load
	Store
	Loadb
	Storeb
	Add
	Sub
	Addc
	Subc
	Xor
	And
	Or
	Shl
	Shr
	Asr
	IfEq
	IfNe
	IfLt
	IfGe
	IfUlt
	IfUge

	// Psuedoinstructions
	Not
	Neg
	Swap
	IfGt
	IfLe
	IfUgt
	IfUle
	Return

	NumOps
)

func (op Opcode) Asm() string {
	return strings.ReplaceAll(op.String(), "_", ".")
}

func (op Opcode) IsMove() bool {
	return op == Move || op == Swap
}

func (op Opcode) IsCall() bool {
	return op == Call
}

func (op Opcode) IsCommutative() bool {
	return opDefs[op]&commutative != 0
}

func (op Opcode) IsCompare() bool {
	return opDefs[op]&compare != 0
}

func (op Opcode) IsBranch() bool {
	return opDefs[op]&compare != 0
}

func (op Opcode) IsCopy() bool {
	return op == Move
}

func (op Opcode) IsSink() bool {
	return opDefs[op]&sink != 0
}

func (op Opcode) ClobbersArg() bool {
	return opDefs[op]&clobbers != 0
}

func (op Opcode) IsReturn() bool {
	return op == Return
}

type flags uint16

const (
	commutative flags = 1 << iota
	compare
	sink
	clobbers
)

var opDefs = [...]flags{
	Nop:    sink,
	Rets:   sink,
	Error:  sink,
	Halt:   sink,
	Jump:   sink,
	Store:  sink,
	Storeb: sink,
	Add:    clobbers,
	Sub:    clobbers,
	Addc:   clobbers,
	Subc:   clobbers,
	Xor:    clobbers,
	And:    clobbers,
	Or:     clobbers,
	Shl:    clobbers,
	Shr:    clobbers,
	Asr:    clobbers,
	IfEq:   commutative | compare | sink,
	IfNe:   commutative | compare | sink,
	IfLt:   compare | sink,
	IfGe:   compare | sink,
	IfUlt:  compare | sink,
	IfUge:  compare | sink,
	Not:    clobbers,
	Neg:    clobbers,
	IfGt:   compare | sink,
	IfLe:   compare | sink,
	IfUgt:  compare | sink,
	IfUle:  compare | sink,
	Return: sink,
}

func (cpuArch) IsTwoOperand() bool {
	return true
}
