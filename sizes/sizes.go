// Copyright (c) 2021 rj45 (github.com/rj45), MIT Licensed, see LICENSE.

// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sizes

import "go/types"

type Arch interface {
	BasicSizes() [17]byte
	RuneSize() int
	MinAddressableBits() int
}

func SetArch(a Arch) {
	basicSizes = a.BasicSizes()
	runeSize = a.RuneSize()
	minAddressableBits = a.MinAddressableBits()
}

// sizes of basic types
var basicSizes = [17]byte{}

// size of rune
var runeSize = 0

// size of min addressable unit in bits
var minAddressableBits = 0

func WordSize() int {
	return int(basicSizes[types.Uintptr])
}

func MinAddressableBits() int {
	return minAddressableBits
}

func Sizeof(T types.Type) int {
	switch t := T.Underlying().(type) {
	case *types.Basic:
		k := t.Kind()
		if k == types.Int32 {
			if t.Name() == "rune" {
				return int(runeSize)
			}
		}
		if int(k) < len(basicSizes) {
			if s := basicSizes[k]; s > 0 {
				return int(s)
			}
		}
		if k == types.String {
			return int(basicSizes[types.Uintptr]) * 2
		}
	case *types.Array:
		n := t.Len()
		if n <= 0 {
			return 0
		}
		// n > 0
		z := Sizeof(t.Elem())
		return z * int(n)
	case *types.Slice:
		return int(basicSizes[types.Uintptr]) * 3
	case *types.Struct:
		fields := Fieldsof(t)
		n := len(fields)
		if n == 0 {
			return 0
		}
		offsets := Offsetsof(fields)
		return offsets[n-1] + Sizeof(fields[n-1].Type())

	case *types.Interface:
		return int(basicSizes[types.Uintptr]) * 2
	}

	return int(basicSizes[types.Int]) // catch-all
}

func Offsetsof(fields []*types.Var) []int {
	offsets := make([]int, len(fields))
	var o int
	for i, f := range fields {
		offsets[i] = o
		o += Sizeof(f.Type())
	}
	return offsets
}

func Fieldsof(t *types.Struct) []*types.Var {
	n := t.NumFields()
	if n == 0 {
		return nil
	}
	fields := make([]*types.Var, n)
	for i := 0; i < n; i++ {
		fields[i] = t.Field(i)
	}
	return fields
}
