package typ

import (
	"fmt"
	"strconv"
)

type Array struct {
	Element Type
	Count   int

	types *Types
}

var _ Type = &Array{}

func (arr *Array) Kind() Kind {
	return ArrayKind
}

func (arr *Array) SizeOf() int {
	return arr.Element.SizeOf() * arr.Count
}

func (arr *Array) String() string {
	return arr.string(make(map[Type]string))
}

func (arr *Array) GoString() string {
	return fmt.Sprintf("types.ArrayType(%s, %d)", arr.Element.GoString(), arr.Count)
}

func (arr *Array) ZeroValue() interface{} {
	zeros := make([]interface{}, arr.Count)
	zero := arr.Element.ZeroValue()
	for i := 0; i < len(zeros); i++ {
		zeros[i] = zero
	}
	return zeros
}

func (arr *Array) private() {}

func (arr Array) string(refs map[Type]string) string {
	return "[" + strconv.Itoa(arr.Count) + "]" + arr.types.string(arr.Element, refs)
}
