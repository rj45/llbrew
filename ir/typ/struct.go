package typ

import (
	"strings"

	"slices"
)

type Struct struct {
	Name     string
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

func (ctx *Context) PartialStructType(name string, packed bool) Type {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()

	for index, value := range ctx.structs {
		if value.Name == name {
			return typeFor(StructKind, index)
		}
	}

	ctx.structs = append(ctx.structs, Struct{name, nil, packed})
	return typeFor(StructKind, len(ctx.structs)-1)
}

func PartialStructType(name string, packed bool) Type {
	return DefaultContext.PartialStructType(name, packed)
}

func (ctx *Context) CompleteStructType(st Type, elems []Type) {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()

	if ctx.structs[st.index()].Elements != nil {
		if !slices.Equal(elems, ctx.structs[st.index()].Elements) {
			panic("two different structs of the same name!")
		}
		return
	}

	ctx.structs[st.index()].Elements = slices.Clone(elems)
}

func CompleteStructType(st Type, elems []Type) {
	DefaultContext.CompleteStructType(st, elems)
}

func (ctx *Context) StructType(name string, elems []Type, packed bool) Type {
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
		if match && value.Packed == packed && value.Name == name {
			return typeFor(StructKind, index)
		}
	}
	ctx.structs = append(ctx.structs, Struct{name, slices.Clone(elems), packed})
	return typeFor(StructKind, len(ctx.structs)-1)
}

func (t Type) Struct() Struct {
	return DefaultContext.Struct(t)
}

func StructType(name string, elems []Type, packed bool) Type {
	return DefaultContext.StructType(name, elems, packed)
}

func (s Struct) String() string {
	return s.string(make(map[Type]string))
}

func (s Struct) string(refs map[Type]string) string {
	strs := make([]string, len(s.Elements))
	for i, elem := range s.Elements {
		strs[i] = elem.string(refs)
	}
	packed := ""
	if s.Packed {
		packed = "packed "
	}
	return "struct " + packed + s.Name + " {" + strings.Join(strs, ",") + "}"
}

func (s Struct) Reference() string {
	return "struct " + s.Name
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
