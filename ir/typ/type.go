package typ

import "github.com/rj45/llbrew/sizes"

type Type uint32

func (t Type) Kind() Kind {
	return Kind(t >> 28)
}

func (t Type) index() int {
	return int(t & 0xfff_ffff)
}

func typeFor(kind Kind, index int) Type {
	return Type(uint32(kind)<<28 | uint32(index)&0xfff_ffff)
}

func (t Type) String() string {
	return t.string(make(map[Type]string))
}

func (t Type) string(refs map[Type]string) string {
	switch t.Kind() {
	case VoidKind:
		return "void"
	case FloatKind:
	case DoubleKind:
	case LabelKind:
	case IntegerKind:
		return t.Integer().String()
	case FunctionKind:
		return t.Function().string(refs)
	case StructKind:
		if ref, found := refs[t]; found {
			return ref
		}
		refs[t] = t.Struct().Reference()
		return t.Struct().string(refs)
	case ArrayKind:
	case PointerKind:
		return t.Pointer().string(refs)
	case VectorKind:
	case MetadataKind:
	case TokenKind:
	}
	return "todo"
}

func (t Type) ZeroValue() interface{} {
	switch t.Kind() {
	case VoidKind:
		return nil
	case FloatKind:
		return float32(0)
	case DoubleKind:
		return float64(0)
	case LabelKind:
		return ""
	case IntegerKind:
		return 0
	case FunctionKind:
		return nil
	case StructKind:
		elem := t.Struct().Elements
		v := make([]interface{}, len(t.Struct().Elements))
		for i, e := range elem {
			v[i] = e.ZeroValue()
		}
		return v
	case ArrayKind:
		return nil
	case PointerKind:
		return nil
	case VectorKind:
		return nil
	case MetadataKind:
		return nil
	case TokenKind:
		return nil
	}
	panic("unknown type zero value")
}

// Sizeof returns the size of the type in min-addressable-units (usually bytes).
func (t Type) SizeOf() int {
	units := sizes.MinAddressableBits()
	switch t.Kind() {
	case VoidKind:
		return 0
	case FloatKind:
		return 32 / units // todo: add to `sizes` so it's configurable
	case DoubleKind:
		return 64 / units // todo: add to `sizes` so it's configurable
	case LabelKind:
		return sizes.WordSize()
	case IntegerKind:
		return t.Integer().Bits() / units
	case FunctionKind:
		return sizes.WordSize()
	case StructKind:
		total := 0
		for _, st := range t.Struct().Elements {
			total += st.SizeOf()
		}
		return total
	case ArrayKind:
		panic("todo")
	case PointerKind:
		return sizes.WordSize() // todo: maybe word size and pointer size are different?
	case VectorKind:
		panic("todo")
	case MetadataKind:
		return 0
	case TokenKind:
		return 0
	}
	panic("unknown type size")
}

func StringType() Type {
	return PointerType(IntegerType(sizes.MinAddressableBits()), 0)
}
