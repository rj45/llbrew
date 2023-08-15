package typ

import (
	"strings"

	"github.com/rj45/llbrew/sizes"
)

type Function struct {
	Results  []Type
	Params   []Type
	IsVarArg bool

	types *Types
}

func (fn *Function) Kind() Kind {
	return FunctionKind
}

func (fn *Function) SizeOf() int {
	return sizes.PointerSize()
}

func (fn *Function) String() string {
	return fn.string(make(map[Type]string))
}

func (fn *Function) ZeroValue() interface{} {
	return 0
}

func (fn *Function) private() {}

func (fn *Function) string(refs map[Type]string) string {
	strs := make([]string, len(fn.Params))
	for i, param := range fn.Params {
		strs[i] = fn.types.string(param, refs)
	}
	rstr := ""
	if len(fn.Results) == 1 {
		rstr = " " + fn.types.string(fn.Results[0], refs)
	} else {
		strs := make([]string, len(fn.Params))
		for i, param := range fn.Params {
			strs[i] = fn.types.string(param, refs)
		}
		rstr = " (" + strings.Join(strs, ",") + ")"
	}
	return "func(" + strings.Join(strs, ",") + ")" + rstr
}
