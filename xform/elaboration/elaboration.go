package elaboration

import (
	"github.com/rj45/llbrew/xform"
)

//go:generate go run github.com/rj45/llbrew/cmd/rewritegen -o rewrite.go -pkg elaboration -fn rewrite rewrite.rules

var _ = xform.Register(rewrite,
	xform.OnlyPass(xform.Elaboration),
)
