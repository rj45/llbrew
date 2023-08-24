package ir

// EliminateDeadCode eliminates dead code until the code stops
// changing.
func (fn *Func) EliminateDeadCode() {
	changed := true
	for changed {
		changed = false
		for it := fn.InstrIter(); it.HasNext(); it.Next() {
			instr := it.Instr()
			hasUse := false
			for _, def := range instr.defs {
				if len(def.uses) > 0 {
					hasUse = true
					break
				}
				if def.InReg() {
					reg := def.Reg()
					if reg.IsSpecialReg() || reg.IsSavedReg() {
						// todo: figure out how to differenciate logue regs
						hasUse = true
						break
					}
				}
			}
			// todo: how to eliminate dead stores?
			if !hasUse && !instr.IsSink() {
				it.Remove()
				it.Prev()
				changed = true
			}
		}
	}
}
