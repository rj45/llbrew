package typ

import "github.com/rj45/llir2asm/sizes"

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
	switch t.Kind() {
	case VoidKind:
		return "void"
	case FloatKind:
	case DoubleKind:
	case LabelKind:
	case IntegerKind:
		return t.Integer().String()
	case FunctionKind:
		return t.Function().String()
	case StructKind:
	case ArrayKind:
	case PointerKind:
		return t.Pointer().String()
	case VectorKind:
	case MetadataKind:
	case TokenKind:
	}
	return "todo"
}

func StringType() Type {
	return PointerType(IntegerType(sizes.MinAddressableBits()), 0)
}
