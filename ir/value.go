package ir

import (
	"log"

	"github.com/rj45/llbrew/ir/reg"
	"github.com/rj45/llbrew/ir/typ"
)

// Value is a single value that may be stored in a
// single place. This may be a constant or variable,
// stored in a temp, register or on the stack.
type Value struct {
	stg

	ID

	// Type is the type of the Value
	Type typ.Type

	def  *User
	uses []*User

	usestorage [2]*User
}

// Location is the location of a Value
type Location uint8

const (
	InTemp Location = iota
	InConst
	InReg

	// Value is stored on the stack
	OnStack
)

// init initializes the Value.
func (val *Value) init(id ID, typ typ.Type) {
	val.uses = val.usestorage[:0]
	val.ID = id
	val.Type = typ
	val.SetTemp()
}

// Func returns the containing Func.
func (val *Value) Func() *Func {
	return val.def.Block().fn
}

// Def returns the Instr defining the Value,
// or nil if it's not defined
func (val *Value) Def() *User {
	return val.def
}

// NumUses returns the number of uses
func (val *Value) NumUses() int {
	return len(val.uses)
}

// Use returns the ith Instr using this Value
func (val *Value) Use(i int) *User {
	return val.uses[i]
}

// ReplaceUsesWith will go through each use of
// val and replace it with other. Does not modify
// any definitions.
func (val *Value) ReplaceUsesWith(other *Value) {
	tries := 0
	for len(val.uses) > 0 {
		tries++
		use := val.uses[len(val.uses)-1]
		if tries > 1000 {
			log.Panicln("bug in uses ", val, other)
		}
		i := use.ArgIndex(val)
		if i < 0 {
			panic("couldn't find use!")
		}
		use.ReplaceArg(i, other)
	}
}

func (val *Value) IsDefinedByOp(op Op) bool {
	if val.def == nil {
		return false
	}
	if !val.def.IsInstr() {
		return false
	}

	return val.def.Instr().Op == op
}

// NeedsReg indicates if this Value should be allocated
// a register
func (val *Value) NeedsReg() bool {
	return !val.IsConst() && !val.IsBlock() && !val.OnStack()
}

// stg is the storage for a value
type stg interface {
	Location() Location
}

// temps

type tempStg struct{}

func (tempStg) Location() Location { return InTemp }

// InTemp indicates the value is in a temp.
func (val *Value) InTemp() bool {
	return val.Location() == InTemp
}

// Temp returns which temp if the value is in a temp.
func (val *Value) Temp() ID {
	if val.Location() == InTemp {
		return val.ID
	}
	return Placeholder
}

// SetTemp turns the value into a temp.
func (val *Value) SetTemp() {
	val.stg = tempStg{}
}

// regs

type regStg struct{ r reg.Reg }

func (regStg) Location() Location { return InReg }

// InReg indicates if the Value is in a register.
func (val *Value) InReg() bool {
	return val.Location() == InReg
}

// Reg returns which register the Value is in,
// otherwise reg.None if its not in a register.
func (val *Value) Reg() reg.Reg {
	if val.Location() == InReg {
		return val.stg.(regStg).r
	}
	return reg.None
}

// SetReg puts the value in the specified register.
func (val *Value) SetReg(reg reg.Reg) {
	if !val.NeedsReg() {
		panic("assigned reg to non-reg value")
	}
	val.stg = regStg{reg}
}

// stack slots

type stackStg SlotID

func (stackStg) Location() Location { return OnStack }

// OnStack returns whether the value is stored on the stack.
func (val *Value) OnStack() bool {
	return val.Location() == OnStack
}

// StackSlotID returns which spill slot the Value is in, or -1 if not in a spill slot.
func (val *Value) StackSlotID() SlotID {
	if val.Location() == OnStack {
		return SlotID(val.stg.(stackStg))
	}
	return 0
}

// SetStackSlot puts the Value on the stack at the specified slot.
func (val *Value) SetStackSlot(slot SlotID) {
	val.stg = stackStg(slot)
}

// SetSlotIndex sets the stack slot to a specific index
func (val *Value) SetSlotIndex(kind SlotKind, index int) {
	val.SetStackSlot(val.def.fn.Frame.SlotID(kind, index))
}

// MoveToStack moves the value onto the stack in the next slot available
func (val *Value) MoveToStack(kind SlotKind) {
	val.SetStackSlot(val.def.fn.Frame.NewSlotID(kind))
}

// const

// IsConst returns if the Value is constant.
func (val *Value) IsConst() bool {
	return val.Location() == InConst
}

// Const returns the constant value of the Value or NotConst if not constant.
func (val *Value) Const() Const {
	if val.Location() == InConst {
		return val.stg.(Const)
	}
	return notConst{}
}

// SetConst makes the Value the specified constant.
func (val *Value) SetConst(con Const) {
	val.stg = con
}

// util funcs

// addUse adds the instr as a use of this value.
func (val *Value) addUse(user *User) {
	val.uses = append(val.uses, user)
}

// removeUse removes the isntr as a use of this value.
func (val *Value) removeUse(user *User) {
	index := -1
	for i, v := range val.uses {
		if v == user {
			index = i
			break
		}
	}
	if index < 0 {
		log.Panicf("%v does not have use %v", val, user)
	}
	val.uses = append(val.uses[:index], val.uses[index+1:]...)
}
