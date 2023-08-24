package xform

type Tag uint8

const (
	Invalid Tag = iota
	HasFramePointer

	// ...

	NumTags
)

var activeTags []bool
