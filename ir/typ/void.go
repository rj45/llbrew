package typ

type Void struct{}

var _ Type = Void{}

func (v Void) Kind() Kind {
	return VoidKind
}

func (v Void) SizeOf() int {
	return 0
}

func (v Void) String() string {
	return "void"
}

func (v Void) GoString() string {
	return "types.VoidType()"
}

func (v Void) ZeroValue() interface{} {
	return 0
}

func (v Void) private() {}
