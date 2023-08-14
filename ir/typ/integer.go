package typ

import (
	"strconv"

	"github.com/rj45/llbrew/sizes"
)

type Integer int

func (ctx *Context) Integer(typ Type) Integer {
	return Integer(typ.index())
}

func (ctx *Context) IntegerType(bits int) Type {
	return typeFor(IntegerKind, bits)
}

func (t Type) Integer() Integer {
	return DefaultContext.Integer(t)
}

func IntegerType(bits int) Type {
	return DefaultContext.IntegerType(bits)
}

func (i Integer) String() string {
	return "i" + strconv.Itoa(int(i))
}

func (i Integer) Bits() int {
	return int(i)
}

func IntegerWordType() Type {
	return DefaultContext.IntegerType(sizes.WordSize() * sizes.MinAddressableBits())
}
