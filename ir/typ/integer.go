package typ

import (
	"strconv"

	"github.com/rj45/llbrew/sizes"
)

type Integer uint8

var _ Type = Integer(0)

func (i Integer) Kind() Kind {
	return IntegerKind
}

func (i Integer) SizeOf() int {
	if int(i) < sizes.MinAddressableBits() {
		return 1
	}
	return int(i) / sizes.MinAddressableBits()
}

func (i Integer) String() string {
	return "i" + strconv.Itoa(int(i))
}

func (i Integer) ZeroValue() interface{} {
	return 0
}

func (i Integer) private() {}

func (i Integer) Bits() int {
	return int(i)
}
