package rj32

import (
	"fmt"
	"strings"

	"github.com/rj45/llir2asm/ir"
)

func (cpuArch) Asm(op ir.Op, defs, args []string, emit func(string)) {
	switch op {
	case Load, Loadb:
		emit(fmt.Sprintf("%s %s, [%s, %s]", op, defs[0], args[0], args[1]))
		return
	case Store, Storeb:
		emit(fmt.Sprintf("%s [%s, %s], %s", op, args[0], args[1], args[2]))
		return
	case Return:
		emit("return")
		return
	case Call:
		emit("call " + args[0])
		return
	case IfEq, IfNe, IfLt, IfLe, IfGt, IfGe, IfUlt, IfUle, IfUge, IfUgt:
		emit(strings.Replace(op.String(), "_", ".", -1) + " " + args[0] + ", " + args[1])
		emit("    jump " + args[2])
		if len(args) > 3 {
			emit("jump " + args[3])
		}
	default:
		if op.ClobbersArg() {
			args = args[1:]
		}
		if len(defs) > 0 && len(args) > 0 {
			emit(op.String() + " " + strings.Join(defs, ", ") + ", " + strings.Join(args, ", "))
			return
		} else if len(defs) > 0 {
			emit(op.String() + " " + strings.Join(defs, ", "))
			return
		}
		emit(op.String() + " " + strings.Join(args, ", "))
		return
	}
}
