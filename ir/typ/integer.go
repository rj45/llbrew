package typ

import (
	"strconv"

	"github.com/rj45/llir2asm/sizes"
)

type Integer struct {
	Bits int
}

func (ctx *Context) Integer(typ Type) Integer {
	if typ.Kind() != IntegerKind {
		return Integer{}
	}
	return ctx.integer[typ.index()]
}

func (ctx *Context) IntegerType(bits int) Type {
	for index, value := range ctx.integer {
		if value.Bits == bits {
			return typeFor(IntegerKind, index)
		}
	}
	ctx.integer = append(ctx.integer, Integer{bits})
	return typeFor(IntegerKind, len(ctx.integer)-1)
}

func (t Type) Integer() Integer {
	return DefaultContext.Integer(t)
}

func IntegerType(bits int) Type {
	return DefaultContext.IntegerType(bits)
}

func (i Integer) String() string {
	return "i" + strconv.Itoa(i.Bits)
}

func IntegerWordType() Type {
	return DefaultContext.IntegerType(sizes.WordSize() * sizes.MinAddressableBits())
}
