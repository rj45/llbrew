package typ

type Type interface {
	Kind() Kind
	String() string
	ZeroValue() interface{}
	SizeOf() int
	GoString() string

	private()
}
