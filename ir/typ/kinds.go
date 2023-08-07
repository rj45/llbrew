package typ

type Kind uint8

const (
	VoidKind Kind = iota
	FloatKind
	DoubleKind
	LabelKind
	IntegerKind
	FunctionKind
	StructKind
	ArrayKind
	PointerKind
	VectorKind
	MetadataKind
	TokenKind
)
