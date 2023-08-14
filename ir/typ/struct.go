package typ

import (
	"strings"
)

type Struct struct {
	Elements []Type
	Packed   bool
}

func (ctx *Context) Struct(typ Type) Struct {
	ctx.lock.RLock()
	defer ctx.lock.RUnlock()

	if typ.Kind() != StructKind {
		return Struct{}
	}
	return ctx.structs[typ.index()]
}

func (ctx *Context) StructType(elems []Type, packed bool) Type {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()

	for index, value := range ctx.structs {
		match := true
		for i, elem := range value.Elements {
			if elem != elems[i] {
				match = false
				break
			}
		}
		if match && value.Packed == packed {
			return typeFor(StructKind, index)
		}
	}
	// todo copy slice?
	ctx.structs = append(ctx.structs, Struct{elems, packed})
	return typeFor(StructKind, len(ctx.structs)-1)
}

func (t Type) Struct() Struct {
	return DefaultContext.Struct(t)
}

func StructType(elems []Type, packed bool) Type {
	return DefaultContext.StructType(elems, packed)
}

func (s Struct) String() string {
	strs := make([]string, len(s.Elements))
	for i, elem := range s.Elements {
		strs[i] = elem.String()
	}
	packed := ""
	if s.Packed {
		packed = "packed"
	}
	return "{" + strings.Join(strs, ",") + "}" + packed
}

func (s Struct) OffsetOf(element int) int {
	total := 0
	for i, t := range s.Elements {
		if element == i {
			return total
		}
		total += t.SizeOf()
	}
	panic("invalid element index")
}
