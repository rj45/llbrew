package ir

// EliminateDeadCode eliminates dead code until the code stops
// changing.
func (fn *Func) EliminateDeadCode() {
	alive := map[*Value]bool{}
	for _, v := range fn.alive {
		alive[v] = true
	}

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
				if alive[def] {
					hasUse = true
				}
			}
			// todo: how to eliminate dead stores?
			if !hasUse && !instr.IsSink() && !instr.IsCall() {
				it.Remove()
				changed = true
			}
		}
	}
}
