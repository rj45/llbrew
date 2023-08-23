package typ

import (
	"slices"

	"github.com/rj45/llbrew/sizes"
)

type Types struct {
	functions []*Function
	pointers  []*Pointer
	structs   []*Struct
	arrays    []*Array
}

func (ctx *Types) IntegerType(bits int) Type {
	return Integer(bits)
}

func (ctx *Types) IntegerWordType() Type {
	return Integer(sizes.WordSize())
}

func (ctx *Types) FunctionType(results []Type, params []Type, isVarArg bool) Type {
next:
	for _, fn := range ctx.functions {
		if len(fn.Results) != len(results) {
			continue
		}
		if len(fn.Params) != len(params) {
			continue
		}
		if fn.IsVarArg != isVarArg {
			continue
		}

		for i, result := range results {
			if fn.Results[i] != result {
				continue next
			}
		}

		for i, param := range params {
			if fn.Params[i] != param {
				continue next
			}
		}

		return fn
	}
	ctx.functions = append(ctx.functions, &Function{
		Results:  slices.Clone(results),
		Params:   slices.Clone(params),
		IsVarArg: isVarArg,
		types:    ctx,
	})
	return ctx.functions[len(ctx.functions)-1]
}

func (ctx *Types) PointerType(elem Type, addrspace int) Type {
	for _, value := range ctx.pointers {
		if value.Element == elem && value.AddrSpace == addrspace {
			return value
		}
	}
	ctx.pointers = append(ctx.pointers, &Pointer{elem, addrspace, ctx})
	return ctx.pointers[len(ctx.pointers)-1]
}

func (ctx *Types) VoidPointer() Type {
	return ctx.PointerType(ctx.VoidType(), 0)
}

func (ctx *Types) StructType(name string, elems []Type, packed bool) Type {
	for _, value := range ctx.structs {
		match := true
		for i, elem := range value.Elements {
			if elem != elems[i] {
				match = false
				break
			}
		}
		if match && value.Packed == packed && value.Name == name {
			return value
		}
	}
	ctx.structs = append(ctx.structs, &Struct{name, slices.Clone(elems), packed, ctx})
	return ctx.structs[len(ctx.structs)-1]
}

func (ctx *Types) PartialStructType(name string, packed bool) Type {
	for _, value := range ctx.structs {
		if value.Name == name {
			return value
		}
	}

	ctx.structs = append(ctx.structs, &Struct{name, nil, packed, ctx})
	return ctx.structs[len(ctx.structs)-1]
}

func (ctx *Types) CompleteStructType(t Type, elems []Type) {
	st := t.(*Struct)

	if st.Elements != nil {
		if !slices.Equal(elems, st.Elements) {
			panic("two different structs of the same name!")
		}
		return
	}

	st.Elements = slices.Clone(elems)
}

func (ctx *Types) ArrayType(elem Type, count int) Type {
	for _, value := range ctx.arrays {
		if value.Element == elem && value.Count == count {
			return value
		}
	}
	ctx.arrays = append(ctx.arrays, &Array{elem, count, ctx})
	return ctx.arrays[len(ctx.arrays)-1]
}

func (ctx *Types) StringType() Type {
	return ctx.PointerType(Integer(sizes.MinAddressableBits()), 0)
}

func (ctx *Types) VoidType() Type {
	return Void{}
}

func (ctx *Types) string(t Type, refs map[Type]string) string {
	switch t.Kind() {
	case VoidKind:
		return "void"
	case FloatKind:
	case DoubleKind:
	case IntegerKind:
		return t.(Integer).String()
	case FunctionKind:
		return t.(*Function).string(refs)
	case StructKind:
		if ref, found := refs[t]; found {
			return ref
		}
		refs[t] = t.(*Struct).Reference()
		return t.(*Struct).string(refs)
	case ArrayKind:
	case PointerKind:
		return t.(*Pointer).string(refs)
	}
	return "todo"
}
