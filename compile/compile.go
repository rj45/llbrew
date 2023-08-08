package compile

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/rj45/llir2asm/arch"
	"github.com/rj45/llir2asm/html"
	"github.com/rj45/llir2asm/ir"
	"github.com/rj45/llir2asm/ir/op"
	"github.com/rj45/llir2asm/ir/typ"
	"github.com/rj45/llir2asm/regalloc"
	"github.com/rj45/llir2asm/regalloc/verify"
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

type Compiler struct {
	Filename string

	OptSpeed int
	OptSize  int

	DumpLL  io.WriteCloser
	DumpIR  io.WriteCloser
	DumpSSA string

	ctx       llvm.Context
	mod       llvm.Module
	initLevel int

	prog *ir.Program
	pkg  *ir.Package

	fn       *ir.Func
	blkmap   map[llvm.BasicBlock]*ir.Block
	valuemap map[llvm.Value]*ir.Value
	instrmap map[llvm.Value]*ir.Instr
}

func (c *Compiler) dispose() {
	// yuck :-/
	if c.initLevel > 1 {
		c.mod.Dispose()
	}

	if c.initLevel > 0 {
		c.ctx.Dispose()
	}
}

func (c *Compiler) Compile(filename string) (*ir.Program, error) {

	c.Filename = filename

	err := c.loadIR()
	if err != nil {
		return nil, err
	}
	defer c.dispose()

	// some function attributes prevent optimizations from running
	c.fixFunctionAttributes()

	// run the LLVM optimizer
	c.optimize()

	fmt.Fprint(c.DumpLL, c.mod.String())

	// convert the LLVM program to our own IR
	c.initProgram()
	c.convertGlobals()
	c.convertFunctions()

	// convert the program into assembly
	err = c.transformProgram()
	if err != nil {
		return nil, err
	}

	c.prog.Emit(c.DumpIR, ir.SSAString{})

	return c.prog, nil
}

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

	passBuilder := llvm.NewPassManagerBuilder()
	defer passBuilder.Dispose()

	passBuilder.SetOptLevel(c.OptSpeed)
	passBuilder.SetSizeLevel(c.OptSize)

	{
		passManager := llvm.NewFunctionPassManagerForModule(c.mod)
		defer passManager.Dispose()
		passBuilder.PopulateFunc(passManager)

		passManager.InitializeFunc()
		for fn := c.mod.FirstFunction(); !fn.IsNil(); fn = llvm.NextFunction(fn) {
			passManager.RunFunc(fn)
		}
		passManager.FinalizeFunc()
	}

	{
		modPasses := llvm.NewPassManager()
		defer modPasses.Dispose()
		passBuilder.Populate(modPasses)
		modPasses.Run(c.mod)
	}

	// err = llvm.VerifyModule(mod, llvm.PrintMessageAction)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func (c *Compiler) initProgram() {
	c.prog = &ir.Program{}
	c.pkg = &ir.Package{
		Name: "main",
	}
	c.prog.AddPackage(c.pkg)
}

func (c *Compiler) convertGlobals() {
	for glob := c.mod.FirstGlobal(); !glob.IsNil(); glob = llvm.NextGlobal(glob) {
		c.pkg.NewGlobal(glob.Name(), convertType(glob.Type()))
		// todo: set global value?
	}
}

func (c *Compiler) convertFunctions() {
	for fn := c.mod.FirstFunction(); !fn.IsNil(); fn = llvm.NextFunction(fn) {
		c.convertFunction(fn)
	}
}

func (c *Compiler) convertFunction(fn llvm.Value) {
	c.fn = c.pkg.NewFunc(fn.Name(), convertFuncType(fn))

	c.fn.Referenced = true

	c.blkmap = make(map[llvm.BasicBlock]*ir.Block)
	c.valuemap = make(map[llvm.Value]*ir.Value)
	c.instrmap = make(map[llvm.Value]*ir.Instr)

	c.convertInstructions(fn)
	c.convertOperands(fn)
}

func (c *Compiler) convertInstructions(fn llvm.Value) {
	first := true
	for blk := fn.FirstBasicBlock(); !blk.IsNil(); blk = llvm.NextBasicBlock(blk) {
		nblk := c.fn.NewBlock()
		c.blkmap[blk] = nblk

		if first {
			first = false
			for param := fn.FirstParam(); !param.IsNil(); param = llvm.NextParam(param) {
				pval := c.fn.NewValue(convertType(param.Type()))
				nblk.AddDef(pval)
				c.valuemap[param] = pval
			}
		}

		c.fn.InsertBlock(-1, nblk)

		for instr := blk.FirstInstruction(); !instr.IsNil(); instr = llvm.NextInstruction(instr) {
			op := translateOpcode(instr.InstructionOpcode())
			ninstr := c.fn.NewInstr(op, convertType(instr.Type()))
			nblk.InsertInstr(-1, ninstr)
			c.instrmap[instr] = ninstr
			if ninstr.NumDefs() > 0 {
				c.valuemap[instr] = ninstr.Def(0)
			}
		}
	}
}

func (c *Compiler) convertOperands(fn llvm.Value) {
	for blk := fn.FirstBasicBlock(); !blk.IsNil(); blk = llvm.NextBasicBlock(blk) {
		// handle preds & succs
		nblk := c.blkmap[blk]
		bval := blk.AsValue()
		for i := 0; i < bval.IncomingCount(); i++ {
			inc := bval.IncomingBlock(i)
			pred := c.blkmap[inc]
			nblk.AddPred(pred)
			pred.AddSucc(nblk)
		}

		for instr := blk.FirstInstruction(); !instr.IsNil(); instr = llvm.NextInstruction(instr) {
			ninstr := c.instrmap[instr]
			for i := 0; i < instr.OperandsCount(); i++ {
				operand := instr.Operand(i)
				oval := c.valuemap[operand]
				if oval == nil {
					if operand.IsConstant() {
						opertyp := operand.Type()
						ntyp := convertType(opertyp)
						switch opertyp.TypeKind() {
						case llvm.IntegerTypeKind:
							ninstr.InsertArg(-1, c.fn.ValueFor(ntyp, operand.SExtValue()))
						case llvm.FunctionTypeKind:
							panic(operand.Name())
						case llvm.PointerTypeKind:
							if ntyp.Pointer().Element.Kind() == typ.FunctionKind {
								otherfn := c.pkg.Func(operand.Name())
								ninstr.InsertArg(0, c.fn.ValueFor(ntyp.Pointer().Element, otherfn))
							}
						default:
							log.Panicf("todo: other constant types: %s", ntyp)
						}

					} else {

						log.Panicf("todo: other operand types %d", operand.IntrinsicID())
					}
				} else {
					ninstr.InsertArg(-1, oval)
				}
			}
		}
	}
}

func (c *Compiler) transformProgram() error {
	arch.SetArch("rj32")

	for _, pkg := range c.prog.Packages() {
		for _, fn := range pkg.Funcs() {
			var w dumper = nopDumper{}
			if c.DumpSSA != "" && strings.Contains(fn.FullName, c.DumpSSA) {
				w = html.NewHTMLWriter("ssa.html", fn)
			}
			defer w.Close()

			w.WritePhase("initial", "initial")

			xform.Transform(xform.Elaboration, fn)
			w.WritePhase("elaboration", "elaboration")

			xform.Transform(xform.Simplification, fn)
			w.WritePhase("simplification", "simplification")

			xform.Transform(xform.Lowering, fn)
			w.WritePhase("lowering", "lowering")

			xform.Transform(xform.Legalization, fn)
			w.WritePhase("legalization", "legalization")

			ra := regalloc.NewRegAlloc(fn)
			err := ra.Allocate()
			// if *debug {
			// 	regalloc.WriteGraphvizCFG(ra)
			// 	regalloc.DumpLivenessChart(ra)
			// 	regalloc.WriteGraphvizInterferenceGraph(ra)
			// 	regalloc.WriteGraphvizLivenessGraph(ra)
			// }
			if err != nil {
				return fmt.Errorf("register allocation failed: %w", err)
			}
			errs := verify.Verify(fn)
			for _, err := range errs {
				log.Printf("verification error: %s\n", err)
			}
			if len(errs) > 0 {
				log.Fatal("verification failed")
			}
			w.WritePhase("regalloc", "regalloc")

			xform.Transform(xform.CleanUp, fn)
			w.WritePhase("cleanup", "cleanup")

			xform.Transform(xform.Finishing, fn)
			w.WritePhase("finishing", "finishing")
		}
	}
	return nil
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
