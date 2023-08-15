package translate

import (
	"log"
	"strings"

	"github.com/rj45/llbrew/ir/typ"
	"tinygo.org/x/go-llvm"
)

func (trans *translator) translateType(t llvm.Type) typ.Type {
	return trans.translatePartialType(t, make(map[llvm.Type]typ.Type))
}

func (trans *translator) translatePartialType(t llvm.Type, incomplete map[llvm.Type]typ.Type) typ.Type {
	if t.IsNil() {
		panic("nil type")
	}

	if it, ok := incomplete[t]; ok {
		return it
	}

	switch t.TypeKind() {
	case llvm.IntegerTypeKind:
		return trans.types.IntegerType(t.IntTypeWidth())
	case llvm.FunctionTypeKind:
		pt := t.ParamTypes()
		params := make([]typ.Type, len(pt))
		for i, pt := range pt {
			params[i] = trans.translatePartialType(pt, incomplete)
		}
		return trans.types.FunctionType([]typ.Type{trans.translatePartialType(t.ReturnType(), incomplete)}, params, t.IsFunctionVarArg())
	case llvm.StructTypeKind:
		name := strings.TrimPrefix(t.StructName(), "struct.")
		ntype := trans.types.PartialStructType(name, t.IsStructPacked())
		incomplete[t] = ntype
		se := t.StructElementTypes()
		elems := make([]typ.Type, len(se))
		for i, s := range se {
			elems[i] = trans.translatePartialType(s, incomplete)
		}
		trans.types.CompleteStructType(ntype, elems)
		return ntype
	case llvm.PointerTypeKind:
		return trans.types.PointerType(trans.translatePartialType(t.ElementType(), incomplete), t.PointerAddressSpace())
	case llvm.VoidTypeKind:
		return trans.types.VoidType()
	case llvm.LabelTypeKind:
		panic("should not be")
	default:
		log.Panicf("Unknown type: %#v (%s)", t, t.TypeKind().String())
		return nil
	}
}

func (trans *translator) translateFuncType(fn llvm.Value) *typ.Function {
	incomplete := make(map[llvm.Type]typ.Type)

	if fn.Type().TypeKind() == llvm.FunctionTypeKind {
		return trans.translatePartialType(fn.Type(), incomplete).(*typ.Function)
	}

	params := make([]typ.Type, 0, fn.ParamsCount())

	for param := fn.FirstParam(); !param.IsNil(); param = llvm.NextParam(param) {
		params = append(params, trans.translatePartialType(param.Type(), incomplete))
	}

	for blk := fn.LastBasicBlock(); !blk.IsNil(); blk = llvm.PrevBasicBlock(blk) {
		for inst := blk.LastInstruction(); !inst.IsNil(); inst = llvm.PrevInstruction(inst) {
			if inst.Opcode() == llvm.Ret {
				if inst.Type().TypeKind() == llvm.VoidTypeKind {
					return trans.types.FunctionType(nil, params, false).(*typ.Function)
				}
				return trans.types.FunctionType([]typ.Type{trans.translatePartialType(inst.Type(), incomplete)}, params, false).(*typ.Function)
			}
		}
	}

	return trans.types.FunctionType(nil, params, false).(*typ.Function)
}
