package typ

import "strconv"

type Pointer struct {
	Element   Type
	AddrSpace int
}

func (ctx *Context) Pointer(typ Type) Pointer {
	ctx.lock.RLock()
	defer ctx.lock.RUnlock()

	if typ.Kind() != PointerKind {
		return Pointer{}
	}
	return ctx.pointer[typ.index()]
}

func (ctx *Context) PointerType(elem Type, addrspace int) Type {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()

	for index, value := range ctx.pointer {
		if value.Element == elem && value.AddrSpace == addrspace {
			return typeFor(PointerKind, index)
		}
	}
	ctx.pointer = append(ctx.pointer, Pointer{elem, addrspace})
	return typeFor(PointerKind, len(ctx.pointer)-1)
}

func (t Type) Pointer() Pointer {
	return DefaultContext.Pointer(t)
}

func PointerType(elem Type, addrspace int) Type {
	return DefaultContext.PointerType(elem, addrspace)
}

func (ptr Pointer) String() string {
	space := ""
	if ptr.AddrSpace != 0 {
		space = "(" + strconv.Itoa(ptr.AddrSpace) + ")"
	}
	return "*" + space + ptr.Element.String()
}
