package typ

type Label struct{}

func (ctx *Context) IsLabel(typ Type) bool {
	return typ.Kind() == LabelKind
}

func (ctx *Context) LabelType() Type {
	return typeFor(LabelKind, 0)
}

func (t Type) IsLabel() bool {
	return t.Kind() == LabelKind
}

func LabelType() Type {
	return DefaultContext.LabelType()
}
