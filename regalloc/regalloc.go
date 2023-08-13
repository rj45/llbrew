package regalloc

import (
	"errors"
	"fmt"
	"log"

	"github.com/rj45/llir2asm/ir"
	"github.com/rj45/llir2asm/ir/reg"
)

type RegAlloc struct {
	fn *ir.Func

	info []blockInfo

	iGraph iGraph
}

type blockInfo struct {
	liveIns  map[ir.ID]struct{}
	liveOuts map[ir.ID]struct{}
}

// NewRegAlloc returns a new register allocator, ready to have
// the registers allocated for the given function.
func NewRegAlloc(fn *ir.Func) *RegAlloc {
	info := make([]blockInfo, fn.NumBlocks())

	return &RegAlloc{
		fn:   fn,
		info: info,
	}
}

var ErrCriticalEdges = errors.New("the CFG has critical edges")

// CheckInput will verify the structure of the
// input code, which is useful in tests and fuzzing.
func (ra *RegAlloc) CheckInput() error {
	fn := ra.fn

	for i := 0; i < fn.NumBlocks(); i++ {
		pred := fn.Block(i)

		// if fn.Block(0).NumInstrs() < 1 {
		// 	// skip empty blocks
		// 	continue
		// }

		if pred.NumSuccs() > 1 {
			for i := 0; i < pred.NumSuccs(); i++ {
				succ := pred.Succ(i)

				// the first block has an implicit pred
				extraPred := 0
				if succ.Index() == 0 {
					extraPred = 1
				}

				if (succ.NumPreds() + extraPred) > 1 {
					fmt.Println("Critical edge:", pred, "->", succ)
					return ErrCriticalEdges
				}
			}
		}
	}
	return nil
}

// Allocate will run the allocator and assign a physical
// register or stack slot to each Value that needs one.
//
// This uses a graph colouring algorithm designed for use on SSA
// code. The SSA code should have copies added for block
// defs/args such that the allocator is free to choose different
// registers as it crosses that boundary. Attempts will be
// made to strongly prefer choosing the same register across
// a copy and across block args/defs to minimize the number
// of copies in the resultant code.
//
// A simpler algorithm is chosen rather than a fast algorithm
// because this can easily be a very complex piece of code and
// simpler code is easier to ensure is correct. It's extremely
// important that the register allocator produces correct code
// or else pretty much anything can happen.
//
// To that end, there is a verifier in the verify sub-package
// which will double check the work of the allocator. See that
// code for more information on how it works.
func (ra *RegAlloc) Allocate() error {
	if err := ra.liveInOutScan(); err != nil {
		return err
	}

	ra.buildInterferenceGraph()
	ra.preColour()
	// ra.iGraph.coalesceMoves()
	ra.iGraph.pickColours()
	if err := ra.assignRegisters(); err != nil {
		return err
	}

	return nil
}

var regList []reg.Reg
var savedStart uint16

const dontColour = 0xffff

// preColour finds all the values with already assigned registers and sets their colour to them
func (ra *RegAlloc) preColour() {
	if regList == nil {
		regList = append(regList, reg.ArgRegs...)
		regList = append(regList, reg.TempRegs...)
		savedStart = uint16(len(regList) + 1)
		regList = append(regList, reg.SavedRegs...)
	}

	for id := range ra.iGraph.nodes {
		node := &ra.iGraph.nodes[id]

		val := node.val.ValueIn(ra.fn)
		if val == nil {
			continue
		}

		if val.InReg() && val.Reg() != reg.None {
			found := false
			for i, reg := range regList {
				if val.Reg() == reg {
					node.colour = uint16(i + 1)
					found = true
					break
				}
			}

			if !found {
				// mark node not to be coloured
				node.colour = dontColour
			}
		}
	}
}

var ErrTooManyRequiredRegisters = errors.New("too many required registers")

func (ra *RegAlloc) assignRegisters() error {
	for id := range ra.iGraph.nodes {
		node := &ra.iGraph.nodes[id]

		val := node.val.ValueIn(ra.fn)

		if val == nil {
			continue
		}

		if node.colour == dontColour {
			continue
		}

		// colour zero is "noColour", so subtract one
		regIndex := int(node.colour - 1)

		if regIndex >= len(regList) {
			return ErrTooManyRequiredRegisters
		}

		if val.InReg() && val.Reg() != reg.None && val.Reg() != regList[regIndex] {
			log.Panicf("setting pre-set %s id %d reg %s to %s", ra.fn.Name, val.ID, val, regList[regIndex])
		}

		val.SetReg(regList[regIndex])

		for _, id := range node.merged {
			val := id.ValueIn(ra.fn)
			if val.NeedsReg() {
				val.SetReg(regList[regIndex])
			}
		}
	}

	return nil
}
