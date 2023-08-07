package typ

type Context struct {
	// todo: add rwlock

	integer  []Integer
	function []Function
	pointer  []Pointer
}

var DefaultContext *Context = &Context{}
