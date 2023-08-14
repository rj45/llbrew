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
	minAddressableBits = a.MinAddressableBits()
}

// sizes of basic types
var basicSizes = [17]byte{}

// size of min addressable unit in bits
var minAddressableBits = 0

func WordSize() int {
	return int(basicSizes[types.Uintptr])
}

func MinAddressableBits() int {
	return minAddressableBits
}
