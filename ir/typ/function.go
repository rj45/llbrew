package typ

import (
	"fmt"
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

func (fn *Function) GoString() string {
	params := make([]string, len(fn.Params))
	for i, elem := range fn.Params {
		params[i] = elem.GoString()
	}
	results := make([]string, len(fn.Results))
	for i, elem := range fn.Results {
		results[i] = elem.GoString()
	}
	return fmt.Sprintf("types.FunctionType([]typ.Type{%s}, []typ.Type{%s}, %v)",
		strings.Join(results, ", "), strings.Join(params, ", "), fn.IsVarArg)
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
