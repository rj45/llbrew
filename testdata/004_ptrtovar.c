// Copyright (c) 2018 Andrew Chambers
// MIT Licensed
// From github.com/c-testsuite/c-testsuite

int x;

int main() {
	int *p;

	x = 4;
	p = &x;
	*p = 42;

	return x;
}