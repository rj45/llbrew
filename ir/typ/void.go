package typ

type Void struct{}

func (ctx *Context) IsVoid(typ Type) bool {
	return typ.Kind() == VoidKind
}

func (ctx *Context) VoidType() Type {
	return typeFor(VoidKind, 0)
}

func (t Type) IsVoid() bool {
	return t.Kind() == VoidKind
}

func VoidType() Type {
	return DefaultContext.VoidType()
}
