package typ

import "strings"

type Function struct {
	Results  []Type
	Params   []Type
	IsVarArg bool
}

func (ctx *Context) Function(typ Type) Function {
	if typ.Kind() != FunctionKind {
		return Function{}
	}

	ctx.lock.RLock()
	defer ctx.lock.RUnlock()
	return ctx.functions[typ.index()]
}

func (ctx *Context) FunctionType(results []Type, params []Type, isVarArg bool) Type {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()

next:
	for index, fn := range ctx.functions {
		if len(fn.Results) != len(results) {
			continue
		}
		if len(fn.Params) != len(params) {
			continue
		}
		if fn.IsVarArg != isVarArg {
			continue
		}

		for i, result := range results {
			if fn.Results[i] != result {
				continue next
			}
		}

		for i, param := range params {
			if fn.Params[i] != param {
				continue next
			}
		}

		return typeFor(FunctionKind, index)
	}
	// todo: copy slices?
	ctx.functions = append(ctx.functions, Function{
		Results:  results,
		Params:   params,
		IsVarArg: isVarArg,
	})
	return typeFor(FunctionKind, len(ctx.functions)-1)
}

func (t Type) Function() Function {
	return DefaultContext.Function(t)
}

func FunctionType(results []Type, params []Type, isVarArg bool) Type {
	return DefaultContext.FunctionType(results, params, isVarArg)
}

func (fn Function) String() string {
	return fn.string(make(map[Type]string))
}

func (fn Function) string(refs map[Type]string) string {
	strs := make([]string, len(fn.Params))
	for i, param := range fn.Params {
		strs[i] = param.string(refs)
	}
	rstr := ""
	if len(fn.Results) == 1 {
		rstr = " " + fn.Results[0].string(refs)
	} else {
		strs := make([]string, len(fn.Params))
		for i, param := range fn.Params {
			strs[i] = param.string(refs)
		}
		rstr = " (" + strings.Join(strs, ",") + ")"
	}
	return "func(" + strings.Join(strs, ",") + ")" + rstr
}
