package typ

import (
	"strings"
)

type Struct struct {
	Name     string
	Elements []Type
	Packed   bool

	types *Types
}

var _ Type = &Struct{}

func (ptr *Struct) Kind() Kind {
	return StructKind
}

func (ptr *Struct) SizeOf() int {
	ttl := 0
	for _, e := range ptr.Elements {
		ttl += e.SizeOf()
	}
	return ttl
}

func (ptr *Struct) String() string {
	return ptr.string(make(map[Type]string))
}

func (ptr *Struct) ZeroValue() interface{} {
	v := make([]interface{}, len(ptr.Elements))
	for i, e := range ptr.Elements {
		v[i] = e.ZeroValue()
	}
	return v
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

func (ptr *Struct) private() {}

func (s Struct) string(refs map[Type]string) string {
	strs := make([]string, len(s.Elements))
	for i, elem := range s.Elements {
		strs[i] = s.types.string(elem, refs)
	}
	packed := ""
	if s.Packed {
		packed = "packed "
	}
	return "struct " + packed + s.Name + " {" + strings.Join(strs, ",") + "}"
}
