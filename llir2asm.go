package main

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/rj45/llir2asm/arch"
	"github.com/rj45/llir2asm/ir"
	"github.com/rj45/llir2asm/ir/op"
	"github.com/rj45/llir2asm/ir/typ"
	"github.com/rj45/llir2asm/xform"
	"tinygo.org/x/go-llvm"

	_ "github.com/rj45/llir2asm/arch/rj32"

	_ "github.com/rj45/llir2asm/xform/cleanup"
	_ "github.com/rj45/llir2asm/xform/elaboration"
	_ "github.com/rj45/llir2asm/xform/finishing"
	_ "github.com/rj45/llir2asm/xform/legalization"
	_ "github.com/rj45/llir2asm/xform/lowering"
	_ "github.com/rj45/llir2asm/xform/simplification"
)

func init() {
	llvm.InitializeAllTargets()
	llvm.InitializeAllTargetMCs()
	llvm.InitializeAllTargetInfos()
	llvm.InitializeAllAsmParsers()
	llvm.InitializeAllAsmPrinters()
}

func main() {
	ctx := llvm.NewContext()
	defer ctx.Dispose()

	optsize := 2
	optspeed := 2

	filename := os.Args[1]

	name := strings.TrimSuffix(path.Base(filename), path.Ext(filename))

	buf, err := llvm.NewMemoryBufferFromFile(filename)
	if err != nil {
		log.Fatalf("unable to read %s: %s", filename, err)
	}
	//defer buf.Dispose() // throws error???

	mod, err := ctx.ParseIR(buf)
	if err != nil {
		log.Fatalf("unable to parse LLVM IR %s: %s", filename, err)
	}
	defer mod.Dispose()

	for fn := mod.FirstFunction(); !fn.IsNil(); fn = llvm.NextFunction(fn) {
		// remove disabling optimizations
		fn.RemoveEnumFunctionAttribute(llvm.AttributeKindID("optnone"))

		// set any size optimizations
		if optsize >= 1 {
			fn.AddFunctionAttr(ctx.CreateEnumAttribute(llvm.AttributeKindID("optsize"), 0))
		}
		if optsize >= 2 {
			fn.AddFunctionAttr(ctx.CreateEnumAttribute(llvm.AttributeKindID("minsize"), 0))
		}
	}

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

	passBuilder := llvm.NewPassManagerBuilder()
	defer passBuilder.Dispose()

	passBuilder.SetOptLevel(optspeed)
	passBuilder.SetSizeLevel(optsize)

	{
		passManager := llvm.NewFunctionPassManagerForModule(mod)
		defer passManager.Dispose()
		passBuilder.PopulateFunc(passManager)

		passManager.InitializeFunc()
		for fn := mod.FirstFunction(); !fn.IsNil(); fn = llvm.NextFunction(fn) {
			passManager.RunFunc(fn)
		}
		passManager.FinalizeFunc()
	}

	{
		modPasses := llvm.NewPassManager()
		defer modPasses.Dispose()
		passBuilder.Populate(modPasses)
		modPasses.Run(mod)
	}

	err = llvm.VerifyModule(mod, llvm.PrintMessageAction)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(mod.String())

	prog := &ir.Program{}

	pkg := &ir.Package{
		Name: name,
	}

	prog.AddPackage(pkg)

	for glob := mod.FirstGlobal(); !glob.IsNil(); glob = llvm.NextGlobal(glob) {
		pkg.NewGlobal(glob.Name(), convertType(glob.Type()))
		// todo: set global value?
	}

	for fn := mod.FirstFunction(); !fn.IsNil(); fn = llvm.NextFunction(fn) {
		nfn := pkg.NewFunc(fn.Name(), convertFuncType(fn))

		nfn.Referenced = true

		blkmap := make(map[llvm.BasicBlock]*ir.Block)
		valuemap := make(map[llvm.Value]*ir.Value)
		instrmap := make(map[llvm.Value]*ir.Instr)

		for blk := fn.FirstBasicBlock(); !blk.IsNil(); blk = llvm.NextBasicBlock(blk) {

			nblk := nfn.NewBlock()
			blkmap[blk] = nblk

			nfn.InsertBlock(-1, nblk)

			for instr := blk.FirstInstruction(); !instr.IsNil(); instr = llvm.NextInstruction(instr) {
				op := translateOpcode(instr.InstructionOpcode())
				ninstr := nfn.NewInstr(op, convertType(instr.Type()))
				nblk.InsertInstr(-1, ninstr)
				instrmap[instr] = ninstr
				if ninstr.NumDefs() > 0 {
					valuemap[instr] = ninstr.Def(0)
				}
			}
		}

		for blk := fn.FirstBasicBlock(); !blk.IsNil(); blk = llvm.NextBasicBlock(blk) {
			// handle preds & succs
			nblk := blkmap[blk]
			bval := blk.AsValue()
			for i := 0; i < bval.IncomingCount(); i++ {
				inc := bval.IncomingBlock(i)
				pred := blkmap[inc]
				nblk.AddPred(pred)
				pred.AddSucc(nblk)
			}

			for instr := blk.FirstInstruction(); !instr.IsNil(); instr = llvm.NextInstruction(instr) {
				ninstr := instrmap[instr]
				for i := 0; i < instr.OperandsCount(); i++ {
					operand := instr.Operand(i)
					oval := valuemap[operand]
					if oval == nil {
						if operand.IsConstant() {
							if operand.Type().TypeKind() == llvm.IntegerTypeKind {
								ninstr.InsertArg(-1, nfn.ValueFor(convertType(operand.Type()), operand.SExtValue()))
							} else {
								panic("todo: other constant types")
							}
						} else {
							panic("todo: other operand types")
						}
					} else {
						ninstr.InsertArg(-1, oval)
					}
				}
			}
		}
	}

	arch.SetArch("rj32")
	for _, pkg := range prog.Packages() {
		for _, fn := range pkg.Funcs() {
			xform.Transform(xform.Elaboration, fn)
			xform.Transform(xform.Simplification, fn)
			xform.Transform(xform.Lowering, fn)
			xform.Transform(xform.Legalization, fn)
		}
	}

	prog.Emit(os.Stdout, ir.SSAString{})

	// mod.Dump()
}

func convertFuncType(fn llvm.Value) typ.Type {
	if fn.Type().TypeKind() == llvm.FunctionTypeKind {
		return convertType(fn.Type())
	}

	// if fn.AllocatedType().TypeKind() == llvm.FunctionTypeKind {
	// 	return convertType(fn.Type())
	// }

	// if fn.CalledFunctionType().TypeKind() == llvm.FunctionTypeKind {
	// 	return convertType(fn.Type())
	// }

	params := make([]typ.Type, 0, fn.ParamsCount())

	for param := fn.FirstParam(); !param.IsNil(); param = llvm.NextParam(param) {
		params = append(params, convertType(param.Type()))
	}

	for blk := fn.LastBasicBlock(); !blk.IsNil(); blk = llvm.PrevBasicBlock(blk) {
		for inst := blk.LastInstruction(); !inst.IsNil(); inst = llvm.PrevInstruction(inst) {
			if inst.Opcode() == llvm.Ret {
				if inst.Type().TypeKind() == llvm.VoidTypeKind {
					return typ.FunctionType(nil, params, false)
				}
				return typ.FunctionType([]typ.Type{convertType(inst.Type())}, params, false)
			}
		}
	}

	return typ.FunctionType(nil, params, false)
}

func convertType(t llvm.Type) typ.Type {
	if t.IsNil() {
		panic("nil type")
	}

	switch t.TypeKind() {
	case llvm.IntegerTypeKind:
		return typ.IntegerType(t.IntTypeWidth())
	case llvm.FunctionTypeKind:
		pt := t.ParamTypes()
		params := make([]typ.Type, len(pt))
		for i, pt := range pt {
			params[i] = convertType(pt)
		}
		return typ.FunctionType([]typ.Type{convertType(t.ReturnType())}, params, t.IsFunctionVarArg())
	case llvm.PointerTypeKind:
		return typ.PointerType(convertType(t.ElementType()), t.PointerAddressSpace())
	case llvm.VoidTypeKind:
		return typ.VoidType()
	default:
		log.Panicf("Unknown type: %#v (%s)", t, t.TypeKind().String())
		return 0
	}
}

var opcodeMap = [...]op.Op{
	llvm.Ret:         op.Ret,
	llvm.Br:          op.Br,
	llvm.Switch:      op.Switch,
	llvm.IndirectBr:  op.IndirectBr,
	llvm.Invoke:      op.Invoke,
	llvm.Unreachable: op.Unreachable,

	// Standard Binary Operators
	llvm.Add:  op.Add,
	llvm.FAdd: op.FAdd,
	llvm.Sub:  op.Sub,
	llvm.FSub: op.FSub,
	llvm.Mul:  op.Mul,
	llvm.FMul: op.FMul,
	llvm.UDiv: op.UDiv,
	llvm.SDiv: op.SDiv,
	llvm.FDiv: op.FDiv,
	llvm.URem: op.URem,
	llvm.SRem: op.SRem,
	llvm.FRem: op.FRem,

	// Logical Operators
	llvm.Shl:  op.Shl,
	llvm.LShr: op.LShr,
	llvm.AShr: op.AShr,
	llvm.And:  op.And,
	llvm.Or:   op.Or,
	llvm.Xor:  op.Xor,

	// Memory Operators
	llvm.Alloca:        op.Alloca,
	llvm.Load:          op.Load,
	llvm.Store:         op.Store,
	llvm.GetElementPtr: op.GetElementPtr,

	// Cast Operators
	llvm.Trunc:    op.Trunc,
	llvm.ZExt:     op.ZExt,
	llvm.SExt:     op.SExt,
	llvm.FPToUI:   op.FPToUI,
	llvm.FPToSI:   op.FPToSI,
	llvm.UIToFP:   op.UIToFP,
	llvm.SIToFP:   op.SIToFP,
	llvm.FPTrunc:  op.FPTrunc,
	llvm.FPExt:    op.FPExt,
	llvm.PtrToInt: op.PtrToInt,
	llvm.IntToPtr: op.IntToPtr,
	llvm.BitCast:  op.BitCast,

	// Other Operators
	llvm.ICmp:   op.ICmp,
	llvm.FCmp:   op.FCmp,
	llvm.PHI:    op.PHI,
	llvm.Call:   op.Call,
	llvm.Select: op.Select,

	// UserOp1
	// UserOp2
	llvm.VAArg:          op.VAArg,
	llvm.ExtractElement: op.ExtractElement,
	llvm.InsertElement:  op.InsertElement,
	llvm.ShuffleVector:  op.ShuffleVector,
	llvm.ExtractValue:   op.ExtractValue,
	llvm.InsertValue:    op.InsertValue,
}

func translateOpcode(opcode llvm.Opcode) ir.Op {
	op2 := opcodeMap[opcode]
	if op2 == op.Invalid {
		log.Panicf("bad opcode %d", opcode)
	}
	return op2
}
