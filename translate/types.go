package translate

import (
	"log"
	"strings"

	"github.com/rj45/llbrew/ir/typ"
	"tinygo.org/x/go-llvm"
)

func translateType(t llvm.Type) typ.Type {
	return translatePartialType(t, make(map[llvm.Type]typ.Type))
}

func translatePartialType(t llvm.Type, incomplete map[llvm.Type]typ.Type) typ.Type {
	if t.IsNil() {
		panic("nil type")
	}

	if it, ok := incomplete[t]; ok {
		return it
	}

	switch t.TypeKind() {
	case llvm.IntegerTypeKind:
		return typ.IntegerType(t.IntTypeWidth())
	case llvm.FunctionTypeKind:
		pt := t.ParamTypes()
		params := make([]typ.Type, len(pt))
		for i, pt := range pt {
			params[i] = translatePartialType(pt, incomplete)
		}
		return typ.FunctionType([]typ.Type{translatePartialType(t.ReturnType(), incomplete)}, params, t.IsFunctionVarArg())
	case llvm.StructTypeKind:
		name := strings.TrimPrefix(t.StructName(), "struct.")
		ntype := typ.PartialStructType(name, t.IsStructPacked())
		incomplete[t] = ntype
		se := t.StructElementTypes()
		elems := make([]typ.Type, len(se))
		for i, s := range se {
			elems[i] = translatePartialType(s, incomplete)
		}
		typ.CompleteStructType(ntype, elems)
		return ntype
	case llvm.PointerTypeKind:
		return typ.PointerType(translatePartialType(t.ElementType(), incomplete), t.PointerAddressSpace())
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
	incomplete := make(map[llvm.Type]typ.Type)

	if fn.Type().TypeKind() == llvm.FunctionTypeKind {
		return translatePartialType(fn.Type(), incomplete)
	}

	params := make([]typ.Type, 0, fn.ParamsCount())

	for param := fn.FirstParam(); !param.IsNil(); param = llvm.NextParam(param) {
		params = append(params, translatePartialType(param.Type(), incomplete))
	}

	for blk := fn.LastBasicBlock(); !blk.IsNil(); blk = llvm.PrevBasicBlock(blk) {
		for inst := blk.LastInstruction(); !inst.IsNil(); inst = llvm.PrevInstruction(inst) {
			if inst.Opcode() == llvm.Ret {
				if inst.Type().TypeKind() == llvm.VoidTypeKind {
					return typ.FunctionType(nil, params, false)
				}
				return typ.FunctionType([]typ.Type{translatePartialType(inst.Type(), incomplete)}, params, false)
			}
		}
	}

	return typ.FunctionType(nil, params, false)
}
