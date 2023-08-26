// Copyright (c) 2018 Andrew Chambers
// MIT Licensed
// From github.com/c-testsuite/c-testsuite

struct foo {
	int i, j, k;
};

int f1(struct foo f, struct foo *p, int n) {
	if (f.i != p->i)
		return 0;
	return p->j + n;
}

int main(void) {
	struct foo f;

	f.i = f.j = 1;

	return f1(f, &f, 2);;
}