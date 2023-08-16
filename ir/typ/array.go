package typ

import (
	"strconv"
)

type Array struct {
	Element Type
	Count   int

	types *Types
}

var _ Type = &Array{}

func (ptr *Array) Kind() Kind {
	return ArrayKind
}

func (ptr *Array) SizeOf() int {
	return ptr.Element.SizeOf() * ptr.Count
}

func (ptr *Array) String() string {
	return ptr.string(make(map[Type]string))
}

func (ptr *Array) ZeroValue() interface{} {
	zeros := make([]interface{}, ptr.Count)
	zero := ptr.Element.ZeroValue()
	for i := 0; i < len(zeros); i++ {
		zeros[i] = zero
	}
	return zeros
}

func (ptr *Array) private() {}

func (ptr Array) string(refs map[Type]string) string {
	return "[" + strconv.Itoa(ptr.Count) + "]" + ptr.types.string(ptr.Element, refs)
}
