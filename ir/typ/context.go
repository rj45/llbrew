package typ

import "sync"

type Context struct {
	// todo: add rwlock

	functions []Function
	pointers  []Pointer
	structs   []Struct

	lock sync.RWMutex
}

var DefaultContext *Context = &Context{}
