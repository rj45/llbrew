package typ

type Kind uint8

const (
	VoidKind Kind = iota
	FloatKind
	DoubleKind
	IntegerKind
	FunctionKind
	StructKind
	ArrayKind
	PointerKind
)
