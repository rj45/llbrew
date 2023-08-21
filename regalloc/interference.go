package regalloc

import (
	"flag"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/rj45/llbrew/ir"
)

type iNodeID uint32

type iGraph struct {
	fn      *ir.Func
	ra      *RegAlloc
	nodes   []iNode
	valNode map[ir.ID]iNodeID

	maxColour uint16
}

type iNode struct {
	val        ir.ID
	interferes map[iNodeID]struct{}
	moves      []iNodeID
	merged     []ir.ID

	colour uint16
	order  uint16

	callerSaved bool
}

var debugalloc = flag.Bool("debugalloc", false, "emit log messages for allocation decisions")

func (ig *iGraph) dbg(format string, args ...interface{}) {
	if *debugalloc {
		newargs := make([]interface{}, len(args))
		for i, arg := range args {
			newargs[i] = arg
			switch arg := arg.(type) {
			case iNodeID:
				if ig.nodes[arg].val == 0 {
					newargs[i] = "<brkn>"
				} else {
					newargs[i] = ig.nodes[arg].val.ValueIn(ig.fn).String()
				}
			case *iNode:
				val := arg.val.ValueIn(ig.fn)
				if val == nil {
					newargs[i] = "<removed>"
				} else {
					newargs[i] = val.String()
				}
			case ir.ID:
				val := arg.ValueIn(ig.fn)
				if val == nil {
					newargs[i] = "<removed>"
				} else {
					newargs[i] = val.String()
				}
			case map[ir.ID]struct{}:
				var list []string
				for id := range arg {
					list = append(list, id.ValueIn(ig.fn).String())
				}
				newargs[i] = strings.Join(list, ",")
			case []iNodeID:
				var list []string
				for _, id := range arg {
					if ig.nodes[id].val != 0 {
						list = append(list, ig.nodes[id].val.ValueIn(ig.fn).String())
					} else {
						list = append(list, "<0>")
					}
				}
				newargs[i] = strings.Join(list, ",")
			}
		}
		fmt.Printf(format+"\n", newargs...)
	}
}

func (ig *iGraph) addNode(id ir.ID) iNodeID {
	if !id.ValueIn(ig.fn).NeedsReg() {
		panic("attempt to add non reg value: " + id.ValueIn(ig.fn).IDString())
	}

	nodeID, found := ig.valNode[id]
	if !found {
		nodeID = iNodeID(len(ig.nodes))
		ig.nodes = append(ig.nodes, iNode{
			val: id,
		})
		ig.valNode[id] = nodeID
		ig.dbg("%s: add interference node %s", ig.fn.Name, nodeID)
	}

	return nodeID
}

func (ig *iGraph) addEdge(var1 ir.ID, var2 ir.ID) {
	node1ID := ig.addNode(var1)
	node2ID := ig.addNode(var2)

	if var1 == var2 {
		// don't add edges between ourself
		return
	}

	for _, pair := range [2][2]iNodeID{{node1ID, node2ID}, {node2ID, node1ID}} {
		node := &ig.nodes[pair[0]]
		neighbor := pair[1]
		if _, found := node.interferes[neighbor]; !found {
			if node.interferes == nil {
				node.interferes = make(map[iNodeID]struct{})
			}

			// add to the interferes map
			node.interferes[neighbor] = struct{}{}

			ig.dbg("%s: add interference edge %s -- %s", ig.fn.Name, node, neighbor)
		}
	}
}

func (ig *iGraph) addMove(var1 ir.ID, var2 ir.ID) {
	node1ID := ig.addNode(var1)
	node2ID := ig.addNode(var2)

	if node1ID == node2ID {
		// don't add moves between ourself
		return
	}

	ig.linkMove(node1ID, node2ID)
	ig.linkMove(node2ID, node1ID)
}

func (ig *iGraph) linkMove(nodeID iNodeID, neighbor iNodeID) {
	node := &ig.nodes[nodeID]

	found := false
	for _, id := range node.moves {
		if id == neighbor {
			found = true
		}
	}

	if !found {
		// add it to the moves list
		ig.dbg("%s: move nodes before: %s: %s -- adding %s", ig.fn.Name, node, node.moves, neighbor)
		node.moves = append(node.moves, neighbor)
		ig.dbg("%s: add move edge %s -- %s", ig.fn.Name, nodeID, neighbor)
		ig.dbg("%s: move nodes now: %s: %s", ig.fn.Name, node, node.moves)
	}
}

func (ig *iGraph) merge(var1 ir.ID, var2 ir.ID) bool {
	node1ID := ig.addNode(var1)
	node2ID := ig.addNode(var2)

	if var1 == var2 {
		// don't merge ourself
		return false
	}

	node1 := &ig.nodes[node1ID]
	node2 := &ig.nodes[node2ID]

	for _, list := range [2][]ir.ID{{node2.val}, node2.merged} {
		for _, val := range list {
			found := false
			for _, mval := range node1.merged {
				if mval == val {
					found = true
					break
				}
			}
			if !found {
				ig.dbg("%s: merging %s & %s -- pulling in %s", ig.fn.Name, node1, node2, val)
				node1.merged = append(node1.merged, val)
			}
		}
	}

	for _, m := range node1.merged {
		ig.valNode[m] = node1ID
	}

	// update the node moves
	for _, move := range node2.moves {
		moves := ig.nodes[move].moves
		idx := slices.Index(moves, node2ID)
		if idx >= 0 {
			moves = slices.Delete(moves, idx, idx)
			if !slices.Contains(moves, node1ID) {
				moves = append(moves, node1ID)
			}
		}

		ig.nodes[move].moves = moves

		if !slices.Contains(node1.moves, move) {
			node1.moves = append(node1.moves, move)
		}
	}

	if node1.interferes == nil && len(node2.interferes) > 0 {
		node1.interferes = make(map[iNodeID]struct{})
	}

	for interferance := range node2.interferes {
		node1.interferes[interferance] = struct{}{}
		delete(ig.nodes[interferance].interferes, node2ID)
	}

	if node2.callerSaved {
		node1.callerSaved = true
	}

	if node2.colour != noColour && node1.colour != noColour && node2.colour != node1.colour {
		log.Panicf("%s: tried to merge two pre-coloured nodes %s and %s", ig.fn.Name, node1.val.InstrIn(ig.fn), node2.val.InstrIn(ig.fn))
	} else if node2.colour != noColour {
		node1.colour = node2.colour
	}

	ig.dbg("%s: merged %s -- %s", ig.fn.Name, node1, node2)

	// clear out the node
	*node2 = iNode{}

	return true
}

// buildInterferenceGraph takes the liveness information and builds a
// graph where nodes in the graph represent variables, and edges between
// the nodes represent variables that are live at the same time, in other
// words, variables that interfere with one another. This is done in order
// to aide in colouring the graph with non-interfering registers.
func (ra *RegAlloc) buildInterferenceGraph() {
	ig := &ra.iGraph
	ig.ra = ra
	ig.fn = ra.fn
	fn := ra.fn

	ig.nodes = nil
	ig.valNode = make(map[ir.ID]iNodeID)

	for i := 0; i < fn.NumBlocks(); i++ {
		blk := fn.Block(i)
		info := ra.info[blk.Index()]

		live := make(map[ir.ID]struct{})
		for k := range info.liveOuts {
			live[k] = struct{}{}
		}

		// block args are live immediately before leaving the block
		// and there is an implicit move between them and the defs of
		// succ blocks
		offset := 0
		for s := 0; s < blk.NumSuccs(); s++ {
			succ := blk.Succ(s)
			for d := 0; d < succ.NumDefs(); d++ {
				def := succ.Def(d)
				arg := blk.Arg(offset + d)

				live[arg.ID] = struct{}{}

				ig.merge(def.ID, arg.ID)
				// ig.addMove(def.ID, arg.ID)
			}

			offset += succ.NumDefs()
		}

		// all currently live variables interfere
		for id1 := range live {
			for id2 := range live {
				if id1 != id2 {
					ig.addEdge(id1, id2)
				}
			}
		}

		for j := blk.NumInstrs() - 1; j >= 0; j-- {
			instr := blk.Instr(j)

			// all defs interfere with one another, so removing it from the
			// live set should be done after adding edges
			for d := 0; d < instr.NumDefs(); d++ {
				def := instr.Def(d)
				if def.NeedsReg() {
					// make sure the node is in the graph, even if there's no
					// other live values at the time
					ig.addNode(def.ID)

					// make sure all live vars are marked as interfering
					for id := range live {
						ig.addEdge(def.ID, id)
					}

					// if it's a move (aka copy)
					if instr.Op.IsCopy() && instr.Arg(d).NeedsReg() {
						// add the move between the corresponding defs and args
						ig.addMove(def.ID, instr.Arg(d).ID)
					}
				}
			}

			// now we can remove each def from the live set
			for d := 0; d < instr.NumDefs(); d++ {
				def := instr.Def(d)
				if def.NeedsReg() {
					// def is now no longer live
					delete(live, def.ID)
				}
			}

			// at a call site, any variables live across the call site must not be
			// assigned to caller saved registers, otherwise the variable should be
			// spilled which is handled separately
			if instr.Op.IsCall() {
				for id := range live {
					node := &ig.nodes[ig.valNode[id]]
					node.callerSaved = true
					ig.dbg("%s: marking val %s in val %s as caller saved", ra.fn.Name, node.val.ValueIn(ra.fn), node.val)
				}
			}

			// if the instruction clobbers its first arg (aka it's two operand) then
			// ensure they are assigned the same register by merging the nodes
			if instr.ClobbersArg() {
				ig.merge(instr.Def(0).ID, instr.Arg(0).ID)
			}

			// mark each used arg as now live
			for u := 0; u < instr.NumArgs(); u++ {
				use := instr.Arg(u)
				if use.NeedsReg() {
					live[use.ID] = struct{}{}
				}
			}
		}

		for d := 0; d < blk.NumDefs(); d++ {
			def := blk.Def(d)
			if def.NeedsReg() {
				// make sure the node is in the graph, even if there's no
				// other live values at the time
				ig.addNode(def.ID)

				// make sure all live vars are marked as interfering
				for id := range live {
					ig.addEdge(def.ID, id)
				}

				delete(live, def.ID)
			}
		}
	}
}

// try to merge moves that don't interfere with each other
func (ig *iGraph) coalesceMoves() {
	changed := true
	for changed {
		changed = false
		for _, nd := range ig.nodes {
			if len(nd.moves) == 0 {
				continue
			}

			if nd.val == 0 {
				// already merged
				continue
			}

			interferes := false
			for _, id1 := range nd.moves {
				if _, found := nd.interferes[id1]; found {
					interferes = true
					break
				}

				for _, id2 := range nd.moves {
					if id1 == id2 {
						continue
					}
					if _, found := ig.nodes[id1].interferes[id2]; found {
						interferes = true
						break
					}
				}
			}

			if !interferes {
				ig.dbg("%s: moves do not interfere: %v", ig.fn.Name, nd.moves)

				for _, id := range nd.moves {
					if ig.nodes[id].val == 0 {
						continue
					}
					// don't coalesce pre-coloured nodes
					// todo: figure out if this is safe and do it anyway
					if ig.nodes[id].colour > 0 || nd.colour > 0 {
						continue
					}

					if ig.merge(nd.val, ig.nodes[id].val) {
						changed = true
						break
					}

				}
				if changed {
					break
				}
			}
		}
	}
}

// findPerfectEliminationOrder finds the perfect elimination order by
// using the max cardinality search algorithm. This is done because
// the graph should be chordal thanks to SSA. Chordal graphs can
// be optimally coloured in reverse perfect elimination order.
// There are other algorithms that could find the PEO as well,
// such as lexicographic breadth first search. This seemed simpler
// though it may be slower (not sure).
func (ig *iGraph) findPerfectEliminationOrder() []iNodeID {
	marked := make(map[iNodeID]struct{})
	output := make([]iNodeID, 0, len(ig.nodes))
	unmarked := make([]iNodeID, len(ig.nodes))
	for i := range ig.nodes {
		unmarked[i] = iNodeID(i)
	}

	// for each unmarked node
	for len(unmarked) > 0 {
		// find the unmarked node with the most marked neighbors
		maxNode := unmarked[0]
		maxI := 0
		maxCard := -1
		for i, cand := range unmarked {
			card := 0
			for neighbor := range ig.nodes[cand].interferes {
				if _, found := marked[neighbor]; found {
					card++
				}
			}
			// hasMoreMoves := len(ig.nodes[cand].moves) > len(ig.nodes[maxI].moves)
			if card > maxCard {
				maxI = i
				maxNode = cand
				maxCard = card
			}
		}

		// remove node from unmarked list. Order doesn't matter
		// so the faster way of removing an item from the slice works.
		unmarked[maxI] = unmarked[len(unmarked)-1]
		unmarked = unmarked[:len(unmarked)-1]

		// mark the node
		marked[maxNode] = struct{}{}

		// add node to output
		output = append(output, maxNode)

		ig.nodes[maxNode].order = uint16(len(output) - 1)
	}

	return output
}

func (ig *iGraph) pickColours() {
	order := ig.findPerfectEliminationOrder()

	// pick colours in reverse perfect elimination order
	for i := len(order) - 1; i >= 0; i-- {
		nodeID := order[i]
		node := &ig.nodes[nodeID]
		if node.val == 0 {
			continue
		}
		node.pickColour(ig)
	}
}

const noColour uint16 = 0

func (nd *iNode) findMostUsedMoveColour(ig *iGraph, moves []iNodeID) (iNodeID, int) {
	var best iNodeID
	uses := -1

	// try to pick a move colour if that colour doesn't
	// interfere with any others
	for _, mv := range moves {
		moveColour := ig.nodes[mv].colour

		if _, ok := nd.interferes[mv]; ok {
			ig.dbg("move interferes: %s -- %s", nd, mv)
			continue
		}

		// skip if the move node has not already been assigned a colour
		if moveColour == noColour || moveColour == dontColour {
			continue
		}

		// check if that colour interferes with any neighbors
		interferes := false
		for nb := range nd.interferes {
			if mv == nb {
				continue
			}
			if ig.nodes[nb].colour == moveColour {
				interferes = true
				break
			}
		}

		// if it doesn't interfere  and the move colour is caller saved if it needs to be
		if !interferes && (!nd.callerSaved || moveColour >= ig.ra.savedStart) {
			val := ig.nodes[mv].val.ValueIn(ig.fn)
			if val.NumUses() > uses {
				uses = val.NumUses()
				best = mv
			}
		}
	}
	return best, uses
}

// findMoveNodes recursively searches for all connected move nodes
func (nd *iNode) findMoveNodes(ig *iGraph, moves []iNodeID) []iNodeID {
	for _, mv := range nd.moves {
		if ig.nodes[mv].val == 0 {
			continue
		}
		if slices.Contains(moves, mv) {
			continue
		}

		moves = append(moves, mv)
		moves = ig.nodes[mv].findMoveNodes(ig, moves)
	}
	return moves
}

func (nd *iNode) checkIfInterferes(ig *iGraph, colour uint16) bool {
	interferes := false
	// for each neighbour in the interferences
	for nb := range nd.interferes {
		if ig.nodes[nb].val == 0 {
			continue
		}

		// if the neighbour already has this colour
		if ig.nodes[nb].colour == colour {
			// then it interferes and we can't use it
			ig.dbg("%s: checking node %s: fail: colour %d already assigned to %s", ig.fn.Name, nd.val, colour, nb)
			interferes = true
			break
		}
		ig.dbg("%s: checking node %s: colour %d not assigned to %s", ig.fn.Name, nd.val, colour, nb)
	}
	return interferes
}

func (nd *iNode) checkIfInterferesWithMove(ig *iGraph, moves []iNodeID, colour uint16) bool {
	for _, mv := range moves {
		// if _, ok := nd.interferes[mv]; ok {
		// 	continue
		// }
		if ig.nodes[mv].checkIfInterferes(ig, colour) {
			return true
		}
	}
	return false
}

func (nd *iNode) pickColour(ig *iGraph) {
	if nd.colour != noColour {
		ig.dbg("%s: %s already has colour %d", ig.fn.Name, nd, nd.colour)
		// already coloured
		return
	}

	moves := nd.findMoveNodes(ig, nil)

	best, uses := nd.findMostUsedMoveColour(ig, moves)

	if uses >= 0 {
		moveColour := ig.nodes[best].colour
		nd.colour = moveColour
		ig.dbg("%s: pick move colour %d for %s", ig.fn.Name, nd.colour, nd)
		return
	}

	// if the node must be in caller saved registers, then start it there rather
	// than at 1 where the callee saved registers are
	start := uint16(1)
	if nd.callerSaved {
		ig.dbg("%s: starting node %s in callee saved regs", ig.fn.Name, nd.val)
		start = ig.ra.savedStart
	}

	// first check for a colour that doesn't interfere with any moves, otherwise
	// take the first non-interfering colour
	for _, checkMoves := range []bool{true, false} {
		// find the lowest numbered colour that doesn't interfere
		for colour := start; ; colour++ {
			// if it doesn't interfere then
			if nd.checkIfInterferes(ig, colour) {
				continue
			}

			if checkMoves && nd.checkIfInterferesWithMove(ig, moves, colour) {
				continue
			}

			// choose the colour
			nd.colour = colour

			ig.dbg("%s: pick colour %d for %s checkMoves: %v", ig.fn.Name, nd.colour, nd, checkMoves)

			// keep track of the largest chosen colour
			if ig.maxColour < colour {
				ig.maxColour = colour
			}

			return
		}
	}
}
