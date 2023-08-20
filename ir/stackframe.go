package ir

import (
	"slices"

	"github.com/rj45/llbrew/sizes"
)

var stackLayout = [...]SlotKind{
	ArgSlot,
	SpillSlot,
	AllocaSlot,
	SavedSlot,
	ParamSlot, // in caller's frame
}

// StackFrame represents the stack frame layout of the current function.
// "Slot" IDs are handed out for offsets to values which are stored on the
// stack, and space is allocated according to the largest value that will
// be assigned to that stack slot.
//
// A stack frame looks like this:
//
//	low mem addresses
//	+------------------------+
//	|                        |   <-- SP
//	+------------------------+ \
//	| callee stack arg 0     |  |
//	+------------------------+  |
//	| callee stack arg 1     |  |
//	+------------------------+  |
//	| spill 0                |  |
//	+------------------------+  |
//	| spill 1                |  |
//	+------------------------+  |
//	| alloca 0               |  |
//	+------------------------+   > stack frame
//	| alloca 1               |  |
//	+------------------------+  |
//	| saved reg 0            |  |
//	+------------------------+  |
//	| saved reg 1            |  |
//	+------------------------+  |
//	| saved RA               |  |
//	+------------------------+  |
//	| saved SP               |  | <-- previous SP
//	+------------------------+ /
//	| stack param 0          | \
//	+------------------------+  |
//	| stack param 1          |   > caller's stack frame
//	+------------------------+  |
//	high mem addresses
//
// The SP is saved instead of a frame pointer to save a register. A frame
// pointer may be required to support dynamic stack allocations, but other
// than that, it is not needed. SP is saved to aide in unwinding the
// stack for debugging purposes.
//
// The stack parameters of a called function reside in the caller's stack
// frame.
//
// The order of items on the stack frame is an attempt to limit the size of
// offsets on load/store instructions that may appear frequently, such as
// variables spilled to the stack during register allocation.
type StackFrame struct {
	fn *Func

	nextIDs [NumStackAreas]SlotID

	slots     [NumStackAreas][]slotValue
	sizes     [NumStackAreas][]int
	offsets   [NumStackAreas][]int
	totals    [NumStackAreas]int
	frameSize int
}

func (frame *StackFrame) Func() *Func {
	return frame.fn
}

// NewSlotID returns the next unused SlotID of the given kind
func (frame *StackFrame) NewSlotID(kind SlotKind) SlotID {
	next := frame.nextIDs[kind]
	frame.nextIDs[kind]++
	return SlotID(uint32(kind)<<24 | uint32(next))
}

// SlotID returns a specific slot, making sure that NewSlotID will
// return the next unused one after this if it's not already been
// given out.
func (frame *StackFrame) SlotID(kind SlotKind, index int) SlotID {
	if frame.nextIDs[kind] < SlotID(index+1) {
		frame.nextIDs[kind] = SlotID(index + 1)
	}
	return SlotID(uint32(kind)<<24 | uint32(index))
}

func (frame *StackFrame) offsetOf(slot SlotID) int {
	return frame.offsets[slot.Kind()][slot.Index()]
}

// FrameSize is the total stack frame size minus
func (frame *StackFrame) FrameSize() int {
	return frame.frameSize
}

// ReplaceOffsets replaces all the stack offset variables with the
// actual calculated stack offsets. `Scan` must be called first.
func (frame *StackFrame) ReplaceOffsets() {
	fn := frame.fn
	for kind := range frame.slots {
		for _, sv := range frame.slots[kind] {
			val := sv.value.ValueIn(fn)
			val.SetConst(ConstFor(frame.offsetOf(sv.slot)))
		}
	}
}

// Scan the function for stack variables and calculate the SP offset for them.
func (frame *StackFrame) Scan() {
	for i := 0; i < int(NumStackAreas); i++ {
		frame.slots[i] = nil
		frame.sizes[i] = nil
		frame.offsets[i] = nil
	}

	var blk *Block = nil
	for it := frame.fn.InstrIter(); it.HasNext(); it.Next() {
		instr := it.Instr()
		if instr.blk != blk {
			blk = instr.blk

			for d := 0; d < blk.NumDefs(); d++ {
				frame.countValue(blk.Def(d))
			}
			for a := 0; a < blk.NumArgs(); a++ {
				frame.countValue(blk.Arg(a))
			}
		}

		for d := 0; d < instr.NumDefs(); d++ {
			frame.countValue(instr.Def(d))
		}
		for a := 0; a < instr.NumArgs(); a++ {
			frame.countValue(instr.Arg(a))
		}
	}

	frame.recalcOffsets()
}

func (frame *StackFrame) countValue(v *Value) {
	if !v.OnStack() {
		return
	}
	slot := v.StackSlotID()
	kind := slot.Kind()
	index := slot.Index()

	for len(frame.sizes[kind]) <= index {
		frame.sizes[kind] = append(frame.sizes[kind], 0)
	}

	for len(frame.offsets[kind]) <= index {
		frame.offsets[kind] = append(frame.offsets[kind], 0)
	}

	size := v.Type.SizeOf()
	if frame.sizes[kind][index] < size {
		frame.sizes[kind][index] = size
	}

	sv := slotValue{slot: slot, value: v.ID}
	if !slices.Contains(frame.slots[kind], sv) {
		frame.slots[kind] = append(frame.slots[kind], sv)
	}
}

func (frame *StackFrame) recalcOffsets() {
	clear(frame.totals[:])
	offset := sizes.WordSize() // reserve space for saving SP
	frame.frameSize = offset
	for _, kind := range stackLayout {
		if kind == ParamSlot && frame.frameSize > sizes.WordSize() {
			// skip over SP on the stack
			offset += sizes.WordSize()
		}
		for i, size := range frame.sizes[kind] {
			frame.offsets[kind][i] = offset
			offset += size
			frame.totals[kind] += size

			// param slots are stored on caller's frame
			if kind != ParamSlot {
				frame.frameSize += size
			}
		}
	}
	if frame.frameSize <= sizes.WordSize() {
		// if we would only store SP on the frame
		// then don't bother making a frame
		frame.frameSize = 0
	}
}
