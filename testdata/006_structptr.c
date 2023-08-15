// Copyright (c) 2018 Andrew Chambers
// MIT Licensed
// From github.com/c-testsuite/c-testsuite

struct S { int x; int y; };

void setS(struct S *p) {
    p->x = 9;
    p->y = 36;
}

int main() {

	struct S s;

	setS(&s);
	return s.y + s.x - 3;
}
