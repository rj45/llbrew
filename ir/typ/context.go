package typ

import "sync"

type Context struct {
	// todo: add rwlock

	function []Function
	pointer  []Pointer

	lock sync.RWMutex
}

var DefaultContext *Context = &Context{}
