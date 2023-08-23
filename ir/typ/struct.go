package typ

import (
	"fmt"
	"strings"
)

type Struct struct {
	Name     string
	Elements []Type
	Packed   bool

	types *Types
}

var _ Type = &Struct{}

func (st *Struct) Kind() Kind {
	return StructKind
}

func (st *Struct) SizeOf() int {
	ttl := 0
	for _, e := range st.Elements {
		ttl += e.SizeOf()
	}
	return ttl
}

func (st *Struct) String() string {
	return st.string(make(map[Type]string))
}

func (st *Struct) GoString() string {
	elems := make([]string, len(st.Elements))
	for i, elem := range st.Elements {
		elems[i] = elem.GoString()
	}
	return fmt.Sprintf("types.StructType(%q, []Type{%s}, %v)", st.Name, strings.Join(elems, ", "), st.Packed)
}

func (st *Struct) ZeroValue() interface{} {
	v := make([]interface{}, len(st.Elements))
	for i, e := range st.Elements {
		v[i] = e.ZeroValue()
	}
	return v
}

func (st *Struct) Reference() string {
	return "struct " + st.Name
}

func (st *Struct) OffsetOf(element int) int {
	total := 0
	for i, t := range st.Elements {
		if element == i {
			return total
		}
		total += t.SizeOf()
	}
	panic("invalid element index")
}

func (st *Struct) private() {}

func (st *Struct) string(refs map[Type]string) string {
	strs := make([]string, len(st.Elements))
	for i, elem := range st.Elements {
		strs[i] = st.types.string(elem, refs)
	}
	packed := ""
	if st.Packed {
		packed = "packed "
	}
	return "struct " + packed + st.Name + " {" + strings.Join(strs, ",") + "}"
}
