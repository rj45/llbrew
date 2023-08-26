package op

//go:generate go run github.com/dmarkham/enumer -type=Op -transform title-lower

type def uint16

const (
	sink def = 1 << iota
	compare
	constant
	move
	commute
	branch
)

type Op uint8

func (op Op) IsCompare() bool {
	return opDefs[op]&compare != 0
}

func (op Op) IsSink() bool {
	return opDefs[op]&sink != 0
}

func (op Op) IsConst() bool {
	return opDefs[op]&constant != 0
}

func (op Op) IsCopy() bool {
	return opDefs[op]&move != 0
}

func (op Op) IsCommutative() bool {
	return opDefs[op]&commute != 0
}

func (op Op) IsCall() bool {
	return op == Call
}

func (op Op) ClobbersArg() bool {
	return false
}

func (op Op) IsBranch() bool {
	return opDefs[op]&branch != 0
}

func (op Op) IsReturn() bool {
	return op == Ret
}

func (op Op) Opposite() Op {
	switch op {
	case Equal:
		return NotEqual
	case NotEqual:
		return Equal
	case Less:
		return GreaterEqual
	case LessEqual:
		return Greater
	case Greater:
		return LessEqual
	case GreaterEqual:
		return Less
	case ULess:
		return UGreaterEqual
	case ULessEqual:
		return UGreater
	case UGreater:
		return ULessEqual
	case UGreaterEqual:
		return ULess
	}
	return op
}

const (
	Invalid Op = iota

	// Control flow
	Ret
	If
	Jump
	Switch
	IndirectBr
	Invoke
	Unreachable

	// Compares
	FCmp
	Equal
	NotEqual
	Less
	LessEqual
	Greater
	GreaterEqual
	ULess
	ULessEqual
	UGreater
	UGreaterEqual

	// Standard Binary Operators
	Add
	FAdd
	Sub
	FSub
	Mul
	FMul
	UDiv
	SDiv
	FDiv
	URem
	SRem
	FRem

	// Logical Operators
	Shl
	LShr
	AShr
	And
	Or
	Xor

	// Memory Operators
	Alloca
	Load
	Store
	GetElementPtr

	// Cast Operators
	Trunc
	ZExt
	SExt
	FPToUI
	FPToSI
	UIToFP
	SIToFP
	FPTrunc
	FPExt
	PtrToInt
	IntToPtr
	BitCast

	// Other Operators
	Call
	Select
	Copy

	// UserOp1
	// UserOp2
	VAArg
	ExtractElement
	InsertElement
	ShuffleVector
	ExtractValue
	InsertValue

	NumOps
)

var opDefs = [...]def{
	Copy:          move,
	Store:         sink,
	Ret:           sink,
	If:            sink,
	Jump:          sink,
	Switch:        sink,
	IndirectBr:    sink,
	Invoke:        sink,
	Unreachable:   sink,
	Add:           commute,
	Mul:           commute,
	And:           commute,
	Or:            commute,
	Xor:           commute,
	Equal:         compare | commute,
	NotEqual:      compare | commute,
	Less:          compare,
	LessEqual:     compare,
	Greater:       compare,
	GreaterEqual:  compare,
	ULess:         compare,
	ULessEqual:    compare,
	UGreater:      compare,
	UGreaterEqual: compare,
	NumOps:        0, // make sure array is large enough
}
