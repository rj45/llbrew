package translate

import (
	"log"

	"github.com/rj45/llbrew/ir/typ"
	"tinygo.org/x/go-llvm"
)

func translateType(t llvm.Type) typ.Type {
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
			params[i] = translateType(pt)
		}
		return typ.FunctionType([]typ.Type{translateType(t.ReturnType())}, params, t.IsFunctionVarArg())
	case llvm.StructTypeKind:
		se := t.StructElementTypes()
		elems := make([]typ.Type, len(se))
		for i, s := range se {
			elems[i] = translateType(s)
		}
		return typ.StructType(elems, t.IsStructPacked())
	case llvm.PointerTypeKind:
		return typ.PointerType(translateType(t.ElementType()), t.PointerAddressSpace())
	case llvm.VoidTypeKind:
		return typ.VoidType()
	case llvm.LabelTypeKind:
		return typ.LabelType()
	default:
		log.Panicf("Unknown type: %#v (%s)", t, t.TypeKind().String())
		return 0
	}
}

func translateFuncType(fn llvm.Value) typ.Type {
	if fn.Type().TypeKind() == llvm.FunctionTypeKind {
		return translateType(fn.Type())
	}

	params := make([]typ.Type, 0, fn.ParamsCount())

	for param := fn.FirstParam(); !param.IsNil(); param = llvm.NextParam(param) {
		params = append(params, translateType(param.Type()))
	}

	for blk := fn.LastBasicBlock(); !blk.IsNil(); blk = llvm.PrevBasicBlock(blk) {
		for inst := blk.LastInstruction(); !inst.IsNil(); inst = llvm.PrevInstruction(inst) {
			if inst.Opcode() == llvm.Ret {
				if inst.Type().TypeKind() == llvm.VoidTypeKind {
					return typ.FunctionType(nil, params, false)
				}
				return typ.FunctionType([]typ.Type{translateType(inst.Type())}, params, false)
			}
		}
	}

	return typ.FunctionType(nil, params, false)
}
