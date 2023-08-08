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
	return ctx.function[typ.index()]
}

func (ctx *Context) FunctionType(results []Type, params []Type, isVarArg bool) Type {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()

next:
	for index, fn := range ctx.function {
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
	ctx.function = append(ctx.function, Function{
		Results:  results,
		Params:   params,
		IsVarArg: isVarArg,
	})
	return typeFor(FunctionKind, len(ctx.function)-1)
}

func (t Type) Function() Function {
	return DefaultContext.Function(t)
}

func FunctionType(results []Type, params []Type, isVarArg bool) Type {
	return DefaultContext.FunctionType(results, params, isVarArg)
}

func (fn Function) String() string {
	strs := make([]string, len(fn.Params))
	for i, param := range fn.Params {
		strs[i] = param.String()
	}
	rstr := ""
	if len(fn.Results) == 1 {
		rstr = " " + fn.Results[0].String()
	} else {
		strs := make([]string, len(fn.Params))
		for i, param := range fn.Params {
			strs[i] = param.String()
		}
		rstr = " (" + strings.Join(strs, ",") + ")"
	}
	return "func(" + strings.Join(strs, ",") + ")" + rstr
}
