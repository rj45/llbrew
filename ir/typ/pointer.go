package typ

import (
	"strconv"

	"github.com/rj45/llbrew/sizes"
)

type Pointer struct {
	Element   Type
	AddrSpace int

	types *Types
}

var _ Type = &Pointer{}

func (ptr *Pointer) Kind() Kind {
	return PointerKind
}

func (ptr *Pointer) SizeOf() int {
	return sizes.PointerSize()
}

func (ptr *Pointer) String() string {
	return ptr.string(make(map[Type]string))
}

func (ptr *Pointer) ZeroValue() interface{} {
	return 0
}

func (ptr *Pointer) private() {}

func (ptr Pointer) string(refs map[Type]string) string {
	space := ""
	if ptr.AddrSpace != 0 {
		space = "(" + strconv.Itoa(ptr.AddrSpace) + ")"
	}
	return "*" + space + ptr.types.string(ptr.Element, refs)
}
