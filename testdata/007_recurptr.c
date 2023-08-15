// Copyright (c) 2018 Andrew Chambers
// MIT Licensed
// From github.com/c-testsuite/c-testsuite

struct S { struct S *p; int x; } s;

int main() {
	s.x = 42;
	s.p = &s;
	return s.p->p->p->p->p->x;
}
