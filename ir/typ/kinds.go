package typ

type Kind uint8

const (
	VoidKind Kind = iota
	FloatKind
	DoubleKind
	IntegerKind
	LabelKind
	FunctionKind
	StructKind
	ArrayKind
	PointerKind
	VectorKind
	MetadataKind
	TokenKind
)
