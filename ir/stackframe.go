package ir

import "slices"

/*
A stack frame looks like this:

low mem addresses
+------------------------+
|                        |  |
+------------------------+  |
| callee stack arg 3     |  |
+------------------------+  |
| callee stack arg 4     |  |
+------------------------+   > callee frame
| callee SP              |  |
+------------------------+  |
| callee RA              | /  <-- SP
+------------------------+ \
| alloca 0               |  |
+------------------------+  |
| alloca 1               |  |
+------------------------+  |
| spill 0                |  |
+------------------------+   > stack frame
| spill 1                |  |
+------------------------+  |
| saved reg 0            |  |
+------------------------+  |
| saved reg 1            |  |
+------------------------+  |
| stack param 3          |  |
+------------------------+  |
| stack param 4          |  |
+------------------------+  |
| saved SP               |  |
+------------------------+  |
| saved RA               |  | <-- previous SP
+------------------------+ /
high mem addresses

*/

// StackFrame represents the stack frame layout of the current function.
// "Slot" IDs are handed out to values which are stored on the stack,
// and space is allocated according to the largest value that will
// be assigned to that stack slot.
type StackFrame struct {
	fn *Func

	nextIDs [NumStackAreas]SlotID

	slots   [NumStackAreas][]slotValue
	sizes   [NumStackAreas][]int
	offsets [NumStackAreas][]int
	totals  [NumStackAreas]int
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

func (frame *StackFrame) OffsetOf(slot SlotID) int {
	return frame.offsets[slot.Kind()][slot.Index()]
}

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
	for kind := 0; kind < int(NumStackAreas); kind++ {
		offset := 0
		for i, size := range frame.sizes[kind] {
			frame.offsets[kind][i] = offset
			offset += size
			frame.totals[kind] += size
		}
	}
}
