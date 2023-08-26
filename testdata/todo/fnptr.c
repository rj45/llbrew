// Copyright (c) 2018 Andrew Chambers
// MIT Licensed
// From github.com/c-testsuite/c-testsuite

struct S {
	int	(*fptr)();
};

int foo() {
	return 42;
}

int main() {
	struct S v;

	v.fptr = foo;
	return v.fptr();
}