package compile

import (
	"fmt"
	"log"
	"strings"

	"tinygo.org/x/go-llvm"
)

func (c *Compiler) loadIR() error {
	c.ctx = llvm.NewContext()
	c.initLevel = 1

	buf, err := llvm.NewMemoryBufferFromFile(c.Filename)
	if err != nil {
		c.dispose()
		return fmt.Errorf("unable to read %s: %w", c.Filename, err)
	}
	//defer buf.Dispose() // throws error???

	mod, err := c.ctx.ParseIR(buf)
	if err != nil {
		c.dispose()
		return fmt.Errorf("unable to parse LLVM IR %s: %w", c.Filename, err)
	}
	c.mod = mod
	c.initLevel = 2

	return nil
}

func (c *Compiler) fixFunctionAttributes() {
	for fn := c.mod.FirstFunction(); !fn.IsNil(); fn = llvm.NextFunction(fn) {
		// remove disabling optimizations
		fn.RemoveEnumFunctionAttribute(llvm.AttributeKindID("optnone"))

		// set any size optimizations
		if c.OptSize >= 1 {
			fn.AddFunctionAttr(c.ctx.CreateEnumAttribute(llvm.AttributeKindID("optsize"), 0))
		}
		if c.OptSize >= 2 {
			fn.AddFunctionAttr(c.ctx.CreateEnumAttribute(llvm.AttributeKindID("minsize"), 0))
		}
	}
}

func (c *Compiler) optimize() {
	// {
	// 	pm := llvm.NewPassManager()
	// 	defer pm.Dispose()
	// 	pm.AddGlobalDCEPass()
	// 	pm.AddGlobalOptimizerPass()
	// 	pm.AddIPSCCPPass()
	// 	pm.AddInstructionCombiningPass()
	// 	pm.AddAggressiveDCEPass()
	// 	pm.AddFunctionAttrsPass()
	// 	pm.AddGlobalDCEPass()
	// 	if !pm.Run(mod) {
	// 		fmt.Println("initial pass did nothing!")
	// 	}
	// }

	// passBuilder := llvm.NewPassManagerBuilder()
	// defer passBuilder.Dispose()

	// passBuilder.SetOptLevel(c.OptSpeed)
	// passBuilder.SetSizeLevel(c.OptSize)
	// passBuilder.SetDisableUnrollLoops(true)

	// {
	// 	passManager := llvm.NewFunctionPassManagerForModule(c.mod)
	// 	defer passManager.Dispose()
	// 	passBuilder.PopulateFunc(passManager)

	// 	passManager.InitializeFunc()
	// 	for fn := c.mod.FirstFunction(); !fn.IsNil(); fn = llvm.NextFunction(fn) {
	// 		passManager.RunFunc(fn)
	// 	}
	// 	passManager.FinalizeFunc()
	// }

	// {
	// 	modPasses := llvm.NewPassManager()
	// 	defer modPasses.Dispose()
	// 	passBuilder.Populate(modPasses)
	// 	modPasses.Run(c.mod)
	// }

	pbo := llvm.NewPassBuilderOptions()
	defer pbo.Dispose()

	pbo.SetLoopInterleaving(false)
	pbo.SetLoopVectorization(false)
	pbo.SetSLPVectorization(false)

	defaultPass := "default<O0>"

	if c.OptSize > 0 {
		switch c.OptSize {
		case 1:
			defaultPass = "default<Os>"
		case 2:
			defaultPass = "default<Oz>"
		}
	} else {
		switch c.OptSpeed {
		case 1:
			defaultPass = "default<O1>"
		case 2:
			defaultPass = "default<O2>"
		case 3:
			defaultPass = "default<O3>"
		}
	}

	targ, err := llvm.GetTargetFromTriple(c.mod.Target())
	if err != nil {
		log.Fatal(err)
	}

	tm := targ.CreateTargetMachine(c.mod.Target(), "", "", llvm.CodeGenLevelAggressive, llvm.RelocStatic, llvm.CodeModelSmall)

	passes := []string{
		defaultPass,
		"break-crit-edges",
	}

	err = c.mod.RunPasses(strings.Join(passes, ","), tm, pbo)
	if err != nil {
		log.Fatal(err)
	}

	// err = llvm.VerifyModule(mod, llvm.PrintMessageAction)
	// if err != nil {
	// 	log.Fatal(err)
	// }

}

func (c *Compiler) splitCriticalEdges() {
	for fn := c.mod.FirstFunction(); !fn.IsNil(); fn = llvm.NextFunction(fn) {

		var splits [][2]llvm.BasicBlock
		for bb := fn.FirstBasicBlock(); !bb.IsNil(); bb = llvm.NextBasicBlock(bb) {
			lastInstr := bb.LastInstruction()
			if !isBranchInstruction(lastInstr) {
				continue
			}

			for i := 0; i < lastInstr.OperandsCount(); i++ {
				operand := lastInstr.Operand(i)

				if operand.Type().TypeKind() != llvm.LabelTypeKind {
					continue
				}
				successor := operand.AsBasicBlock()

				if hasCriticalEdge(bb, successor) {
					splits = append(splits, [2]llvm.BasicBlock{bb, successor})
				}
			}
		}

		for _, split := range splits {
			splitCriticalEdge(split[0], split[1])
		}

		// llvm.VerifyFunction(fn, llvm.PrintMessageAction)
	}

}

func isBranchInstruction(instr llvm.Value) bool {
	return instr.InstructionOpcode() == llvm.Br
}

func hasCriticalEdge(source, dest llvm.BasicBlock) bool {
	lastInstr := source.LastInstruction()

	operCount := lastInstr.OperandsCount()
	if lastInstr.OperandsCount() > 0 && lastInstr.Operand(0).Type().TypeKind() != llvm.LabelTypeKind {
		operCount--
	}

	numSuccessors := operCount
	numPredecessors := len(getPredecessors(dest))

	if dest == dest.Parent().FirstBasicBlock() {
		numPredecessors++
	}

	return numSuccessors > 1 && numPredecessors > 1
}

func splitCriticalEdge(source, dest llvm.BasicBlock) {
	lastInstr := source.LastInstruction()

	// Create a new block
	builder := llvm.NewBuilder()
	defer builder.Dispose()

	newBlock := llvm.AddBasicBlock(source.Parent(), dest.AsValue().Name()+".split")
	newBlock.MoveAfter(source)

	// Redirect the edge
	builder.SetInsertPointAtEnd(newBlock)
	builder.CreateBr(dest)
	for i := 0; i < lastInstr.OperandsCount(); i++ {
		if lastInstr.Operand(i).AsBasicBlock() == dest {
			lastInstr.SetOperand(i, newBlock.AsValue())
		}
	}

	// Update PHI nodes
	for instr := dest.FirstInstruction(); !instr.IsNil(); instr = llvm.NextInstruction(instr) {
		phi := instr.IsAPHINode()

		if phi.IsNil() {
			break
		}

		builder.SetInsertPoint(dest, phi)
		newPhi := builder.CreatePHI(phi.Type(), "")
		incVal := make([]llvm.Value, phi.IncomingCount())
		incBB := make([]llvm.BasicBlock, phi.IncomingCount())
		for i := 0; i < phi.IncomingCount(); i++ {
			incVal[i] = phi.IncomingValue(i)
			incBB[i] = phi.IncomingBlock(i)

			if phi.IncomingBlock(i) == source {
				incBB[i] = newBlock
			}
		}
		newPhi.AddIncoming(incVal, incBB)
		phi.ReplaceAllUsesWith(newPhi)
		phi.RemoveFromParentAsInstruction()
		instr = newPhi
	}
}

func getPredecessors(bb llvm.BasicBlock) []llvm.BasicBlock {
	var predecessors []llvm.BasicBlock
	parentFunc := bb.Parent()
	for iter := parentFunc.FirstBasicBlock(); !iter.IsNil(); iter = llvm.NextBasicBlock(iter) {
		terminator := iter.LastInstruction()
		if terminator.InstructionOpcode() == llvm.Br {
			for i := 0; i < terminator.OperandsCount(); i++ {
				if terminator.Operand(i).Type().TypeKind() == llvm.LabelTypeKind &&
					terminator.Operand(i).AsBasicBlock() == bb {
					predecessors = append(predecessors, iter)
					break
				}
			}
		}
	}
	return predecessors
}
