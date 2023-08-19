package ir

import "strconv"

// SlotID represents a variable width "slot" where a value is stored on
// the stack frame.
type SlotID uint32

func (sid SlotID) Kind() SlotKind {
	return SlotKind(sid >> 24)
}

func (sid SlotID) Index() int {
	return int(sid & 0xffffff)
}

func (sid SlotID) String() string {
	kind := sid.Kind()
	if kind == InvalidSlot || kind > NumStackAreas {
		return "s<inv>"
	}
	return strPrefixes[kind] + strconv.Itoa(sid.Index())
}

type SlotKind uint8

const (
	// InvalidSlot is an invalid unassigned slot
	InvalidSlot SlotKind = iota

	// Param slots are where a function's parameters live, which
	// are Args passed in from a calling function
	ParamSlot

	// SavedSlot is where callee saved registers live on the stack
	SavedSlot

	// AllocaSlot is where stack allocated data lives on the stack
	AllocaSlot

	// SpillSlot is where the register allocator stores spilled variables
	// These slots are reused when their previous values are no longer required
	SpillSlot

	// ArgSlot is where a caller stores a callee's parameters before a function
	// call. In other words they are "Args" before a function is called, and
	// become "Params" once the called function starts executing
	ArgSlot

	NumStackAreas
)

var strPrefixes = [...]string{
	InvalidSlot: "!",
	ParamSlot:   "sp",
	SavedSlot:   "sr",
	AllocaSlot:  "sm",
	SpillSlot:   "ss",
	ArgSlot:     "sa",
}

type slotValue struct {
	slot  SlotID
	value ID
}
